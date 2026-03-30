package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var logo = `
 ██╗   ██╗ █████╗ ██╗     ██╗  ██╗██╗   ██╗██████╗ ██╗███████╗
 ██║   ██║██╔══██╗██║     ██║ ██╔╝╚██╗ ██╔╝██╔══██╗██║██╔════╝
 ██║   ██║███████║██║     █████╔╝  ╚████╔╝ ██████╔╝██║█████╗
 ╚██╗ ██╔╝██╔══██║██║     ██╔═██╗   ╚██╔╝  ██╔══██╗██║██╔══╝
  ╚████╔╝ ██║  ██║███████╗██║  ██╗   ██║   ██║  ██║██║███████╗
   ╚═══╝  ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚══════╝`

var (
	purple     = lipgloss.Color("#7B2FBE")
	darkPurple = lipgloss.Color("#3B1F6E")
	dim        = lipgloss.Color("#555555")
	white      = lipgloss.Color("#FAFAFA")
	cyan       = lipgloss.Color("#00CED1")

	logoStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true)

	taglineStyle = lipgloss.NewStyle().
			Foreground(dim).
			Italic(true)

	frameStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(darkPurple).
			Padding(1, 4)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	helpTextStyle = lipgloss.NewStyle().
			Foreground(dim)

	helpBarStyle = lipgloss.NewStyle().
			MarginTop(1)
)

// ExitMsg is sent when the TUI exits with a message to print.
type ExitMsg struct {
	Message string
}

type Model struct {
	width  int
	height int
}

func New() Model { return Model{} }

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m, tea.Sequence(
				tea.Println("run: valkyrie features create --name <name> --description <desc>"),
				tea.Quit,
			)
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	content := logoStyle.Render(logo) + "\n\n" +
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
