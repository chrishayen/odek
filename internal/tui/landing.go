package tui

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

var logoBig = `
  ██████╗ ██████╗ ███████╗██╗  ██╗
 ██╔═══██╗██╔══██╗██╔════╝██║ ██╔╝
 ██║   ██║██║  ██║█████╗  █████╔╝
 ██║   ██║██║  ██║██╔══╝  ██╔═██╗
 ╚██████╔╝██████╔╝███████╗██║  ██╗
  ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝`

var (
	border = lipgloss.Color("#666666")
	dim    = lipgloss.Color("#888888")
)

var taglineStyle = lipgloss.NewStyle().
	Foreground(dim).
	Italic(true)

var frameStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(border).
	Padding(1, 4)

// 70s VHS tape stripe palette
var gradStops []colorful.Color

func init() {
	hexes := []string{
		"#F9D050", // yellow
		"#F6B830", // amber
		"#F29A1E", // orange
		"#EC7A1E", // deep orange
		"#E05E22", // red-orange
		"#D4443C", // red
		"#B53A2E", // deep red
		"#863520", // rust
		"#5A2A10", // brown
	}
	for _, h := range hexes {
		c, _ := colorful.Hex(h)
		gradStops = append(gradStops, c)
	}
}

type model struct {
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Create - for now just quit
			return m, tea.Quit
		case tea.KeyEsc, tea.KeyBackspace:
			os.Exit(0)
		}
	}
	return m, nil
}

func renderStripes(text string, stops []colorful.Color) string {
	lines := strings.Split(text, "\n")
	n := len(stops) - 1

	var out strings.Builder
	for i, line := range lines {
		t := float64(i) / float64(len(lines)-1) * float64(n)
		idx := int(t)
		if idx >= n {
			idx = n - 1
		}
		frac := t - float64(idx)
		c := stops[idx].BlendLuv(stops[idx+1], frac)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex())).Bold(true)
		out.WriteString(style.Render(line))
		if i < len(lines)-1 {
			out.WriteRune('\n')
		}
	}
	return out.String()
}

func (m model) View() string {
	gradientLogo := renderStripes(logoBig, gradStops)

	content := gradientLogo + "\n\n" +
		taglineStyle.Render("Tree Composition CLI and Rune Server")

	framed := frameStyle.Render(content)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		framed,
	)
}

func Run() {
	teaModel := model{}
	p := tea.NewProgram(teaModel)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
