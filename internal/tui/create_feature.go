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

	"sort"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	runepkg "github.com/chrishayen/valkyrie/internal/rune"
)

type formState int

const (
	stateIdle formState = iota
	stateDecomposing
	stateDone
	stateError
	stateAuthError
	stateRefining
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

	namespaceStyle = lipgloss.NewStyle().
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
	Refs          []string `json:"refs"`
}

type existingMatch struct {
	Name   string `json:"name"`
	Covers string `json:"covers"`
}

type decomposeResult struct {
	FeatureName   string          `json:"feature_name,omitempty"`
	Summary       string          `json:"summary,omitempty"`
	FlowDiagram   string          `json:"flow_diagram,omitempty"`
	NewRunes      []proposedRune  `json:"new_runes"`
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

type rightPane int

const (
	paneDetail rightPane = iota
	paneTree
	paneFlow
)

type createFeatureModel struct {
	descInput   textarea.Model
	refineInput textinput.Model
	state       formState
	port        int
	width       int
	jobID       string
	spinner     spinner.Model
	result      *decomposeResult
	errMsg      string
	authURL     string
	runeCursor  int
	height      int
	rightPane   rightPane
	leftScroll  int
	requirement string // accumulated requirements text
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
		height:    height,
		spinner:   s,
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

func (m *createFeatureModel) decompose(req string) tea.Cmd {
	m.state = stateDecomposing
	port := m.port
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{"requirement": req})
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
		if m.state == stateRefining {
			switch msg.String() {
			case "esc":
				m.state = stateDone
				return nil
			case "enter":
				refinement := strings.TrimSpace(m.refineInput.Value())
				if refinement == "" {
					m.state = stateDone
					return nil
				}
				m.requirement = m.requirement + "\n\n" + refinement
				m.runeCursor = 0
				m.leftScroll = 0
				m.rightPane = paneDetail
				return m.decompose(m.requirement)
			}
			var cmd tea.Cmd
			m.refineInput, cmd = m.refineInput.Update(msg)
			return cmd
		}
		if m.state == stateDone {
			switch msg.String() {
			case "j", "down":
				if m.result != nil && m.runeCursor < len(m.result.NewRunes)-1 {
					m.runeCursor++
				}
				return nil
			case "k", "up":
				if m.runeCursor > 0 {
					m.runeCursor--
				}
				return nil
			case "t":
				m.rightPane = (m.rightPane + 1) % 3
				return nil
			case "r":
				m.state = stateRefining
				m.refineInput = textinput.New()
				m.refineInput.Placeholder = "Add refinement..."
				m.refineInput.Width = m.width - 4
				m.refineInput.Focus()
				return m.refineInput.Cursor.BlinkCmd()
			case "enter":
				m.state = stateIdle
				m.result = nil
				m.requirement = ""
				m.runeCursor = 0
				m.leftScroll = 0
				m.rightPane = paneDetail
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
		m.runeCursor = 0
		m.leftScroll = 0
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

	if m.state == stateRefining {
		var cmd tea.Cmd
		m.refineInput, cmd = m.refineInput.Update(msg)
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
	case stateRefining:
		return m.viewRefining(width)
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

func (m *createFeatureModel) viewRefining(width int) string {
	var b strings.Builder
	b.WriteString(inputLabel.Render("Refine") + " ")
	b.WriteString(m.refineInput.View())
	b.WriteString("\n\n")
	b.WriteString(m.viewResult(width))
	return b.String()
}

// runeGroup holds runes under a common namespace.
type runeGroup struct {
	namespace string
	indices   []int // indices into NewRunes
}

// groupRunesByNamespace groups runes by their top-level namespace (first dot segment).
func groupRunesByNamespace(runes []proposedRune) []runeGroup {
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
	result := make([]runeGroup, len(order))
	for i, ns := range order {
		result[i] = runeGroup{namespace: ns, indices: groups[ns]}
	}
	return result
}

// featureName returns the API-provided name or derives one from the rune namespaces.
func (m *createFeatureModel) featureName() string {
	if m.result.FeatureName != "" {
		return m.result.FeatureName
	}
	// Derive: use the first non-std namespace, or fall back to "feature"
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

// leafName returns the part of a dot-path after the top-level namespace.
func leafName(fullPath string) string {
	if dot := strings.IndexByte(fullPath, '.'); dot > 0 {
		return fullPath[dot+1:]
	}
	return fullPath
}

var (
	treeRefStyle  = lipgloss.NewStyle().Foreground(dim)
	treeLineStyle = lipgloss.NewStyle().Foreground(border)
)

func renderCompositionTree(runes []proposedRune) string {
	// Build refs map
	refsMap := map[string][]string{}
	for _, r := range runes {
		if len(r.Refs) > 0 {
			refsMap[r.Name] = r.Refs
		}
	}

	// Collect all unique path prefixes to synthesize intermediate nodes
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

	// Find roots (paths with no dot = top-level namespaces), preserve order from runes
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
	// Extract leaf segment
	leaf := path
	if dot := strings.LastIndexByte(path, '.'); dot >= 0 {
		leaf = path[dot+1:]
	}

	// Render this node's line
	if depth == 0 {
		b.WriteString(namespaceStyle.Render(leaf) + "\n")
	} else {
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		b.WriteString(prefix + treeLineStyle.Render(connector) + runeLeafStyle.Render(leaf) + "\n")
	}

	// Child prefix for deeper levels
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

	leftWidth := width / 3
	if leftWidth < 25 {
		leftWidth = 25
	}
	rightWidth := width - leftWidth - 3

	groups := groupRunesByNamespace(m.result.NewRunes)

	// Left pane: feature header + grouped rune list
	var left strings.Builder
	cursorLine := 0
	lineNum := 0

	// Feature name
	left.WriteString(featureNameStyle.Render(m.featureName()) + "\n")
	lineNum++
	if m.result.Summary != "" {
		left.WriteString(featureSummaryStyle.Render(m.result.Summary) + "\n")
		lineNum++
	}
	left.WriteString("\n")
	lineNum++

	// Grouped rune list
	for gi, g := range groups {
		left.WriteString(namespaceStyle.Render(g.namespace) +
			runeSigStyle.Render(fmt.Sprintf(" (%d)", len(g.indices))) + "\n")
		lineNum++

		for _, idx := range g.indices {
			r := m.result.NewRunes[idx]
			leaf := leafName(r.Name)
			if len(leaf) > leftWidth-6 {
				leaf = leaf[:leftWidth-7] + "~"
			}
			if idx == m.runeCursor {
				cursorLine = lineNum
				left.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FFFFFF")).
					Background(lipgloss.Color("#333333")).
					Width(leftWidth - 2).
					Render("  > " + leaf) + "\n")
			} else {
				left.WriteString(runeLeafStyle.Render("    "+leaf) + "\n")
			}
			lineNum++
		}

		if gi < len(groups)-1 {
			left.WriteString("\n")
			lineNum++
		}
	}

	if len(m.result.ExistingRunes) > 0 {
		left.WriteString("\n" + lipgloss.NewStyle().Foreground(dim).Italic(true).Render("existing") + "\n")
		lineNum += 2
		for _, r := range m.result.ExistingRunes {
			name := r.Name
			if len(name) > leftWidth-6 {
				name = name[:leftWidth-7] + "~"
			}
			left.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).
				Render("    " + name) + "\n")
			lineNum++
		}
	}

	// Right pane: detail, tree, or flow diagram
	var right strings.Builder
	switch m.rightPane {
	case paneTree:
		right.WriteString(renderCompositionTree(m.result.NewRunes))
	case paneFlow:
		if m.result.FlowDiagram != "" {
			right.WriteString(m.result.FlowDiagram)
		} else {
			right.WriteString(runeSigStyle.Render("No flow diagram available"))
		}
	default:
	}
	if m.rightPane == paneDetail && m.runeCursor < len(m.result.NewRunes) {
		r := m.result.NewRunes[m.runeCursor]
		right.WriteString(runeNameStyle.Render(r.Name) + "\n\n")

		if r.Signature != "" {
			right.WriteString(runeSigStyle.Render(r.Signature) + "\n\n")
		}

		if r.Description != "" {
			right.WriteString(lipgloss.NewStyle().Width(rightWidth - 2).Render(r.Description) + "\n")
		}

		if len(r.PositiveTests) > 0 || len(r.NegativeTests) > 0 {
			right.WriteString("\n")
			for _, t := range r.PositiveTests {
				right.WriteString(testPassStyle.Render("+ ") + lipgloss.NewStyle().Width(rightWidth - 4).Render(t) + "\n")
			}
			for _, t := range r.NegativeTests {
				right.WriteString(testFailStyle.Render("- ") + lipgloss.NewStyle().Width(rightWidth - 4).Render(t) + "\n")
			}
		}
	}

	// Separator
	sep := lipgloss.NewStyle().Foreground(border).Render("│")

	// Join panes line by line
	leftLines := strings.Split(left.String(), "\n")
	rightLines := strings.Split(right.String(), "\n")

	paneHeight := m.height - 10
	if paneHeight < 5 {
		paneHeight = 5
	}

	// Scroll left pane to keep cursor visible
	if cursorLine < m.leftScroll {
		m.leftScroll = cursorLine
	}
	if cursorLine >= m.leftScroll+paneHeight {
		m.leftScroll = cursorLine - paneHeight + 1
	}
	if m.leftScroll < 0 {
		m.leftScroll = 0
	}

	// Apply scroll to left pane
	if m.leftScroll < len(leftLines) {
		leftLines = leftLines[m.leftScroll:]
	}
	if len(leftLines) > paneHeight {
		leftLines = leftLines[:paneHeight]
	}

	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}
	if maxLines > paneHeight {
		maxLines = paneHeight
	}

	var b strings.Builder
	for i := 0; i < maxLines; i++ {
		l := ""
		if i < len(leftLines) {
			l = leftLines[i]
		}
		r := ""
		if i < len(rightLines) {
			r = rightLines[i]
		}
		lRendered := lipgloss.NewStyle().Width(leftWidth).Render(l)
		b.WriteString(lRendered + " " + sep + " " + r + "\n")
	}

	b.WriteString("\n" + statusOk.Render(fmt.Sprintf("%d runes proposed", len(m.result.NewRunes))))

	return b.String()
}
