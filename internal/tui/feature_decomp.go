package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"shotgun.dev/odek/internal/decomposer"
)

const (
	blinkInterval = 500 * time.Millisecond
	// blinkHoldAfterMove is how long the cursor stays steadily "on" after the
	// user moves the selection. Makes rapid navigation visually stable.
	blinkHoldAfterMove = 500 * time.Millisecond
)

type blinkTickMsg struct{}

func blinkTick() tea.Cmd {
	return tea.Tick(blinkInterval, func(time.Time) tea.Msg {
		return blinkTickMsg{}
	})
}

type decompKanjiTickMsg struct{}

func decompKanjiTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg {
		return decompKanjiTickMsg{}
	})
}

var (
	mockSep   = lipgloss.Color("#3e3e3e")
	mockHot   = lipgloss.Color("#f74e82")
	mockFocus = lipgloss.Color("#9d7cf3")
)

const (
	decompMinColW    = 14 // hard floor for a navigable column
	decompColW       = 22 // nominal navigable column width before growth
	decompMinRightW  = 24 // detail pane floor
	decompMaxVisCols = 5  // max navigable columns onscreen; extras horizontally shift
)

const emptyTopCopy = "Describe your feature in the chat. This pane updates when you change scope."

type featureDecompModel struct {
	width       int
	height      int
	vp          viewport.Model
	selPath     []string // fully-qualified path selected in each open navigable column
	focusedCol  int      // index into selPath
	colScroll   int      // leftmost visible navigable column
	inputActive bool
	active      bool
	blinkOn     bool
	inSplit     bool
	steadyUntil time.Time
	input       textinput.Model

	sess        *decomposer.Session
	state       *decomposeState
	snap        decomposer.Snapshot
	decomposing bool // mirrored from state, read under the same lock
	kanjiOffset int
}

func newFeatureDecompModel(width, height int, sess *decomposer.Session, state *decomposeState) featureDecompModel {
	ti := textinput.New()
	ti.Prompt = "> "
	s := ti.Styles()
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)
	s.Focused.Text = lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)
	s.Blurred.Prompt = lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)
	s.Blurred.Text = lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)
	ti.SetStyles(s)

	m := featureDecompModel{
		width:  width,
		height: height,
		vp:     viewport.New(),
		input:  ti,
		active: true,
		sess:   sess,
		state:  state,
	}
	m.resize(width, height)
	return m
}

func (m *featureDecompModel) SetActive(active bool) {
	if m.active == active {
		return
	}
	m.active = active
	if active {
		// Start the blink in the "off" phase so the first visible change
		// after activation is dot → invisible, not dot → (same) bright pink.
		m.blinkOn = false
		m.steadyUntil = time.Time{}
	}
	m.refreshColumns()
}

func (m *featureDecompModel) SetInSplit(v bool) {
	m.inSplit = v
}

// decompBottomChromeRows is the number of rows consumed below the viewport:
// blank (above vp) + input + blank + help = 4. No trailing blank below help
// because the chat pane's body overflows by 1 row so its own trailing blank
// is clipped; keeping the decomp help on the last row lines them up.
const decompBottomChromeRows = 4

func (m *featureDecompModel) resize(w, h int) {
	m.width = w
	m.height = h
	if w <= 0 || h <= 0 {
		return
	}
	topH := h / 3
	bottomH := max(h-topH-decompBottomChromeRows, 1)
	m.vp.SetWidth(w)
	m.vp.SetHeight(bottomH)
	m.refreshColumns()
}

// refreshColumns pulls the latest session snapshot (via the shared state
// pointer so post-transition /decompose runs are picked up) and re-renders
// the column viewport.
func (m *featureDecompModel) refreshColumns() {
	if m.state != nil {
		m.state.mu.Lock()
		m.sess = m.state.session
		m.decomposing = m.state.decomposing
		m.state.mu.Unlock()
	}
	if m.sess != nil {
		m.snap = m.sess.Snapshot()
	} else {
		m.snap = decomposer.Snapshot{}
	}
	m.normalizeSelection()
	m.vp.SetContent(buildColumns(m.vp.Width(), m.vp.Height(), m.selPath, m.focusedCol, m.colScroll, m.active, m.blinkOn, m.decomposing, m.snap))
}

// parentOfCol returns the snapshot key whose ChildrenByName entry feeds
// column colIdx. Column 0's parent is the synthetic "root"; deeper
// columns inherit the selection of the column to their left.
func (m *featureDecompModel) parentOfCol(colIdx int) string {
	if colIdx <= 0 {
		return "root"
	}
	if colIdx-1 >= len(m.selPath) {
		return ""
	}
	return m.selPath[colIdx-1]
}

// normalizeSelection keeps selPath consistent with the current snapshot.
// It truncates stale entries and clamps focusedCol/colScroll into range.
// On initial selection (when selPath is empty) it auto-extends one level
// deeper so the first item's children appear in column 1 immediately —
// the user sees the hierarchy without having to press right.
func (m *featureDecompModel) normalizeSelection() {
	rootKids := m.snap.ChildrenByName["root"]
	if len(rootKids) == 0 {
		m.selPath = nil
		m.focusedCol = 0
		m.colScroll = 0
		return
	}
	freshSelection := len(m.selPath) == 0
	if freshSelection {
		m.selPath = []string{rootKids[0]}
	}

	// Validate each entry against its parent's current child list; trim
	// the cascade at the first entry that no longer belongs.
	for i := 0; i < len(m.selPath); i++ {
		parent := "root"
		if i > 0 {
			parent = m.selPath[i-1]
		}
		kids := m.snap.ChildrenByName[parent]
		if len(kids) == 0 {
			m.selPath = m.selPath[:i]
			break
		}
		found := false
		for _, k := range kids {
			if k == m.selPath[i] {
				found = true
				break
			}
		}
		if !found {
			m.selPath = append(m.selPath[:i], kids[0])
		}
	}

	if len(m.selPath) == 0 {
		m.selPath = []string{rootKids[0]}
		freshSelection = true
	}

	// Auto-extend one level so the first item shows its children on
	// initial render. Only on a fresh selection — we don't want to
	// re-extend every time the user moves up/down in column 0, because
	// that would clobber any deeper selection they've drilled into.
	if freshSelection && len(m.selPath) == 1 {
		first := m.selPath[0]
		if kids := m.snap.ChildrenByName[first]; len(kids) > 0 {
			m.selPath = append(m.selPath, kids[0])
		}
	}

	if m.focusedCol < 0 {
		m.focusedCol = 0
	}
	if m.focusedCol >= len(m.selPath) {
		m.focusedCol = len(m.selPath) - 1
	}
	m.ensureColVisible()
}

// ensureColVisible adjusts colScroll so focusedCol is within the visible
// window [colScroll, colScroll+decompMaxVisCols).
func (m *featureDecompModel) ensureColVisible() {
	if m.focusedCol < m.colScroll {
		m.colScroll = m.focusedCol
	}
	if m.focusedCol >= m.colScroll+decompMaxVisCols {
		m.colScroll = m.focusedCol - decompMaxVisCols + 1
	}
	if m.colScroll < 0 {
		m.colScroll = 0
	}
	maxScroll := max(len(m.selPath)-decompMaxVisCols, 0)
	if m.colScroll > maxScroll {
		m.colScroll = maxScroll
	}
}

// steadyCursor holds the focused-column selection bullet visibly "on"
// for a short grace period after a navigation keystroke.
func (m *featureDecompModel) steadyCursor() {
	m.blinkOn = true
	m.steadyUntil = time.Now().Add(blinkHoldAfterMove)
}

// moveSelection shifts the selection in the focused column by delta,
// truncates deeper selections, and auto-extends back to a leaf.
func (m *featureDecompModel) moveSelection(delta int) {
	if m.focusedCol >= len(m.selPath) {
		return
	}
	parent := m.parentOfCol(m.focusedCol)
	kids := m.snap.ChildrenByName[parent]
	if len(kids) == 0 {
		return
	}
	cur := m.selPath[m.focusedCol]
	idx := -1
	for i, k := range kids {
		if k == cur {
			idx = i
			break
		}
	}
	next := idx + delta
	if next < 0 || next >= len(kids) {
		return
	}
	m.selPath = append(m.selPath[:m.focusedCol:m.focusedCol], kids[next])
	m.normalizeSelection()
	m.steadyCursor()
	m.refreshColumns()
}

func (m featureDecompModel) Init() tea.Cmd { return tea.Batch(blinkTick(), decompKanjiTick()) }

func (m featureDecompModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil
	case blinkTickMsg:
		if m.active {
			if time.Now().Before(m.steadyUntil) {
				// Hold the cursor on while the user is navigating.
				m.blinkOn = true
			} else {
				m.blinkOn = !m.blinkOn
			}
			m.refreshColumns()
		}
		return m, blinkTick()
	case decompKanjiTickMsg:
		if m.state.WorkInProgress() {
			m.kanjiOffset += 2
		}
		return m, decompKanjiTick()
	case kanjiTickMsg:
		// Chat pane's ticker reaches us via the split's broadcast —
		// consume it so it doesn't bleed into the viewport component.
		return m, nil
	case expansionEventMsg:
		if msg.ok {
			m.refreshColumns()
		}
		return m, nil
	case decomposeDoneMsg:
		// Chat just installed a new session in shared state; re-snapshot
		// so the right pane shows the fresh top-level runes instead of
		// whatever it was displaying before (empty state, a prior
		// decomposition, etc.).
		m.refreshColumns()
		return m, nil
	case decomposeStartedMsg:
		// The user just sent a message and the decomposer is running.
		// Re-snapshot so we pick up state.decomposing=true and render
		// the loading indicator.
		m.refreshColumns()
		return m, nil
	case tea.KeyPressMsg:
		if m.inputActive {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter", "esc":
				m.inputActive = false
				m.input.SetValue("")
				m.input.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "space":
			m.inputActive = true
			return m, m.input.Focus()
		case "left", "h":
			if m.focusedCol > 0 {
				m.focusedCol--
				m.steadyCursor()
				m.ensureColVisible()
				m.refreshColumns()
			}
			return m, nil
		case "right", "l":
			if m.focusedCol+1 < len(m.selPath) {
				m.focusedCol++
				m.steadyCursor()
				m.ensureColVisible()
				m.refreshColumns()
			} else if m.focusedCol < len(m.selPath) {
				cur := m.selPath[m.focusedCol]
				kids := m.snap.ChildrenByName[cur]
				if len(kids) > 0 {
					m.selPath = append(m.selPath, kids[0])
					m.focusedCol++
					m.normalizeSelection()
					m.steadyCursor()
					m.refreshColumns()
				}
			}
			return m, nil
		case "up", "k":
			m.moveSelection(-1)
			return m, nil
		case "down", "j":
			m.moveSelection(+1)
			return m, nil
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	cmds = append(cmds, cmd)
	if m.inputActive {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m featureDecompModel) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true
	v.BackgroundColor = bgMain

	if m.width <= 0 || m.height <= 0 {
		return v
	}

	topH := m.height / 3
	top := renderDecompTop(m.width, topH, m.focusedCol*4+m.kanjiOffset, m.snap)
	blank := lipgloss.NewStyle().Background(bgMain).Render(strings.Repeat(" ", m.width))
	bottom := m.vp.View()

	lastLine := blank
	if m.inputActive {
		rendered := m.input.View()
		pad := max(m.width-lipgloss.Width(rendered), 0)
		lastLine = rendered + lipgloss.NewStyle().Background(bgMain).Render(strings.Repeat(" ", pad))
	}

	help := renderDecompHelp(m.width, m.inputActive, m.inSplit)

	v.Content = lipgloss.JoinVertical(lipgloss.Left,
		top,
		blank,
		bottom,
		lastLine,
		blank,
		help,
	)
	return v
}

func renderDecompHelp(width int, inputActive, showTabSwitch bool) string {
	keyStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(fgBody).Background(bgMain)
	sepStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
	padStyle := lipgloss.NewStyle().Background(bgMain)

	type binding struct{ key, desc string }
	var bindings []binding
	if inputActive {
		bindings = []binding{
			{"enter", "submit"},
			{"esc", "cancel"},
			{"ctrl+c", "quit"},
		}
	} else {
		bindings = []binding{
			{"↑/↓", "navigate"},
			{"←/→", "navigate"},
		}
		if showTabSwitch {
			bindings = append(bindings, binding{"tab", "chat"})
		}
	}

	var b strings.Builder
	for i, bind := range bindings {
		if i > 0 {
			b.WriteString(sepStyle.Render("  •  "))
		}
		b.WriteString(keyStyle.Render(bind.key))
		b.WriteString(padStyle.Render(" "))
		b.WriteString(descStyle.Render(bind.desc))
	}
	content := b.String()
	pad := max(width-lipgloss.Width(content)-2, 0)
	return padStyle.Render(" ") + content + padStyle.Render(strings.Repeat(" ", pad)+" ")
}

func renderDecompTop(width, height, kanjiOffset int, snap decomposer.Snapshot) string {
	const logoText = "ODEK "
	logoCells := len(logoText)
	logoStyled := lipgloss.NewStyle().
		Foreground(accent).
		Background(bgMain).
		Bold(true).
		Render(logoText)

	bgStyle := lipgloss.NewStyle().Background(bgMain)
	textStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)
	blankLine := bgStyle.Render(strings.Repeat(" ", width))

	var logoRow strings.Builder
	logoRow.WriteString(logoStyled)
	if logoCells < width {
		logoRow.WriteString(bgStyle.Render(strings.Repeat(" ", width-logoCells)))
	}

	renderTextLine := func(s string) string {
		styled := textStyle.Render(s)
		if pad := width - len(s); pad > 0 {
			styled += bgStyle.Render(strings.Repeat(" ", pad))
		}
		return styled
	}

	textMax := max(min(width/2, 70), 20)
	topCopy := strings.TrimSpace(snap.Summary)
	if topCopy == "" {
		topCopy = snap.Requirement
	}
	if !snap.HasSession || topCopy == "" {
		topCopy = emptyTopCopy
	}
	rawLines := wrapDecompText(topCopy, textMax)

	runesText := fmt.Sprintf("%d runes", snap.TotalRunes)
	switch snap.Phase {
	case decomposer.PhaseContract:
		runesText = "designing contract…"
	case decomposer.PhaseExtraction:
		runesText = fmt.Sprintf("extracting runes… %d bytes", snap.ExtractionBytes)
	case "error":
		runesText = fmt.Sprintf("error: %s", snap.ErrorMsg)
	}

	lines := make([]string, 0, height)
	appendLine := func(s string) {
		if len(lines) < height {
			lines = append(lines, s)
		}
	}

	appendLine(logoRow.String())
	appendLine(blankLine)
	appendLine(renderKanjiLine(width, 2, kanjiOffset))
	appendLine(renderKanjiLine(width, 3, -kanjiOffset))
	appendLine(blankLine)
	for _, raw := range rawLines {
		appendLine(renderTextLine(raw))
	}
	appendLine(blankLine)
	appendLine(renderTextLine(runesText))
	for len(lines) < height {
		lines = append(lines, blankLine)
	}
	return strings.Join(lines, "\n")
}

func wrapDecompText(text string, width int) []string {
	words := strings.Fields(text)
	var lines []string
	var cur strings.Builder
	curLen := 0
	for _, w := range words {
		wLen := len(w)
		switch {
		case curLen == 0:
			cur.WriteString(w)
			curLen = wLen
		case curLen+1+wLen <= width:
			cur.WriteByte(' ')
			cur.WriteString(w)
			curLen += 1 + wLen
		default:
			lines = append(lines, cur.String())
			cur.Reset()
			cur.WriteString(w)
			curLen = wLen
		}
	}
	if curLen > 0 {
		lines = append(lines, cur.String())
	}
	return lines
}

// leafTag returns a short tag shown in the detail pane's right-aligned
// slot. The 2-pass pipeline produces the whole tree at once, so there is
// no per-rune lifecycle — the only useful tag is whether the rune is a
// leaf (no children).
func leafTag(isLeaf bool) string {
	if isLeaf {
		return "leaf"
	}
	return ""
}

// statusGlyph returns the bullet for a rune row in the column list.
// The glyph is always a small pink dot; only the style varies with
// selection state (blink / focused cursor / unfocused rose).
func statusGlyph(selected, active, blinkOn bool) (string, lipgloss.Style) {
	base := "• "
	switch {
	case selected && active && blinkOn:
		return base, lipgloss.NewStyle().Foreground(mockHot).Background(bgMain).Bold(true)
	case selected && active && !blinkOn:
		return base, lipgloss.NewStyle().Foreground(bgMain).Background(bgMain)
	case selected:
		return base, lipgloss.NewStyle().Foreground(lipgloss.Color("#9a3050")).Background(bgMain)
	}
	return base, lipgloss.NewStyle().Foreground(accent).Background(bgMain)
}

func renderRuneInfo(name string, r decomposer.Rune, children []string, maxW int) string {
	summaryW := min(max(maxW-4, 20), 72)
	summaryStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain).Italic(true).Width(summaryW)
	headingStyle := lipgloss.NewStyle().Foreground(fgBody).Background(bgMain).Bold(true)
	sigStyle := lipgloss.NewStyle().Foreground(accent).Background(bgMain)
	plusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Background(bgMain)
	minusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Background(bgMain)
	questionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Background(bgMain).Italic(true)
	depStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Background(bgMain)
	bodyStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)

	var lines []string
	if r.Description != "" {
		lines = append(lines, summaryStyle.Render(r.Description))
		lines = append(lines, "")
	}
	if sig := decomposer.NormalizeFunctionSig(r.FunctionSig); sig != "" {
		lines = append(lines, sigStyle.Render("fn "+sig))
	}
	if len(r.PositiveTests) > 0 || len(r.NegativeTests) > 0 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, headingStyle.Render("Behavior"))
		for _, p := range r.PositiveTests {
			lines = append(lines, plusStyle.Render("  + ")+bodyStyle.Render(p))
		}
		for _, mn := range r.NegativeTests {
			lines = append(lines, minusStyle.Render("  - ")+bodyStyle.Render(mn))
		}
	}
	for _, q := range r.Assumptions {
		lines = append(lines, questionStyle.Render("? "+q))
	}
	if len(children) > 0 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, headingStyle.Render("Helpers"))
		for _, d := range children {
			lines = append(lines, depStyle.Render("  -> ")+bodyStyle.Render(d))
		}
	}
	if len(r.Dependencies) > 0 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, headingStyle.Render("Dependencies"))
		for _, d := range r.Dependencies {
			lines = append(lines, depStyle.Render("  -> ")+bodyStyle.Render(d))
		}
	}
	_ = name
	return strings.Join(lines, "\n")
}

func buildColumns(innerW, innerH int, selPath []string, focusedCol, colScroll int, active, blinkOn, decomposing bool, snap decomposer.Snapshot) string {
	const sepWidth = 3

	textStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)
	dimStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
	ruleFocusedStyle := lipgloss.NewStyle().Foreground(mockFocus).Background(bgMain)
	iconStyle := lipgloss.NewStyle().Foreground(mockHot).Background(bgMain).Bold(true)
	tagStyle := lipgloss.NewStyle().Foreground(fgBody).Background(bgMain)
	placeholderStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain).Italic(true)
	bodyDimStyle := lipgloss.NewStyle().Foreground(fgDim).Background(bgMain).Italic(true)
	bgSpace := lipgloss.NewStyle().Background(bgMain)

	titled := func(name, tag string, w int, focused bool) string {
		contentW := max(w-2, 1)
		nameSeg := iconStyle.Render("◆ ") + textStyle.Render(name)
		header := nameSeg
		if tag != "" {
			tagSeg := tagStyle.Render("# " + tag)
			pad := max(contentW-lipgloss.Width(nameSeg)-lipgloss.Width(tagSeg), 1)
			spacer := bgSpace.Render(strings.Repeat(" ", pad))
			header = nameSeg + spacer + tagSeg
		}
		ruleStyle := dimStyle
		if focused {
			ruleStyle = ruleFocusedStyle
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			header,
			ruleStyle.Render(strings.Repeat("─", contentW)),
		)
	}

	renderColumnBox := func(content string, w int) string {
		return lipgloss.NewStyle().
			Background(bgMain).
			Padding(0, 1).
			Width(w).
			Height(innerH).
			Render(content)
	}

	sepLine := bgSpace.Render(" ") + lipgloss.NewStyle().Foreground(mockSep).Background(bgMain).Render("▏") + bgSpace.Render(" ")
	var sepBuilder strings.Builder
	for i := range innerH {
		if i > 0 {
			sepBuilder.WriteString("\n")
		}
		sepBuilder.WriteString(sepLine)
	}
	sep := sepBuilder.String()

	// Empty-state early return: no decomposition yet. Preserve the copy
	// the 2-column version used so the chat/decomp split keeps its
	// familiar feel.
	if len(selPath) == 0 {
		leftPlaceholder := "send a message to start"
		if decomposing {
			leftPlaceholder = "decomposing…"
		}
		leftContentW := max(decompColW-2, 1)
		leftContent := lipgloss.JoinVertical(lipgloss.Left,
			iconStyle.Render("◆ ")+textStyle.Render("root"),
			dimStyle.Render(strings.Repeat("─", leftContentW)),
			bgSpace.Render(strings.Repeat(" ", leftContentW)),
			placeholderStyle.Render(leftPlaceholder),
		)
		leftCol := renderColumnBox(leftContent, decompColW)

		detailW := max(innerW-decompColW-sepWidth, decompMinRightW)
		titleText := "(no decomposition)"
		bodyText := "Describe your feature in the chat. Scope changes refine the tree; questions stay in chat."
		if decomposing {
			titleText = "decomposing…"
			bodyText = "Asking the model to decompose the feature. This usually takes a few seconds."
		}
		detailContent := titled(titleText, "", detailW, false) + "\n\n" + bodyDimStyle.Render(bodyText)
		detailCol := renderColumnBox(detailContent, detailW)
		return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, sep, detailCol)
	}

	// Decide how many navigable columns to show on screen.
	nNav := len(selPath)
	visible := nNav
	if visible > decompMaxVisCols {
		visible = decompMaxVisCols
	}
	if colScroll < 0 {
		colScroll = 0
	}
	if colScroll > nNav-visible {
		colScroll = nNav - visible
	}

	// Width budget: each navigable column grows to fit its own content
	// (no minimum share of the terminal). Any leftover width goes to
	// the detail pane, not to inflate the nav columns.
	colWidths := make([]int, visible)
	for vi := 0; vi < visible; vi++ {
		colIdx := colScroll + vi
		parent := "root"
		if colIdx > 0 {
			parent = selPath[colIdx-1]
		}
		kids := snap.ChildrenByName[parent]

		want := decompMinColW
		if colIdx == 0 {
			// Section headers (std / project name) participate in
			// the width calc for column 0.
			for _, pkg := range snap.PackagePaths {
				if n := len(displayName(snap, pkg)) + 2; n > want {
					want = n
				}
			}
		} else {
			headerLabel := displayName(snap, selPath[colIdx-1])
			// "◆ " + label + " ›" shift hint + padding
			if n := 2 + len(headerLabel) + 2 + 2; n > want {
				want = n
			}
		}
		for _, p := range kids {
			// glyph "• " (2) + name + " ›" hint (2) + padding (2)
			if n := 2 + len(displayName(snap, p)) + 2 + 2; n > want {
				want = n
			}
		}
		colWidths[vi] = want
	}

	// Now allocate width. Reserve the nominal detail pane width, then
	// see how much is left over for nav columns.
	detailW := max(innerW/3, decompMinRightW)
	navArea := innerW - detailW - sepWidth*visible
	total := 0
	for _, w := range colWidths {
		total += w
	}
	// If nav columns don't fit, steal from the detail pane down to its
	// minimum; if they still don't fit, shrink the deepest navigable
	// column first so column 0 (the entry-point into the tree) keeps its
	// full content width.
	if total > navArea {
		grab := min(total-navArea, detailW-decompMinRightW)
		if grab > 0 {
			detailW -= grab
			navArea += grab
		}
	}
	for total > navArea {
		shrunk := false
		for i := len(colWidths) - 1; i >= 0; i-- {
			if colWidths[i] > decompMinColW {
				colWidths[i]--
				total--
				shrunk = true
				break
			}
		}
		if !shrunk {
			break
		}
	}
	// Any slack left over goes to the detail pane so nav columns stay
	// at their natural content width instead of being stretched out.
	if total < navArea {
		detailW += navArea - total
	}

	// renderRuneRow produces one list row for a rune path — status glyph
	// plus the short display name, with an optional "›" hint when the
	// rune has children to drill into. Used by both column 0 (sectioned)
	// and deeper columns (flat).
	renderRuneRow := func(p, selected string, focused bool, contentW int) string {
		name := displayName(snap, p)
		isSel := p == selected
		prefix, prefixStyle := statusGlyph(isSel, focused && active, blinkOn)
		hint := ""
		if len(snap.ChildrenByName[p]) > 0 {
			hint = " ›"
		}
		// Use visual column width (lipgloss.Width) rather than byte length —
		// the glyph and hint are multi-byte UTF-8 but each render to 2 cols.
		prefixW := lipgloss.Width(prefix)
		hintW := lipgloss.Width(hint)
		label := name
		maxName := contentW - prefixW - hintW
		if maxName < 1 {
			maxName = 1
		}
		if lipgloss.Width(label) > maxName {
			if maxName > 1 {
				label = label[:maxName-1] + "…"
			} else {
				label = "…"
			}
		}
		line := prefixStyle.Render(prefix) + textStyle.Render(label)
		if hint != "" {
			line += dimStyle.Render(hint)
		}
		if gap := contentW - prefixW - lipgloss.Width(label) - hintW; gap > 0 {
			line += bgSpace.Render(strings.Repeat(" ", gap))
		}
		return line
	}

	// applyScroll windows parts[] around selRow so it fits in innerH
	// rows, adding "↑"/"↓" chrome rows when content is clipped. Every
	// entry in parts[] must be exactly one terminal row.
	arrowStyle := lipgloss.NewStyle().Foreground(fgDim).Background(bgMain)
	applyScroll := func(parts []string, selRow, contentW int) string {
		total := len(parts)
		if total <= innerH {
			return lipgloss.JoinVertical(lipgloss.Left, parts...)
		}
		// Iteratively pick yOffset + chrome. Chrome rows shrink the
		// visible content area, which may push yOffset further to keep
		// the selected row in view. Two passes is always enough: one
		// for initial chrome, one for the adjusted window.
		visible := innerH
		yOffset := 0
		for iter := 0; iter < 2; iter++ {
			if selRow < yOffset {
				yOffset = selRow
			}
			if selRow >= yOffset+visible {
				yOffset = selRow - visible + 1
			}
			if yOffset < 0 {
				yOffset = 0
			}
			if yOffset+visible > total {
				yOffset = total - visible
				if yOffset < 0 {
					yOffset = 0
				}
			}
			needsUp := yOffset > 0
			needsDown := yOffset+visible < total
			chrome := 0
			if needsUp {
				chrome++
			}
			if needsDown {
				chrome++
			}
			newVisible := innerH - chrome
			if newVisible == visible {
				break
			}
			visible = newVisible
		}
		if yOffset+visible > total {
			yOffset = total - visible
		}
		if yOffset < 0 {
			yOffset = 0
		}
		needsUp := yOffset > 0
		needsDown := yOffset+visible < total

		arrowRow := func(glyph string) string {
			pad := max(contentW-1, 0)
			return bgSpace.Render(strings.Repeat(" ", pad)) + arrowStyle.Render(glyph)
		}
		out := make([]string, 0, innerH)
		if needsUp {
			out = append(out, arrowRow("↑"))
		}
		end := yOffset + visible
		if end > total {
			end = total
		}
		out = append(out, parts[yOffset:end]...)
		if needsDown {
			out = append(out, arrowRow("↓"))
		}
		return lipgloss.JoinVertical(lipgloss.Left, out...)
	}

	rendered := make([]string, 0, visible+1)
	for vi := 0; vi < visible; vi++ {
		colIdx := colScroll + vi
		focused := colIdx == focusedCol
		w := colWidths[vi]
		contentW := max(w-2, 1)

		var parts []string
		selRow := 0

		if colIdx == 0 {
			// Sectioned layout: std header + its runes, then project
			// header + its runes. Section headers are pure display and
			// don't consume a nav slot — up/down steps through the flat
			// list in ChildrenByName["root"].
			selected := ""
			if len(selPath) > 0 {
				selected = selPath[0]
			}
			pkgHeaderStyle := lipgloss.NewStyle().Foreground(accentSoft).Background(bgMain).Bold(true)
			for pi, pkg := range snap.PackagePaths {
				if pi > 0 {
					parts = append(parts, bgSpace.Render(strings.Repeat(" ", contentW)))
				}
				parts = append(parts, pkgHeaderStyle.Render(displayName(snap, pkg)))
				for _, p := range snap.ChildrenByName[pkg] {
					if p == selected {
						selRow = len(parts)
					}
					parts = append(parts, renderRuneRow(p, selected, focused, contentW))
				}
			}
			if len(parts) == 0 {
				parts = append(parts, placeholderStyle.Render("(no runes)"))
			}
		} else {
			parent := selPath[colIdx-1]
			kids := snap.ChildrenByName[parent]
			headerLabel := displayName(snap, selPath[colIdx-1])
			if vi == 0 && colScroll > 0 {
				headerLabel = "‹ " + headerLabel
			}
			if vi == visible-1 && colScroll+visible < nNav {
				headerLabel = headerLabel + " ›"
			}
			// Decompose the titled() output into per-row entries so
			// scrolling math can count rows accurately.
			headerRow := iconStyle.Render("◆ ") + textStyle.Render(headerLabel)
			ruleStyle := dimStyle
			if focused {
				ruleStyle = ruleFocusedStyle
			}
			ruleRow := ruleStyle.Render(strings.Repeat("─", contentW))
			parts = append(parts, headerRow, ruleRow, bgSpace.Render(strings.Repeat(" ", contentW)))

			if len(kids) == 0 {
				parts = append(parts, placeholderStyle.Render("(no children)"))
			} else {
				selected := ""
				if colIdx < len(selPath) {
					selected = selPath[colIdx]
				}
				for _, p := range kids {
					if p == selected {
						selRow = len(parts)
					}
					parts = append(parts, renderRuneRow(p, selected, focused, contentW))
				}
			}
		}
		content := applyScroll(parts, selRow, contentW)
		rendered = append(rendered, renderColumnBox(content, w))
	}

	// Detail pane: anchored to the deepest selected entry. Handles
	// both package-root rows (which have no Rune spec) and real runes.
	detailCol := renderDetailPane(selPath, snap, detailW, titled)
	rendered = append(rendered, renderColumnBox(detailCol, detailW))

	pieces := make([]string, 0, 2*len(rendered)-1)
	for i, col := range rendered {
		if i > 0 {
			pieces = append(pieces, sep)
		}
		pieces = append(pieces, col)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, pieces...)
}

// displayName returns the short column label for a fully-qualified path.
// Falls back to the last dot-segment when the snapshot has no entry.
func displayName(snap decomposer.Snapshot, path string) string {
	if n, ok := snap.DisplayNameByPath[path]; ok && n != "" {
		return n
	}
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return path
}

// renderDetailPane builds the rightmost column's content. The deepest
// selected entry is always a rune path — packages are never selectable,
// they're section headers in column 0.
func renderDetailPane(selPath []string, snap decomposer.Snapshot, w int, titled func(string, string, int, bool) string) string {
	if len(selPath) == 0 {
		return ""
	}
	path := selPath[len(selPath)-1]
	name := displayName(snap, path)

	r := snap.RuneByPath[path]
	childPaths := snap.ChildrenByName[path]
	shortChildren := make([]string, 0, len(childPaths))
	for _, cp := range childPaths {
		shortChildren = append(shortChildren, displayName(snap, cp))
	}
	return titled(name, leafTag(len(childPaths) == 0), w, false) + "\n\n" + renderRuneInfo(name, r, shortChildren, w)
}

