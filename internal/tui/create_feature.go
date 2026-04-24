package tui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"shotgun.dev/odek/internal/decomposer"
	openai "shotgun.dev/odek/openai"
)

type kanjiTickMsg struct{}

func kanjiTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg { return kanjiTickMsg{} })
}

const createFeatureChromeHeight = 6

// decomposeState is a heap-allocated holder for the active decomposition
// session, shared between the chat model and the decomposition page so both
// sides of the split pane see the same Session.
//
// The `decomposing` flag is set synchronously when a decompose operation
// starts — before the cmd is returned to Tea — so the right pane can show
// an immediate loading indicator when the user sends a message, without
// waiting for the LLM call to finish.
//
// The thinking buffer is a shared log of reasoning_content deltas streamed
// from the LLM. Both panes render it as a scrolling marquee in the kanji
// area while work is in flight.
type decomposeState struct {
	mu          sync.Mutex
	session     *decomposer.Session
	decomposing bool
	thinking    []rune
}

const maxThinkingRunes = 12000

func (s *decomposeState) get() *decomposer.Session {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.session
}

// Decomposing reports whether a decompose is in flight. Read by the right
// pane's render path to show a loading indicator.
func (s *decomposeState) Decomposing() bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.decomposing
}

func (s *decomposeState) setDecomposing(v bool) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.decomposing = v
}

// AppendThinking adds a streamed reasoning_content chunk to the shared
// thinking buffer. Safe to call from any goroutine — the LLM client invokes
// this from its streaming read loop.
func (s *decomposeState) AppendThinking(chunk string) {
	if s == nil || chunk == "" {
		return
	}
	s.mu.Lock()
	s.thinking = append(s.thinking, []rune(chunk)...)
	if len(s.thinking) > maxThinkingRunes {
		s.thinking = s.thinking[len(s.thinking)-maxThinkingRunes:]
	}
	s.mu.Unlock()
}

// ThinkingSnapshot returns a copy of the current thinking buffer. Safe to
// call from the render path.
func (s *decomposeState) ThinkingSnapshot() []rune {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.thinking) == 0 {
		return nil
	}
	out := make([]rune, len(s.thinking))
	copy(out, s.thinking)
	return out
}

// ClearThinking drops the accumulated thinking buffer. Called when a
// fresh request starts or after the scroll-out animation has pushed the
// last thinking rune off-screen.
func (s *decomposeState) ClearThinking() {
	if s == nil {
		return
	}
	s.mu.Lock()
	s.thinking = nil
	s.mu.Unlock()
}

// WorkInProgress reports whether any background work is currently active:
// either a decompose call is in flight, or the session is streaming
// expansion events. Both panes drive their kanji animation off this so
// the chat and decompose views animate together.
func (s *decomposeState) WorkInProgress() bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	dec := s.decomposing
	sess := s.session
	s.mu.Unlock()
	if dec {
		return true
	}
	if sess != nil {
		return sess.Snapshot().Expanding
	}
	return false
}

type createFeatureModel struct {
	ctx         context.Context
	client      *openai.Client
	decomposer  *decomposer.Decomposer
	state       *decomposeState
	width       int
	height      int
	chat        chatModel
	kanjiOffset int
	inSplit     bool
}

func newCreateFeatureModel(ctx context.Context, client *openai.Client, dec *decomposer.Decomposer, width, height int) createFeatureModel {
	state := &decomposeState{}
	m := createFeatureModel{
		ctx:        ctx,
		client:     client,
		decomposer: dec,
		state:      state,
		width:      width,
		height:     height,
	}

	sendHandler := makeFeatureSendHandler(ctx, client, dec, state)
	chatWidth := max(width-viewPadX*2, 20)
	chatHeight := max(height-createFeatureChromeHeight, 5)
	m.chat = newChatModel(
		chatWidth,
		chatHeight,
		withChatPlaceholder("Describe the feature..."),
		withChatSendHandler(sendHandler),
		withChatWelcome("Describe your feature. Chat freely — I'll update the rune tree when you change scope."),
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
		if m.state.WorkInProgress() {
			m.kanjiOffset += 2
		}
		if snap := m.state.ThinkingSnapshot(); len(snap) > 0 {
			m.chat.SetPendingThinking(string(snap))
		}
		return m, kanjiTick()

	case decompKanjiTickMsg:
		// Decompose pane's ticker reaches us via the split's broadcast —
		// consume it so it doesn't bleed into the chat component.
		return m, nil

	case decomposeDoneMsg:
		reply := chatReplyMsg{id: msg.id, content: msg.content, headline: msg.headline}
		var cmd tea.Cmd
		m.chat, cmd = m.chat.Update(reply)
		if msg.events != nil {
			return m, tea.Batch(cmd, pumpExpansionCmd(msg.events))
		}
		return m, cmd

	case expansionEventMsg:
		if !msg.ok {
			// Channel closed; pump for this session is done.
			return m, nil
		}
		// Producer has already called sess.Apply; we just re-pump the
		// same channel to read the next event. Using msg.source (not
		// state.session.Events) ensures a rapid re-decompose doesn't
		// hijack the old pump onto the new session's channel.
		return m, pumpExpansionCmd(msg.source)

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
				width:      m.width,
				height:     m.height,
				help:       newHelpModel(),
				ctx:        m.ctx,
				client:     m.client,
				decomposer: m.decomposer,
			}, landingTick()
		case "ctrl+enter", "ctrl+s":
			if m.chat.Busy() {
				return m, nil
			}
			pin := renderFeaturePin()
			sess := m.state.get()
			if m.width >= splitPaneMinWidth {
				leftW, rightW := splitWidths(m.width)
				left := m
				left.resize(leftW, m.height)
				right := newFeatureDecompModel(rightW, m.height, sess, m.state)
				split := newSplitPaneModel(left, right, m.width, m.height)
				return split, split.Init()
			}
			dest := newFeatureDecompModel(m.width, m.height, sess, m.state)
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
	bodyBlock := lipgloss.NewStyle().Height(m.height - 3).Render(body)
	blankRow := padStyle.Render(strings.Repeat(" ", innerWidth))

	content := bodyBlock + "\n" + blankRow + "\n" + helpBar + "\n" + blankRow
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
	sepStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
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

// makeFeatureSendHandler returns a SendHandler that routes every user turn
// through the chat LLM. The LLM replies in plain text for discussion, or
// calls the `decompose` tool when the user's message revises the feature
// spec; tool calls run the decomposer and the resulting summary is shown
// as the assistant reply.
func makeFeatureSendHandler(ctx context.Context, client *openai.Client, dec *decomposer.Decomposer, state *decomposeState) SendHandler {
	return func(history []chatMessage, userInput string, id int) tea.Cmd {
		// Set the decomposing flag synchronously (before returning a Cmd)
		// so the very next render — which happens right after this
		// Update returns — sees decomposing=true. Pair it with an
		// instant decomposeStartedMsg so the right pane's Update is
		// triggered to re-snapshot the state.
		state.setDecomposing(true)
		// Wipe any leftover thinking from the prior turn so the new
		// marquee starts fresh.
		state.ClearThinking()
		startCmd := func() tea.Msg { return decomposeStartedMsg{id: id} }

		if client == nil {
			state.setDecomposing(false)
			return offlineReply(id)
		}
		// Attach a thinking callback so every streamed reasoning_content
		// delta from any LLM call under this ctx appends to the shared
		// buffer. Both panes read from it on every kanji tick.
		thinkingCtx := openai.WithThinkingCallback(ctx, state.AppendThinking)
		return tea.Batch(startCmd, chatHandler(thinkingCtx, client, dec, state, id, history))
	}
}

// buildDiscussion formats a chat transcript suitable for inclusion in a
// refinement-pass user message. Rules:
//   - System notes are skipped (never shown to the model).
//   - The first non-empty user message is skipped — it's the requirement
//     and already appears at the top of the refinement prompt.
//   - Assistant turns with an "Effort:" headline are skipped because
//     their content is a decomposition reply, already carried structurally
//     in SessionContext.Prior.
//   - Everything else is rendered as "you: ..." / "clank: ...".
//
// Returns empty string when there is nothing to say beyond the requirement.
func buildDiscussion(history []chatMessage) string {
	var b strings.Builder
	skippedRequirement := false
	for _, msg := range history {
		switch msg.role {
		case roleSystemNote:
			continue
		case roleUser:
			content := strings.TrimSpace(msg.content)
			if content == "" {
				continue
			}
			if !skippedRequirement {
				skippedRequirement = true
				continue
			}
			fmt.Fprintf(&b, "you: %s\n", content)
		case roleAssistant:
			content := strings.TrimSpace(msg.content)
			if content == "" && msg.headline == "" {
				continue
			}
			if strings.HasPrefix(msg.headline, "Effort:") {
				// Prior decomposition reply — structural, redundant with
				// SessionContext.Prior.
				continue
			}
			if msg.headline != "" {
				fmt.Fprintf(&b, "clank [%s]: %s\n", msg.headline, content)
			} else {
				fmt.Fprintf(&b, "clank: %s\n", content)
			}
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

// extractRequirement returns the first non-empty user message in history.
// That message is treated as the feature requirement; subsequent user
// messages become discussion.
func extractRequirement(history []chatMessage) string {
	for _, msg := range history {
		if msg.role != roleUser {
			continue
		}
		content := strings.TrimSpace(msg.content)
		if content != "" {
			return content
		}
	}
	return ""
}

// expansionEventMsg carries a single ExpansionEvent off the channel into
// the bubbletea Update loop. Carries the source channel so the handler
// can re-pump the exact same channel without re-reading state (which may
// already point at a newer session after a rapid /decompose).
type expansionEventMsg struct {
	event  decomposer.ExpansionEvent
	source <-chan decomposer.ExpansionEvent
	ok     bool
}

// decomposeDoneMsg is emitted by decomposeHandler once NewSession has
// returned. It carries both the chat reply payload and (optionally) the
// live event channel, so the parent Update can forward the reply to the
// chat AND start pumping events off the channel.
type decomposeDoneMsg struct {
	id       int
	content  string
	headline string
	events   <-chan decomposer.ExpansionEvent
}

// decomposeStartedMsg is emitted immediately when the user sends a
// message (before the LLM call completes). It lets the right pane refresh
// so it can show a loading indicator while the decomposer runs, instead
// of the user staring at a stale empty-state pane for several seconds.
type decomposeStartedMsg struct {
	id int
}

// pumpExpansionCmd returns a tea.Cmd that blocks on one event from ch.
// The model re-schedules the pump on every ok=true event, and stops when
// the channel closes (ok=false). Call sites must check ok before reading
// event.
func pumpExpansionCmd(ch <-chan decomposer.ExpansionEvent) tea.Cmd {
	if ch == nil {
		return nil
	}
	return func() tea.Msg {
		evt, ok := <-ch
		return expansionEventMsg{event: evt, source: ch, ok: ok}
	}
}

// chatHandler runs on every user turn. It forwards the chat history to the
// chat LLM with the `decompose` tool available. If the model calls the tool,
// we run the decomposer with the tool's arguments. Otherwise the model's
// text reply is shown as a normal chat reply.
func chatHandler(ctx context.Context, client *openai.Client, dec *decomposer.Decomposer, state *decomposeState, id int, history []chatMessage) tea.Cmd {
	return func() tea.Msg {
		defer state.setDecomposing(false)

		messages := buildChatMessages(history)
		resp, err := client.Chat(ctx, &openai.ChatCompletionRequest{
			Model:      openai.DefaultModel,
			Messages:   messages,
			Tools:      []openai.Tool{chatDecomposeTool},
			ToolChoice: "auto",
		})
		if err != nil {
			return chatErrMsg{id: id, err: err}
		}
		if len(resp.Choices) == 0 {
			return chatErrMsg{id: id, err: fmt.Errorf("no response from chat model")}
		}
		choice := resp.Choices[0].Message

		for _, call := range choice.ToolCalls {
			if call.Function.Name != "decompose" {
				continue
			}
			var args struct {
				Levels int `json:"levels"`
				Effort int `json:"effort"`
			}
			_ = json.Unmarshal([]byte(call.Function.Arguments), &args)
			if args.Levels <= 0 {
				args.Levels = defaultAutoLevels
			}
			if args.Levels > 10 {
				args.Levels = 10
			}
			if args.Effort <= 0 || args.Effort > 5 {
				args.Effort = defaultAutoEffort
			}
			req := extractRequirement(history)
			discussion := buildDiscussion(history)
			return runDecompose(ctx, dec, state, id, req, discussion, args.Levels, args.Effort)
		}

		// No decompose tool call; chat LLM replied with plain text.
		content := strings.TrimSpace(choice.Content)
		if content == "" {
			content = "(no reply)"
		}
		return chatReplyMsg{id: id, content: content}
	}
}

// runDecompose is the shared backend for both auto and manual decompose.
// It pulls the prior session + any live expansion check, builds a
// SessionContext from discussion + prior, runs NewSession, and returns the
// appropriate tea.Msg (decomposeDoneMsg on success, chatReplyMsg with
// "clarification" headline when the model asked for clarification,
// chatErrMsg otherwise).
//
// On every exit path this clears state.decomposing so the right pane's
// loading indicator goes away.
func runDecompose(ctx context.Context, dec *decomposer.Decomposer, state *decomposeState, id int, requirement, discussion string, levels, effortLvl int) tea.Msg {
	defer state.setDecomposing(false)

	if dec == nil {
		return chatErrMsg{id: id, err: fmt.Errorf("decomposer unavailable")}
	}
	if requirement == "" {
		return chatErrMsg{id: id, err: fmt.Errorf("describe the feature first")}
	}

	state.mu.Lock()
	prior := state.session
	state.mu.Unlock()
	if prior != nil && prior.Events != nil {
		return chatErrMsg{id: id, err: fmt.Errorf("a decomposition is already in progress")}
	}

	var priorResp *decomposer.DecompositionResponse
	if prior != nil && prior.Root != nil {
		priorResp = prior.Root.Response
	}

	if levels < 1 {
		levels = 1
	}
	if effortLvl < 1 || effortLvl > 5 {
		effortLvl = defaultAutoEffort
	}
	cfg := decomposer.ConfigForEffort(effortLvl)
	// levels shadows MaxDepth; effort still controls ParallelInitial/RuneCap.
	if levels-1 < cfg.MaxDepth {
		cfg.MaxDepth = levels - 1
	}
	cfg.Recurse = levels > 1

	sessCtx := decomposer.SessionContext{
		Discussion: discussion,
		Prior:      priorResp,
	}
	sess, err := dec.NewSession(ctx, requirement, effortLvl, "", cfg, sessCtx)
	if err != nil {
		var clar *decomposer.ClarificationNeeded
		if errors.As(err, &clar) {
			return chatReplyMsg{
				id:       id,
				content:  clar.Message,
				headline: "clarification",
			}
		}
		return chatErrMsg{id: id, err: err}
	}

	state.mu.Lock()
	state.session = sess
	state.mu.Unlock()

	content := renderDecompositionSummary(sess)
	headline := fmt.Sprintf("Effort: %d/5", effortLvl)

	var events <-chan decomposer.ExpansionEvent
	if levels > 1 {
		events = dec.ExpandStreaming(ctx, sess, cfg)
	}

	return decomposeDoneMsg{id: id, content: content, headline: headline, events: events}
}

// defaultAutoLevels / defaultAutoEffort are the parameters used for auto
// mode and as the fallback when the chat LLM calls the decompose tool
// without specifying them. Kept conservative to keep turnaround quick.
const (
	defaultAutoLevels = 1
	defaultAutoEffort = 2
)

// chatDecomposeTool is the tool definition given to the chat LLM. The LLM
// calls this when the user's latest message revises the feature spec —
// i.e. a scope or requirement change, not a question or discussion.
var chatDecomposeTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "decompose",
		Description: "Apply a scope or requirement change to the feature spec. Call this whenever the user's latest message revises what the library should do — adds, removes, or renames capabilities, tightens or relaxes scope, or otherwise changes the requirements. Do NOT call it for questions, clarifications, or design discussion that doesn't change what the library does. The tool runs the decomposition pipeline using the full conversation history as context.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"levels": map[string]any{
					"type":        "integer",
					"minimum":     1,
					"maximum":     10,
					"description": "Recursion depth. 1 = top-level only (default). Use higher when the user asks for a deeper tree.",
				},
				"effort": map[string]any{
					"type":        "integer",
					"minimum":     1,
					"maximum":     5,
					"description": "Effort level 1-5, controls parallelism and rune cap. Default 2. Raise when the user asks for high quality or mentions parallel attempts.",
				},
			},
		},
	},
}

// buildChatMessages builds the conversation sent to the chat LLM. Includes
// a system prompt that explains the LLM's role and the decomposition-vs-
// discussion classification it's responsible for.
func buildChatMessages(history []chatMessage) []openai.ChatMessage {
	system := `You are Odek, a software library design collaborator. The user is iterating on a library spec, and on every turn you either discuss it or update it.

Classify the user's latest message by intent:

1. **Spec change** — the user is revising what the library should do: adding/removing/renaming capabilities, tightening or relaxing scope, changing the library's purpose. Call the ` + "`decompose`" + ` tool. Do not reply with text before calling it; the tool's output becomes the reply.

2. **Discussion** — the user is asking a question, exploring tradeoffs, clarifying how something works, or giving feedback that doesn't change what the library does. Reply in plain text. Be concise and practical. Do not paste rune trees yourself.

Examples:
- "make it a scientific calculator" → spec change, call the tool.
- "also support matrices" → spec change, call the tool.
- "drop the divide function" → spec change, call the tool.
- "how does logarithm work?" → discussion, reply in text.
- "what should divide-by-zero return?" → discussion, reply in text (until the user picks a behavior).
- "what's the tradeoff between X and Y?" → discussion, reply in text.

When in doubt, prefer discussion: it's cheap to follow up with a spec change, but a spurious rewrite costs the user their prior structure.`

	msgs := []openai.ChatMessage{{Role: openai.RoleSystem, Content: system}}
	for _, msg := range history {
		if msg.role == roleSystemNote {
			continue
		}
		role := openai.RoleUser
		content := msg.content
		if msg.role == roleAssistant {
			role = openai.RoleAssistant
			if msg.headline != "" {
				content = "[" + msg.headline + "] " + msg.content
			}
		}
		msgs = append(msgs, openai.ChatMessage{Role: role, Content: content})
	}
	return msgs
}

func offlineReply(id int) tea.Cmd {
	return func() tea.Msg {
		return chatErrMsg{id: id, err: fmt.Errorf("chat offline: no client configured")}
	}
}
