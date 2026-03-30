package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chrishayen/valkyrie/internal/feature"
	"github.com/lucasb-eyer/go-colorful"
)

type screen int

const (
	screenSplash screen = iota
	screenCreateFeature
)

var logoBig = `
  ██████╗ ██████╗ ███████╗██╗  ██╗
 ██╔═══██╗██╔══██╗██╔════╝██║ ██╔╝
 ██║   ██║██║  ██║█████╗  █████╔╝
 ██║   ██║██║  ██║██╔══╝  ██╔═██╗
 ╚██████╔╝██████╔╝███████╗██║  ██╗
  ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝`

var logoSmall = "O D E K"

var (
	border  = lipgloss.Color("#666666")
	dim     = lipgloss.Color("#888888")
	helpKey = lipgloss.Color("#6A9FD9")

	taglineStyle = lipgloss.NewStyle().
			Foreground(dim).
			Italic(true)

	frameStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(1, 4)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(helpKey).
			Bold(true)

	helpTextStyle = lipgloss.NewStyle().
			Foreground(dim)

	helpBarStyle = lipgloss.NewStyle()

	_ = lipgloss.NewStyle()
)

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
		// Pad to full width
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

func renderGradient(text string, stops []colorful.Color) string {
	lines := strings.Split(text, "\n")

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
	var out strings.Builder
	for i, line := range lines {
		runes := []rune(line)
		for j, r := range runes {
			if r == ' ' {
				out.WriteRune(r)
				continue
			}
			t := float64(j) / float64(maxLen) * float64(n)
			idx := int(t)
			if idx >= n {
				idx = n - 1
			}
			frac := t - float64(idx)
			c := stops[idx].BlendLuv(stops[idx+1], frac)
			out.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex())).Bold(true).Render(string(r)))
		}
		if i < len(lines)-1 {
			out.WriteRune('\n')
		}
	}
	return out.String()
}

type Model struct {
	width        int
	height       int
	screen       screen
	featureStore *feature.Store
	createForm   createFeatureModel
}

func New(featureStore *feature.Store) Model {
	return Model{
		screen:       screenSplash,
		featureStore: featureStore,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.screen == screenCreateFeature {
			m.createForm.resize(msg.Width, msg.Height)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.screen == screenSplash {
				return m, tea.Quit
			}
		case "esc":
			if m.screen == screenCreateFeature {
				m.screen = screenSplash
				return m, nil
			}
		case "enter":
			if m.screen == screenSplash {
				m.createForm = newCreateFeatureModel(m.featureStore, m.width, m.height)
				m.screen = screenCreateFeature
				return m, m.createForm.descInput.Focus()
			}
		}
	}

	if m.screen == screenCreateFeature {
		cmd := m.createForm.update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	switch m.screen {
	case screenCreateFeature:
		return m.viewCreateFeature()
	default:
		return m.viewSplash()
	}
}

func (m Model) viewSplash() string {
	gradientLogo := renderStripes(logoBig, gradStops)

	content := gradientLogo + "\n\n" +
		taglineStyle.Render("agentic code orchestration")

	framed := frameStyle.Render(content)

	help := helpBarStyle.Render(
		fmt.Sprintf("%s %s    %s %s",
			helpKeyStyle.Render("enter"),
			helpTextStyle.Render("create feature"),
			helpKeyStyle.Render("q"),
			helpTextStyle.Render("quit"),
		),
	)

	block := lipgloss.JoinVertical(lipgloss.Center, framed, help)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		block,
	)
}

func (m Model) viewCreateFeature() string {
	// Header: gradient logo on full-width dark bar
	header := renderGradientOnBg(" "+logoSmall, gradStops, "#1A1A1A", m.width)

	// Help bar pinned to bottom
	help := helpBarStyle.Render(
		fmt.Sprintf("%s %s    %s %s    %s %s",
			helpKeyStyle.Render("enter"),
			helpTextStyle.Render("create"),
			helpKeyStyle.Render("alt+enter"),
			helpTextStyle.Render("new line"),
			helpKeyStyle.Render("esc"),
			helpTextStyle.Render("back"),
		),
	)

	// Form in the middle
	form := m.createForm.view(m.width)

	headerLines := strings.Count(header, "\n") + 1
	formLines := strings.Count(form, "\n") + 1
	helpHeight := 1

	var b strings.Builder
	b.WriteString(header + "\n")
	b.WriteString("\n")
	b.WriteString(form)

	// Push help bar to bottom
	usedLines := headerLines + 1 + formLines + helpHeight
	gap := m.height - usedLines
	if gap > 0 {
		b.WriteString(strings.Repeat("\n", gap))
	}
	b.WriteString(help)

	return b.String()
}
