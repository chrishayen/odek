package tui

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"shotgun.dev/odek/internal/effort"
	openai "shotgun.dev/odek/openai"
)

type effortState int

const (
	effortIdle effortState = iota
	effortEstimating
	effortDone
	effortFailed
)

type effortDoneMsg struct{ result effort.Result }
type effortErrMsg struct{ err error }

type createFeatureModel struct {
	ctx        context.Context
	client     *openai.Client
	width      int
	height     int
	descInput  textarea.Model
	help       help.Model
	spinner    spinner.Model
	state      effortState
	result     effort.Result
	estimateOf string
	errMsg     string
}

func newCreateFeatureModel(ctx context.Context, client *openai.Client, width, height int) createFeatureModel {
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

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(accent)

	m := createFeatureModel{
		ctx:       ctx,
		client:    client,
		width:     width,
		height:    height,
		descInput: ta,
		help:      h,
		spinner:   sp,
	}
	m.resize(width, height)
	return m
}

func (m *createFeatureModel) resize(width, height int) {
	m.width = width
	m.height = height
	m.descInput.SetWidth(max(width-viewPadX*2-4, 40))
	m.descInput.SetHeight(max(height-12, 3))
	m.help.SetWidth(width - viewPadX*2)
}

func (m createFeatureModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}

func (m createFeatureModel) startEstimate(req string) tea.Cmd {
	ctx := m.ctx
	client := m.client
	return func() tea.Msg {
		res, err := effort.Estimate(ctx, client, req)
		if err != nil {
			return effortErrMsg{err: err}
		}
		return effortDoneMsg{result: res}
	}
}

func (m createFeatureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case effortDoneMsg:
		m.state = effortDone
		m.result = msg.result
		return m, nil

	case effortErrMsg:
		m.state = effortFailed
		m.errMsg = msg.err.Error()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return model{
				width:  m.width,
				height: m.height,
				help:   newHelpModel(),
				ctx:    m.ctx,
				client: m.client,
			}, nil
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			req := strings.TrimSpace(m.descInput.Value())
			if req == "" || m.state == effortEstimating || m.client == nil {
				return m, nil
			}
			m.state = effortEstimating
			m.estimateOf = req
			m.errMsg = ""
			return m, m.startEstimate(req)
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
	form.WriteString("\n")
	form.WriteString(m.renderEffortBlock())

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

var (
	effortLabelStyle  = lipgloss.NewStyle().Foreground(fgBright).Background(accent).Bold(true).Padding(0, 1)
	effortReasonStyle = lipgloss.NewStyle().Foreground(fgBody).Padding(0, 1)
	effortErrorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func (m createFeatureModel) renderEffortBlock() string {
	switch m.state {
	case effortEstimating:
		return "\n" + m.spinner.View() + " estimating effort..."
	case effortDone:
		label := effortLabelStyle.Render(fmt.Sprintf("effort %d/5", m.result.Level))
		return "\n" + label + effortReasonStyle.Render(m.result.Reason)
	case effortFailed:
		return "\n" + effortErrorStyle.Render("effort estimate failed: "+m.errMsg)
	default:
		return ""
	}
}
