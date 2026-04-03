package tui

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/chrishayen/odek/internal/draft"
	"github.com/chrishayen/odek/internal/feature"
)

type draftSelectedMsg struct{ draft draft.Draft }
type featureSelectedMsg struct{ feature feature.Feature }
type newFeatureMsg struct{}

// listItem wraps either a Draft or Feature for bubbles/list.
type listItem struct {
	draft   *draft.Draft
	feature *feature.Feature
}

func (i listItem) FilterValue() string {
	if i.draft != nil {
		return i.draft.FeatureName
	}
	return i.feature.Name
}

func (i listItem) isDraft() bool { return i.draft != nil }

// featureDelegate renders items in the feature list with color-coded drafts and features.
type featureDelegate struct{}

func (d featureDelegate) Height() int                             { return 2 }
func (d featureDelegate) Spacing() int                            { return 1 }
func (d featureDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

var (
	fdNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5A623")).
			Bold(true)

	fdIDStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))

	fdDraftTagStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D4A843"))

	fdDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))

	fdTimeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))

	fdStatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))

	fdSelectedBorder = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color("#F5A623"))
)

func (d featureDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	li, ok := item.(listItem)
	if !ok {
		return
	}

	selected := index == m.Index()
	availWidth := m.Width()

	var title, desc string

	if li.isDraft() {
		name := li.draft.FeatureName
		if name == "" {
			name = "untitled"
		}
		if selected {
			title = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render(name) +
				"  " + fdIDStyle.Render(li.draft.ID)
		} else {
			title = fdNameStyle.Render(name) + "  " + fdIDStyle.Render(li.draft.ID)
		}

		desc = fdDraftTagStyle.Render("draft")
		if li.draft.Summary != "" {
			desc += "  " + fdDescStyle.Render(li.draft.Summary)
		} else if li.draft.Requirement != "" {
			desc += "  " + fdDescStyle.Render(li.draft.Requirement)
		}
		desc += "  " + fdTimeStyle.Render(timeAgo(li.draft.UpdatedAt))
	} else {
		f := li.feature
		if selected {
			title = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render(f.Name)
		} else {
			title = fdNameStyle.Render(f.Name)
		}

		status := f.Status
		if status == "" {
			status = "unknown"
		}
		version := f.Version
		if version == "" {
			version = "-"
		}
		desc = fdStatusStyle.Render(status + "  " + version)
	}

	// Truncate if needed
	if lipgloss.Width(title) > availWidth-4 {
		title = title[:availWidth-5] + "~"
	}
	if lipgloss.Width(desc) > availWidth-4 {
		desc = desc[:availWidth-5] + "~"
	}

	if selected {
		block := fdSelectedBorder.Width(availWidth).
			Padding(0, 0, 0, 1).
			Render(title + "\n" + desc)
		fmt.Fprint(w, block)
	} else {
		pad := lipgloss.NewStyle().Padding(0, 0, 0, 2)
		fmt.Fprint(w, pad.Render(title)+"\n"+pad.Render(desc))
	}
}

// featureListModel wraps bubbles/list.
type featureListModel struct {
	list         list.Model
	draftStore   *draft.Store
	featureStore *feature.Store
	drafts       []draft.Draft
	features     []feature.Feature
	err          string
}

func newFeatureListModel(draftStore *draft.Store, featureStore *feature.Store, width, height int) featureListModel {
	l := list.New(nil, featureDelegate{}, width, height)
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
