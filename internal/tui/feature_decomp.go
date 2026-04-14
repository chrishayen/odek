package tui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

var (
	mockSep   = lipgloss.Color("#3e3e3e")
	mockHot   = lipgloss.Color("#f74e82")
	mockFocus = lipgloss.Color("#9d7cf3")
)

const (
	mockNumCols       = 2
	mockNavigableCols = mockNumCols - 1 // rightmost column is the summary/detail pane, not navigable
	mockColW          = 30
)

var mockTomlItems = []string{
	"decode",
	"encode",
	"encode_value",
	"tokenize",
	"get_array",
	"get_bool",
	"get_int",
	"get_string",
	"parse",
}

type mockRuneInfo struct {
	summary   string
	sig       string
	pluses    []string
	minuses   []string
	questions []string
	tag       string
	deps      []string
}

var mockRuneInfos = []mockRuneInfo{
	{
		summary: "One-shot helper that chains tokenize + parse so callers can hand in source text and get back a ready-to-query value tree.",
		sig:     "(source: string) -> result[toml_value, string]",
		pluses:  []string{"tokenizes then parses, returning the root table"},
		minuses: []string{"propagates tokenize and parse errors"},
		tag:     "decoding",
		deps:    []string{"toml.tokenize", "toml.parse"},
	},
	{
		summary: "Serializes a full toml_value tree back into TOML text, choosing header grouping and key order for round-trip fidelity.",
		sig:     "(root: toml_value) -> result[string, string]",
		pluses: []string{
			"returns a TOML document whose decode round-trips to the same value tree",
			`emits nested tables as "[section.sub]" headers ordered by appearance`,
		},
		minuses: []string{"returns error when the root is not a table"},
		tag:     "encoding",
		deps:    []string{"toml.encode_value"},
	},
	{
		summary:   "Renders a single TOML value — scalar, string, number, bool, or inline array — in its literal source form.",
		sig:       "(value: toml_value) -> string",
		pluses:    []string{"renders a scalar value or inline array in TOML literal syntax"},
		questions: []string{"used internally by encode; exposed for callers who build values directly"},
		tag:       "encoding",
	},
	{
		summary: "Turns raw TOML source into a flat stream of tokens — keys, values, brackets, equals signs, newlines — ready for the parser to consume.",
		sig:     "(source: string) -> result[list[toml_token], string]",
		pluses:  []string{"splits the input into keys, equals, strings, numbers, booleans, brackets, and newlines"},
		minuses: []string{"returns error on an unterminated string literal", "returns error on an invalid number literal"},
		tag:     "lexing",
	},
	{
		summary: "Fetches an array at a dotted path within a decoded TOML value tree.",
		sig:     "(root: toml_value, path: string) -> optional[list[toml_value]]",
		pluses:  []string{"returns the array at a dotted path"},
		minuses: []string{"returns none when missing or not an array"},
		tag:     "lookup",
	},
	{
		summary: "Looks up a boolean leaf at a dotted path within a decoded TOML value tree.",
		sig:     "(root: toml_value, path: string) -> optional[bool]",
		pluses:  []string{"returns the boolean value at a dotted path"},
		minuses: []string{"returns none when missing or not a boolean"},
		tag:     "lookup",
	},
	{
		summary: "Looks up an integer leaf at a dotted path within a decoded TOML value tree.",
		sig:     "(root: toml_value, path: string) -> optional[i64]",
		pluses:  []string{"returns the integer value at a dotted path"},
		minuses: []string{"returns none when missing or not an integer"},
		tag:     "lookup",
	},
	{
		summary: "Looks up a string leaf at a dotted path within a decoded TOML value tree.",
		sig:     "(root: toml_value, path: string) -> optional[string]",
		pluses:  []string{`returns the string value at a dotted path like "server.host"`},
		minuses: []string{"returns none when any segment is missing or the leaf is not a string"},
		tag:     "lookup",
	},
	{
		summary: "Walks a token stream and builds the root table value: nested sections, inline arrays, and scalar leaves.",
		sig:     "(tokens: list[toml_token]) -> result[toml_value, string]",
		pluses: []string{
			"returns the root table value containing every top-level key",
			`supports nested tables via "[section.sub]" headers`,
			"supports inline arrays of homogeneous type",
		},
		minuses: []string{
			"returns error on duplicate keys within the same table",
			"returns error on unclosed section headers",
		},
		tag: "parsing",
	},
}

const mockTopCopy = "A full subsystem: tokenize, parse, build a typed value tree, and emit it back as text. Values can be strings, integers, floats, booleans, arrays, or nested tables."

type featureDecompModel struct {
	width       int
	height      int
	pin         string
	vp          viewport.Model
	selectedIdx int
	focusedCol  int
	inputActive bool
	active      bool
	blinkOn     bool
	inSplit     bool
	steadyUntil time.Time
	input       textinput.Model
}

func newFeatureDecompModel(width, height int, pin string) featureDecompModel {
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
		pin:    pin,
		vp:     viewport.New(),
		input:  ti,
		active: true,
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

// mockBottomChromeRows is the number of rows consumed below the viewport:
// blank + pin + blank + input + help = 5, plus 1 spacer above the viewport.
const mockBottomChromeRows = 6

func (m *featureDecompModel) resize(w, h int) {
	m.width = w
	m.height = h
	if w <= 0 || h <= 0 {
		return
	}
	topH := h / 3
	bottomH := max(h-topH-mockBottomChromeRows, 1)
	m.vp.SetWidth(w)
	m.vp.SetHeight(bottomH)
	m.refreshColumns()
}

func (m *featureDecompModel) refreshColumns() {
	m.vp.SetContent(buildMockColumns(m.vp.Width(), m.vp.Height(), m.selectedIdx, m.focusedCol, m.active, m.blinkOn))
}

// PinRow reports the row at which the feature pin renders, so the slide
// transition can land its end-state icon at the same position.
func (m featureDecompModel) PinRow() int {
	if m.height <= 0 {
		return 0
	}
	topH := m.height / 3
	bottomH := max(m.height-topH-mockBottomChromeRows, 1)
	// JoinVertical order: top(topH) + blank(1) + bottom(bottomH) + blank(1) = pin row.
	return topH + 1 + bottomH + 1
}

func (m featureDecompModel) Init() tea.Cmd { return blinkTick() }

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
			if mockNavigableCols > 1 {
				m.focusedCol = (m.focusedCol - 1 + mockNavigableCols) % mockNavigableCols
				m.refreshColumns()
			}
			return m, nil
		case "right", "l":
			if mockNavigableCols > 1 {
				m.focusedCol = (m.focusedCol + 1) % mockNavigableCols
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
			if m.focusedCol == 0 && m.selectedIdx < len(mockTomlItems)-1 {
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
	top := renderMockTop(m.width, topH, m.selectedIdx*4)
	blank := lipgloss.NewStyle().Background(bgMain).Render(strings.Repeat(" ", m.width))
	bottom := m.vp.View()
	pinRow := renderMockPinRow(m.width, m.pin)

	lastLine := blank
	if m.inputActive {
		rendered := m.input.View()
		pad := max(m.width-lipgloss.Width(rendered), 0)
		lastLine = rendered + lipgloss.NewStyle().Background(bgMain).Render(strings.Repeat(" ", pad))
	}

	help := renderMockHelp(m.width, m.inputActive, m.inSplit)

	v.Content = lipgloss.JoinVertical(lipgloss.Left,
		top,
		blank,
		bottom,
		blank,
		pinRow,
		blank,
		lastLine,
		help,
	)
	return v
}

// renderMockPinRow places the feature pin at the left edge of a full-width
// row, padding the rest with bgMain so the row has the same width as the
// surrounding layout.
func renderMockPinRow(width int, pin string) string {
	bgPad := lipgloss.NewStyle().Background(bgMain)
	pinW := lipgloss.Width(pin)
	if pinW >= width {
		return pin
	}
	return pin + bgPad.Render(strings.Repeat(" ", width-pinW))
}

func renderMockHelp(width int, inputActive, showTabSwitch bool) string {
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

func renderMockTop(width, height, kanjiOffset int) string {
	const (
		logoText  = "ODEK "
		runesText = "9 runes"
	)
	logoCells := len(logoText)
	logoStyled := lipgloss.NewStyle().
		Foreground(accent).
		Background(bgMain).
		Bold(true).
		Render(logoText)

	bgStyle := lipgloss.NewStyle().Background(bgMain)
	textStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)
	blankLine := bgStyle.Render(strings.Repeat(" ", width))

	// Row 0: logo + bg padding.
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
	rawLines := wrapMockText(mockTopCopy, textMax)

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

func wrapMockText(text string, width int) []string {
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

func renderMockRuneInfo(info mockRuneInfo, maxW int) string {
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
	if info.summary != "" {
		lines = append(lines, summaryStyle.Render(info.summary))
		lines = append(lines, "")
	}
	if info.sig != "" {
		lines = append(lines, sigStyle.Render("@ "+info.sig))
	}
	if len(info.pluses) > 0 || len(info.minuses) > 0 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, headingStyle.Render("Behavior"))
		for _, p := range info.pluses {
			lines = append(lines, plusStyle.Render("  + ")+bodyStyle.Render(p))
		}
		for _, mn := range info.minuses {
			lines = append(lines, minusStyle.Render("  - ")+bodyStyle.Render(mn))
		}
	}
	for _, q := range info.questions {
		lines = append(lines, questionStyle.Render("? "+q))
	}
	if len(info.deps) > 0 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, headingStyle.Render("References"))
		for _, d := range info.deps {
			lines = append(lines, depStyle.Render("  -> ")+bodyStyle.Render(d))
		}
	}
	return strings.Join(lines, "\n")
}

func buildMockColumns(innerW, innerH, selectedIdx, focusedCol int, active, blinkOn bool) string {
	const sepWidth = 3
	widths := []int{mockColW, max(innerW-mockColW-sepWidth, mockColW)}

	textStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain)
	dimStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
	ruleFocusedStyle := lipgloss.NewStyle().Foreground(mockFocus).Background(bgMain)
	iconStyle := lipgloss.NewStyle().Foreground(mockHot).Background(bgMain).Bold(true)
	var cursorStyle lipgloss.Style
	switch {
	case active && blinkOn:
		cursorStyle = lipgloss.NewStyle().
			Foreground(mockHot).
			Background(bgMain).
			Bold(true)
	case active && !blinkOn:
		cursorStyle = lipgloss.NewStyle().
			Foreground(bgMain).
			Background(bgMain)
	default:
		cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9a3050")).
			Background(bgMain)
	}
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

	cols := make([]string, mockNumCols)
	for i := range mockNumCols {
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
			parts := []string{
				iconStyle.Render("◆ ") + textStyle.Render("toml"),
				ruleStyle.Render(strings.Repeat("─", contentW)),
				bgSpace.Render(strings.Repeat(" ", contentW)),
			}

			for j, name := range mockTomlItems {
				prefix := "• "
				circleStyle := dimStyle
				if focused && j == selectedIdx {
					circleStyle = cursorStyle
				}
				line := circleStyle.Render(prefix) + textStyle.Render(name)
				gap := max(contentW-len(prefix)-len(name), 0)
				if gap > 0 {
					line += bgSpace.Render(strings.Repeat(" ", gap))
				}
				parts = append(parts, line)
			}
			content = lipgloss.JoinVertical(lipgloss.Left, parts...)
		case 1:
			info := mockRuneInfos[selectedIdx]
			base := titled(mockTomlItems[selectedIdx], info.tag, widths[i])
			content = base + "\n\n" + renderMockRuneInfo(info, widths[i])
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

	pieces := make([]string, 0, 2*mockNumCols-1)
	for i, col := range cols {
		if i > 0 {
			pieces = append(pieces, sep)
		}
		pieces = append(pieces, col)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, pieces...)
}
