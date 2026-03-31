package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/chrishayen/odek/internal/chat"
)

type screen int

const (
	screenSplash screen = iota
	screenFeatureList
	screenCreateFeature
)

var logoBig = `
  ██████╗ ██████╗ ███████╗██╗  ██╗
 ██╔═══██╗██╔══██╗██╔════╝██║ ██╔╝
 ██║   ██║██║  ██║█████╗  █████╔╝
 ██║   ██║██║  ██║██╔══╝  ██╔═██╗
 ╚██████╔╝██████╔╝███████╗██║  ██╗
  ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝`

var logoSmall = "ODEK"

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
	port         int
	registryPath string
	draftStore   *DraftStore
	chatStore    *chat.Store
	createForm   createFeatureModel
	featureList  featureListModel
}

func New(port int, registryPath string, chatStore *chat.Store) Model {
	return Model{
		screen:       screenSplash,
		port:         port,
		registryPath: registryPath,
		draftStore:   NewDraftStore(registryPath),
		chatStore:    chatStore,
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
		if m.screen == screenFeatureList {
			m.featureList.width = msg.Width
			m.featureList.height = msg.Height
		}
		return m, nil

	case goBackMsg:
		switch m.screen {
		case screenCreateFeature:
			if draft := m.createForm.toDraft(); draft != nil {
				_ = m.draftStore.Save(*draft)
			}
			m.screen = screenSplash
			return m, nil
		case screenFeatureList:
			m.screen = screenSplash
			return m, nil
		}

	case draftSelectedMsg:
		m.createForm = newCreateFeatureModelFromDraft(m.port, m.width, m.height, m.draftStore, m.chatStore, msg.draft)
		m.screen = screenCreateFeature
		if m.createForm.state == stateDone {
			return m, nil
		}
		return m, m.createForm.descInput.Focus()

	case newFeatureMsg:
		m.createForm = newCreateFeatureModel(m.port, m.width, m.height, m.draftStore, m.chatStore)
		m.screen = screenCreateFeature
		return m, m.createForm.descInput.Focus()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.screen == screenSplash {
				return m, tea.Quit
			}
		case "l":
			if m.screen == screenSplash {
				m.featureList = newFeatureListModel(m.draftStore, m.width, m.height)
				m.screen = screenFeatureList
				return m, nil
			}
		case "enter":
			if m.screen == screenSplash {
				m.createForm = newCreateFeatureModel(m.port, m.width, m.height, m.draftStore, m.chatStore)
				m.screen = screenCreateFeature
				return m, m.createForm.descInput.Focus()
			}
		}
	}

	if m.screen == screenCreateFeature {
		cmd := m.createForm.update(msg)
		return m, cmd
	}

	if m.screen == screenFeatureList {
		cmd := m.featureList.update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	switch m.screen {
	case screenFeatureList:
		return m.viewFeatureList()
	case screenCreateFeature:
		return m.viewCreateFeature()
	default:
		return m.viewSplash()
	}
}

func (m Model) viewSplash() string {
	gradientLogo := renderStripes(logoBig, gradStops)

	content := gradientLogo + "\n\n" +
		taglineStyle.Render("Tree Composition CLI and Rune Server")

	framed := frameStyle.Render(content)

	help := helpBarStyle.Render(
		fmt.Sprintf("%s %s    %s %s    %s %s",
			helpKeyStyle.Render("enter"),
			helpTextStyle.Render("create feature"),
			helpKeyStyle.Render("l"),
			helpTextStyle.Render("drafts"),
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

	// Help bar pinned to bottom — context-sensitive
	var helpText string
	switch m.createForm.state {
	case stateDone:
		switch m.createForm.focus {
		case focusLeft:
			helpText = fmt.Sprintf("%s    %s %s    %s %s    %s %s    %s %s    %s %s    %s %s",
				featureNameStyle.Render("feature"),
				helpKeyStyle.Render("j/k"),
				helpTextStyle.Render("navigate"),
				helpKeyStyle.Render("r"),
				helpTextStyle.Render("refine"),
				helpKeyStyle.Render("c"),
				helpTextStyle.Render("chat"),
				helpKeyStyle.Render("tab"),
				helpTextStyle.Render("next"),
				helpKeyStyle.Render("enter"),
				helpTextStyle.Render("new"),
				helpKeyStyle.Render("bksp"),
				helpTextStyle.Render("back"),
			)
		case focusMiddle:
			helpText = fmt.Sprintf("%s    %s %s    %s %s    %s %s    %s %s    %s %s",
				featureNameStyle.Render("rune"),
				helpKeyStyle.Render("r"),
				helpTextStyle.Render("refine"),
				helpKeyStyle.Render("c"),
				helpTextStyle.Render("chat"),
				helpKeyStyle.Render("tab"),
				helpTextStyle.Render("next"),
				helpKeyStyle.Render("enter"),
				helpTextStyle.Render("new"),
				helpKeyStyle.Render("bksp"),
				helpTextStyle.Render("back"),
			)
		case focusRight:
			helpText = fmt.Sprintf("%s    %s %s    %s %s    %s %s    %s %s",
				featureNameStyle.Render("chat"),
				helpKeyStyle.Render("c"),
				helpTextStyle.Render("resume"),
				helpKeyStyle.Render("tab"),
				helpTextStyle.Render("next"),
				helpKeyStyle.Render("enter"),
				helpTextStyle.Render("new"),
				helpKeyStyle.Render("bksp"),
				helpTextStyle.Render("back"),
			)
		}
	case stateChat:
		helpText = fmt.Sprintf("%s %s    %s %s",
			helpKeyStyle.Render("enter"),
			helpTextStyle.Render("send"),
			helpKeyStyle.Render("bksp"),
			helpTextStyle.Render("back"),
		)
	case stateRefining:
		helpText = fmt.Sprintf("%s %s    %s %s",
			helpKeyStyle.Render("enter"),
			helpTextStyle.Render("submit"),
			helpKeyStyle.Render("esc"),
			helpTextStyle.Render("cancel"),
		)
	case stateDecomposing:
		helpText = fmt.Sprintf("%s %s",
			helpKeyStyle.Render("bksp"),
			helpTextStyle.Render("back"),
		)
	case stateAuthError:
		helpText = fmt.Sprintf("%s %s    %s %s",
			helpKeyStyle.Render("l"),
			helpTextStyle.Render("login"),
			helpKeyStyle.Render("bksp"),
			helpTextStyle.Render("back"),
		)
	default:
		helpText = fmt.Sprintf("%s %s    %s %s",
			helpKeyStyle.Render("enter"),
			helpTextStyle.Render("create"),
			helpKeyStyle.Render("alt+enter"),
			helpTextStyle.Render("new line"),
		)
	}
	help := helpBarStyle.Render(helpText)

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

func (m Model) viewFeatureList() string {
	header := renderGradientOnBg(" "+logoSmall, gradStops, "#1A1A1A", m.width)

	form := m.featureList.view()

	var helpText string
	if len(m.featureList.drafts) > 0 {
		helpText = fmt.Sprintf("%s %s    %s %s    %s %s    %s %s    %s %s",
			helpKeyStyle.Render("j/k"),
			helpTextStyle.Render("navigate"),
			helpKeyStyle.Render("enter"),
			helpTextStyle.Render("select"),
			helpKeyStyle.Render("n"),
			helpTextStyle.Render("new"),
			helpKeyStyle.Render("d"),
			helpTextStyle.Render("delete"),
			helpKeyStyle.Render("bksp"),
			helpTextStyle.Render("back"),
		)
	} else {
		helpText = fmt.Sprintf("%s %s    %s %s",
			helpKeyStyle.Render("n"),
			helpTextStyle.Render("new"),
			helpKeyStyle.Render("bksp"),
			helpTextStyle.Render("back"),
		)
	}
	help := helpBarStyle.Render(helpText)

	headerLines := strings.Count(header, "\n") + 1
	formLines := strings.Count(form, "\n") + 1
	helpHeight := 1

	var b strings.Builder
	b.WriteString(header + "\n")
	b.WriteString("\n")
	b.WriteString(form)

	usedLines := headerLines + 1 + formLines + helpHeight
	gap := m.height - usedLines
	if gap > 0 {
		b.WriteString(strings.Repeat("\n", gap))
	}
	b.WriteString(help)

	return b.String()
}
