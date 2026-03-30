package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type formState int

const (
	stateIdle formState = iota
	stateDecomposing
	stateDone
	stateError
	stateAuthError
)

var (
	inputLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Bold(true)

	statusOk = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66CC66"))

	statusErr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CC6666"))

	runeNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5A623")).
			Bold(true)

	runeSigStyle = lipgloss.NewStyle().
			Foreground(dim)

	testPassStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66CC66"))

	testFailStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CC6666"))
)

// API response types

type decomposeResponse struct {
	JobID string `json:"job_id"`
}

type jobResponse struct {
	ID     string          `json:"id"`
	Status string          `json:"status"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}

type proposedRune struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Signature     string   `json:"signature"`
	PositiveTests []string `json:"positive_tests"`
	NegativeTests []string `json:"negative_tests"`
}

type existingMatch struct {
	Name   string `json:"name"`
	Covers string `json:"covers"`
}

type decomposeResult struct {
	NewRunes      []proposedRune `json:"new_runes"`
	ExistingRunes []existingMatch `json:"existing_runes"`
	TreeOutput    string          `json:"tree_output"`
}

// Messages

type decomposeStartedMsg struct {
	jobID string
}

type decomposeErrorMsg struct {
	err error
}

type pollTickMsg struct{}

type decomposeDoneMsg struct {
	result decomposeResult
}

type loginDoneMsg struct {
	err error
}

type createFeatureModel struct {
	descInput textarea.Model
	state     formState
	port      int
	width     int
	jobID     string
	spinner   spinner.Model
	result    *decomposeResult
	errMsg  string
	authURL string
}

func newCreateFeatureModel(port, width, height int) createFeatureModel {
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

	inputWidth := width - 4
	if inputWidth < 40 {
		inputWidth = 40
	}
	taHeight := height - 20
	if taHeight < 3 {
		taHeight = 3
	}
	ta.SetWidth(inputWidth)
	ta.SetHeight(taHeight)
	ta.CharLimit = 2000

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623"))

	return createFeatureModel{
		descInput: ta,
		state:     stateIdle,
		port:      port,
		width:     width,
		spinner:   s,
	}
}

func (m *createFeatureModel) resize(width, height int) {
	m.width = width
	inputWidth := width - 4
	if inputWidth < 40 {
		inputWidth = 40
	}
	m.descInput.SetWidth(inputWidth)
	taHeight := height - 20
	if taHeight < 3 {
		taHeight = 3
	}
	m.descInput.SetHeight(taHeight)
}

func (m *createFeatureModel) submit() tea.Cmd {
	desc := strings.TrimSpace(m.descInput.Value())
	if desc == "" {
		m.errMsg = "description is required"
		m.state = stateError
		return nil
	}

	m.state = stateDecomposing
	m.descInput.Blur()

	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{"requirement": desc})
		resp, err := http.Post(
			fmt.Sprintf("http://localhost:%d/api/decompose", m.port),
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			return decomposeErrorMsg{err: err}
		}
		defer resp.Body.Close()

		var dr decomposeResponse
		if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
			return decomposeErrorMsg{err: err}
		}
		return decomposeStartedMsg{jobID: dr.JobID}
	}
}

func (m *createFeatureModel) pollJob() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return pollTickMsg{}
	})
}

func (m *createFeatureModel) checkJob() tea.Cmd {
	jobID := m.jobID
	port := m.port
	return func() tea.Msg {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/decompose/%s", port, jobID))
		if err != nil {
			return decomposeErrorMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		var job jobResponse
		if err := json.Unmarshal(data, &job); err != nil {
			return decomposeErrorMsg{err: err}
		}

		switch job.Status {
		case "completed":
			var result decomposeResult
			json.Unmarshal(job.Result, &result)
			return decomposeDoneMsg{result: result}
		case "failed":
			return decomposeErrorMsg{err: fmt.Errorf("%s", job.Error)}
		default:
			return pollTickMsg{}
		}
	}
}

func (m *createFeatureModel) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == stateIdle || m.state == stateError {
			if msg.String() == "enter" {
				return m.submit()
			}
		}
		if m.state == stateDone {
			if msg.String() == "enter" {
				m.state = stateIdle
				m.result = nil
				m.descInput.Reset()
				m.descInput.Focus()
				return nil
			}
		}
		if m.state == stateAuthError {
			if msg.String() == "l" {
				exe, _ := os.Executable()
				c := exec.Command(exe, "login")
				c.Stdin = os.Stdin
				return tea.ExecProcess(c, func(err error) tea.Msg {
					return loginDoneMsg{err: err}
				})
			}
		}

	case decomposeStartedMsg:
		m.jobID = msg.jobID
		return tea.Batch(m.spinner.Tick, m.pollJob())

	case pollTickMsg:
		if m.state == stateDecomposing {
			return m.checkJob()
		}

	case decomposeDoneMsg:
		m.state = stateDone
		m.result = &msg.result
		return nil

	case decomposeErrorMsg:
		if strings.Contains(msg.err.Error(), "auth error") || strings.Contains(msg.err.Error(), "token expired") {
			m.state = stateAuthError
			m.errMsg = msg.err.Error()
			return nil
		}
		m.state = stateError
		m.errMsg = msg.err.Error()
		m.descInput.Focus()
		return nil

	case loginDoneMsg:
		if msg.err != nil {
			m.state = stateAuthError
			m.errMsg = fmt.Sprintf("login failed: %v", msg.err)
			return nil
		}
		// Give the proxy's file watcher time to detect the new token
		m.state = stateIdle
		m.errMsg = ""
		m.authURL = ""
		m.descInput.Focus()
		return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return nil })

	case spinner.TickMsg:
		if m.state == stateDecomposing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return cmd
		}
	}

	if m.state == stateIdle || m.state == stateError {
		var cmd tea.Cmd
		m.descInput, cmd = m.descInput.Update(msg)
		return cmd
	}

	return nil
}

func (m *createFeatureModel) view(width int) string {
	switch m.state {
	case stateDecomposing:
		return m.viewDecomposing()
	case stateDone:
		return m.viewResult(width)
	case stateAuthError:
		return m.viewAuthError()
	default:
		return m.viewForm()
	}
}

func (m *createFeatureModel) viewForm() string {
	var b strings.Builder
	b.WriteString(inputLabel.Render("Describe your feature") + "\n\n")
	b.WriteString(m.descInput.View())

	if m.state == stateError && m.errMsg != "" {
		b.WriteString("\n\n")
		b.WriteString(statusErr.Render(m.errMsg))
	}

	return b.String()
}

func (m *createFeatureModel) viewAuthError() string {
	var b strings.Builder
	b.WriteString(statusErr.Render("Authentication required") + "\n\n")
	if m.authURL != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(dim).Render(m.errMsg) + "\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6A9FD9")).Bold(true).Render(m.authURL) + "\n")
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(dim).Render(m.errMsg) + "\n")
	}
	return b.String()
}

func (m *createFeatureModel) viewDecomposing() string {
	return m.spinner.View() + " Decomposing into runes..."
}

func (m *createFeatureModel) viewResult(width int) string {
	if m.result == nil {
		return ""
	}

	var b strings.Builder

	if len(m.result.NewRunes) > 0 {
		b.WriteString(inputLabel.Render(fmt.Sprintf("New Runes (%d)", len(m.result.NewRunes))) + "\n\n")
		for _, r := range m.result.NewRunes {
			b.WriteString(runeNameStyle.Render(r.Name))
			if r.Signature != "" {
				b.WriteString("  " + runeSigStyle.Render(r.Signature))
			}
			b.WriteString("\n")
			if r.Description != "" {
				b.WriteString("  " + r.Description + "\n")
			}
			for _, t := range r.PositiveTests {
				b.WriteString(testPassStyle.Render("  + ") + t + "\n")
			}
			for _, t := range r.NegativeTests {
				b.WriteString(testFailStyle.Render("  - ") + t + "\n")
			}
			b.WriteString("\n")
		}
	}

	if len(m.result.ExistingRunes) > 0 {
		b.WriteString(inputLabel.Render("Existing Runes") + "\n\n")
		for _, r := range m.result.ExistingRunes {
			b.WriteString("  " + runeNameStyle.Render(r.Name))
			if r.Covers != "" {
				b.WriteString("  " + runeSigStyle.Render(r.Covers))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(statusOk.Render(fmt.Sprintf("%d new runes proposed", len(m.result.NewRunes))))

	return b.String()
}
