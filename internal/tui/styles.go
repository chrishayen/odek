package tui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

const viewPadX = 1

var logoSmall = "ODEK"

var (
	border  = lipgloss.Color("241")
	dim     = lipgloss.Color("241")
	helpKey = lipgloss.Color("99")

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(helpKey).
			Bold(true)

	helpTextStyle = lipgloss.NewStyle().
			Foreground(dim)

	helpBarStyle = lipgloss.NewStyle().PaddingLeft(5)
)

var (
	keyNew           = key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new"))
	keySend          = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "send"))
	keyNewLine       = key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "new line"))
	keyScroll        = key.NewBinding(key.WithKeys("pgup", "pgdown"), key.WithHelp("pgup/pgdn", "scroll"))
	keyCancel        = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
	keyQuit          = key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit"))
	keyDecomposition = key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "decomposition"))
)

type splashKeyMap struct{}

func (splashKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keyNew, keyDecomposition, keyQuit}
}
func (splashKeyMap) FullHelp() [][]key.Binding { return nil }

type formKeyMap struct{}

func (formKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keySend, keyNewLine, keyScroll, keyCancel}
}
func (formKeyMap) FullHelp() [][]key.Binding { return nil }

func newHelpModel() help.Model {
	h := help.New()
	h.Styles.ShortKey = helpKeyStyle
	h.Styles.ShortDesc = helpTextStyle
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(dim)
	h.ShortSeparator = "    "
	return h
}
