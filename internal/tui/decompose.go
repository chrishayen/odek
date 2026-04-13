package tui

import (
	"context"
	"fmt"
	"os"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"shotgun.dev/odek/decompose"
	openai "shotgun.dev/odek/openai"
)

type formState int

const (
	stateIdle formState = iota
	stateDecomposing
	stateError
)

var statusErr = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

type decomposeDoneMsg struct{ result *decompose.Decomposition }
type decomposeErrorMsg struct{ err error }

type decomposeModel struct {
	ctx          context.Context
	client       *openai.Client
	systemPrompt string
	feature      string
	width        int
	height       int
	state        formState
	errMsg       string
	spinner      spinner.Model
}

func newDecomposeModel(ctx context.Context, client *openai.Client, systemPrompt, feature string) decomposeModel {
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	return decomposeModel{
		ctx:          ctx,
		client:       client,
		systemPrompt: systemPrompt,
		feature:      feature,
		width:        120,
		height:       40,
		state:        stateIdle,
		spinner:      sp,
	}
}

func RunDecompose(ctx context.Context, client *openai.Client, systemPrompt, feature string) {
	m := newDecomposeModel(ctx, client, systemPrompt, feature)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m decomposeModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.startDecompose())
}

func (m decomposeModel) startDecompose() tea.Cmd {
	feature := m.feature
	client := m.client
	systemPrompt := m.systemPrompt
	ctx := m.ctx

	return func() tea.Msg {
		result, err := decompose.DecomposeStructured(ctx, client, systemPrompt, feature)
		if err != nil {
			return decomposeErrorMsg{err: err}
		}
		return decomposeDoneMsg{result: result.Decomposition}
	}
}

func (m decomposeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case decomposeDoneMsg:
		// TODO: real decompositions don't yet carry responsibility data;
		// for now we hand the mock groups to the new screen regardless.
		_ = msg.result
		d, groups := mockDecomposition()
		next := newDecompositionModel(d, groups)
		next.width = m.width
		next.height = m.height
		return next, next.Init()

	case decomposeErrorMsg:
		m.state = stateError
		m.errMsg = msg.err.Error()
		return m, nil

	case spinner.TickMsg:
		if m.state == stateDecomposing || m.state == stateIdle {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	switch m.state {
	case stateIdle, stateDecomposing:
		return m.updateLoadingState(msg)
	case stateError:
		return m.updateErrorState(msg)
	}
	return m, nil
}

func (m decomposeModel) updateLoadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyPressMsg); ok {
		switch km.String() {
		case "esc", "backspace", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m decomposeModel) updateErrorState(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyPressMsg); ok {
		switch km.String() {
		case "esc", "backspace", "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m decomposeModel) View() tea.View {
	var content string
	switch m.state {
	case stateIdle, stateDecomposing:
		content = m.viewLoading()
	case stateError:
		content = m.viewError()
	}
	v := tea.NewView(content)
	v.AltScreen = true
	v.BackgroundColor = bgMain
	return v
}

func (m decomposeModel) viewLoading() string {
	text := "Decomposing into runes..."
	return m.spinner.View() + " " + text
}

func (m decomposeModel) viewError() string {
	return statusErr.Render("Error: "+m.errMsg) + "\n\nPress backspace to exit"
}

