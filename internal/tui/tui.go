package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/chrishayen/odek/internal/draft"
	runepkg "github.com/chrishayen/odek/internal/rune"
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

	viewPadX = 1
)

// Key bindings
var (
	keyCreate     = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "create"))
	keyNewLine    = key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "new line"))
	keySubmit     = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit"))
	keyCancel     = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
	keyBack       = key.NewBinding(key.WithKeys("backspace"), key.WithHelp("bksp", "back"))
	keyQuit       = key.NewBinding(key.WithKeys("backspace"), key.WithHelp("bksp", "quit"))
	keyNavigate        = key.NewBinding(key.WithKeys("j", "k"), key.WithHelp("j/k", "navigate"))
	keyComment         = key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "comment"))
	keyApprove         = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "approve"))
	keySubmitRefine    = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit"))
	keyHydrate         = key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "hydrate"))
	keyLogin           = key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "login"))
	keyFeatures        = key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "features"))
)

// Help keymaps per state
type splashKeyMap struct{}

func (splashKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keyCreate, keyFeatures, keyQuit}
}
func (splashKeyMap) FullHelp() [][]key.Binding { return nil }

type formKeyMap struct{}

func (formKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keyCreate, keyNewLine, keyCancel}
}
func (formKeyMap) FullHelp() [][]key.Binding { return nil }

type doneKeyMap struct{ hasComments bool }

func (k doneKeyMap) ShortHelp() []key.Binding {
	bindings := []key.Binding{keyComment, keyApprove}
	if k.hasComments {
		bindings = append(bindings, keySubmitRefine)
	}
	bindings = append(bindings, keyBack)
	return bindings
}
func (doneKeyMap) FullHelp() [][]key.Binding { return nil }

type approvedKeyMap struct{}

func (approvedKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keyHydrate, keyBack}
}
func (approvedKeyMap) FullHelp() [][]key.Binding { return nil }

type refiningKeyMap struct{}

func (refiningKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keySubmit, keyCancel}
}
func (refiningKeyMap) FullHelp() [][]key.Binding { return nil }

type decomposingKeyMap struct{}

func (decomposingKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keyBack}
}
func (decomposingKeyMap) FullHelp() [][]key.Binding { return nil }

type authErrorKeyMap struct{}

func (authErrorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keyLogin, keyBack}
}
func (authErrorKeyMap) FullHelp() [][]key.Binding { return nil }

func newHelpModel() help.Model {
	h := help.New()
	h.Styles.ShortKey = helpKeyStyle
	h.Styles.ShortDesc = helpTextStyle
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(dim)
	h.ShortSeparator = "    "
	return h
}

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

type Model struct {
	width        int
	height       int
	screen       screen
	port         int
	registryPath string
	draftStore   *draft.Store
	runeStore    *runepkg.Store
	createForm   createFeatureModel
	featureList  featureListModel
	help         help.Model
}

func New(port int, registryPath string, runeStore *runepkg.Store, draftStore *draft.Store) Model {
	return Model{
		screen:       screenSplash,
		port:         port,
		registryPath: registryPath,
		draftStore:   draftStore,
		runeStore:    runeStore,
		help:         newHelpModel(),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		innerWidth := msg.Width - viewPadX*2
		m.help.Width = innerWidth
		if m.screen == screenCreateFeature {
			m.createForm.resize(innerWidth, msg.Height)
		}
		if m.screen == screenFeatureList {
			m.featureList.list.SetSize(innerWidth, msg.Height-3)
		}
		return m, nil

	case goBackMsg:
		switch m.screen {
		case screenCreateFeature:
			if m.createForm.state != stateApproved {
				m.createForm.saveDraft()
			}
			m.screen = screenSplash
			return m, nil
		case screenFeatureList:
			m.screen = screenSplash
			return m, nil
		}

	case draftSelectedMsg:
		m.createForm = newCreateFeatureModelFromDraft(m.port, m.width-viewPadX*2, m.height, m.draftStore, msg.draft)

		m.screen = screenCreateFeature
		if m.createForm.state == stateDone {
			return m, nil
		}
		return m, m.createForm.descInput.Focus()

	case featureSelectedMsg:
		m.createForm = newCreateFeatureModel(m.port, m.width-viewPadX*2, m.height, m.draftStore)
		m.screen = screenCreateFeature
		// Load runes directly — featureLoadedMsg will set stateApproved
		return m, loadFeatureRunes(msg.name, m.port)

	case newFeatureMsg:
		m.createForm = newCreateFeatureModel(m.port, m.width-viewPadX*2, m.height, m.draftStore)
		m.screen = screenCreateFeature
		return m, m.createForm.descInput.Focus()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			if m.screen == screenSplash {
				return m, tea.Quit
			}
		case "f":
			if m.screen == screenSplash {
				m.featureList = newFeatureListModel(m.draftStore, m.runeStore, m.width-viewPadX*2, m.height)
				m.screen = screenFeatureList
				return m, nil
			}
		case "enter":
			if m.screen == screenSplash {
				m.createForm = newCreateFeatureModel(m.port, m.width-viewPadX*2, m.height, m.draftStore)
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
	helpBar := helpBarStyle.Render(m.help.View(splashKeyMap{}))

	block := lipgloss.JoinVertical(lipgloss.Center, framed, helpBar)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		block,
	)
}

func (m Model) createFeatureKeyMap() help.KeyMap {
	switch m.createForm.state {
	case stateDone:
		return doneKeyMap{hasComments: m.createForm.hasComments()}
	case stateRefining:
		return refiningKeyMap{}
	case stateDecomposing, stateHydrating:
		return decomposingKeyMap{}
	case stateApproved:
		return approvedKeyMap{}
	case stateAuthError:
		return authErrorKeyMap{}
	default:
		return formKeyMap{}
	}
}

func (m Model) viewCreateFeature() string {
	innerWidth := m.width - viewPadX*2
	header := renderGradientOnBg(" "+logoSmall, gradStops, "#1A1A1A", innerWidth)
	form := m.createForm.view(innerWidth)
	helpBar := helpBarStyle.Render(m.help.View(m.createFeatureKeyMap()))

	body := header + "\n\n" + form
	bodyBlock := lipgloss.NewStyle().Height(m.height - 1).Render(body)

	content := bodyBlock + "\n" + helpBar
	return lipgloss.NewStyle().PaddingLeft(viewPadX).PaddingRight(viewPadX).Render(content)
}

func (m Model) viewFeatureList() string {
	innerWidth := m.width - viewPadX*2
	header := renderGradientOnBg(" "+logoSmall, gradStops, "#1A1A1A", innerWidth)
	content := m.featureList.view()

	body := header + "\n\n" + content
	return lipgloss.NewStyle().PaddingLeft(viewPadX).PaddingRight(viewPadX).Height(m.height).Render(body)
}
