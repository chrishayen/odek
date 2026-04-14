package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"shotgun.dev/odek/internal/effort"
	openai "shotgun.dev/odek/openai"
)

type kanjiTickMsg struct{}

func kanjiTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg { return kanjiTickMsg{} })
}

const createFeatureChromeHeight = 6

type createFeatureModel struct {
	ctx         context.Context
	client      *openai.Client
	width       int
	height      int
	chat        chatModel
	kanjiOffset int
	inSplit     bool
}

func newCreateFeatureModel(ctx context.Context, client *openai.Client, width, height int) createFeatureModel {
	m := createFeatureModel{
		ctx:    ctx,
		client: client,
		width:  width,
		height: height,
	}

	sendHandler := makeFeatureSendHandler(ctx, client)
	chatWidth := max(width-viewPadX*2, 20)
	chatHeight := max(height-createFeatureChromeHeight, 5)
	m.chat = newChatModel(
		chatWidth,
		chatHeight,
		withChatPlaceholder("Describe the feature..."),
		withChatSendHandler(sendHandler),
		withChatWelcome("Describe your feature. I'll estimate effort first, then we can iterate."),
	)
	return m
}

func (m *createFeatureModel) resize(width, height int) {
	m.width = width
	m.height = height
	chatWidth := max(width-viewPadX*2, 20)
	chatHeight := max(height-createFeatureChromeHeight, 5)
	m.chat.SetSize(chatWidth, chatHeight)
}

func (m *createFeatureModel) Focus() tea.Cmd {
	return m.chat.Focus()
}

func (m *createFeatureModel) Blur() {
	m.chat.Blur()
}

func (m *createFeatureModel) SetChatInput(s string) {
	m.chat.SetInput(s)
}

func (m *createFeatureModel) SetInSplit(v bool) {
	m.inSplit = v
}

func (m createFeatureModel) Init() tea.Cmd {
	return tea.Batch(m.chat.Init(), kanjiTick())
}

func (m createFeatureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case kanjiTickMsg:
		if m.chat.Busy() {
			m.kanjiOffset += 2
		}
		return m, kanjiTick()

	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			m.chat = m.chat.ScrollLines(-1)
			return m, nil
		case "down":
			m.chat = m.chat.ScrollLines(1)
			return m, nil
		case "esc":
			if m.chat.Busy() {
				return m, nil
			}
			return model{
				width:  m.width,
				height: m.height,
				help:   newHelpModel(),
				ctx:    m.ctx,
				client: m.client,
			}, nil
		case "ctrl+enter", "ctrl+s":
			if m.chat.Busy() {
				return m, nil
			}
			pin := renderFeaturePin()
			if m.width >= splitPaneMinWidth {
				leftW, rightW := splitWidths(m.width)
				left := m
				left.resize(leftW, m.height)
				right := newFeatureDecompModel(rightW, m.height, pin)
				split := newSplitPaneModel(left, right, m.width, m.height)
				return split, split.Init()
			}
			dest := newFeatureDecompModel(m.width, m.height, pin)
			t := newTransition(m, dest, m.width, m.height, pin)
			return t, t.Init()
		default:
			m.kanjiOffset += 2
		}
	}

	var cmd tea.Cmd
	m.chat, cmd = m.chat.Update(msg)
	return m, cmd
}

func (m createFeatureModel) View() tea.View {
	innerWidth := m.width - viewPadX*2

	logoText := "ODEK "
	logoStyle := lipgloss.NewStyle().Foreground(accent).Background(bgMain).Bold(true)
	padStyle := lipgloss.NewStyle().Background(bgMain)
	arrowStyle := lipgloss.NewStyle().Foreground(fgDim).Background(bgMain)
	statusStr := m.chat.StatusView()
	statusWidth := lipgloss.Width(statusStr)
	upArrow := ""
	if m.chat.CanScrollUp() {
		upArrow = arrowStyle.Render("↑")
	}
	logoPad := max(innerWidth-len(logoText)-statusWidth-lipgloss.Width(upArrow), 0)
	header := logoStyle.Render(logoText) + padStyle.Render(strings.Repeat(" ", logoPad)) + statusStr + upArrow

	helpBar := renderFormHelpBar(innerWidth, m.chat.CanScrollDown(), m.inSplit)

	scrollOff := m.chat.ViewportYOffset() * 2
	kanjiLine1 := renderKanjiLine(innerWidth, 2, m.kanjiOffset+scrollOff)
	kanjiLine2 := renderKanjiLine(innerWidth, 3, -(m.kanjiOffset + scrollOff))
	body := header + "\n\n" + kanjiLine1 + "\n" + kanjiLine2 + "\n" + m.chat.View()
	bodyBlock := lipgloss.NewStyle().Height(m.height - 1).Render(body)

	content := bodyBlock + "\n" + helpBar
	rendered := lipgloss.NewStyle().PaddingLeft(viewPadX).PaddingRight(viewPadX).Render(content)
	v := tea.NewView(rendered)
	v.AltScreen = true
	v.BackgroundColor = bgMain
	return v
}

func renderKanjiLine(width, row, offset int) string {
	kanjiStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
	bgStyle := lipgloss.NewStyle().Background(bgMain)
	var kb strings.Builder
	cells := 0
	for cells+2 <= width {
		kb.WriteRune(kanjiAt(row, cells+offset))
		cells += 2
	}
	line := kanjiStyle.Render(kb.String())
	if cells < width {
		line += bgStyle.Render(" ")
	}
	return line
}

func renderFormHelpBar(width int, scrollDown, showTabSwitch bool) string {
	keyStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(fgBody).Background(bgMain)
	sepStyle := lipgloss.NewStyle().Foreground(fgDim).Background(bgMain)
	padStyle := lipgloss.NewStyle().Background(bgMain)

	type binding struct{ key, desc string }
	bindings := []binding{
		{"enter", "send"},
		{"alt+enter", "new line"},
		{"↑/↓", "scroll"},
	}
	if showTabSwitch {
		bindings = append(bindings, binding{"tab", "runes"})
	}

	var b strings.Builder
	for i, bind := range bindings {
		if i > 0 {
			b.WriteString(sepStyle.Render("  •  "))
		}
		b.WriteString(keyStyle.Render(bind.key))
		b.WriteString(padStyle.Render(" "))
		b.WriteString(descStyle.Render(bind.desc))
	}
	content := b.String()
	trailing := padStyle.Render(" ")
	if scrollDown {
		trailing = lipgloss.NewStyle().Foreground(fgDim).Background(bgMain).Render("↓")
	}
	pad := max(width-lipgloss.Width(content)-1-lipgloss.Width(trailing), 0)
	return padStyle.Render(" ") + content + padStyle.Render(strings.Repeat(" ", pad)) + trailing
}

// makeFeatureSendHandler returns a SendHandler that routes the first submission
// through effort.Estimate and subsequent submissions through client.Chat,
// carrying the feature description and the initial effort estimate forward as
// system context.
func makeFeatureSendHandler(ctx context.Context, client *openai.Client) SendHandler {
	return func(history []chatMessage, userInput string, id int) tea.Cmd {
		if client == nil {
			return func() tea.Msg {
				return chatErrMsg{id: id, err: fmt.Errorf("chat offline: no client configured")}
			}
		}
		if isFirstTurn(history) {
			return func() tea.Msg {
				res, err := effort.Estimate(ctx, client, userInput)
				if err != nil {
					return chatErrMsg{id: id, err: err}
				}
				return chatReplyMsg{
					id:       id,
					content:  res.Reason,
					headline: fmt.Sprintf("Effort: %d/5", res.Level),
				}
			}
		}
		return func() tea.Msg {
			req := buildFollowupRequest(history)
			resp, err := client.Chat(ctx, req)
			if err != nil {
				return chatErrMsg{id: id, err: err}
			}
			if len(resp.Choices) == 0 {
				return chatErrMsg{id: id, err: fmt.Errorf("empty response")}
			}
			return chatReplyMsg{
				id:      id,
				content: resp.Choices[0].Message.Content,
			}
		}
	}
}

// isFirstTurn reports whether the just-submitted message is the very first
// user turn — i.e. history contains only that user message, with no prior
// assistant reply.
func isFirstTurn(history []chatMessage) bool {
	userCount := 0
	for _, msg := range history {
		if msg.role == roleUser {
			userCount++
		}
		if msg.role == roleAssistant {
			return false
		}
	}
	return userCount == 1
}

func buildFollowupRequest(history []chatMessage) *openai.ChatCompletionRequest {
	featureDesc := ""
	effortHeadline := ""
	effortReason := ""
	for _, msg := range history {
		if featureDesc == "" && msg.role == roleUser {
			featureDesc = msg.content
			continue
		}
		if effortHeadline == "" && msg.role == roleAssistant && msg.headline != "" {
			effortHeadline = msg.headline
			effortReason = msg.content
		}
	}

	system := "You are Odek, a software feature design collaborator. "
	system += fmt.Sprintf("The user is iterating on this feature: %q. ", featureDesc)
	if effortHeadline != "" {
		system += fmt.Sprintf("Prior effort estimate: %s — %s. ", effortHeadline, effortReason)
	}
	system += "Respond concisely and practically."

	msgs := []openai.ChatMessage{{Role: "system", Content: system}}
	for _, msg := range history {
		role := "user"
		content := msg.content
		if msg.role == roleAssistant {
			role = "assistant"
			if msg.headline != "" {
				content = msg.headline + " — " + msg.content
			}
		}
		msgs = append(msgs, openai.ChatMessage{Role: role, Content: content})
	}

	return &openai.ChatCompletionRequest{
		Model:    "default",
		Messages: msgs,
	}
}
