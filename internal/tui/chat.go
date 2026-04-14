package tui

import (
	"regexp"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

type chatRole int

const (
	roleUser chatRole = iota
	roleAssistant
	roleSystemNote
)

type chatStatus int

const (
	chatIdle chatStatus = iota
	chatSending
	chatError
)

type chatMessage struct {
	role     chatRole
	headline string
	content  string
}

// SendHandler is invoked by the chat whenever the user submits a turn. It
// returns a tea.Cmd that must eventually produce either a chatReplyMsg or a
// chatErrMsg carrying the provided id, so the chat can match the response back
// to its pending turn.
type SendHandler func(history []chatMessage, userInput string, id int) tea.Cmd

type chatReplyMsg struct {
	id       int
	content  string
	headline string
}

type chatErrMsg struct {
	id  int
	err error
}

const (
	chatInputHeight   = 3
	chatScrollReserve = 3 // topInd + blank below top + botInd
)

var (
	chatUserLabel      = lipgloss.NewStyle().Foreground(accent).Bold(true)
	chatAssistantLabel = lipgloss.NewStyle().Foreground(accentSoft).Bold(true)
	chatBodyStyle      = lipgloss.NewStyle().Foreground(fgBody)
	chatAssistBodyStyle = lipgloss.NewStyle().Foreground(fgBright)
	chatSystemStyle    = lipgloss.NewStyle().Foreground(fgBody).Italic(true)
	chatErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	chatPendingStyle   = lipgloss.NewStyle().Foreground(accent).Italic(true)

	chatHeadlineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(accent).
				Bold(true).
				Padding(0, 1)

	chatUserBlockStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder(), false, false, false, true).
				BorderForeground(accent).
				PaddingLeft(1)

	chatAssistantBlockStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder(), false, false, false, true).
				BorderForeground(fgDim).
				PaddingLeft(1)

	chatInputFrame = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderForeground(accent).
			PaddingLeft(1)
)

type chatModel struct {
	width       int
	height      int
	input       textarea.Model
	viewport    viewport.Model
	spinner     spinner.Model
	messages    []chatMessage
	status      chatStatus
	errMsg      string
	pendingID   int
	nextID      int
	send        SendHandler
	placeholder string
}

type chatOption func(*chatModel)

func withChatPlaceholder(s string) chatOption {
	return func(m *chatModel) { m.placeholder = s }
}

func withChatSendHandler(h SendHandler) chatOption {
	return func(m *chatModel) { m.send = h }
}

func withChatWelcome(text string) chatOption {
	return func(m *chatModel) {
		m.messages = append(m.messages, chatMessage{role: roleSystemNote, content: text})
	}
}

func newChatModel(width, height int, opts ...chatOption) chatModel {
	ta := textarea.New()
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
	ta.CharLimit = 2000
	ta.Focus()

	vp := viewport.New()
	vp.SoftWrap = true

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(accent)

	m := chatModel{
		width:    width,
		height:   height,
		input:    ta,
		viewport: vp,
		spinner:  sp,
	}
	for _, opt := range opts {
		opt(&m)
	}
	if m.placeholder != "" {
		m.input.Placeholder = m.placeholder
	}
	m.SetSize(width, height)
	return m
}

func (m *chatModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	inputInner := max(width-2, 10)
	m.input.SetWidth(inputInner)
	m.input.SetHeight(chatInputHeight)
	m.viewport.SetWidth(width)
	m.viewport.SetHeight(max(height-chatInputHeight-chatScrollReserve, 3))
	m.refreshContent()
}

func (m *chatModel) refreshContent() {
	wasAtBottom := m.viewport.AtBottom()
	content := m.renderHistory(m.viewport.Width())
	h := lipgloss.Height(content)
	vh := m.viewport.Height()
	if h < vh {
		content = strings.Repeat("\n", vh-h) + content
	}
	m.viewport.SetContent(content)
	if wasAtBottom {
		m.viewport.GotoBottom()
	}
}

func (m chatModel) renderHistory(width int) string {
	if len(m.messages) == 0 {
		return ""
	}
	blocks := make([]string, 0, len(m.messages))
	for _, msg := range m.messages {
		blocks = append(blocks, m.renderMessage(msg, width))
	}
	return strings.Join(blocks, "\n\n")
}

func (m chatModel) renderMessage(msg chatMessage, width int) string {
	innerWidth := max(width-4, 10)
	switch msg.role {
	case roleSystemNote:
		return chatSystemStyle.Width(max(width, 10)).Render("— " + msg.content)
	case roleUser:
		label := chatUserLabel.Render("you")
		body := chatBodyStyle.Width(innerWidth).Render(msg.content)
		block := lipgloss.JoinVertical(lipgloss.Left, label, body)
		return chatUserBlockStyle.Render(block)
	case roleAssistant:
		label := chatAssistantLabel.Render("clank")
		if msg.headline != "" {
			pill := chatHeadlineStyle.Render(msg.headline)
			label = lipgloss.JoinHorizontal(lipgloss.Top, label, "  ", pill)
		}
		body := renderMarkdown(msg.content, innerWidth)
		block := lipgloss.JoinVertical(lipgloss.Left, label, body)
		return chatAssistantBlockStyle.Render(block)
	}
	return ""
}

// renderMarkdown renders assistant message content, applying chroma syntax
// highlighting to fenced code blocks and plain styling to surrounding text.
func renderMarkdown(content string, width int) string {
	const fence = "```"
	var out strings.Builder
	s := content
	for {
		idx := strings.Index(s, fence)
		if idx == -1 {
			if t := strings.TrimRight(s, "\n"); t != "" {
				out.WriteString(chatAssistBodyStyle.Width(width).Render(renderInlineCode(t)))
			}
			break
		}
		// text before the fence
		if before := strings.TrimRight(s[:idx], "\n"); before != "" {
			out.WriteString(chatAssistBodyStyle.Width(width).Render(renderInlineCode(before)))
			out.WriteString("\n")
		}
		rest := s[idx+3:]
		end := strings.Index(rest, fence)
		if end == -1 {
			// unclosed fence — treat remainder as plain text
			out.WriteString(chatAssistBodyStyle.Width(width).Render(strings.TrimRight(s[idx:], "\n")))
			break
		}
		block := rest[:end]
		lang, code := "", block
		if nl := strings.Index(block, "\n"); nl >= 0 {
			lang = strings.TrimSpace(block[:nl])
			code = block[nl+1:]
		}
		out.WriteString(highlightCode(lang, strings.TrimRight(code, "\n"), width))
		out.WriteString("\n")
		s = rest[end+3:]
	}
	return strings.TrimRight(out.String(), "\n")
}

var (
	inlineCodeRe    = regexp.MustCompile("`([^`\n]+)`")
	inlineCodeStyle = lipgloss.NewStyle().
			Foreground(fgBright).
			Background(lipgloss.Color("#1e1e1e"))
)

func renderInlineCode(s string) string {
	return inlineCodeRe.ReplaceAllStringFunc(s, func(match string) string {
		return inlineCodeStyle.Render(match[1 : len(match)-1])
	})
}

// codeBgAnsi is the ANSI true-colour escape for the code block background:
// #1e1e1e (rgb 30,30,30), very slightly lighter than bgMain #171717 (rgb 23,23,23).
const codeBgAnsi = "\x1b[48;2;30;30;30m"

// chromaBgRe matches any ANSI background-colour escape so we can strip chroma's
// own background (dracula uses #282a36) before injecting ours.
var chromaBgRe = regexp.MustCompile(`\x1b\[48[;0-9]*m`)

func highlightCode(lang, code string, width int) string {
	lex := lexers.Get(lang)
	if lex == nil {
		lex = lexers.Analyse(code)
	}
	if lex == nil {
		lex = lexers.Fallback
	}
	lex = chroma.Coalesce(lex)

	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal16m")
	iter, err := lex.Tokenise(nil, code)
	if err != nil {
		return " " + code
	}
	var buf strings.Builder
	if err = formatter.Format(&buf, style, iter); err != nil {
		return " " + code
	}

	// Strip chroma's own background (dracula uses #282a36) so only ours shows.
	raw := chromaBgRe.ReplaceAllString(strings.TrimRight(buf.String(), "\n"), "")

	// Process line-by-line: prefix each line with our background + 3-space margin,
	// and re-inject the background after every reset within the line so token
	// resets don't clobber it mid-line. Pad each line to full width so the
	// background covers the entire row.
	const reset = "\x1b[0m"
	const margin = 3
	blank := codeBgAnsi + strings.Repeat(" ", width) + reset
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		line = strings.ReplaceAll(line, reset, reset+codeBgAnsi)
		// Measure visible width of the ANSI-decorated content after the margin.
		visibleWidth := lipgloss.Width(codeBgAnsi + "   " + line + reset)
		pad := max(width-visibleWidth, 0)
		lines[i] = codeBgAnsi + "   " + line + strings.Repeat(" ", pad) + reset
	}
	return blank + "\n" + strings.Join(lines, "\n") + "\n" + blank
}

// History returns a snapshot of the real chat turns (user + assistant),
// excluding system notes, for the SendHandler to build a request payload.
func (m *chatModel) History() []chatMessage {
	out := make([]chatMessage, 0, len(m.messages))
	for _, msg := range m.messages {
		if msg.role == roleSystemNote {
			continue
		}
		out = append(out, msg)
	}
	return out
}

func (m chatModel) Busy() bool {
	return m.status == chatSending
}

func (m chatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m chatModel) Update(msg tea.Msg) (chatModel, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.status != chatSending {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case chatReplyMsg:
		if msg.id != m.pendingID {
			return m, nil
		}
		m.messages = append(m.messages, chatMessage{
			role:     roleAssistant,
			headline: msg.headline,
			content:  msg.content,
		})
		m.status = chatIdle
		m.errMsg = ""
		m.refreshContent()
		m.viewport.GotoBottom()
		return m, nil

	case chatErrMsg:
		if msg.id != m.pendingID {
			return m, nil
		}
		m.status = chatError
		m.errMsg = msg.err.Error()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if m.status == chatSending {
				return m, nil
			}
			content := strings.TrimSpace(m.input.Value())
			if content == "" {
				return m, nil
			}
			m.messages = append(m.messages, chatMessage{role: roleUser, content: content})
			m.input.Reset()
			m.nextID++
			id := m.nextID
			m.pendingID = id
			m.status = chatSending
			m.errMsg = ""
			m.refreshContent()
			m.viewport.GotoBottom()
			history := m.History()
			cmds := []tea.Cmd{m.spinner.Tick}
			if m.send != nil {
				cmds = append(cmds, m.send(history, content, id))
			}
			return m, tea.Batch(cmds...)
		case "pgup", "pgdown", "ctrl+u", "ctrl+d":
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m chatModel) ViewportYOffset() int {
	return m.viewport.YOffset()
}

func (m chatModel) ScrollLines(n int) chatModel {
	if n > 0 {
		m.viewport.ScrollDown(n)
	} else if n < 0 {
		m.viewport.ScrollUp(-n)
	}
	return m
}

func (m chatModel) StatusView() string {
	switch m.status {
	case chatSending:
		return chatPendingStyle.Render(m.spinner.View())
	case chatError:
		return chatErrorStyle.Render("err: " + m.errMsg)
	}
	return ""
}

func (m chatModel) View() string {
	var topInd, botInd string
	if !m.viewport.AtBottom() {
		topInd = chatScrollIndicator(true, m.viewport.Width())
	}
	if !m.viewport.AtTop() {
		botInd = chatScrollIndicator(false, m.viewport.Width())
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		botInd,
		"",
		m.viewport.View(),
		topInd,
		chatInputFrame.Render(m.input.View()),
	)
}

func chatScrollIndicator(up bool, width int) string {
	arrow := "↑"
	if up {
		arrow = "↓"
	}
	return lipgloss.NewStyle().Foreground(fgDim).Width(width).Align(lipgloss.Right).Render(arrow)
}
