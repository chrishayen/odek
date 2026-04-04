package tui

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	runepkg "github.com/chrishayen/odek/internal/rune"
)

type formState int

const (
	stateIdle formState = iota
	stateDecomposing
	stateDone
	stateError
	stateAuthError
	stateRefining
	stateApproved
	stateHydrating
)

var (
	inputLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Bold(true)

	statusOk = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66CC66"))

	statusErr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CC6666"))

	featureNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F5A623")).
				Bold(true)

	featureSummaryStyle = lipgloss.NewStyle().
				Foreground(dim).
				Italic(true)

	pkgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6A9FD9")).
			Bold(true)

	runeNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5A623")).
			Bold(true)

	runeLeafStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC"))

	runeSigStyle = lipgloss.NewStyle().
			Foreground(dim)

	testPassStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66CC66"))

	testFailStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CC6666"))

	assumptionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D4A843"))

	paneHeaderActive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F5A623")).
				Bold(true)

	paneHeaderInactive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555"))
)

// --- draft suffix helpers ---

// shortID returns 6 random hex chars for draft uniqueness.
func shortID() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// addDraftSuffix appends "_suffix" to the top-level segment of a dot-path.
// e.g. addDraftSuffix("write_bing_bong.write", "a8f3b2") → "write_bing_bong_a8f3b2.write"
func addDraftSuffix(name, suffix string) string {
	parts := strings.SplitN(name, ".", 2)
	parts[0] = parts[0] + "_" + suffix
	return strings.Join(parts, ".")
}

// removeDraftSuffix strips "_suffix" from the top-level segment of a dot-path.
func removeDraftSuffix(name, suffix string) string {
	parts := strings.SplitN(name, ".", 2)
	parts[0] = strings.TrimSuffix(parts[0], "_"+suffix)
	return strings.Join(parts, ".")
}

// extractDraftSuffix returns the 6-char hex suffix from a name like "foo_a8f3b2", or "".
func extractDraftSuffix(name string) string {
	top := strings.SplitN(name, ".", 2)[0]
	if len(top) < 8 { // minimum: "x_a8f3b2"
		return ""
	}
	if top[len(top)-7] != '_' {
		return ""
	}
	suffix := top[len(top)-6:]
	for _, c := range suffix {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return ""
		}
	}
	return suffix
}

// hasDraftSuffix checks if a rune name's top-level segment ends with "_" + suffix.
func hasDraftSuffix(name, suffix string) bool {
	top := strings.SplitN(name, ".", 2)[0]
	return strings.HasSuffix(top, "_"+suffix)
}

// runeListItem is an item in the decomposition rune list.
type runeListItem struct {
	runeIdx      int    // index into NewRunes; -1 for non-rune items
	name         string // display name
	isHeader     bool   // top-level package header
	isExisting   bool   // existing rune section
	isSpacer     bool   // empty visual separator between groups
	count        int    // child count for package headers
	covers       string // what existing rune covers
	hasComment   bool   // has a review comment
}

func (i runeListItem) Title() string       { return i.name }
func (i runeListItem) Description() string { return "" }
func (i runeListItem) FilterValue() string { return i.name }

// runeListDelegate renders items in the rune list with left-border selection.
type runeListDelegate struct{}

func (d runeListDelegate) Height() int                             { return 1 }
func (d runeListDelegate) Spacing() int                            { return 0 }
func (d runeListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d runeListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	ri, ok := item.(runeListItem)
	if !ok {
		return
	}
	if ri.isSpacer {
		fmt.Fprint(w, strings.Repeat(" ", m.Width()))
		return
	}

	selected := index == m.Index()
	availWidth := m.Width()

	selectedBorder := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#F5A623")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Width(availWidth)

	var str string
	switch {
	case ri.isHeader && ri.isExisting:
		str = lipgloss.NewStyle().
			Foreground(dim).
			Italic(true).
			Width(availWidth).
			Render(" " + ri.name)
	case ri.isExisting:
		name := ri.name
		if len(name) > availWidth-6 {
			name = name[:availWidth-7] + "~"
		}
		if selected {
			str = selectedBorder.
				Foreground(lipgloss.Color("#777777")).
				Padding(0, 0, 0, 3).
				Render(name)
		} else {
			str = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555")).
				Width(availWidth).
				Render("    " + name)
		}
	case ri.isHeader:
		countStr := fmt.Sprintf(" (%d)", ri.count)
		commentMarker := ""
		if ri.hasComment {
			commentMarker = " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623")).Render("●")
		}
		if selected {
			str = selectedBorder.
				Bold(true).
				Render(ri.name + countStr + commentMarker)
		} else {
			str = lipgloss.NewStyle().
				Width(availWidth).
				Render(" " + pkgStyle.Render(ri.name) + runeSigStyle.Render(countStr) + commentMarker)
		}
	default:
		name := ri.name
		commentMarker := ""
		if ri.hasComment {
			commentMarker = " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623")).Render("●")
		}
		if len(name) > availWidth-6 {
			name = name[:availWidth-7] + "~"
		}
		if selected {
			str = selectedBorder.
				Padding(0, 0, 0, 3).
				Render(name + commentMarker)
		} else {
			str = runeLeafStyle.Width(availWidth).Render("    " + name + commentMarker)
		}
	}

	fmt.Fprint(w, str)
}

func renderPaneHeader(label string, width int, active bool) string {
	style := paneHeaderInactive
	if active {
		style = paneHeaderActive
	}
	prefix := "── "
	suffix := " "
	inner := prefix + label + suffix
	remaining := width - lipgloss.Width(inner)
	if remaining < 0 {
		remaining = 0
	}
	return style.Render(inner + strings.Repeat("─", remaining))
}

// API response types

type decomposeResponse struct {
	JobID string `json:"job_id"`
}

type jobResponse struct {
	ID       string          `json:"id"`
	Status   string          `json:"status"`
	Progress string          `json:"progress,omitempty"`
	Result   json.RawMessage `json:"result,omitempty"`
	Error    string          `json:"error,omitempty"`
}

type proposedRune struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Signature     string   `json:"signature"`
	PositiveTests []string `json:"positive_tests"`
	NegativeTests []string `json:"negative_tests"`
	Assumptions   []string `json:"assumptions,omitempty"`
	Refs          []string `json:"refs"`
}

type existingMatch struct {
	Name   string `json:"name"`
	Covers string `json:"covers"`
}

type decomposeResult struct {
	FeatureName      string            `json:"feature_name,omitempty"`
	Summary          string            `json:"summary,omitempty"`
	FlowDiagram      string            `json:"flow_diagram,omitempty"`
	NewRunes         []proposedRune    `json:"new_runes"`
	ExistingRunes    []existingMatch   `json:"existing_runes"`
	TreeOutput       string            `json:"tree_output"`
	PackageSummaries map[string]string `json:"package_summaries,omitempty"`
}

// Messages

type decomposeStartedMsg struct {
	jobID string
}

type decomposeErrorMsg struct {
	err error
}

type pollTickMsg struct{}

type progressMsg struct {
	text string
}

type decomposeDoneMsg struct {
	result decomposeResult
}

type loginDoneMsg struct {
	err error
}

type goBackMsg struct{}

type commitDoneMsg struct{}

type commitErrorMsg struct {
	err error
}

type hydrateDoneMsg struct{}
type hydrateErrorMsg struct{ err error }

type featureLoadedMsg struct {
	feature string
	runes   []proposedRune
	summary string
}

type featureLoadErrorMsg struct{ err error }

type inputMode int

const (
	inputRefineFeature inputMode = iota
	inputRefineRune
)

type createFeatureModel struct {
	descInput    textarea.Model
	refineInput  textinput.Model
	state        formState
	port         int
	width        int
	jobID        string
	spinner      spinner.Model
	result       *decomposeResult
	errMsg       string
	authURL      string
	height       int
	requirement  string
	inputMode    inputMode
	progressText string
	runeList     list.Model
	midVP        viewport.Model

	// Review comments (ephemeral, not persisted to draft)
	runeComments   map[int]string // rune index → comment
	featureComment string

	// Drafts
	draftName   string // top-level suffixed rune name for this draft (empty = new)
	draftSuffix string // 6-char hex suffix shared by all runes in this draft
	runeStore   *runepkg.Store
}

func newCreateFeatureModel(port, width, height int, runeStore *runepkg.Store) createFeatureModel {
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

	rl := list.New(nil, runeListDelegate{}, width/3, height-6)
	rl.SetShowTitle(false)
	rl.SetShowStatusBar(false)
	rl.SetFilteringEnabled(false)
	rl.SetShowHelp(false)
	rl.SetShowPagination(false)
	rl.KeyMap.Quit.Unbind()
	rl.KeyMap.ShowFullHelp.Unbind()
	rl.KeyMap.CloseFullHelp.Unbind()
	rl.Styles.TitleBar = lipgloss.NewStyle()
	rl.Styles.NoItems = lipgloss.NewStyle()

	midVP := viewport.New(width-width/3-3, height-6)
	midVP.KeyMap = viewport.KeyMap{}

	return createFeatureModel{
		descInput:    ta,
		state:        stateIdle,
		port:         port,
		width:        width,
		height:       height,
		spinner:      s,
		runeList:     rl,
		midVP:        midVP,
		runeStore:    runeStore,
		runeComments: map[int]string{},
	}
}

func newCreateFeatureModelFromDraft(port, width, height int, runeStore *runepkg.Store, r runepkg.Rune) createFeatureModel {
	m := newCreateFeatureModel(port, width, height, runeStore)
	m.draftName = r.Name
	m.draftSuffix = extractDraftSuffix(r.Name)
	m.requirement = r.Description

	// Load all runes belonging to this draft (same suffix)
	allDrafts, err := runeStore.ListByStatus("draft")
	if err == nil {
		var draftRunes []runepkg.Rune
		for _, dr := range allDrafts {
			if m.draftSuffix != "" && hasDraftSuffix(dr.Name, m.draftSuffix) {
				// Strip suffix for display
				stripped := dr
				stripped.Name = removeDraftSuffix(dr.Name, m.draftSuffix)
				draftRunes = append(draftRunes, stripped)
			}
		}
		if len(draftRunes) > 0 {
			featureName := removeDraftSuffix(r.Name, m.draftSuffix)
			m.result = runesToResult(featureName, r.Description, draftRunes)
			sortRunesStdLast(m.result.NewRunes)
			m.buildRuneListItems()
			m.state = stateDone
			m.descInput.Blur()
		} else if m.requirement != "" {
			m.descInput.SetValue(m.requirement)
		}
	} else if m.requirement != "" {
		m.descInput.SetValue(m.requirement)
	}
	return m
}

// runesToResult converts loaded runes into a decomposeResult for display.
func runesToResult(featureName, summary string, runes []runepkg.Rune) *decomposeResult {
	proposed := make([]proposedRune, len(runes))
	for i, r := range runes {
		proposed[i] = proposedRune{
			Name:          r.Name,
			Description:   r.Description,
			Signature:     r.Signature,
			PositiveTests: r.PositiveTests,
			NegativeTests: r.NegativeTests,
			Assumptions:   r.Assumptions,
			Refs:          r.Dependencies,
		}
	}
	return &decomposeResult{
		FeatureName: featureName,
		Summary:     summary,
		NewRunes:    proposed,
	}
}

// saveDraft persists the current state as draft runes with suffixed names.
// Returns the draft name (top-level suffixed rune name).
func (m *createFeatureModel) saveDraft() string {
	if m.runeStore == nil || (m.requirement == "" && m.result == nil) {
		return m.draftName
	}

	if m.result == nil || len(m.result.NewRunes) == 0 {
		return m.draftName
	}

	name := m.result.FeatureName
	if name == "" {
		return m.draftName
	}

	// Generate suffix once per draft
	if m.draftSuffix == "" {
		m.draftSuffix = shortID()
	}

	// Clean up old draft runes (same suffix)
	m.deleteDraftRunes()

	suffixedName := addDraftSuffix(name, m.draftSuffix)

	for _, pr := range m.result.NewRunes {
		r := runepkg.Rune{
			Name:          addDraftSuffix(pr.Name, m.draftSuffix),
			Description:   pr.Description,
			Signature:     pr.Signature,
			PositiveTests: pr.PositiveTests,
			NegativeTests: pr.NegativeTests,
			Assumptions:   pr.Assumptions,
			Dependencies:  pr.Refs,
			Status:        "draft",
		}
		_ = m.runeStore.Create(r)
	}

	m.draftName = suffixedName
	return m.draftName
}

// deleteDraftRunes removes all runes belonging to the current draft (by suffix).
func (m *createFeatureModel) deleteDraftRunes() {
	if m.draftSuffix == "" || m.runeStore == nil {
		return
	}
	allDrafts, err := m.runeStore.ListByStatus("draft")
	if err != nil {
		return
	}
	for _, r := range allDrafts {
		if hasDraftSuffix(r.Name, m.draftSuffix) {
			_ = m.runeStore.Delete(r.Name)
		}
	}
}

func (m *createFeatureModel) resize(width, height int) {
	m.width = width
	m.height = height
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

	leftWidth := width / 3
	availHeight := height - 6
	if availHeight < 5 {
		availHeight = 5
	}
	m.runeList.SetSize(leftWidth, availHeight)
	m.midVP.Width = width - leftWidth - 3
	m.midVP.Height = availHeight
}

func (m *createFeatureModel) submit() tea.Cmd {
	desc := strings.TrimSpace(m.descInput.Value())
	if desc == "" {
		m.errMsg = "description is required"
		m.state = stateError
		return nil
	}

	m.requirement = desc
	m.descInput.Blur()
	return m.decompose(desc)
}

func (m *createFeatureModel) selectedRuneName() string {
	if item, ok := m.runeList.SelectedItem().(runeListItem); ok && item.runeIdx >= 0 {
		return m.result.NewRunes[item.runeIdx].Name
	}
	return "rune"
}

func (m *createFeatureModel) hasComments() bool {
	return len(m.runeComments) > 0 || m.featureComment != ""
}

func (m *createFeatureModel) commentCount() int {
	n := len(m.runeComments)
	if m.featureComment != "" {
		n++
	}
	return n
}

func (m *createFeatureModel) openRefine(mode inputMode, placeholder string) tea.Cmd {
	m.state = stateRefining
	m.inputMode = mode
	m.refineInput = textinput.New()
	m.refineInput.Placeholder = placeholder
	m.refineInput.Width = m.width - 4

	// Pre-fill with existing comment if editing
	if mode == inputRefineRune {
		if item, ok := m.runeList.SelectedItem().(runeListItem); ok && item.runeIdx >= 0 {
			if existing, ok := m.runeComments[item.runeIdx]; ok {
				m.refineInput.SetValue(existing)
			}
		}
	} else if mode == inputRefineFeature {
		if m.featureComment != "" {
			m.refineInput.SetValue(m.featureComment)
		}
	}

	m.refineInput.Focus()
	return m.refineInput.Cursor.BlinkCmd()
}

func (m *createFeatureModel) decompose(req string) tea.Cmd {
	m.state = stateDecomposing
	m.progressText = ""
	port := m.port
	prevDecomp := ""
	prevName := ""
	prevSummary := ""
	if m.result != nil {
		prevDecomp = m.result.TreeOutput
		prevName = m.result.FeatureName
		prevSummary = m.result.Summary
	}
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{"requirement": req, "decomposition": prevDecomp, "feature_name": prevName, "summary": prevSummary})
		resp, err := http.Post(
			fmt.Sprintf("http://localhost:%d/api/decompose", port),
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
	isHydrating := m.state == stateHydrating

	// Pick the right status endpoint based on job type
	endpoint := "decompose"
	if isHydrating {
		endpoint = "hydrate"
	}

	return func() tea.Msg {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/%s/%s", port, endpoint, jobID))
		if err != nil {
			if isHydrating {
				return hydrateErrorMsg{err: err}
			}
			return decomposeErrorMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		var job jobResponse
		if err := json.Unmarshal(data, &job); err != nil {
			if isHydrating {
				return hydrateErrorMsg{err: err}
			}
			return decomposeErrorMsg{err: err}
		}

		switch job.Status {
		case "completed":
			if isHydrating {
				return hydrateDoneMsg{}
			}
			var result decomposeResult
			json.Unmarshal(job.Result, &result)
			return decomposeDoneMsg{result: result}
		case "failed":
			if isHydrating {
				return hydrateErrorMsg{err: fmt.Errorf("%s", job.Error)}
			}
			return decomposeErrorMsg{err: fmt.Errorf("%s", job.Error)}
		default:
			if job.Progress != "" {
				return progressMsg{text: job.Progress}
			}
			return pollTickMsg{}
		}
	}
}

func (m *createFeatureModel) update(msg tea.Msg) tea.Cmd {
	// Handle non-key messages first.
	switch msg := msg.(type) {
	case decomposeStartedMsg:
		m.jobID = msg.jobID
		return tea.Batch(m.spinner.Tick, m.pollJob())
	case progressMsg:
		m.progressText = msg.text
		return m.pollJob()
	case pollTickMsg:
		if m.state == stateDecomposing || m.state == stateHydrating {
			return m.checkJob()
		}
		return nil
	case decomposeDoneMsg:
		m.state = stateDone
		m.result = &msg.result
		m.runeComments = map[int]string{}
		m.featureComment = ""
		sortRunesStdLast(m.result.NewRunes)
		m.buildRuneListItems()
		m.runeList.Select(0)
		m.saveDraft()
		return nil
	case commitDoneMsg:
		m.state = stateApproved
		m.draftName = ""
		return nil
	case commitErrorMsg:
		m.state = stateDone
		m.errMsg = msg.err.Error()
		return nil
	case hydrateDoneMsg:
		// Return to splash after hydration
		return func() tea.Msg { return goBackMsg{} }
	case hydrateErrorMsg:
		m.state = stateError
		m.errMsg = msg.err.Error()
		return nil
	case featureLoadedMsg:
		m.result = &decomposeResult{
			FeatureName: msg.feature,
			Summary:     msg.summary,
			NewRunes:    msg.runes,
		}
		sortRunesStdLast(m.result.NewRunes)
		m.buildRuneListItems()
		m.state = stateApproved
		return nil
	case featureLoadErrorMsg:
		m.state = stateError
		m.errMsg = msg.err.Error()
		return nil
	case decomposeErrorMsg:
		return m.handleDecomposeError(msg)
	case loginDoneMsg:
		return m.handleLoginDone(msg)
	case spinner.TickMsg:
		if m.state == stateDecomposing || m.state == stateHydrating {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return cmd
		}
		return nil
	}

	// Dispatch key messages by state.
	switch m.state {
	case stateIdle, stateError:
		return m.updateFormState(msg)
	case stateRefining:
		return m.updateRefiningState(msg)
	case stateDone:
		return m.updateDoneState(msg)
	case stateApproved:
		return m.updateApprovedState(msg)
	case stateAuthError:
		return m.updateAuthErrorState(msg)
	}
	return nil
}

func (m *createFeatureModel) updateFormState(msg tea.Msg) tea.Cmd {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "enter":
			return m.submit()
		case "esc":
			return func() tea.Msg { return goBackMsg{} }
		}
	}
	var cmd tea.Cmd
	m.descInput, cmd = m.descInput.Update(msg)
	return cmd
}

func (m *createFeatureModel) updateRefiningState(msg tea.Msg) tea.Cmd {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "esc":
			m.state = stateDone
			return nil
		case "enter":
			return m.submitRefine()
		}
	}
	var cmd tea.Cmd
	m.refineInput, cmd = m.refineInput.Update(msg)
	return cmd
}

func (m *createFeatureModel) submitRefine() tea.Cmd {
	text := strings.TrimSpace(m.refineInput.Value())
	m.state = stateDone

	// Save identity of selected item (not raw index, which shifts with spacers)
	var savedRuneIdx int = -999
	var savedName string
	if item, ok := m.runeList.SelectedItem().(runeListItem); ok {
		savedRuneIdx = item.runeIdx
		savedName = item.name
	}

	if text == "" {
		// Empty comment clears an existing comment
		if m.inputMode == inputRefineRune {
			if item, ok := m.runeList.SelectedItem().(runeListItem); ok && item.runeIdx >= 0 {
				delete(m.runeComments, item.runeIdx)
			}
		} else {
			m.featureComment = ""
		}
	} else if m.inputMode == inputRefineRune {
		if item, ok := m.runeList.SelectedItem().(runeListItem); ok && item.runeIdx >= 0 {
			m.runeComments[item.runeIdx] = text
		}
	} else {
		m.featureComment = text
	}

	// Rebuild list to update comment markers
	m.buildRuneListItems()
	// Restore selection by identity
	for i, item := range m.runeList.Items() {
		if ri, ok := item.(runeListItem); ok && ri.runeIdx == savedRuneIdx && ri.name == savedName {
			m.runeList.Select(i)
			break
		}
	}
	return nil
}

func (m *createFeatureModel) submitAllRefinements() tea.Cmd {
	if !m.hasComments() {
		return nil
	}

	// Build refinement text from all collected comments
	var parts []string
	if m.featureComment != "" {
		parts = append(parts, m.featureComment)
	}
	for idx, comment := range m.runeComments {
		if idx >= 0 && idx < len(m.result.NewRunes) {
			name := m.result.NewRunes[idx].Name
			parts = append(parts, "For "+name+": "+comment)
		}
	}

	m.requirement = m.requirement + "\n\n" + strings.Join(parts, "\n")
	m.runeComments = map[int]string{}
	m.featureComment = ""
	m.runeList.Select(0)
	return m.decompose(m.requirement)
}

func (m *createFeatureModel) updateDoneState(msg tea.Msg) tea.Cmd {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "c":
			if item, ok := m.runeList.SelectedItem().(runeListItem); ok && item.isHeader {
				return m.openRefine(inputRefineFeature, "Comment on feature...")
			}
			name := m.selectedRuneName()
			return m.openRefine(inputRefineRune, "Comment on "+name+"...")
		case "a":
			return m.commitRunes()
		case "enter":
			return m.submitAllRefinements()
		case "backspace":
			return func() tea.Msg { return goBackMsg{} }
		}
	}
	prevIdx := m.runeList.Index()
	var cmd tea.Cmd
	m.runeList, cmd = m.runeList.Update(msg)
	m.skipSpacers(prevIdx)
	return cmd
}

// skipSpacers nudges the cursor off spacer items after navigation.
func (m *createFeatureModel) skipSpacers(prevIdx int) {
	items := m.runeList.Items()
	idx := m.runeList.Index()
	if idx < 0 || idx >= len(items) {
		return
	}
	ri, ok := items[idx].(runeListItem)
	if !ok || !ri.isSpacer {
		return
	}
	if idx >= prevIdx {
		if idx+1 < len(items) {
			m.runeList.Select(idx + 1)
		} else if idx-1 >= 0 {
			m.runeList.Select(idx - 1)
		}
	} else {
		if idx-1 >= 0 {
			m.runeList.Select(idx - 1)
		} else if idx+1 < len(items) {
			m.runeList.Select(idx + 1)
		}
	}
}

func (m *createFeatureModel) updateApprovedState(msg tea.Msg) tea.Cmd {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "h":
			return m.hydrateAll()
		case "backspace":
			return func() tea.Msg { return goBackMsg{} }
		}
	}
	return nil
}

func (m *createFeatureModel) updateAuthErrorState(msg tea.Msg) tea.Cmd {
	km, ok := msg.(tea.KeyMsg)
	if !ok || km.String() != "l" {
		return nil
	}
	exe, _ := os.Executable()
	c := exec.Command(exe, "login")
	c.Stdin = os.Stdin
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return loginDoneMsg{err: err}
	})
}

func (m *createFeatureModel) commitRunes() tea.Cmd {
	if m.result == nil || m.runeStore == nil {
		return nil
	}

	// Ensure draft is saved first
	m.saveDraft()
	if m.draftSuffix == "" {
		return func() tea.Msg {
			return commitErrorMsg{err: fmt.Errorf("no draft to commit")}
		}
	}

	suffix := m.draftSuffix
	runeStore := m.runeStore
	return func() tea.Msg {
		// Find all runes with this draft suffix
		allDrafts, err := runeStore.ListByStatus("draft")
		if err != nil {
			return commitErrorMsg{err: err}
		}
		for _, r := range allDrafts {
			if !hasDraftSuffix(r.Name, suffix) {
				continue
			}
			// Create the rune with the clean name (no suffix, no draft status)
			clean := r
			clean.Name = removeDraftSuffix(r.Name, suffix)
			clean.Status = ""
			// Delete existing if present, then create
			_ = runeStore.Delete(clean.Name)
			if err := runeStore.Create(clean); err != nil {
				// Skip if already exists (e.g. shared std rune)
				if !strings.Contains(err.Error(), "already exists") {
					return commitErrorMsg{err: err}
				}
			}
			// Delete the suffixed draft version
			_ = runeStore.Delete(r.Name)
		}
		return commitDoneMsg{}
	}
}

func (m *createFeatureModel) hydrateAll() tea.Cmd {
	port := m.port
	m.state = stateHydrating
	m.progressText = ""
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]any{"verify": true, "concurrency": 0})
		resp, err := http.Post(
			fmt.Sprintf("http://localhost:%d/api/hydrate", port),
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			return hydrateErrorMsg{err: err}
		}
		defer resp.Body.Close()

		var dr decomposeResponse
		if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
			return hydrateErrorMsg{err: err}
		}
		return decomposeStartedMsg{jobID: dr.JobID}
	}
}

func loadFeatureRunes(featureName string, port int) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/runes", port))
		if err != nil {
			return featureLoadErrorMsg{err: err}
		}
		defer resp.Body.Close()

		var allRunes []struct {
			Name          string   `json:"name"`
			Description   string   `json:"description"`
			Signature     string   `json:"signature"`
			PositiveTests []string `json:"positive_tests"`
			NegativeTests []string `json:"negative_tests"`
			Assumptions   []string `json:"assumptions"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&allRunes); err != nil {
			return featureLoadErrorMsg{err: err}
		}

		// Filter runes belonging to this feature (matching top-level package)
		var matched []proposedRune
		var summary string
		for _, r := range allRunes {
			pkg := r.Name
			if dot := strings.IndexByte(pkg, '.'); dot > 0 {
				pkg = pkg[:dot]
			}
			if pkg == featureName {
				if r.Name == featureName {
					summary = r.Description
				}
				matched = append(matched, proposedRune{
					Name:          r.Name,
					Description:   r.Description,
					Signature:     r.Signature,
					PositiveTests: r.PositiveTests,
					NegativeTests: r.NegativeTests,
					Assumptions:   r.Assumptions,
				})
			}
		}

		return featureLoadedMsg{feature: featureName, runes: matched, summary: summary}
	}
}

func (m *createFeatureModel) handleDecomposeError(msg decomposeErrorMsg) tea.Cmd {
	if strings.Contains(msg.err.Error(), "auth error") || strings.Contains(msg.err.Error(), "token expired") {
		m.state = stateAuthError
		m.errMsg = msg.err.Error()
		return nil
	}
	m.state = stateError
	m.errMsg = msg.err.Error()
	m.descInput.Focus()
	return nil
}

func (m *createFeatureModel) handleLoginDone(msg loginDoneMsg) tea.Cmd {
	if msg.err != nil {
		m.state = stateAuthError
		m.errMsg = fmt.Sprintf("login failed: %v", msg.err)
		return nil
	}
	m.state = stateIdle
	m.errMsg = ""
	m.authURL = ""
	m.descInput.Focus()
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return nil })
}

func (m *createFeatureModel) view(width int) string {
	switch m.state {
	case stateDecomposing:
		return m.viewDecomposing()
	case stateHydrating:
		return m.viewHydrating()
	case stateDone:
		return m.viewResult(width)
	case stateRefining:
		return m.viewRefining(width)
	case stateApproved:
		return m.viewApproved(width)
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
	text := "Decomposing into runes..."
	if m.progressText != "" {
		text = m.progressText
	}
	return m.spinner.View() + " " + text
}

func (m *createFeatureModel) viewHydrating() string {
	text := "Hydrating runes..."
	if m.progressText != "" {
		text = m.progressText
	}
	return m.spinner.View() + " " + text
}

func (m *createFeatureModel) viewApproved(width int) string {
	if m.result == nil {
		return ""
	}

	leftWidth := width / 3
	midWidth := width - leftWidth - 3
	if leftWidth < 20 {
		leftWidth = 20
	}
	if midWidth < 20 {
		midWidth = 20
	}

	availHeight := m.height - 6

	// Left pane: feature header + committed count
	var left strings.Builder
	featureHdr := paneHeaderActive.Render("── feature ") + featureNameStyle.Render(m.featureName()) + " "
	if remaining := leftWidth - lipgloss.Width(featureHdr); remaining > 0 {
		featureHdr += paneHeaderActive.Render(strings.Repeat("─", remaining))
	}
	left.WriteString(featureHdr + "\n")
	if m.result.Summary != "" {
		left.WriteString(featureSummaryStyle.Width(leftWidth - 2).Render(m.result.Summary) + "\n")
	}
	left.WriteString("\n")
	left.WriteString(statusOk.Render(fmt.Sprintf("  ✓ %d runes committed", len(m.result.NewRunes))) + "\n")

	// Right pane: composition tree
	var mid strings.Builder
	treeHdr := paneHeaderInactive.Render("── composition ") + " "
	if remaining := midWidth - lipgloss.Width(treeHdr); remaining > 0 {
		treeHdr += paneHeaderInactive.Render(strings.Repeat("─", remaining))
	}
	mid.WriteString(treeHdr + "\n\n")
	mid.WriteString(renderCompositionTree(m.result.NewRunes))

	m.midVP.Width = midWidth
	m.midVP.Height = availHeight
	m.midVP.SetContent(mid.String())

	sepChar := lipgloss.NewStyle().Foreground(border).Render("│")
	sepLines := make([]string, availHeight)
	for i := range sepLines {
		sepLines[i] = " " + sepChar + " "
	}

	leftContent := lipgloss.NewStyle().Width(leftWidth).Height(availHeight).Render(left.String())
	layout := lipgloss.JoinHorizontal(lipgloss.Top,
		leftContent,
		strings.Join(sepLines, "\n"),
		m.midVP.View(),
	)

	return layout
}

func (m *createFeatureModel) viewRefining(width int) string {
	var label string
	switch m.inputMode {
	case inputRefineFeature:
		label = "Comment on feature"
	case inputRefineRune:
		label = "Comment on " + m.selectedRuneName()
	}
	var b strings.Builder
	b.WriteString(m.viewResult(width))
	b.WriteString("\n")
	b.WriteString(inputLabel.Render(label) + " ")
	b.WriteString(m.refineInput.View())
	return b.String()
}

// sortRunesStdLast reorders runes so std package runes come after feature runes.
func sortRunesStdLast(runes []proposedRune) {
	sort.SliceStable(runes, func(i, j int) bool {
		ni, nj := runes[i].Name, runes[j].Name
		if dot := strings.IndexByte(ni, '.'); dot > 0 {
			ni = ni[:dot]
		}
		if dot := strings.IndexByte(nj, '.'); dot > 0 {
			nj = nj[:dot]
		}
		if ni == "std" && nj != "std" {
			return false
		}
		if ni != "std" && nj == "std" {
			return true
		}
		return false
	})
}

// runeGroup holds runes under a common top-level package.
type runeGroup struct {
	pkg string
	indices   []int // indices into NewRunes
}

// groupRunesByPackage groups runes by their top-level package (first dot segment).
func groupRunesByPackage(runes []proposedRune) []runeGroup {
	order := []string{}
	groups := map[string][]int{}
	for i, r := range runes {
		ns := r.Name
		if dot := strings.IndexByte(ns, '.'); dot > 0 {
			ns = ns[:dot]
		}
		if _, ok := groups[ns]; !ok {
			order = append(order, ns)
		}
		groups[ns] = append(groups[ns], i)
	}
	sort.SliceStable(order, func(i, j int) bool {
		if order[i] == "std" {
			return false
		}
		if order[j] == "std" {
			return true
		}
		return false
	})
	result := make([]runeGroup, len(order))
	for i, ns := range order {
		result[i] = runeGroup{pkg: ns, indices: groups[ns]}
	}
	return result
}

// featureName returns the API-provided name or derives one from the rune packages.
func (m *createFeatureModel) featureName() string {
	if m.result.FeatureName != "" {
		return m.result.FeatureName
	}
	for _, r := range m.result.NewRunes {
		ns := r.Name
		if dot := strings.IndexByte(ns, '.'); dot > 0 {
			ns = ns[:dot]
		}
		if ns != "std" {
			return ns
		}
	}
	return "feature"
}

// leafName returns the part of a dot-path after the top-level package.
func leafName(fullPath string) string {
	if dot := strings.IndexByte(fullPath, '.'); dot > 0 {
		return fullPath[dot+1:]
	}
	return fullPath
}

// buildRuneListItems populates the rune list from the current decomposition result.
func (m *createFeatureModel) buildRuneListItems() {
	if m.result == nil {
		m.runeList.SetItems(nil)
		return
	}

	groups := groupRunesByPackage(m.result.NewRunes)
	var items []list.Item

	for gi, g := range groups {
		if gi > 0 {
			items = append(items, runeListItem{isSpacer: true, runeIdx: -1})
		}
		selfIdx := -1
		for _, idx := range g.indices {
			if leafName(m.result.NewRunes[idx].Name) == g.pkg {
				selfIdx = idx
				break
			}
		}

		// Count visible children (non-empty runes)
		visibleCount := 0
		for _, idx := range g.indices {
			if idx == selfIdx {
				continue
			}
			r := m.result.NewRunes[idx]
			if r.Description != "" || r.Signature != "" || len(r.Assumptions) > 0 {
				visibleCount++
			}
		}

		_, selfHasComment := m.runeComments[selfIdx]
		items = append(items, runeListItem{
			runeIdx:    selfIdx,
			name:       g.pkg,
			isHeader:   true,
			count:      visibleCount,
			hasComment: selfHasComment,
		})

		for _, idx := range g.indices {
			if idx == selfIdx {
				continue
			}
			r := m.result.NewRunes[idx]
			// Skip empty structural entries (intermediate path nodes)
			if r.Description == "" && r.Signature == "" && len(r.Assumptions) == 0 {
				continue
			}
			_, hasComment := m.runeComments[idx]
			items = append(items, runeListItem{
				runeIdx:    idx,
				name:       leafName(r.Name),
				hasComment: hasComment,
			})
		}
	}

	if len(m.result.ExistingRunes) > 0 {
		if len(items) > 0 {
			items = append(items, runeListItem{isSpacer: true, runeIdx: -1})
		}
		items = append(items, runeListItem{
			name:       "existing",
			isHeader:   true,
			isExisting: true,
			runeIdx:    -1,
		})
		for _, r := range m.result.ExistingRunes {
			items = append(items, runeListItem{
				runeIdx:    -1,
				name:       r.Name,
				isExisting: true,
				covers:     r.Covers,
			})
		}
	}

	m.runeList.SetItems(items)
}

// isPackage returns true if the named rune has children (other runes prefixed by name+".").
func isPackage(name string, runes []proposedRune) bool {
	prefix := name + "."
	for _, r := range runes {
		if strings.HasPrefix(r.Name, prefix) {
			return true
		}
	}
	return false
}

// childRunes returns all runes whose name is prefixed by name+".".
func childRunes(name string, runes []proposedRune) []proposedRune {
	prefix := name + "."
	var children []proposedRune
	for _, r := range runes {
		if strings.HasPrefix(r.Name, prefix) {
			children = append(children, r)
		}
	}
	return children
}

var (
	treeRefStyle  = lipgloss.NewStyle().Foreground(dim)
	treeLineStyle = lipgloss.NewStyle().Foreground(border)
)

func renderCompositionTree(runes []proposedRune) string {
	refsMap := map[string][]string{}
	for _, r := range runes {
		if len(r.Refs) > 0 {
			refsMap[r.Name] = r.Refs
		}
	}

	allPaths := map[string]bool{}
	for _, r := range runes {
		parts := strings.Split(r.Name, ".")
		for i := 1; i <= len(parts); i++ {
			allPaths[strings.Join(parts[:i], ".")] = true
		}
	}
	pathSlice := make([]string, 0, len(allPaths))
	for p := range allPaths {
		pathSlice = append(pathSlice, p)
	}
	sort.Strings(pathSlice)

	childrenMap := runepkg.BuildChildrenMap(pathSlice)

	seen := map[string]bool{}
	var roots []string
	for _, r := range runes {
		ns := r.Name
		if dot := strings.IndexByte(ns, '.'); dot > 0 {
			ns = ns[:dot]
		}
		if !seen[ns] {
			seen[ns] = true
			roots = append(roots, ns)
		}
	}

	var b strings.Builder
	for i, root := range roots {
		renderNode(&b, root, childrenMap, refsMap, "", true, 0)
		if i < len(roots)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func renderNode(b *strings.Builder, path string, childrenMap map[string][]string, refsMap map[string][]string, prefix string, isLast bool, depth int) {
	leaf := path
	if dot := strings.LastIndexByte(path, '.'); dot >= 0 {
		leaf = path[dot+1:]
	}

	if depth == 0 {
		b.WriteString(pkgStyle.Render(leaf) + "\n")
	} else {
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		b.WriteString(prefix + treeLineStyle.Render(connector) + runeLeafStyle.Render(leaf) + "\n")
	}

	var childPrefix string
	if depth == 0 {
		childPrefix = ""
	} else if isLast {
		childPrefix = prefix + "    "
	} else {
		childPrefix = prefix + treeLineStyle.Render("│") + "   "
	}

	refs := refsMap[path]
	children := childrenMap[path]
	total := len(refs) + len(children)

	for i, ref := range refs {
		c := "├── "
		if i+len(children) == total-1 {
			c = "└── "
		}
		b.WriteString(childPrefix + treeLineStyle.Render(c) + treeRefStyle.Render("-> "+ref) + "\n")
	}

	for i, child := range children {
		childIsLast := len(refs)+i == total-1
		renderNode(b, child, childrenMap, refsMap, childPrefix, childIsLast, depth+1)
	}
}

func (m *createFeatureModel) viewResult(width int) string {
	if m.result == nil {
		return ""
	}

	// Width allocation
	leftWidth := width / 3
	midWidth := width - leftWidth - 3
	if leftWidth < 20 {
		leftWidth = 20
	}
	if midWidth < 20 {
		midWidth = 20
	}

	footerLines := 2 // status line + help bar
	availHeight := m.height - 4 - footerLines
	if m.state == stateRefining {
		availHeight -= 2
	}
	if availHeight < 5 {
		availHeight = 5
	}

	// Left pane: feature header + rune list
	var header strings.Builder
	hdrStyle := paneHeaderActive
	featureHdr := hdrStyle.Render("── feature ") + featureNameStyle.Render(m.featureName()) + " "
	if remaining := leftWidth - lipgloss.Width(featureHdr); remaining > 0 {
		featureHdr += hdrStyle.Render(strings.Repeat("─", remaining))
	}
	header.WriteString(featureHdr + "\n")
	headerLines := 1
	if m.result.Summary != "" {
		wrapped := featureSummaryStyle.Width(leftWidth - 2).Render(m.result.Summary)
		headerLines += strings.Count(wrapped, "\n") + 1
		header.WriteString(wrapped + "\n")
	}
	header.WriteString("\n")
	headerLines++

	m.runeList.SetSize(leftWidth, availHeight-headerLines)
	leftContent := header.String() + m.runeList.View()

	// Middle pane: rune or package detail
	var mid strings.Builder

	selectedRuneIdx := -1
	var selectedExisting *existingMatch
	if item, ok := m.runeList.SelectedItem().(runeListItem); ok {
		if item.runeIdx >= 0 && item.runeIdx < len(m.result.NewRunes) {
			selectedRuneIdx = item.runeIdx
		} else if item.isExisting && !item.isHeader {
			for i := range m.result.ExistingRunes {
				if m.result.ExistingRunes[i].Name == item.name {
					selectedExisting = &m.result.ExistingRunes[i]
					break
				}
			}
		}
	}

	if selectedRuneIdx >= 0 {
		r := m.result.NewRunes[selectedRuneIdx]

		if isPackage(r.Name, m.result.NewRunes) {
			// Package view
			pkgHdr := paneHeaderInactive.Render("── package ") + pkgStyle.Render(r.Name) + " "
			if remaining := midWidth - lipgloss.Width(pkgHdr); remaining > 0 {
				pkgHdr += paneHeaderInactive.Render(strings.Repeat("─", remaining))
			}
			mid.WriteString(pkgHdr + "\n")

			children := childRunes(r.Name, m.result.NewRunes)

			if summary := m.result.PackageSummaries[r.Name]; summary != "" {
				mid.WriteString("\n" + lipgloss.NewStyle().Width(midWidth-2).Render(summary) + "\n")
			}

			if len(children) > 0 {
				mid.WriteString("\n")
				for _, child := range children {
					if child.Description == "" && child.Signature == "" && len(child.Assumptions) == 0 {
						continue
					}
					leaf := leafName(child.Name)
					title := runeNameStyle.Render(leaf)
					if child.Signature != "" {
						title += " " + runeSigStyle.Render(child.Signature)
					}
					mid.WriteString(lipgloss.NewStyle().Width(midWidth-2).Render("  "+title) + "\n")
					if child.Description != "" {
						mid.WriteString(lipgloss.NewStyle().Foreground(dim).Width(midWidth-2).PaddingLeft(4).Render(child.Description) + "\n")
					}
					mid.WriteString("\n")
				}
			}

			// Aggregate assumptions from package and all children
			var allAssumptions []string
			allAssumptions = append(allAssumptions, r.Assumptions...)
			for _, child := range children {
				allAssumptions = append(allAssumptions, child.Assumptions...)
			}
			if len(allAssumptions) > 0 {
				mid.WriteString("\n" + runeSigStyle.Render("assumes:") + "\n")
				for _, a := range allAssumptions {
					mid.WriteString(assumptionStyle.Render("? ") + lipgloss.NewStyle().Width(midWidth-4).Render(a) + "\n")
				}
			}

			// Aggregate links from package and all children
			var allRefs []string
			allRefs = append(allRefs, r.Refs...)
			for _, child := range children {
				allRefs = append(allRefs, child.Refs...)
			}
			if len(allRefs) > 0 {
				mid.WriteString("\n" + runeSigStyle.Render("dependencies:") + "\n")
				for _, ref := range allRefs {
					mid.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6A9FD9")).Render("-> ") + lipgloss.NewStyle().Foreground(dim).Width(midWidth-4).Render(ref) + "\n")
				}
			}

			// Show comment if present
			if comment, ok := m.runeComments[selectedRuneIdx]; ok {
				commentHdr := "\n" + paneHeaderInactive.Render("── your comment ") + " "
				if remaining := midWidth - lipgloss.Width(commentHdr); remaining > 0 {
					commentHdr += paneHeaderInactive.Render(strings.Repeat("─", remaining))
				}
				mid.WriteString(commentHdr + "\n")
				mid.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623")).Width(midWidth-2).Render(comment) + "\n")
			}
		} else {
			// Rune view
			runeHdr := paneHeaderInactive.Render("── rune ") + runeNameStyle.Render(r.Name) + " "
			if remaining := midWidth - lipgloss.Width(runeHdr); remaining > 0 {
				runeHdr += paneHeaderInactive.Render(strings.Repeat("─", remaining))
			}
			mid.WriteString(runeHdr + "\n")

			if r.Signature != "" {
				mid.WriteString("\n" + runeSigStyle.Width(midWidth-2).Render(r.Signature) + "\n")
			}

			if r.Description != "" {
				mid.WriteString(lipgloss.NewStyle().Width(midWidth-2).Render(r.Description) + "\n")
			}

			if len(r.PositiveTests) > 0 || len(r.NegativeTests) > 0 {
				mid.WriteString("\n")
				for _, t := range r.PositiveTests {
					mid.WriteString(testPassStyle.Render("+ ") + lipgloss.NewStyle().Width(midWidth-4).Render(t) + "\n")
				}
				for _, t := range r.NegativeTests {
					mid.WriteString(testFailStyle.Render("- ") + lipgloss.NewStyle().Width(midWidth-4).Render(t) + "\n")
				}
			}

			if len(r.Assumptions) > 0 {
				mid.WriteString("\n" + runeSigStyle.Render("assumes:") + "\n")
				for _, a := range r.Assumptions {
					mid.WriteString(assumptionStyle.Render("? ") + lipgloss.NewStyle().Width(midWidth-4).Render(a) + "\n")
				}
			}

			if len(r.Refs) > 0 {
				mid.WriteString("\n" + runeSigStyle.Render("dependencies:") + "\n")
				for _, ref := range r.Refs {
					mid.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6A9FD9")).Render("-> ") + lipgloss.NewStyle().Foreground(dim).Width(midWidth-4).Render(ref) + "\n")
				}
			}

			// Show comment if present
			if comment, ok := m.runeComments[selectedRuneIdx]; ok {
				commentHdr := "\n" + paneHeaderInactive.Render("── your comment ") + " "
				if remaining := midWidth - lipgloss.Width(commentHdr); remaining > 0 {
					commentHdr += paneHeaderInactive.Render(strings.Repeat("─", remaining))
				}
				mid.WriteString(commentHdr + "\n")
				mid.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623")).Width(midWidth-2).Render(comment) + "\n")
			}
		}
	} else if selectedExisting != nil {
		hdr := paneHeaderInactive.Render("── existing ") + lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render(selectedExisting.Name) + " "
		if remaining := midWidth - lipgloss.Width(hdr); remaining > 0 {
			hdr += paneHeaderInactive.Render(strings.Repeat("─", remaining))
		}
		mid.WriteString(hdr + "\n")
		if selectedExisting.Covers != "" {
			mid.WriteString("\n" + lipgloss.NewStyle().Width(midWidth-2).Render(selectedExisting.Covers) + "\n")
		}
	} else {
		mid.WriteString(renderPaneHeader("rune", midWidth, false) + "\n")
	}

	// Feature-level comment and pending count (shown in description pane)
	if m.featureComment != "" {
		featureCommentHdr := "\n" + paneHeaderInactive.Render("── feature comment ") + " "
		if remaining := midWidth - lipgloss.Width(featureCommentHdr); remaining > 0 {
			featureCommentHdr += paneHeaderInactive.Render(strings.Repeat("─", remaining))
		}
		mid.WriteString(featureCommentHdr + "\n")
		mid.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623")).Width(midWidth-2).Render(m.featureComment) + "\n")
	}
	if count := m.commentCount(); count > 0 {
		mid.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623")).Render(
			fmt.Sprintf("%d comment%s pending", count, map[bool]string{true: "", false: "s"}[count == 1])) + "\n")
	}

	// Update description pane viewport
	m.midVP.Width = midWidth
	m.midVP.Height = availHeight
	m.midVP.SetContent(mid.String())

	// Build separator column
	sepChar := lipgloss.NewStyle().Foreground(border).Render("│")
	sepLines := make([]string, availHeight)
	for i := range sepLines {
		sepLines[i] = " " + sepChar + " "
	}

	// Join panes
	layout := lipgloss.JoinHorizontal(lipgloss.Top,
		leftContent,
		strings.Join(sepLines, "\n"),
		m.midVP.View(),
	)

	// Footer with status and comments
	var footer strings.Builder
	if m.errMsg != "" {
		footer.WriteString(statusErr.Render(m.errMsg))
		m.errMsg = ""
	} else {
		footer.WriteString(statusOk.Render(fmt.Sprintf("%d runes proposed", len(m.result.NewRunes))))
	}

	return layout + "\n\n" + footer.String()
}
