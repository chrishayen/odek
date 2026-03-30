package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chrishayen/valkyrie/internal/feature"
)

var (
	inputLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Bold(true)

	statusOk = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66CC66"))

	statusErr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CC6666"))

	featureRow = lipgloss.NewStyle().
			Foreground(dim)

	sectionTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Bold(true)
)

type createFeatureModel struct {
	descInput textarea.Model
	features  []feature.Feature
	status    string
	statusErr bool
	store     *feature.Store
	width     int
}

func newCreateFeatureModel(store *feature.Store, width, height int) createFeatureModel {
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
	// Reserve space for: header(3) + label(2) + features(~8) + help(2) + padding
	taHeight := height - 20
	if taHeight < 3 {
		taHeight = 3
	}
	ta.SetWidth(inputWidth)
	ta.SetHeight(taHeight)
	ta.CharLimit = 2000

	m := createFeatureModel{
		descInput: ta,
		store:     store,
		width:     width,
	}
	m.loadFeatures()
	return m
}

func (m *createFeatureModel) loadFeatures() {
	if m.store == nil {
		return
	}
	features, err := m.store.List()
	if err != nil {
		return
	}
	m.features = features
}

func (m *createFeatureModel) submit() {
	desc := strings.TrimSpace(m.descInput.Value())
	if desc == "" {
		m.status = "description is required"
		m.statusErr = true
		return
	}

	name := fmt.Sprintf("feature_%d", time.Now().Unix())

	if err := m.store.Create(name, desc); err != nil {
		m.status = err.Error()
		m.statusErr = true
		return
	}

	m.status = fmt.Sprintf("created %q", name)
	m.statusErr = false
	m.descInput.Reset()
	m.descInput.Focus()
	m.loadFeatures()
}

func (m *createFeatureModel) resize(width, height int) {
	m.width = width
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

func (m *createFeatureModel) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.submit()
			return nil
		}
	}

	var cmd tea.Cmd
	m.descInput, cmd = m.descInput.Update(msg)
	return cmd
}

func (m *createFeatureModel) view(width int) string {
	var b strings.Builder

	b.WriteString(inputLabel.Render("Describe your feature") + "\n\n")
	b.WriteString(m.descInput.View())

	if m.status != "" {
		b.WriteString("\n\n")
		if m.statusErr {
			b.WriteString(statusErr.Render(m.status))
		} else {
			b.WriteString(statusOk.Render(m.status))
		}
	}

	return b.String()
}
