package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type draftSelectedMsg struct{ draft Draft }
type newFeatureMsg struct{}

type featureListModel struct {
	drafts     []Draft
	cursor     int
	width      int
	height     int
	draftStore *DraftStore
	err        string
}

func newFeatureListModel(store *DraftStore, width, height int) featureListModel {
	m := featureListModel{
		draftStore: store,
		width:      width,
		height:     height,
	}
	m.reload()
	return m
}

func (m *featureListModel) reload() {
	drafts, err := m.draftStore.List()
	if err != nil {
		m.err = err.Error()
		return
	}
	m.drafts = drafts
	m.err = ""
	if m.cursor >= len(m.drafts) {
		m.cursor = len(m.drafts) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *featureListModel) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < len(m.drafts)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if len(m.drafts) > 0 && m.cursor < len(m.drafts) {
				return func() tea.Msg {
					return draftSelectedMsg{draft: m.drafts[m.cursor]}
				}
			}
		case "n":
			return func() tea.Msg { return newFeatureMsg{} }
		case "d", "x":
			if len(m.drafts) > 0 && m.cursor < len(m.drafts) {
				_ = m.draftStore.Delete(m.drafts[m.cursor].ID)
				m.reload()
			}
		case "backspace":
			return func() tea.Msg { return goBackMsg{} }
		}
	}
	return nil
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1h ago"
		}
		return fmt.Sprintf("%dh ago", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	}
}

func (m *featureListModel) view() string {
	var b strings.Builder

	b.WriteString(renderPaneHeader("drafts", m.width, true) + "\n\n")

	if m.err != "" {
		b.WriteString(statusErr.Render(m.err) + "\n")
		return b.String()
	}

	if len(m.drafts) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(dim).Render("  No drafts yet. Press n to create a new feature.") + "\n")
		return b.String()
	}

	nameWidth := 24
	ageWidth := 10
	summaryWidth := m.width - nameWidth - ageWidth - 8
	if summaryWidth < 20 {
		summaryWidth = 20
	}

	for i, d := range m.drafts {
		name := d.FeatureName
		if name == "" {
			name = "untitled"
		}
		if len(name) > nameWidth-2 {
			name = name[:nameWidth-3] + "~"
		}

		summary := d.Summary
		if summary == "" && d.Requirement != "" {
			summary = strings.ReplaceAll(d.Requirement, "\n", " ")
		}
		if len(summary) > summaryWidth {
			summary = summary[:summaryWidth-3] + "..."
		}

		age := timeAgo(d.UpdatedAt)

		nameCol := lipgloss.NewStyle().Width(nameWidth).Render(name)
		summaryCol := lipgloss.NewStyle().Width(summaryWidth).Foreground(dim).Render(summary)
		ageCol := lipgloss.NewStyle().Width(ageWidth).Foreground(dim).Render(age)

		row := nameCol + summaryCol + ageCol

		if i == m.cursor {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#333333")).
				Width(m.width - 2).
				Render("  "+row) + "\n")
		} else {
			b.WriteString("  " + row + "\n")
		}
	}

	return b.String()
}
