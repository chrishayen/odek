package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/chrishayen/odek/internal/feature"
)

type draftSelectedMsg struct{ draft Draft }
type featureSelectedMsg struct{ feature feature.Feature }
type newFeatureMsg struct{}

// listItem wraps either a Draft or Feature for bubbles/list.
type listItem struct {
	draft   *Draft
	feature *feature.Feature
}

func (i listItem) Title() string {
	if i.draft != nil {
		if i.draft.FeatureName != "" {
			return i.draft.FeatureName
		}
		return "untitled"
	}
	return i.feature.Name
}

func (i listItem) Description() string {
	if i.draft != nil {
		desc := i.draft.Summary
		if desc == "" && i.draft.Requirement != "" {
			desc = i.draft.Requirement
		}
		return desc + "  " + timeAgo(i.draft.UpdatedAt)
	}
	status := i.feature.Status
	if status == "" {
		status = "unknown"
	}
	version := i.feature.Version
	if version == "" {
		version = "-"
	}
	return status + "  " + version
}

func (i listItem) FilterValue() string { return i.Title() }

func (i listItem) isDraft() bool { return i.draft != nil }

// featureListModel wraps bubbles/list.
type featureListModel struct {
	list         list.Model
	draftStore   *DraftStore
	featureStore *feature.Store
	drafts       []Draft
	features     []feature.Feature
	err          string
}

func newFeatureListModel(draftStore *DraftStore, featureStore *feature.Store, width, height int) featureListModel {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F5A623")).
		Padding(0, 0, 0, 2)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2)
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#F5A623")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#F5A623")).
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 1)

	l := list.New(nil, delegate, width, height)
	l.Title = "features"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F5A623")).
		Bold(true).
		Padding(0, 1)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.KeyMap.Quit.Unbind()

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		}
	}

	m := featureListModel{
		list:         l,
		draftStore:   draftStore,
		featureStore: featureStore,
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
	items := make([]list.Item, 0, len(m.drafts)+len(m.features))
	for i := range m.drafts {
		items = append(items, listItem{draft: &m.drafts[i]})
	}
	for i := range m.features {
		items = append(items, listItem{feature: &m.features[i]})
	}
	m.list.SetItems(items)
}

func (m *featureListModel) totalItems() int {
	return len(m.drafts) + len(m.features)
}

func (m *featureListModel) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item, ok := m.list.SelectedItem().(listItem); ok {
				if item.isDraft() {
					return func() tea.Msg { return draftSelectedMsg{draft: *item.draft} }
				}
				return func() tea.Msg { return featureSelectedMsg{feature: *item.feature} }
			}
		case "n":
			return func() tea.Msg { return newFeatureMsg{} }
		case "d", "x":
			if item, ok := m.list.SelectedItem().(listItem); ok && item.isDraft() {
				_ = m.draftStore.Delete(item.draft.ID)
				m.reload()
				return nil
			}
		case "backspace":
			return func() tea.Msg { return goBackMsg{} }
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *featureListModel) view() string {
	if m.err != "" {
		return statusErr.Render(m.err) + "\n"
	}
	return m.list.View()
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
