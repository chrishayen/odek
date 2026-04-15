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
// session plus the auto-decompose toggle. Shared between the chat model
// and the decomposition page so both sides of the split pane see the same
// Session, and so the send handler (a closure captured at construction
// time) can read the toggle on every turn.
//
// The `decomposing` flag is set synchronously when a /decompose-equivalent
// operation starts — before the cmd is returned to Tea — so the right pane
// can show an immediate loading indicator when the user sends a message,
// without waiting for the LLM call to finish.
type decomposeState struct {
	mu          sync.Mutex
	session     *decomposer.Session
	auto        bool
	decomposing bool
}

func (s *decomposeState) get() *decomposer.Session {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.session
}

// AutoDecompose reports whether auto-decompose is currently on.
func (s *decomposeState) AutoDecompose() bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.auto
}

// toggleAuto flips the auto-decompose flag and returns the new value.
func (s *decomposeState) toggleAuto() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.auto = !s.auto
	return s.auto
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
	state := &decomposeState{auto: true}
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
		withChatWelcome("Describe your feature. Auto-decompose is on — each message updates the rune tree. Ctrl+T to toggle."),
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
		case "ctrl+t":
			if m.state != nil {
				m.state.toggleAuto()
			}
			return m, nil
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

	helpBar := renderFormHelpBar(innerWidth, m.chat.CanScrollDown(), m.inSplit, m.state.AutoDecompose())

	scrollOff := m.chat.ViewportYOffset() * 2
	kanjiLine1 := renderKanjiLine(innerWidth, 2, m.kanjiOffset+scrollOff)
	kanjiLine2 := renderKanjiLine(innerWidth, 3, -(m.kanjiOffset + scrollOff))
	body := header + "\n\n" + kanjiLine1 + "\n" + kanjiLine2 + "\n" + m.chat.View()
	bodyBlock := lipgloss.NewStyle().Height(m.height - 2).Render(body)
	blankRow := padStyle.Render(strings.Repeat(" ", innerWidth))

	content := bodyBlock + "\n" + helpBar + "\n" + blankRow
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

func renderFormHelpBar(width int, scrollDown, showTabSwitch, autoOn bool) string {
	keyStyle := lipgloss.NewStyle().Foreground(fgBright).Background(bgMain).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(fgBody).Background(bgMain)
	sepStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
	padStyle := lipgloss.NewStyle().Background(bgMain)

	type binding struct{ key, desc string }
	autoLabel := "live: off"
	if autoOn {
		autoLabel = "live: on"
	}
	bindings := []binding{
		{"enter", "send"},
		{"alt+enter", "new line"},
		{"↑/↓", "scroll"},
		{"ctrl+t", autoLabel},
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

// makeFeatureSendHandler returns a SendHandler whose routing depends on
// the auto-decompose toggle:
//   - auto on  -> every turn runs the decomposer with the full chat
//     history as refinement context; no chat LLM is involved.
//   - auto off -> each turn goes through the chat LLM, which can respond
//     in plain text or call the `decompose` tool. Tool calls run the
//     decomposer; plain replies are shown as normal chat messages.
func makeFeatureSendHandler(ctx context.Context, client *openai.Client, dec *decomposer.Decomposer, state *decomposeState) SendHandler {
	return func(history []chatMessage, userInput string, id int) tea.Cmd {
		// Set the decomposing flag synchronously (before returning a Cmd)
		// so the very next render — which happens right after this
		// Update returns — sees decomposing=true. Pair it with an
		// instant decomposeStartedMsg so the right pane's Update is
		// triggered to re-snapshot the state.
		state.setDecomposing(true)
		startCmd := func() tea.Msg { return decomposeStartedMsg{id: id} }

		if state.AutoDecompose() {
			return tea.Batch(startCmd, autoDecomposeHandler(ctx, dec, state, id, history))
		}
		if client == nil {
			state.setDecomposing(false)
			return offlineReply(id)
		}
		return tea.Batch(startCmd, manualChatHandler(ctx, client, dec, state, id, history))
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

// autoDecomposeHandler runs on every user turn when auto-decompose is on.
// It bypasses the chat LLM entirely: the user's message is added to the
// refinement discussion and the decomposer produces a fresh session, which
// replaces the prior one.
func autoDecomposeHandler(ctx context.Context, dec *decomposer.Decomposer, state *decomposeState, id int, history []chatMessage) tea.Cmd {
	return func() tea.Msg {
		req := extractRequirement(history)
		discussion := buildDiscussion(history)
		return runDecompose(ctx, dec, state, id, req, discussion, defaultAutoLevels, defaultAutoEffort)
	}
}

// manualChatHandler runs on every user turn when auto-decompose is off.
// It forwards the chat history to the chat LLM with the `decompose` tool
// available. If the model calls the tool, we run the decomposer with the
// tool's arguments. Otherwise the model's text reply is shown as a normal
// chat reply.
func manualChatHandler(ctx context.Context, client *openai.Client, dec *decomposer.Decomposer, state *decomposeState, id int, history []chatMessage) tea.Cmd {
	return func() tea.Msg {
		messages := buildManualChatMessages(history)
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
		// Clear the decomposing flag here (runDecompose would have done it
		// on the tool-call path).
		state.setDecomposing(false)
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

	content := renderDecompositionSummary(sess, priorResp)
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

// chatDecomposeTool is the tool definition given to the chat LLM in manual
// mode. The LLM calls this when the user's request implies they want a
// decomposition.
var chatDecomposeTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "decompose",
		Description: "Run the real decomposition pipeline on the current feature. Call this when the user indicates they want a decomposition — e.g. 'decompose this', 'break it down', 'show me the structure', 'let's build it', 'generate the runes'. The tool uses the full conversation history as context. Do not call it for general design discussion.",
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

// buildManualChatMessages builds the conversation sent to the chat LLM in
// manual mode. Includes a system prompt that explains the LLM's role and
// how to trigger decomposition via the tool.
func buildManualChatMessages(history []chatMessage) []openai.ChatMessage {
	system := `You are Odek, a software feature design collaborator. The user is iterating on a feature, and you help them refine scope, names, tradeoffs, and edge cases in plain text.

When the user indicates they want a decomposition — e.g. "decompose this", "break it down", "show me the structure", "let's build it", "generate the runes" — call the ` + "`decompose`" + ` tool. You do not need to reply with text before calling it. The tool runs the real decomposition pipeline on the current feature using the full conversation history as context.

Otherwise, respond in plain text. Be concise and practical. Do not paste rune trees yourself; let the tool produce them.`

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
