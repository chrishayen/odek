package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"shotgun.dev/odek/decompose"
	openai "shotgun.dev/odek/openai"
)

type formState int

const (
	stateIdle formState = iota
	stateDecomposing
	stateDone
	stateError
)

var (
	statusOk  = lipgloss.NewStyle().Foreground(lipgloss.Color("#66CC66"))
	statusErr = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC6666"))

	featureNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F5A623")).Bold(true)
	featureSummaryStyle = lipgloss.NewStyle().
				Foreground(dim).Italic(true)

	pkgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6A9FD9")).Bold(true)
	runeNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5A623")).Bold(true)
	runeLeafStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
	runeSigStyle  = lipgloss.NewStyle().Foreground(dim)

	testPassStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#66CC66"))
	assumptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D4A843"))

	paneHeaderActive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F5A623")).Bold(true)
	paneHeaderInactive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555"))

	treeRefStyle  = lipgloss.NewStyle().Foreground(dim)
	treeLineStyle = lipgloss.NewStyle().Foreground(border)
)

type runeListItem struct {
	runeIdx  int
	name     string
	fullPath string
	isHeader bool
	isRef    bool
	refName  string
	isSpacer bool
	count    int
	refCount int
}

func (i runeListItem) Title() string       { return i.name }
func (i runeListItem) Description() string { return "" }
func (i runeListItem) FilterValue() string { return i.name }

type runeListDelegate struct{}

func (d runeListDelegate) Height() int                             { return 1 }
func (d runeListDelegate) Spacing() int                            { return 0 }
func (d runeListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d runeListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	ri, ok := item.(runeListItem)
	if !ok || ri.isSpacer {
		return
	}

	availWidth := m.Width()
	selected := index == m.Index()

	selectedBorder := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#F5A623")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Width(availWidth)

	var str string
	if ri.isHeader {
		countStr := ""
		if ri.count > 0 {
			countStr = fmt.Sprintf(" (%d)", ri.count)
		}
		if selected {
			str = selectedBorder.Bold(true).Render(ri.name + countStr)
		} else {
			str = lipgloss.NewStyle().Width(availWidth).
				Render(" " + pkgStyle.Render(ri.name) + runeSigStyle.Render(countStr))
		}
	} else {
		name := ri.name
		if len(name) > availWidth-6 {
			name = name[:availWidth-7] + "~"
		}
		if selected {
			str = selectedBorder.Padding(0, 0, 0, 3).Render(name)
		} else {
			str = runeLeafStyle.Width(availWidth).Render("    " + name)
		}
	}
	fmt.Fprint(w, str)
}

type decomposeDoneMsg struct{ result *decompose.Decomposition }
type decomposeErrorMsg struct{ err error }

type decomposeModel struct {
	ctx            context.Context
	client         *openai.Client
	systemPrompt   string
	feature        string
	width          int
	height         int
	state          formState
	result         *decompose.Decomposition
	errMsg         string
	spinner        spinner.Model
	runeList       list.Model
	midVP          viewport.Model
	groupByPackage bool
}

func newDecomposeModel(ctx context.Context, client *openai.Client, systemPrompt, feature string) decomposeModel {
	width, height := 120, 40

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#F5A623"))

	runeList := list.New(nil, runeListDelegate{}, width/3, height-6)
	runeList.SetShowTitle(false)
	runeList.SetShowStatusBar(false)
	runeList.SetFilteringEnabled(false)
	runeList.SetShowHelp(false)
	runeList.SetShowPagination(false)

	midVP := viewport.New(width-width/3-3, height-6)
	midVP.KeyMap = viewport.KeyMap{}

	return decomposeModel{
		ctx:            ctx,
		client:         client,
		systemPrompt:   systemPrompt,
		feature:        feature,
		width:          width,
		height:         height,
		state:          stateIdle,
		spinner:        sp,
		runeList:       runeList,
		midVP:          midVP,
		groupByPackage: true,
	}
}

func RunDecompose(ctx context.Context, client *openai.Client, systemPrompt, feature string) {
	m := newDecomposeModel(ctx, client, systemPrompt, feature)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m decomposeModel) Init() tea.Cmd {
	return m.startDecompose()
}

func (m decomposeModel) startDecompose() tea.Cmd {
	feature := m.feature
	client := m.client
	systemPrompt := m.systemPrompt
	ctx := m.ctx

	return func() tea.Msg {
		result, err := decompose.DecomposeStructured(ctx, client, systemPrompt, feature)
		if err != nil {
			return decomposeErrorMsg{err: err}
		}
		return decomposeDoneMsg{result: result.Decomposition}
	}
}

func (m decomposeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case decomposeDoneMsg:
		m.state = stateDone
		m.result = msg.result
		m.buildRuneListItems()
		if len(m.runeList.Items()) > 0 {
			m.runeList.Select(0)
		}
		return m, nil

	case decomposeErrorMsg:
		m.state = stateError
		m.errMsg = msg.err.Error()
		return m, nil

	case spinner.TickMsg:
		if m.state == stateDecomposing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	switch m.state {
	case stateIdle, stateDecomposing:
		return m.updateLoadingState(msg)
	case stateDone:
		return m.updateDoneState(msg)
	case stateError:
		return m.updateErrorState(msg)
	}
	return m, nil
}

func (m decomposeModel) updateLoadingState(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "esc", "backspace":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m decomposeModel) updateDoneState(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "backspace":
			return m, tea.Quit
		case "g":
			m.groupByPackage = !m.groupByPackage
			m.buildRuneListItems()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.runeList, cmd = m.runeList.Update(msg)
	return m, cmd
}

func (m decomposeModel) updateErrorState(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "esc", "backspace":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m decomposeModel) View() string {
	switch m.state {
	case stateIdle, stateDecomposing:
		return m.viewLoading()
	case stateDone:
		return m.viewResult(m.width)
	case stateError:
		return m.viewError()
	}
	return ""
}

func (m decomposeModel) viewLoading() string {
	text := "Decomposing into runes..."
	return m.spinner.View() + " " + text
}

func (m decomposeModel) viewError() string {
	return statusErr.Render("Error: "+m.errMsg) + "\n\nPress backspace to exit"
}

func (m *decomposeModel) buildRuneListItems() {
	if m.result == nil || m.result.RuneTree == nil {
		m.runeList.SetItems(nil)
		return
	}

	var runes []*decompose.Rune
	collectRunes(m.result.RuneTree, &runes)

	groups := m.groupRunes(runes)
	var items []list.Item

	pkgOrder := 0
	for gi, g := range groups {
		if gi > 0 {
			items = append(items, runeListItem{isSpacer: true})
		}

		selfIdx := -1
		for _, idx := range g.indices {
			if m.leafName(runes[idx].Path) == g.pkg {
				selfIdx = idx
				break
			}
		}

		hasContent := false
		visibleCount := 0
		for _, idx := range g.indices {
			if idx == selfIdx {
				continue
			}
			r := runes[idx]
			if r.Description != "" || r.Signature != "" || len(r.Tests) > 0 {
				hasContent = true
				visibleCount++
			}
		}

		if !hasContent && g.pkg == "std" {
			continue
		}

		items = append(items, runeListItem{
			runeIdx:  selfIdx,
			name:     g.pkg,
			isHeader: true,
			count:    visibleCount,
		})

		for _, idx := range g.indices {
			if idx == selfIdx {
				continue
			}
			r := runes[idx]
			if r.Description == "" && r.Signature == "" && len(r.Tests) == 0 {
				continue
			}
			items = append(items, runeListItem{
				runeIdx:  idx,
				name:     m.leafName(r.Path),
				fullPath: r.Path,
			})
		}

		pkgOrder++
	}

	m.runeList.SetItems(items)
}

func collectRunes(rune *decompose.Rune, out *[]*decompose.Rune) {
	if rune == nil {
		return
	}
	*out = append(*out, rune)
	for _, child := range rune.Children {
		collectRunes(child, out)
	}
}

type runeGroup struct {
	pkg     string
	indices []int
}

func (m *decomposeModel) groupRunes(runes []*decompose.Rune) []runeGroup {
	order := []string{}
	groups := map[string][]int{}

	for i, r := range runes {
		pkg := m.getPackageName(r.Path)
		if pkg == "" {
			continue
		}
		if _, ok := groups[pkg]; !ok {
			order = append(order, pkg)
		}
		groups[pkg] = append(groups[pkg], i)
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
	for i, pkg := range order {
		result[i] = runeGroup{pkg: pkg, indices: groups[pkg]}
	}
	return result
}

func (m *decomposeModel) getFeatureName() string {
	if m.result.FeatureName != "" {
		return m.result.FeatureName
	}
	if m.result.RuneTree != nil && m.result.RuneTree.Path != "" {
		parts := strings.Split(m.result.RuneTree.Path, "/")
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] != "std" {
				return parts[i]
			}
		}
	}
	return "feature"
}

func (m *decomposeModel) leafName(fullPath string) string {
	parts := strings.Split(fullPath, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullPath
}

func (m *decomposeModel) getPackageName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 1 && parts[0] != "" {
		return parts[0]
	}
	return ""
}

func (m *decomposeModel) getRuneByIndex(idx int) *decompose.Rune {
	var allRunes []*decompose.Rune
	if m.result.RuneTree != nil {
		collectRunes(m.result.RuneTree, &allRunes)
	}
	if idx < 0 || idx >= len(allRunes) {
		return nil
	}
	return allRunes[idx]
}

func (m *decomposeModel) countVisibleRunes() int {
	count := 0
	for _, item := range m.runeList.Items() {
		if ri, ok := item.(runeListItem); ok && !ri.isHeader && !ri.isSpacer {
			count++
		}
	}
	return count
}

func renderPaneHeader(label string, width int, active bool) string {
	style := paneHeaderInactive
	if active {
		style = paneHeaderActive
	}
	prefix := "── "
	inner := prefix + label
	remaining := width - lipgloss.Width(inner)
	if remaining < 0 {
		remaining = 0
	}
	return style.Render(inner + strings.Repeat("─", remaining))
}

func (m *decomposeModel) viewResult(width int) string {
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

	footerLines := 1
	availHeight := m.height - 4 - footerLines
	if availHeight < 5 {
		availHeight = 5
	}

	var header strings.Builder
	hdrStyle := paneHeaderActive
	featureName := m.getFeatureName()
	featureHdr := hdrStyle.Render("── feature ") + featureNameStyle.Render(featureName) + " "
	if remaining := leftWidth - lipgloss.Width(featureHdr); remaining > 0 {
		featureHdr += hdrStyle.Render(strings.Repeat("─", remaining))
	}
	header.WriteString(featureHdr + "\n")

	headerLines := 1
	if m.result.Description != "" {
		wrapped := featureSummaryStyle.Width(leftWidth - 2).Render(m.result.Description)
		headerLines += strings.Count(wrapped, "\n") + 1
		header.WriteString(wrapped + "\n")
	}
	header.WriteString("\n")
	headerLines++

	m.runeList.SetSize(leftWidth, availHeight-headerLines)
	leftContent := header.String() + m.runeList.View()

	var mid strings.Builder

	selectedRuneIdx := -1
	selectedIsPackage := false
	selectedIsRef := false

	if item, ok := m.runeList.SelectedItem().(runeListItem); ok {
		if item.isHeader {
			selectedIsPackage = true
			selectedRuneIdx = item.runeIdx
		} else if item.runeIdx >= 0 {
			selectedRuneIdx = item.runeIdx
			selectedIsRef = item.isRef
		}
	}

	if selectedRuneIdx >= 0 && m.result.RuneTree != nil {
		r := m.getRuneByIndex(selectedRuneIdx)
		if r == nil {
			mid.WriteString(renderPaneHeader("rune", midWidth, false) + "\n")
		} else if selectedIsPackage {
			mid.WriteString(m.renderPackageView(r, midWidth))
		} else {
			mid.WriteString(m.renderRuneView(r, selectedIsRef, midWidth))
		}
	} else {
		mid.WriteString(renderPaneHeader("rune", midWidth, false) + "\n")
	}

	m.midVP.Width = midWidth
	m.midVP.Height = availHeight
	m.midVP.SetContent(mid.String())

	sepChar := lipgloss.NewStyle().Foreground(border).Render("│")
	sepLines := make([]string, availHeight)
	for i := range sepLines {
		sepLines[i] = " " + sepChar + " "
	}

	layout := lipgloss.JoinHorizontal(lipgloss.Top,
		leftContent,
		strings.Join(sepLines, "\n"),
		m.midVP.View(),
	)

	var footer strings.Builder
	visibleRunes := m.countVisibleRunes()
	footer.WriteString(statusOk.Render(fmt.Sprintf("%d runes proposed", visibleRunes)))

	if m.groupByPackage {
		footer.WriteString(runeSigStyle.Render("  g: group by responsibility"))
	} else {
		footer.WriteString(runeSigStyle.Render("  g: group by package"))
	}

	return layout + "\n\n" + footer.String()
}

func (m *decomposeModel) renderPackageView(pkgRune *decompose.Rune, width int) string {
	var b strings.Builder

	pkgHdr := paneHeaderInactive.Render("── package ") + pkgStyle.Render(m.getPackageName(pkgRune.Path)) + " "
	if remaining := width - lipgloss.Width(pkgHdr); remaining > 0 {
		pkgHdr += paneHeaderInactive.Render(strings.Repeat("─", remaining))
	}
	b.WriteString(pkgHdr + "\n")

	children := m.getChildRunes(pkgRune)
	if len(children) > 0 {
		b.WriteString("\n")
		for _, child := range children {
			leaf := m.leafName(child.Path)
			title := runeNameStyle.Render(leaf)
			if child.Signature != "" {
				title += " " + runeSigStyle.Render(child.Signature)
			}
			b.WriteString(lipgloss.NewStyle().Width(width-2).Render("  "+title) + "\n")

			if child.Description != "" {
				b.WriteString(lipgloss.NewStyle().Foreground(dim).Width(width-2).PaddingLeft(4).Render(child.Description) + "\n")
			}
		}
	}

	var allAssumptions []string
	allAssumptions = append(allAssumptions, pkgRune.Assumptions...)
	for _, child := range children {
		allAssumptions = append(allAssumptions, child.Assumptions...)
	}

	if len(allAssumptions) > 0 {
		b.WriteString("\n" + runeSigStyle.Render("assumes:") + "\n")
		for _, a := range allAssumptions {
			b.WriteString(assumptionStyle.Render("? ") + lipgloss.NewStyle().Width(width-4).Render(a) + "\n")
		}
	}

	var allDeps []string
	allDeps = append(allDeps, pkgRune.Dependencies...)
	for _, child := range children {
		allDeps = append(allDeps, child.Dependencies...)
	}

	if len(allDeps) > 0 {
		b.WriteString("\n" + runeSigStyle.Render("dependencies:") + "\n")
		for _, dep := range allDeps {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6A9FD9")).Render("-> ") + lipgloss.NewStyle().Foreground(dim).Width(width-4).Render(dep) + "\n")
		}
	}

	return b.String()
}

func (m *decomposeModel) renderRuneView(rune *decompose.Rune, isRef bool, width int) string {
	var b strings.Builder

	hdrLabel := "── rune "
	nameStyle := runeNameStyle
	if isRef {
		hdrLabel = "── ref "
		nameStyle = treeRefStyle
	}

	runeHdr := paneHeaderInactive.Render(hdrLabel) + nameStyle.Render(rune.Path) + " "
	if remaining := width - lipgloss.Width(runeHdr); remaining > 0 {
		runeHdr += paneHeaderInactive.Render(strings.Repeat("─", remaining))
	}
	b.WriteString(runeHdr + "\n")

	if rune.Signature != "" {
		b.WriteString("\n" + runeSigStyle.Width(width-2).Render(rune.Signature))
	}

	if rune.Description != "" {
		if rune.Signature == "" {
			b.WriteString("\n")
		}
		wrapped := lipgloss.NewStyle().Foreground(dim).Width(width - 2).Render(rune.Description)
		b.WriteString(wrapped)
	}

	if len(rune.Tests) > 0 {
		b.WriteString("\n")
		for _, test := range rune.Tests {
			b.WriteString(testPassStyle.Render("+ ") + lipgloss.NewStyle().Width(width-4).Render(test.Name) + "\n")
		}
	}

	if len(rune.Assumptions) > 0 {
		b.WriteString("\n" + runeSigStyle.Render("assumes:") + "\n")
		for _, a := range rune.Assumptions {
			b.WriteString(assumptionStyle.Render("? ") + lipgloss.NewStyle().Width(width-4).Render(a) + "\n")
		}
	}

	if len(rune.Dependencies) > 0 {
		b.WriteString("\n" + runeSigStyle.Render("dependencies:") + "\n")
		for _, dep := range rune.Dependencies {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6A9FD9")).Render("-> ") + lipgloss.NewStyle().Foreground(dim).Width(width-4).Render(dep) + "\n")
		}
	}

	return b.String()
}

func (m *decomposeModel) getChildRunes(pkgRune *decompose.Rune) []*decompose.Rune {
	var children []*decompose.Rune
	pkgName := m.getPackageName(pkgRune.Path)
	if pkgName == "" {
		return children
	}

	for _, child := range pkgRune.Children {
		childPkg := m.getPackageName(child.Path)
		if childPkg == pkgName {
			children = append(children, child)
		} else if childPkg != "" {
			collectRunes(child, &children)
		}
	}

	return children
}
