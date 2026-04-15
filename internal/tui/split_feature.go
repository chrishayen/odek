package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

const (
	splitPaneMinWidth = 150
	splitSepCols      = 4
	splitMinSubW      = 20
	splitResizeStep   = 4
)

type splitPaneModel struct {
	left   createFeatureModel
	right  featureDecompModel
	width  int
	height int
	focus  int
	leftW  int
}

func newSplitPaneModel(left createFeatureModel, right featureDecompModel, width, height int) splitPaneModel {
	m := splitPaneModel{
		left:   left,
		right:  right,
		width:  width,
		height: height,
	}
	m.left.SetInSplit(true)
	m.right.SetInSplit(true)
	m.right.SetActive(m.focus == 1)
	m.resize(width, height)
	return m
}

func (m *splitPaneModel) resize(w, h int) {
	oldW := m.width
	m.width = w
	m.height = h
	avail := max(w-splitSepCols, 2)
	switch {
	case m.leftW <= 0:
		// First sizing: use the default 1/3 ratio.
		m.leftW, _ = splitWidths(w)
	case oldW > 0 && oldW != w:
		// Window resized: scale the current left width proportionally so the
		// user's manual adjustment is preserved as a ratio.
		oldAvail := max(oldW-splitSepCols, 2)
		m.leftW = m.leftW * avail / oldAvail
	}
	if m.leftW < splitMinSubW {
		m.leftW = splitMinSubW
	}
	if max := avail - splitMinSubW; m.leftW > max {
		m.leftW = max
	}
	rightW := avail - m.leftW
	m.left.resize(m.leftW, h)
	m.right.resize(rightW, h)
}

func splitWidths(total int) (left, right int) {
	avail := max(total-splitSepCols, 2)
	left = avail / 3
	right = avail - left
	return
}

func (m splitPaneModel) Init() tea.Cmd {
	return tea.Batch(m.left.Init(), m.right.Init())
}

func (m splitPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width < splitPaneMinWidth {
			if m.focus == 0 {
				m.left.resize(msg.Width, msg.Height)
				m.left.SetInSplit(false)
				return m.left, nil
			}
			m.right.resize(msg.Width, msg.Height)
			m.right.SetActive(true)
			m.right.SetInSplit(false)
			return m.right, nil
		}
		m.resize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			m.focus = 1 - m.focus
			if m.focus == 0 {
				m.right.SetActive(false)
				return m, m.left.Focus()
			}
			m.left.Blur()
			m.right.SetActive(true)
			return m, nil
		case "ctrl+]":
			m.leftW += splitResizeStep
			m.resize(m.width, m.height)
			return m, nil
		case "ctrl+[":
			m.leftW -= splitResizeStep
			m.resize(m.width, m.height)
			return m, nil
		case "ctrl+enter", "ctrl+s":
			return m, nil
		case "enter":
			if m.focus == 1 && !m.right.inputActive && m.right.focusedCol == 0 {
				snap := m.right.snap
				if m.right.selectedIdx >= 0 && m.right.selectedIdx < len(snap.TopLevelNames) {
					name := snap.TopLevelNames[m.right.selectedIdx]
					fqn := qualifiedPath(snap.PackageName, name) + ": "
					m.left.SetChatInput(fqn)
					m.focus = 0
					m.right.SetActive(false)
					return m, m.left.Focus()
				}
			}
		}
		if m.focus == 0 {
			next, cmd := m.left.Update(msg)
			if updated, ok := next.(createFeatureModel); ok {
				m.left = updated
				return m, cmd
			}
			return next, tea.Batch(cmd, resizeCmd(m.width, m.height))
		}
		next, cmd := m.right.Update(msg)
		if updated, ok := next.(featureDecompModel); ok {
			m.right = updated
			return m, cmd
		}
		return next, tea.Batch(cmd, resizeCmd(m.width, m.height))
	}

	var cmds []tea.Cmd
	nextLeft, cmdLeft := m.left.Update(msg)
	if u, ok := nextLeft.(createFeatureModel); ok {
		m.left = u
	}
	cmds = append(cmds, cmdLeft)

	nextRight, cmdRight := m.right.Update(msg)
	if u, ok := nextRight.(featureDecompModel); ok {
		m.right = u
	}
	cmds = append(cmds, cmdRight)

	return m, tea.Batch(cmds...)
}

func (m splitPaneModel) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true
	v.BackgroundColor = bgMain

	if m.width <= 0 || m.height <= 0 {
		return v
	}

	avail := max(m.width-splitSepCols, 2)
	leftW := m.leftW
	rightW := avail - leftW
	leftView := clipToBox(m.left.View().Content, leftW, m.height)
	rightView := clipToBox(m.right.View().Content, rightW, m.height)
	if m.focus == 0 {
		rightView = dimUnfocused(rightView, rightW)
	} else {
		leftView = dimUnfocused(leftView, leftW)
	}
	sep := renderSplitSeparator(m.height)

	v.Content = lipgloss.JoinHorizontal(lipgloss.Top, leftView, sep, rightView)
	return v
}

// dimUnfocused flattens the first row (ODEK logo) and the help-bar row
// (second from the bottom — the very bottom is a blank spacer) of a clipped
// pane into plain fgDim text on bgMain. Used on the pane that doesn't have
// focus — the contrast against the focused pane's accent-colored logo and
// styled help bar becomes the focus indicator.
func dimUnfocused(content string, width int) string {
	if width <= 0 {
		return content
	}
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return content
	}
	bgPad := lipgloss.NewStyle().Background(bgMain)
	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("237")).
		Background(bgMain).
		Faint(true)
	dim := func(row string) string {
		plain := ansi.Strip(row)
		styled := dimStyle.Render(plain)
		w := ansi.StringWidth(styled)
		switch {
		case w > width:
			return ansi.Truncate(styled, width, "")
		case w < width:
			return styled + bgPad.Render(strings.Repeat(" ", width-w))
		}
		return styled
	}
	lines[0] = dim(lines[0])
	if len(lines) >= 2 {
		lines[len(lines)-2] = dim(lines[len(lines)-2])
	}
	return strings.Join(lines, "\n")
}

// clipToBox forces `content` into an exact `width` × `height` rectangle,
// truncating over-wide rows via ansi.Truncate (ANSI-aware) and padding short
// rows or missing rows with bgMain spaces. Needed because the child panes'
// help bars overflow their set width; the split enforces the boundary.
func clipToBox(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	bgStyle := lipgloss.NewStyle().Background(bgMain)
	blank := bgStyle.Render(strings.Repeat(" ", width))

	lines := strings.Split(content, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}

	out := make([]string, height)
	for i := range height {
		if i >= len(lines) {
			out[i] = blank
			continue
		}
		line := lines[i]
		w := ansi.StringWidth(line)
		switch {
		case w == width:
			out[i] = line
		case w > width:
			out[i] = ansi.Truncate(line, width, "")
		default:
			out[i] = line + bgStyle.Render(strings.Repeat(" ", width-w))
		}
	}
	return strings.Join(out, "\n")
}

// resizeCmd emits a synthetic WindowSizeMsg carrying the full terminal
// dimensions. Used when the split pane hands off to a sibling model (e.g.
// a transition back to the landing page) so the new model replaces the
// stale per-pane width with the real terminal size on its next Update.
func resizeCmd(w, h int) tea.Cmd {
	return func() tea.Msg { return tea.WindowSizeMsg{Width: w, Height: h} }
}

func renderSplitSeparator(height int) string {
	barStyle := lipgloss.NewStyle().Foreground(mockSep).Background(bgMain)
	bgPad := lipgloss.NewStyle().Background(bgMain).Render(" ")
	// Leading: 1 col bgPad. Trailing: 2 cols bgPad. This matches the left
	// pane's viewPadX right-padding (1 col) + 1 col of separator bgPad on
	// the left side of "│", so the right side has symmetric total padding.
	row := bgPad + barStyle.Render("│") + bgPad + bgPad

	rows := make([]string, height)
	for i := range rows {
		rows[i] = row
	}
	return strings.Join(rows, "\n")
}
