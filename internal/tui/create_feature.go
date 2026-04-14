package tui

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"shotgun.dev/odek/internal/effort"
	openai "shotgun.dev/odek/openai"
)

const createFeatureChromeHeight = 10

type createFeatureModel struct {
	ctx    context.Context
	client *openai.Client
	width  int
	height int
	help   help.Model
	chat   chatModel
}

func newCreateFeatureModel(ctx context.Context, client *openai.Client, width, height int) createFeatureModel {
	m := createFeatureModel{
		ctx:    ctx,
		client: client,
		width:  width,
		height: height,
		help:   newHelpModel(),
	}
	m.help.SetWidth(width - viewPadX*2)

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
	m.help.SetWidth(width - viewPadX*2)
	chatWidth := max(width-viewPadX*2, 20)
	chatHeight := max(height-createFeatureChromeHeight, 5)
	m.chat.SetSize(chatWidth, chatHeight)
}

func (m createFeatureModel) Init() tea.Cmd {
	return m.chat.Init()
}

func (m createFeatureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
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
		}
	}

	var cmd tea.Cmd
	m.chat, cmd = m.chat.Update(msg)
	return m, cmd
}

func (m createFeatureModel) View() tea.View {
	innerWidth := m.width - viewPadX*2
	header := renderGradientOnBg(" "+logoSmall, gradStops, "#1A1A1A", innerWidth)

	helpBar := helpBarStyle.Render(m.help.View(formKeyMap{}))

	body := header + "\n\n" + m.chat.View()
	bodyBlock := lipgloss.NewStyle().Height(m.height - 1).Render(body)

	content := bodyBlock + "\n" + helpBar
	rendered := lipgloss.NewStyle().PaddingLeft(viewPadX).PaddingRight(viewPadX).Render(content)
	v := tea.NewView(rendered)
	v.AltScreen = true
	v.BackgroundColor = bgMain
	return v
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
