package tui

import (
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type createFeatureModel struct {
	width     int
	height    int
	descInput textarea.Model
	help      help.Model
}

func newCreateFeatureModel(width, height int) createFeatureModel {
	ta := textarea.New()
	ta.Placeholder = "Describe the feature..."
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.EndOfBufferCharacter = ' '
	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	s.Focused.Prompt = lipgloss.NewStyle()
	s.Focused.EndOfBuffer = lipgloss.NewStyle()
	s.Focused.Base = lipgloss.NewStyle()
	s.Blurred.CursorLine = lipgloss.NewStyle()
	s.Blurred.Prompt = lipgloss.NewStyle()
	s.Blurred.EndOfBuffer = lipgloss.NewStyle()
	s.Blurred.Base = lipgloss.NewStyle()
	ta.SetStyles(s)
	ta.KeyMap.InsertNewline.SetKeys("alt+enter")
	ta.Focus()
	ta.CharLimit = 2000

	h := newHelpModel()
	h.SetWidth(width - viewPadX*2)

	m := createFeatureModel{
		width:     width,
		height:    height,
		descInput: ta,
		help:      h,
	}
	m.resize(width, height)
	return m
}

func (m *createFeatureModel) resize(width, height int) {
	m.width = width
	m.height = height
	m.descInput.SetWidth(max(width-viewPadX*2-4, 40))
	m.descInput.SetHeight(max(height-8, 3))
	m.help.SetWidth(width - viewPadX*2)
}

func (m createFeatureModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m createFeatureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return model{width: m.width, height: m.height, help: newHelpModel()}, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.descInput, cmd = m.descInput.Update(msg)
	return m, cmd
}

func (m createFeatureModel) View() tea.View {
	innerWidth := m.width - viewPadX*2
	header := renderGradientOnBg(" "+logoSmall, gradStops, "#1A1A1A", innerWidth)

	var form strings.Builder
	form.WriteString(inputLabel.Render("Describe your feature") + "\n\n")
	form.WriteString(m.descInput.View())

	helpBar := helpBarStyle.Render(m.help.View(formKeyMap{}))

	body := header + "\n\n" + form.String()
	bodyBlock := lipgloss.NewStyle().Height(m.height - 1).Render(body)

	content := bodyBlock + "\n" + helpBar
	rendered := lipgloss.NewStyle().PaddingLeft(viewPadX).PaddingRight(viewPadX).Render(content)
	v := tea.NewView(rendered)
	v.AltScreen = true
	v.BackgroundColor = bgMain
	return v
}
