package tui

import (
	"context"
	"strings"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/lucasb-eyer/go-colorful"

	openai "shotgun.dev/odek/openai"
)

var logoBig = `
  ██████╗ ██████╗ ███████╗██╗  ██╗
 ██╔═══██╗██╔══██╗██╔════╝██║ ██╔╝
 ██║   ██║██║  ██║█████╗  █████╔╝
 ██║   ██║██║  ██║██╔══╝  ██╔═██╗
 ╚██████╔╝██████╔╝███████╗██║  ██╗
  ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝`

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
	help   help.Model
	ctx    context.Context
	client *openai.Client
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "n", "N":
			next := newCreateFeatureModel(m.ctx, m.client, m.width, m.height)
			return next, next.Init()
		case "d", "D":
			d, groups := mockDecomposition()
			next := newDecompositionModel(d, groups)
			next.width = m.width
			next.height = m.height
			return next, next.Init()
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
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

func renderGradientOnBg(text string, stops []colorful.Color, bg string, totalWidth int) string {
	lines := strings.Split(text, "\n")
	bgColor := lipgloss.Color(bg)

	maxLen := 0
	for _, line := range lines {
		runes := []rune(line)
		if len(runes) > maxLen {
			maxLen = len(runes)
		}
	}
	if maxLen == 0 {
		return text
	}

	n := len(stops) - 1
	bgStyle := lipgloss.NewStyle().Background(bgColor)
	var out strings.Builder
	for i, line := range lines {
		runes := []rune(line)
		for j, r := range runes {
			if r == ' ' {
				out.WriteString(bgStyle.Render(" "))
				continue
			}
			t := float64(j) / float64(maxLen) * float64(n)
			idx := int(t)
			if idx >= n {
				idx = n - 1
			}
			frac := t - float64(idx)
			c := stops[idx].BlendLuv(stops[idx+1], frac)
			out.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex())).Background(bgColor).Bold(true).Render(string(r)))
		}
		pad := totalWidth - len(runes)
		if pad > 0 {
			out.WriteString(bgStyle.Render(strings.Repeat(" ", pad)))
		}
		if i < len(lines)-1 {
			out.WriteRune('\n')
		}
	}
	return out.String()
}

func (m model) View() tea.View {
	gradientLogo := renderStripes(logoBig, gradStops)

	content := gradientLogo + "\n\n" +
		taglineStyle.Render("Tree Composition CLI and Rune Server")

	framed := frameStyle.Render(content)
	helpBar := helpBarStyle.Render(m.help.View(splashKeyMap{}))

	block := lipgloss.JoinVertical(lipgloss.Left, framed, helpBar)

	placed := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		block,
	)
	v := tea.NewView(placed)
	v.AltScreen = true
	v.BackgroundColor = bgMain
	return v
}

func Run(ctx context.Context, client *openai.Client) {
	teaModel := model{help: newHelpModel(), ctx: ctx, client: client}
	p := tea.NewProgram(teaModel)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
