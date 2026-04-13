package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestDecompositionSmoke(t *testing.T) {
	d, groups := mockDecomposition()
	m := newDecompositionModel(d, groups)

	// Simulate a window size so View has reasonable dimensions.
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 140, Height: 44})
	m = updated.(decompositionModel)

	out := m.View().Content
	if out == "" {
		t.Fatal("empty view")
	}

	// Left pane should show each responsibility name.
	for _, g := range groups {
		if !strings.Contains(out, g.name) {
			t.Errorf("view missing responsibility %q", g.name)
		}
	}

	// Right pane should show the first selected rune's leaf name.
	if m.selectedRow < 0 {
		t.Fatal("no row selected on startup")
	}
	first := m.flatRunes[m.rowEntries[m.selectedRow].runeIdx]
	leaf := leafNameFromPath(first.Path)
	if !strings.Contains(out, leaf) {
		t.Errorf("view missing first selected rune %q", leaf)
	}

	// Navigation: press down a few times, ensure selection advances to another selectable row.
	start := m.selectedRow
	for range 3 {
		next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		m = next.(decompositionModel)
	}
	if m.selectedRow == start {
		t.Error("down key did not advance selection")
	}
	if !m.rowEntries[m.selectedRow].selectable {
		t.Error("advanced selection landed on non-selectable row")
	}

	// Re-render after navigation to catch late panics.
	_ = m.View()
}
