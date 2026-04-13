package tui

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"

	"shotgun.dev/odek/decompose"
)

type responsibilityGroup struct {
	name        string
	description string
	runePaths   []string
}

type rowKind int

const (
	rowHeader rowKind = iota
	rowRune
	rowSpacer
)

type leftRow struct {
	kind       rowKind
	group      *responsibilityGroup
	runeIdx    int
	selectable bool
}

type decompositionModel struct {
	width, height int
	decomposition *decompose.Decomposition
	groups        []responsibilityGroup
	flatRunes     []*decompose.Rune
	pathToIdx     map[string]int
	rowEntries    []leftRow
	selectedRow   int
	help          help.Model
}

var (
	accent     = lipgloss.Color("212")     // pink
	accentSoft = lipgloss.Color("99")      // purple
	bgMain     = lipgloss.Color("#171717") // matches lipgloss gallery terminal bg
	fgBright   = lipgloss.Color("15")
	fgBody     = lipgloss.Color("245")
	fgDim      = lipgloss.Color("241")
)

var (
	accentStyle    = lipgloss.NewStyle().Foreground(accent).Background(bgMain).Bold(true)
	breadcrumbSep  = lipgloss.NewStyle().Foreground(fgDim).Background(bgMain).Render(" › ")
	breadcrumbLogo = accentStyle
	breadcrumbText = lipgloss.NewStyle().Foreground(fgBody).Background(bgMain)

	selectedMarker = accentStyle.Render("❯ ")
	selectedName   = lipgloss.NewStyle().Foreground(fgBright).Background(bgMain).Bold(true)
	unselected     = lipgloss.NewStyle().Foreground(fgBody).Background(bgMain)

	statusLabelStyle = lipgloss.NewStyle().
				Foreground(fgBright).
				Background(accent).
				Bold(true).
				Padding(0, 1)
	statusCountStyle = lipgloss.NewStyle().
				Foreground(fgBright).
				Background(accentSoft).
				Bold(true).
				Padding(0, 1)
	statusMidStyle = lipgloss.NewStyle().
			Foreground(fgBody).
			Background(bgMain).
			Padding(0, 1)
	statusHintStyle = lipgloss.NewStyle().
			Foreground(fgDim).
			Background(bgMain).
			Padding(0, 1)
	statusSpacer = lipgloss.NewStyle().Background(bgMain)
)

func newDecompositionModel(d *decompose.Decomposition, groups []responsibilityGroup) decompositionModel {
	var flat []*decompose.Rune
	flattenRunes(d.RuneTree, &flat)

	pathToIdx := make(map[string]int, len(flat))
	for i, r := range flat {
		pathToIdx[r.Path] = i
	}

	m := decompositionModel{
		width:         120,
		height:        40,
		decomposition: d,
		groups:        groups,
		flatRunes:     flat,
		pathToIdx:     pathToIdx,
		help:          newHelpModel(),
	}
	m.buildRowEntries()
	m.selectedRow = m.firstSelectableRow()
	return m
}

func flattenRunes(r *decompose.Rune, out *[]*decompose.Rune) {
	if r == nil {
		return
	}
	*out = append(*out, r)
	for _, c := range r.Children {
		flattenRunes(c, out)
	}
}

func (m *decompositionModel) buildRowEntries() {
	m.rowEntries = nil
	for gi := range m.groups {
		g := &m.groups[gi]
		if gi > 0 {
			m.rowEntries = append(m.rowEntries, leftRow{kind: rowSpacer})
		}
		m.rowEntries = append(m.rowEntries, leftRow{kind: rowHeader, group: g})
		for _, path := range g.runePaths {
			idx, ok := m.pathToIdx[path]
			if !ok {
				continue
			}
			m.rowEntries = append(m.rowEntries, leftRow{kind: rowRune, runeIdx: idx, selectable: true})
		}
	}
}

func (m decompositionModel) firstSelectableRow() int {
	for i, row := range m.rowEntries {
		if row.selectable {
			return i
		}
	}
	return -1
}

func (m decompositionModel) nextSelectableRow(from, dir int) int {
	n := len(m.rowEntries)
	if n == 0 {
		return -1
	}
	i := from
	for range n {
		i = (i + dir + n) % n
		if m.rowEntries[i].selectable {
			return i
		}
	}
	return from
}

func RunDecomposition() {
	d, groups := mockDecomposition()
	m := newDecompositionModel(d, groups)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m decompositionModel) Init() tea.Cmd { return nil }

func (m decompositionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc", "q", "backspace":
			return model{width: m.width, height: m.height, help: newHelpModel()}, nil
		case "up", "k":
			if m.selectedRow >= 0 {
				m.selectedRow = m.nextSelectableRow(m.selectedRow, -1)
			}
			return m, nil
		case "down", "j":
			if m.selectedRow >= 0 {
				m.selectedRow = m.nextSelectableRow(m.selectedRow, 1)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m decompositionModel) View() tea.View {
	innerWidth := max(m.width-viewPadX*2, 72)
	innerHeight := max(m.height, 24)

	header := m.renderBreadcrumb()
	footer := m.renderStatusBar(innerWidth)

	// body area: total - header(1) - blankBelowHeader(1) - blankAboveFooter(1) - footer(1)
	bodyHeight := max(innerHeight-4, 12)

	leftWidth := innerWidth / 4
	leftWidth = max(leftWidth, 28)
	leftWidth = min(leftWidth, 36)
	rightWidth := innerWidth - leftWidth - 1

	leftBox := m.renderLeftPane(leftWidth, bodyHeight)
	rightCol := m.renderRightPane(rightWidth, bodyHeight)

	gap := lipgloss.NewStyle().Background(bgMain).Render(" ")
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, gap, rightCol)

	content := header + "\n\n" + body + "\n" + footer
	rendered := lipgloss.NewStyle().
		Background(bgMain).
		Padding(0, viewPadX).
		Width(m.width).
		Height(m.height).
		Render(content)
	v := tea.NewView(rendered)
	v.AltScreen = true
	v.BackgroundColor = bgMain
	return v
}

func (m decompositionModel) renderBreadcrumb() string {
	name := ""
	if m.decomposition != nil {
		name = m.decomposition.FeatureName
	}
	return " " + breadcrumbLogo.Render(logoSmall) + breadcrumbSep + breadcrumbText.Render(name)
}

func (m decompositionModel) renderLeftPane(width, height int) string {
	innerW := max(width-4, 1)
	innerH := max(height-2, 0)

	type line struct {
		text     string
		selected bool
	}
	var lines []line

	descStyle := lipgloss.NewStyle().Foreground(fgDim).Italic(true).Width(innerW)

	for i, row := range m.rowEntries {
		switch row.kind {
		case rowSpacer:
			lines = append(lines, line{})
		case rowHeader:
			lines = append(lines, line{text: accentStyle.Render(row.group.name)})
			if row.group.description != "" {
				for _, dl := range strings.Split(descStyle.Render(row.group.description), "\n") {
					lines = append(lines, line{text: dl})
				}
			}
		case rowRune:
			r := m.flatRunes[row.runeIdx]
			leaf := leafNameFromPath(r.Path)
			isSel := i == m.selectedRow
			var txt string
			if isSel {
				txt = "  " + selectedMarker + selectedName.Render(leaf)
			} else {
				txt = "    " + unselected.Render(leaf)
			}
			lines = append(lines, line{text: txt, selected: isSel})
		}
	}

	offset := 0
	if len(lines) > innerH && innerH > 0 {
		selIdx := -1
		for i, l := range lines {
			if l.selected {
				selIdx = i
				break
			}
		}
		if selIdx >= 0 {
			offset = selIdx - innerH/2
			offset = max(offset, 0)
			offset = min(offset, len(lines)-innerH)
		}
	}
	end := min(offset+innerH, len(lines))
	visible := lines[offset:end]

	var b strings.Builder
	for i, l := range visible {
		b.WriteString(l.text)
		if i < len(visible)-1 {
			b.WriteString("\n")
		}
	}

	return titledBox("responsibilities", b.String(), width, height, fgDim, accentSoft)
}

func (m decompositionModel) renderRightPane(width, height int) string {
	sel := m.selectedRune()

	runeBoxH := max(min(height/3, 12), 8)
	runeBox := m.renderRuneBox(sel, width, runeBoxH)

	colsH := max(height-runeBoxH-1, 6)
	cols := m.renderFeatureColumns(width, colsH, sel)

	return runeBox + "\n" + cols
}

func (m decompositionModel) renderRuneBox(r *decompose.Rune, width, height int) string {
	if r == nil {
		return titledBox("selected", "", width, height, fgDim, accentSoft)
	}
	innerW := max(width-4, 1)

	var b strings.Builder
	b.WriteString(accentStyle.Render(leafNameFromPath(r.Path)) + "\n")
	if r.Signature != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(fgDim).Render(r.Signature) + "\n")
	}
	if r.Description != "" {
		b.WriteString("\n")
		wrapped := lipgloss.NewStyle().
			Foreground(fgBody).
			Width(innerW).
			Render(r.Description)
		b.WriteString(wrapped)
	}

	return titledBox("selected", b.String(), width, height, accent, accent)
}

type branchNode struct {
	leaf     string
	selected bool
}

func (b branchNode) String() string {
	if b.selected {
		return lipgloss.NewStyle().Foreground(accent).Bold(true).Render("❯ " + b.leaf)
	}
	return lipgloss.NewStyle().Foreground(fgBody).Render(b.leaf)
}

func buildBranchTree(r *decompose.Rune, selectedPath string) *tree.Tree {
	t := tree.Root(branchNode{
		leaf:     leafNameFromPath(r.Path),
		selected: r.Path == selectedPath,
	})
	for _, c := range r.Children {
		if len(c.Children) == 0 {
			t.Child(branchNode{
				leaf:     leafNameFromPath(c.Path),
				selected: c.Path == selectedPath,
			})
		} else {
			t.Child(buildBranchTree(c, selectedPath))
		}
	}
	return t
}

func (m decompositionModel) renderFeatureColumns(width, height int, sel *decompose.Rune) string {
	root := m.decomposition.RuneTree
	if root == nil || len(root.Children) == 0 {
		return titledBox("composition", "", width, height, fgDim, accentSoft)
	}

	selectedPath := ""
	if sel != nil {
		selectedPath = sel.Path
	}

	branches := root.Children
	n := len(branches)
	gap := 1
	totalGap := gap * (n - 1)
	colW := (width - totalGap) / n
	remainder := (width - totalGap) - colW*n

	enumStyle := lipgloss.NewStyle().Foreground(fgDim)
	rootStyleDim := lipgloss.NewStyle().Foreground(accentSoft).Bold(true)
	rootStyleSel := lipgloss.NewStyle().Foreground(accent).Bold(true)

	var boxes []string
	for i, branch := range branches {
		w := colW
		if i < remainder {
			w++
		}

		hasSelection := containsPath(branch, selectedPath)

		t := buildBranchTree(branch, selectedPath)
		t.EnumeratorStyle(enumStyle)
		if hasSelection {
			t.RootStyle(rootStyleSel)
		} else {
			t.RootStyle(rootStyleDim)
		}

		borderCol := fgDim
		if hasSelection {
			borderCol = accent
		}

		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderCol).
			BorderBackground(bgMain).
			Background(bgMain).
			Padding(0, 1).
			Width(w - 4).
			Height(height - 2).
			MaxHeight(height).
			Render(t.String())
		boxes = append(boxes, box)
	}

	pieces := make([]string, 0, n*2-1)
	for i, box := range boxes {
		if i > 0 {
			pieces = append(pieces, " ")
		}
		pieces = append(pieces, box)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, pieces...)
}

func containsPath(n *decompose.Rune, path string) bool {
	if n == nil || path == "" {
		return false
	}
	if n.Path == path {
		return true
	}
	for _, c := range n.Children {
		if containsPath(c, path) {
			return true
		}
	}
	return false
}

func (m decompositionModel) renderStatusBar(width int) string {
	statusLabel := statusLabelStyle.Render("STATUS")

	featureName := ""
	if m.decomposition != nil {
		featureName = m.decomposition.FeatureName
	}
	midText := statusMidStyle.Render(featureName)

	hints := statusHintStyle.Render("↑/↓ navigate  esc back  q quit")
	count := statusCountStyle.Render(fmt.Sprintf("%d runes", len(m.flatRunes)))

	used := lipgloss.Width(statusLabel) + lipgloss.Width(midText) + lipgloss.Width(hints) + lipgloss.Width(count)
	fill := max(width-used, 0)
	spacer := statusSpacer.Render(strings.Repeat(" ", fill))

	return statusLabel + midText + spacer + hints + count
}

func (m decompositionModel) selectedRune() *decompose.Rune {
	if m.selectedRow < 0 || m.selectedRow >= len(m.rowEntries) {
		return nil
	}
	row := m.rowEntries[m.selectedRow]
	if !row.selectable {
		return nil
	}
	return m.flatRunes[row.runeIdx]
}

func leafNameFromPath(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}

// titledBox renders a rounded-border box of exactly width × height chars.
// If title is non-empty it's rendered as an accent header at the top of the
// content area, followed by a blank line. Content is wrapped to the inner
// width, truncated to fit the inner height, and padded with blank lines.
func titledBox(title, content string, width, height int, borderCol, titleCol color.Color) string {
	width = max(width, 4)
	height = max(height, 3)

	innerW := width - 4
	innerH := height - 2

	var inner string
	if title != "" {
		titleLine := lipgloss.NewStyle().Foreground(titleCol).Bold(true).Render(title)
		if content == "" {
			inner = titleLine
		} else {
			inner = titleLine + "\n\n" + content
		}
	} else {
		inner = content
	}

	pre := lipgloss.NewStyle().Width(innerW).Render(inner)
	lines := strings.Split(pre, "\n")

	truncated := false
	if len(lines) > innerH {
		lines = lines[:innerH]
		truncated = true
	}
	for len(lines) < innerH {
		lines = append(lines, strings.Repeat(" ", innerW))
	}
	if truncated && innerH > 0 {
		marker := lipgloss.NewStyle().Foreground(borderCol).Render("…")
		lines[innerH-1] = marker + strings.Repeat(" ", max(innerW-1, 0))
	}
	fit := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderCol).
		BorderBackground(bgMain).
		Background(bgMain).
		Padding(0, 1).
		Render(fit)
}

var keyNav = key.NewBinding(
	key.WithKeys("up", "down", "j", "k"),
	key.WithHelp("↑/↓", "navigate"),
)

type decompositionKeyMap struct{}

func (decompositionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keyNav, keyCancel, keyQuit}
}
func (decompositionKeyMap) FullHelp() [][]key.Binding { return nil }
