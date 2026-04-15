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
	decompNumCols       = 2
	decompNavigableCols = decompNumCols - 1 // rightmost column is the summary/detail pane, not navigable
	decompColW          = 30                // minimum width of the left column; it grows to fit the longest name
	decompMinRightW     = 24                // keep this much room for the detail pane even when names are long
)

const emptyTopCopy = "Describe your feature in the chat. This pane updates when you change scope."

type featureDecompModel struct {
	width       int
	height      int
	vp          viewport.Model
	selectedIdx int
	focusedCol  int
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
	if n := len(m.snap.TopLevelNames); n == 0 {
		m.selectedIdx = 0
	} else if m.selectedIdx >= n {
		m.selectedIdx = n - 1
	}
	m.vp.SetContent(buildColumns(m.vp.Width(), m.vp.Height(), m.selectedIdx, m.focusedCol, m.active, m.blinkOn, m.decomposing, m.snap))
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
			if decompNavigableCols > 1 {
				m.focusedCol = (m.focusedCol - 1 + decompNavigableCols) % decompNavigableCols
				m.refreshColumns()
			}
			return m, nil
		case "right", "l":
			if decompNavigableCols > 1 {
				m.focusedCol = (m.focusedCol + 1) % decompNavigableCols
				m.refreshColumns()
			}
			return m, nil
		case "up", "k":
			if m.focusedCol == 0 && m.selectedIdx > 0 {
				m.selectedIdx--
				m.blinkOn = true
				m.steadyUntil = time.Now().Add(blinkHoldAfterMove)
				m.refreshColumns()
			}
			return m, nil
		case "down", "j":
			if m.focusedCol == 0 && m.selectedIdx < len(m.snap.TopLevelNames)-1 {
				m.selectedIdx++
				m.blinkOn = true
				m.steadyUntil = time.Now().Add(blinkHoldAfterMove)
				m.refreshColumns()
			}
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
	top := renderDecompTop(m.width, topH, m.selectedIdx*4+m.kanjiOffset, m.snap)
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
	if snap.Expanding {
		runesText += fmt.Sprintf(" · expanding (depth %d)", snap.MaxDepthReached+1)
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

// statusTag returns a short human label for a rune status, shown in the
// detail pane's right-aligned tag slot.
func statusTag(st decomposer.RuneStatus) string {
	switch st {
	case decomposer.StatusInFlight:
		return "expanding…"
	case decomposer.StatusDone:
		return "expanded"
	case decomposer.StatusLeaf:
		return "leaf"
	case decomposer.StatusError:
		return "failed"
	}
	return ""
}

// statusGlyph returns the bullet used in the left column list for a rune
// at the given status, plus the lipgloss style to render it with.
func statusGlyph(st decomposer.RuneStatus, selected, active, blinkOn bool) (string, lipgloss.Style) {
	base := "• "
	switch st {
	case decomposer.StatusInFlight:
		base = "◯ "
	case decomposer.StatusDone:
		base = "● "
	case decomposer.StatusLeaf:
		base = "· "
	case decomposer.StatusError:
		base = "✗ "
	}

	switch {
	case selected && active && blinkOn:
		return base, lipgloss.NewStyle().Foreground(mockHot).Background(bgMain).Bold(true)
	case selected && active && !blinkOn:
		return base, lipgloss.NewStyle().Foreground(bgMain).Background(bgMain)
	case selected:
		return base, lipgloss.NewStyle().Foreground(lipgloss.Color("#9a3050")).Background(bgMain)
	case st == decomposer.StatusDone:
		return base, lipgloss.NewStyle().Foreground(accent).Background(bgMain)
	case st == decomposer.StatusInFlight:
		return base, lipgloss.NewStyle().Foreground(mockHot).Background(bgMain)
	case st == decomposer.StatusError:
		return base, lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Background(bgMain)
	case st == decomposer.StatusLeaf:
		return base, lipgloss.NewStyle().Foreground(fgDim).Background(bgMain)
	}
	return base, lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
}

func renderRuneInfo(name string, r decomposer.Rune, status decomposer.RuneStatus, children []string, maxW int) string {
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
	_ = status
	_ = name
	return strings.Join(lines, "\n")
}

func buildColumns(innerW, innerH, selectedIdx, focusedCol int, active, blinkOn, decomposing bool, snap decomposer.Snapshot) string {
	const sepWidth = 3
	// Grow the left column to fit the longest header/name so lipgloss's
	// Width() never soft-wraps a row. Glyph prefix ("◆ "/"• ") is 2 cells;
	// add 2 more for the column's outer Padding(0, 1).
	leftContentW := 2 + len(nonEmptyOr(snap.PackageName, "(empty)"))
	for _, name := range snap.TopLevelNames {
		if need := 2 + len(name); need > leftContentW {
			leftContentW = need
		}
	}
	leftW := max(leftContentW+2, decompColW)
	if ceiling := innerW - sepWidth - decompMinRightW; ceiling >= decompColW && leftW > ceiling {
		leftW = ceiling
	}
	widths := []int{leftW, max(innerW-leftW-sepWidth, decompMinRightW)}

	textStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)
	dimStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
	ruleFocusedStyle := lipgloss.NewStyle().Foreground(mockFocus).Background(bgMain)
	iconStyle := lipgloss.NewStyle().Foreground(mockHot).Background(bgMain).Bold(true)
	tagStyle := lipgloss.NewStyle().Foreground(fgBody).Background(bgMain)

	titled := func(name, tag string, w int) string {
		contentW := w - 2
		nameSeg := iconStyle.Render("◆ ") + textStyle.Render(name)
		header := nameSeg
		if tag != "" {
			tagSeg := tagStyle.Render("# " + tag)
			pad := max(contentW-lipgloss.Width(nameSeg)-lipgloss.Width(tagSeg), 1)
			spacer := lipgloss.NewStyle().Background(bgMain).Render(strings.Repeat(" ", pad))
			header = nameSeg + spacer + tagSeg
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			header,
			dimStyle.Render(strings.Repeat("─", contentW)),
		)
	}

	cols := make([]string, decompNumCols)
	for i := range decompNumCols {
		var content string
		focused := i == focusedCol
		switch i {
		case 0:
			contentW := widths[i] - 2
			ruleStyle := dimStyle
			if focused {
				ruleStyle = ruleFocusedStyle
			}
			bgSpace := lipgloss.NewStyle().Background(bgMain)
			header := iconStyle.Render("◆ ") + textStyle.Render(nonEmptyOr(snap.PackageName, "(empty)"))
			parts := []string{
				header,
				ruleStyle.Render(strings.Repeat("─", contentW)),
				bgSpace.Render(strings.Repeat(" ", contentW)),
			}

			if len(snap.TopLevelNames) == 0 {
				placeholderText := "send a message to start"
				if decomposing {
					placeholderText = "decomposing…"
				}
				placeholder := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain).Italic(true).
					Render(placeholderText)
				parts = append(parts, placeholder)
			}

			for j, name := range snap.TopLevelNames {
				qualified := qualifiedPath(snap.PackageName, name)
				st := snap.StatusByName[qualified]
				prefix, prefixStyle := statusGlyph(st, focused && j == selectedIdx, active, blinkOn)
				line := prefixStyle.Render(prefix) + textStyle.Render(name)
				gap := max(contentW-len(prefix)-len(name), 0)
				if gap > 0 {
					line += bgSpace.Render(strings.Repeat(" ", gap))
				}
				parts = append(parts, line)
			}
			content = lipgloss.JoinVertical(lipgloss.Left, parts...)

		case 1:
			if len(snap.TopLevelNames) == 0 {
				titleText := "(no decomposition)"
				bodyText := "Describe your feature in the chat. Scope changes refine the tree; questions stay in chat."
				if decomposing {
					titleText = "decomposing…"
					bodyText = "Asking the model to decompose the feature. This usually takes a few seconds."
				}
				content = titled(titleText, "", widths[i]) + "\n\n" +
					lipgloss.NewStyle().Foreground(fgDim).Background(bgMain).Italic(true).
						Render(bodyText)
				break
			}
			name := snap.TopLevelNames[selectedIdx]
			qualified := qualifiedPath(snap.PackageName, name)
			r := snap.RunesByName[name]
			status := snap.StatusByName[qualified]
			children := snap.ChildrenByName[qualified]
			base := titled(name, statusTag(status), widths[i])
			content = base + "\n\n" + renderRuneInfo(name, r, status, children, widths[i])
		}
		cols[i] = lipgloss.NewStyle().
			Background(bgMain).
			Padding(0, 1).
			Width(widths[i]).
			Height(innerH).
			Render(content)
	}

	pad := lipgloss.NewStyle().Background(bgMain).Render(" ")
	line := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain).Render("▏")
	sepLine := pad + line + pad
	var sepBuilder strings.Builder
	for i := range innerH {
		if i > 0 {
			sepBuilder.WriteString("\n")
		}
		sepBuilder.WriteString(sepLine)
	}
	sep := sepBuilder.String()

	pieces := make([]string, 0, 2*decompNumCols-1)
	for i, col := range cols {
		if i > 0 {
			pieces = append(pieces, sep)
		}
		pieces = append(pieces, col)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, pieces...)
}

func nonEmptyOr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// qualifiedPath returns "{pkg}.{rune}" when rune is not already prefixed
// with pkg, matching the session's Status map keys.
func qualifiedPath(pkg, rune string) string {
	if pkg == "" {
		return rune
	}
	if strings.HasPrefix(rune, pkg+".") {
		return rune
	}
	return pkg + "." + rune
}
