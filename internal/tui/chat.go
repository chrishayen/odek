package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/chrishayen/odek/internal/chat"
	"github.com/chrishayen/odek/internal/claude"
)

// chatSignal is emitted by the chat to its parent.
type chatSignal struct {
	Type string // "refine_feature", "refine_rune"
	Data string
}

// Chat messages
type chatResponseMsg struct{ answer string }
type chatErrorMsg struct{ err error }
type chatStartedMsg struct{ jobID string }
type chatPollTickMsg struct{}

var (
	chatUserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6A9FD9")).
			Bold(true)

	chatAssistantStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC"))

	chatContextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Italic(true)

	chatDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#444444"))
)

type chatModel struct {
	session  *chat.Session
	store    *chat.Store
	input    textinput.Model
	spinner  spinner.Model
	port     int
	width    int
	height   int
	scroll   int
	waiting  bool
	jobID    string
	errMsg   string
	signals  []chatSignal // pending signals for parent to consume
}

func newChatModel(store *chat.Store, session *chat.Session, port, width, height int) chatModel {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Width = width - 4
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623"))

	return chatModel{
		session: session,
		store:   store,
		input:   ti,
		spinner: s,
		port:    port,
		width:   width,
		height:  height,
	}
}

func (m *chatModel) resize(width, height int) {
	m.width = width
	m.height = height
	m.input.Width = width - 4
}

// consumeSignals returns and clears pending signals.
func (m *chatModel) consumeSignals() []chatSignal {
	s := m.signals
	m.signals = nil
	return s
}

func (m *chatModel) sendMessage() tea.Cmd {
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		return nil
	}

	m.input.SetValue("")
	m.session.AddUser(text)
	m.waiting = true

	// Save session after adding user message
	if m.store != nil {
		m.store.Save(m.session)
	}

	// Build API messages from session history
	msgs := make([]claude.ChatMessage, len(m.session.Messages))
	for i, msg := range m.session.Messages {
		msgs[i] = claude.ChatMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	ctx := m.session.Context.FormatSystem()
	port := m.port

	return func() tea.Msg {
		body, _ := json.Marshal(map[string]any{
			"messages": msgs,
			"context":  ctx,
		})
		resp, err := http.Post(
			fmt.Sprintf("http://localhost:%d/api/chat", port),
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			return chatErrorMsg{err: err}
		}
		defer resp.Body.Close()
		var dr decomposeResponse
		if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
			return chatErrorMsg{err: err}
		}
		return chatStartedMsg{jobID: dr.JobID}
	}
}

func (m *chatModel) pollChat() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return chatPollTickMsg{}
	})
}

func (m *chatModel) checkChat() tea.Cmd {
	jobID := m.jobID
	port := m.port
	return func() tea.Msg {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/chat/%s", port, jobID))
		if err != nil {
			return chatErrorMsg{err: err}
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		var job jobResponse
		if err := json.Unmarshal(data, &job); err != nil {
			return chatErrorMsg{err: err}
		}
		switch job.Status {
		case "completed":
			var answer string
			json.Unmarshal(job.Result, &answer)
			return chatResponseMsg{answer: answer}
		case "failed":
			return chatErrorMsg{err: fmt.Errorf("%s", job.Error)}
		default:
			return chatPollTickMsg{}
		}
	}
}

func (m *chatModel) parseSignals(answer string) {
	for _, line := range strings.Split(answer, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "SIGNAL:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "SIGNAL:"))
			sigType := "refine_feature"
			if m.session.Context.RuneName != "" {
				sigType = "refine_rune"
			}
			m.signals = append(m.signals, chatSignal{Type: sigType, Data: data})
		}
	}
}

func (m *chatModel) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.waiting {
			return nil
		}
		switch msg.String() {
		case "enter":
			return tea.Batch(m.spinner.Tick, m.sendMessage())
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return cmd

	case chatStartedMsg:
		m.jobID = msg.jobID
		return tea.Batch(m.spinner.Tick, m.pollChat())

	case chatPollTickMsg:
		if m.waiting {
			return m.checkChat()
		}

	case chatResponseMsg:
		m.session.AddAssistant(msg.answer)
		m.waiting = false
		m.errMsg = ""

		// Save session after receiving response
		if m.store != nil {
			m.store.Save(m.session)
		}

		// Check for signals
		m.parseSignals(msg.answer)

		// Scroll to bottom
		m.scrollToBottom()
		return nil

	case chatErrorMsg:
		m.session.AddAssistant("Error: " + msg.err.Error())
		m.waiting = false
		m.errMsg = msg.err.Error()
		if m.store != nil {
			m.store.Save(m.session)
		}
		m.scrollToBottom()
		return nil

	case spinner.TickMsg:
		if m.waiting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return cmd
		}
	}

	return nil
}

func (m *chatModel) scrollToBottom() {
	lines := m.renderMessageLines()
	viewHeight := m.height - 6 // context header + input + borders
	if len(lines) > viewHeight {
		m.scroll = len(lines) - viewHeight
	} else {
		m.scroll = 0
	}
}

func (m *chatModel) renderMessageLines() []string {
	var lines []string
	maxW := m.width - 4

	for _, msg := range m.session.Messages {
		if msg.Role == chat.RoleUser {
			wrapped := wrapText("You: "+msg.Content, maxW)
			for _, l := range strings.Split(wrapped, "\n") {
				lines = append(lines, chatUserStyle.Render(l))
			}
		} else {
			wrapped := wrapText(msg.Content, maxW)
			for _, l := range strings.Split(wrapped, "\n") {
				lines = append(lines, chatAssistantStyle.Render(l))
			}
		}
		lines = append(lines, "") // spacing between messages
	}

	if m.waiting {
		lines = append(lines, m.spinner.View()+" thinking...")
	}

	return lines
}

func (m *chatModel) view() string {
	var b strings.Builder

	// Context header
	ctx := m.session.Context
	var ctxParts []string
	if ctx.FeatureName != "" {
		ctxParts = append(ctxParts, "feature:"+ctx.FeatureName)
	}
	if ctx.RuneName != "" {
		ctxParts = append(ctxParts, "rune:"+ctx.RuneName)
	}
	if len(ctxParts) > 0 {
		b.WriteString(chatContextStyle.Render("  "+strings.Join(ctxParts, "  ")) + "\n")
	}
	b.WriteString(chatDividerStyle.Render(strings.Repeat("─", m.width)) + "\n")

	// Messages area
	msgLines := m.renderMessageLines()
	viewHeight := m.height - 6
	if viewHeight < 3 {
		viewHeight = 3
	}

	// Apply scroll
	start := m.scroll
	if start < 0 {
		start = 0
	}
	if start > len(msgLines) {
		start = len(msgLines)
	}
	visible := msgLines[start:]
	if len(visible) > viewHeight {
		visible = visible[:viewHeight]
	}

	for _, line := range visible {
		b.WriteString("  " + line + "\n")
	}

	// Pad to fill space
	used := len(visible)
	for i := used; i < viewHeight; i++ {
		b.WriteString("\n")
	}

	// Input bar
	b.WriteString(chatDividerStyle.Render(strings.Repeat("─", m.width)) + "\n")
	b.WriteString("  " + m.input.View())

	return b.String()
}

// wrapText wraps text to the given width.
func wrapText(s string, width int) string {
	if width <= 0 {
		return s
	}
	var result strings.Builder
	for _, paragraph := range strings.Split(s, "\n") {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			continue
		}
		lineLen := 0
		for i, w := range words {
			wl := len(w)
			if i > 0 && lineLen+1+wl > width {
				result.WriteString("\n")
				lineLen = 0
			} else if i > 0 {
				result.WriteString(" ")
				lineLen++
			}
			result.WriteString(w)
			lineLen += wl
		}
	}
	return result.String()
}
