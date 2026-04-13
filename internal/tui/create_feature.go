package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Prompt = lipgloss.NewStyle()
	ta.FocusedStyle.EndOfBuffer = lipgloss.NewStyle()
	ta.FocusedStyle.Base = lipgloss.NewStyle()
	ta.BlurredStyle.CursorLine = lipgloss.NewStyle()
	ta.BlurredStyle.Prompt = lipgloss.NewStyle()
	ta.BlurredStyle.EndOfBuffer = lipgloss.NewStyle()
	ta.BlurredStyle.Base = lipgloss.NewStyle()
	ta.KeyMap.InsertNewline.SetKeys("alt+enter")
	ta.Focus()
	ta.CharLimit = 2000

	h := newHelpModel()
	h.Width = width - viewPadX*2

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
	m.help.Width = width - viewPadX*2
}

func (m createFeatureModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m createFeatureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
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

func (m createFeatureModel) View() string {
	innerWidth := m.width - viewPadX*2
	header := renderGradientOnBg(" "+logoSmall, gradStops, "#1A1A1A", innerWidth)

	var form strings.Builder
	form.WriteString(inputLabel.Render("Describe your feature") + "\n\n")
	form.WriteString(m.descInput.View())

	helpBar := helpBarStyle.Render(m.help.View(formKeyMap{}))

	body := header + "\n\n" + form.String()
	bodyBlock := lipgloss.NewStyle().Height(m.height - 1).Render(body)

	content := bodyBlock + "\n" + helpBar
	return lipgloss.NewStyle().PaddingLeft(viewPadX).PaddingRight(viewPadX).Render(content)
}
