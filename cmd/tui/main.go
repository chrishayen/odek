package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	bg        = lipgloss.Color("#1c1c1c")
	accent    = lipgloss.Color("212") // pink
	sepColor  = lipgloss.Color("#3e3e3e")
	textColor = lipgloss.Color("15") // white
)

var kanjiPool = []rune("日月火水木金土山川風花雪心愛空東西南北春夏秋冬父母兄弟姉妹雨雲雷電林森竹松梅桜龍虎鳥魚馬犬猫")

const (
	numCols = 2
	colW    = 30
)

var tomlItems = []string{
	"tokenize",
	"parse",
	"decode",
	"get_string",
	"get_int",
	"get_bool",
	"get_array",
	"encode_value",
	"encode",
}

type model struct {
	width, height int
	vp            viewport.Model
	selectedIdx   int
	focusedCol    int
	inputActive   bool
	input         textinput.Model
}

func newModel() model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.Placeholder = ""
	s := ti.Styles()
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(textColor).Background(bg)
	s.Focused.Text = lipgloss.NewStyle().Foreground(textColor).Background(bg)
	s.Blurred.Prompt = lipgloss.NewStyle().Foreground(textColor).Background(bg)
	s.Blurred.Text = lipgloss.NewStyle().Foreground(textColor).Background(bg)
	ti.SetStyles(s)
	return model{vp: viewport.New(), input: ti}
}

func (m *model) refreshColumns() {
	m.vp.SetContent(buildColumns(m.vp.Height(), m.selectedIdx, m.focusedCol))
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		topH := m.height / 3
		// layout below viewport: blank + bar + blank + blank = 4 rows; plus 1 row spacer above viewport.
		bottomH := max(m.height-topH-5, 1)
		m.vp.SetWidth(m.width)
		m.vp.SetHeight(bottomH)
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
		case "tab":
			m.focusedCol = (m.focusedCol + 1) % numCols
			m.refreshColumns()
			return m, nil
		case "shift+tab":
			m.focusedCol = (m.focusedCol - 1 + numCols) % numCols
			m.refreshColumns()
			return m, nil
		case "up", "k":
			if m.focusedCol == 0 && m.selectedIdx > 0 {
				m.selectedIdx--
				m.refreshColumns()
			}
			return m, nil
		case "down", "j":
			if m.focusedCol == 0 && m.selectedIdx < len(tomlItems)-1 {
				m.selectedIdx++
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

func (m model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true
	v.BackgroundColor = bg

	if m.width == 0 || m.height == 0 {
		return v
	}

	topH := m.height / 3
	top := renderTop(m.width, topH, m.selectedIdx*4)
	blank := lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", m.width))
	bottom := m.vp.View()
	bar := renderBottomBar(m.width)
	lastLine := blank
	if m.inputActive {
		rendered := m.input.View()
		pad := max(m.width-lipgloss.Width(rendered), 0)
		lastLine = rendered + lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", pad))
	}

	v.Content = lipgloss.JoinVertical(lipgloss.Left,
		top,
		blank,
		bottom,
		blank,
		bar,
		blank,
		lastLine,
	)
	return v
}

func renderBottomBar(width int) string {
	const (
		tipW    = 2
		sidePad = 1
	)
	barFill := max(width-tipW-2*sidePad, 0)

	pad := lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", sidePad))
	tip := lipgloss.NewStyle().Background(lipgloss.Color("#f74e82")).Render(strings.Repeat(" ", tipW))
	bar := lipgloss.NewStyle().Background(lipgloss.Color("#686867")).Render(strings.Repeat(" ", barFill))

	return pad + tip + bar + pad
}

const topCopy = "A full subsystem: tokenize, parse, build a typed value tree, and emit it back as text. Values can be strings, integers, floats, booleans, arrays, or nested tables."

func renderTop(width, height, kanjiOffset int) string {
	const logoText = "ODEK "
	logoCells := len(logoText)
	logoStyled := lipgloss.NewStyle().
		Foreground(accent).
		Background(bg).
		Bold(true).
		Render(logoText)

	kanjiStyle := lipgloss.NewStyle().Foreground(sepColor).Background(bg)
	bgStyle := lipgloss.NewStyle().Background(bg)

	blankLine := bgStyle.Render(strings.Repeat(" ", width))

	// Ragged-right word-wrap. Each line's visible width is only the text it
	// actually contains — the remainder of the row fills with kanji.
	textMax := max(min(width/2, 70), 20)
	rawLines := wrapText(topCopy, textMax)
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(bg)
	const textStartRow = 2
	const runesText = "9 runes"
	runesRow := textStartRow + len(rawLines) + 2 // third line below the last text line

	lines := make([]string, height)
	for row := range height {
		switch row {
		case 0:
			var b strings.Builder
			b.WriteString(logoStyled)
			if logoCells < width {
				b.WriteString(bgStyle.Render(strings.Repeat(" ", width-logoCells)))
			}
			lines[row] = b.String()
			continue
		case 1:
			lines[row] = blankLine
			continue
		}

		var b strings.Builder
		cells := 0
		if idx := row - textStartRow; idx >= 0 && idx < len(rawLines) {
			line := rawLines[idx]
			b.WriteString(textStyle.Render(line))
			cells = len(line)
		} else if row == runesRow {
			b.WriteString(textStyle.Render(runesText))
			cells = len(runesText)
		}

		var kb strings.Builder
		for cells+2 <= width {
			kb.WriteRune(kanjiAt(row, cells+kanjiOffset))
			cells += 2
		}
		if kb.Len() > 0 {
			b.WriteString(kanjiStyle.Render(kb.String()))
		}
		if cells < width {
			b.WriteString(bgStyle.Render(strings.Repeat(" ", width-cells)))
		}
		lines[row] = b.String()
	}
	return strings.Join(lines, "\n")
}

func wrapText(text string, width int) []string {
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

func kanjiAt(row, col int) rune {
	h := (row*31 + col*17 + row*col*7) % len(kanjiPool)
	if h < 0 {
		h += len(kanjiPool)
	}
	return kanjiPool[h]
}

func buildColumns(innerH, selectedIdx, focusedCol int) string {
	contentW := colW - 2 // inside horizontal padding

	textStyle := lipgloss.NewStyle().Foreground(textColor).Background(bg)
	dimStyle := lipgloss.NewStyle().Foreground(sepColor).Background(bg)
	ruleFocusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9d7cf3")).Background(bg)
	iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f74e82")).Background(bg).Bold(true)
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("#f74e82")).
		Bold(true)

	titled := func(name string, items []string, sel int, focused bool) string {
		header := iconStyle.Render("◆ ") + textStyle.Render(name)
		ruleStyle := dimStyle
		if focused {
			ruleStyle = ruleFocusedStyle
		}
		parts := []string{
			header,
			ruleStyle.Render(strings.Repeat("─", contentW)),
		}
		for i, it := range items {
			prefix := "• "
			raw := prefix + it
			if i == sel {
				pad := contentW - len(raw)
				if pad < 0 {
					pad = 0
				}
				parts = append(parts, selectedStyle.Render(raw+strings.Repeat(" ", pad)))
			} else {
				parts = append(parts, dimStyle.Render(prefix)+textStyle.Render(it))
			}
		}
		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	}

	cols := make([]string, numCols)
	for i := range numCols {
		var content string
		sel := -1
		focused := i == focusedCol
		switch i {
		case 0:
			if focused {
				sel = selectedIdx
			}
			content = titled("toml", tomlItems, sel, focused)
		case 1:
			content = titled(tomlItems[selectedIdx], nil, -1, focused)
		}
		cols[i] = lipgloss.NewStyle().
			Background(bg).
			Padding(0, 1).
			Width(colW).
			Height(innerH).
			Render(content)
	}

	pad := lipgloss.NewStyle().Background(bg).Render(" ")
	line := lipgloss.NewStyle().Foreground(sepColor).Background(bg).Render("▏")
	sepLine := pad + line + pad
	var sepBuilder strings.Builder
	for i := range innerH {
		if i > 0 {
			sepBuilder.WriteString("\n")
		}
		sepBuilder.WriteString(sepLine)
	}
	sep := sepBuilder.String()

	pieces := make([]string, 0, 2*numCols-1)
	for i, col := range cols {
		if i > 0 {
			pieces = append(pieces, sep)
		}
		pieces = append(pieces, col)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, pieces...)
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
