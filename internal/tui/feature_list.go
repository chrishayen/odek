package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/chrishayen/odek/internal/feature"
)

type draftSelectedMsg struct{ draft Draft }
type featureSelectedMsg struct{ feature feature.Feature }
type newFeatureMsg struct{}

type featureListModel struct {
	drafts       []Draft
	features     []feature.Feature
	cursor       int
	width        int
	height       int
	draftStore   *DraftStore
	featureStore *feature.Store
	err          string
}

func newFeatureListModel(draftStore *DraftStore, featureStore *feature.Store, width, height int) featureListModel {
	m := featureListModel{
		draftStore:   draftStore,
		featureStore: featureStore,
		width:        width,
		height:       height,
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

	if m.featureStore != nil {
		features, err := m.featureStore.List()
		if err != nil {
			m.err = err.Error()
			return
		}
		m.features = features
	}

	m.err = ""
	total := m.totalItems()
	if m.cursor >= total {
		m.cursor = total - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *featureListModel) totalItems() int {
	return len(m.drafts) + len(m.features)
}

func (m *featureListModel) cursorOnDraft() bool {
	return m.cursor < len(m.drafts)
}

func (m *featureListModel) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		total := m.totalItems()
		switch msg.String() {
		case "j", "down":
			if m.cursor < total-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if total > 0 && m.cursor < total {
				if m.cursorOnDraft() {
					return func() tea.Msg {
						return draftSelectedMsg{draft: m.drafts[m.cursor]}
					}
				}
				idx := m.cursor - len(m.drafts)
				return func() tea.Msg {
					return featureSelectedMsg{feature: m.features[idx]}
				}
			}
		case "n":
			return func() tea.Msg { return newFeatureMsg{} }
		case "d", "x":
			if m.cursorOnDraft() && len(m.drafts) > 0 {
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

	if m.err != "" {
		b.WriteString(statusErr.Render(m.err) + "\n")
		return b.String()
	}

	// ── drafts section ──
	b.WriteString(renderPaneHeader("drafts", m.width, true) + "\n\n")

	if len(m.drafts) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(dim).Render("  No drafts yet.") + "\n")
	} else {
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
	}

	b.WriteString("\n")

	// ── features section ──
	b.WriteString(renderPaneHeader("features", m.width, true) + "\n\n")

	if len(m.features) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(dim).Render("  No features in registry.") + "\n")
	} else {
		nameWidth := 24
		statusWidth := 12
		versionWidth := 12

		for i, f := range m.features {
			globalIdx := len(m.drafts) + i

			name := f.Name
			if len(name) > nameWidth-2 {
				name = name[:nameWidth-3] + "~"
			}

			status := f.Status
			if status == "" {
				status = "unknown"
			}

			version := f.Version
			if version == "" {
				version = "-"
			}

			nameCol := lipgloss.NewStyle().Width(nameWidth).Render(name)
			statusCol := lipgloss.NewStyle().Width(statusWidth).Foreground(dim).Render(status)
			versionCol := lipgloss.NewStyle().Width(versionWidth).Foreground(dim).Render(version)

			row := nameCol + statusCol + versionCol

			if globalIdx == m.cursor {
				b.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FFFFFF")).
					Background(lipgloss.Color("#333333")).
					Width(m.width - 2).
					Render("  "+row) + "\n")
			} else {
				b.WriteString("  " + row + "\n")
			}
		}
	}

	return b.String()
}
