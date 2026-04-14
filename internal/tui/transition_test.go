package tui

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// stubModel is a test-only tea.Model that renders a fixed multi-line block.
// Its block is the same string repeated across m.height rows, clipped to
// m.width cells per row. Used to stand in for from/to during transition tests
// so we don't depend on createFeatureModel wiring.
type stubModel struct {
	width  int
	height int
	fill   string
}

func (m stubModel) Init() tea.Cmd { return nil }

func (m stubModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sz, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = sz.Width
		m.height = sz.Height
	}
	return m, nil
}

func (m stubModel) View() tea.View {
	if m.width <= 0 || m.height <= 0 {
		return tea.NewView("")
	}
	row := ansi.Truncate(strings.Repeat(m.fill, m.width), m.width, "")
	if ansi.StringWidth(row) < m.width {
		row += strings.Repeat(" ", m.width-ansi.StringWidth(row))
	}
	rows := make([]string, m.height)
	for i := range m.height {
		rows[i] = row
	}
	v := tea.NewView(strings.Join(rows, "\n"))
	v.BackgroundColor = bgMain
	return v
}

func TestTransitionSmoke(t *testing.T) {
	const (
		w = 80
		h = 24
	)
	from := stubModel{width: w, height: h, fill: "F"}
	to := stubModel{width: w, height: h, fill: "T"}
	pin := renderFeaturePin()

	tr := newTransition(from, to, w, h, pin)

	if tr.fromSnap == "" || tr.toSnap == "" {
		t.Fatal("snapshots not captured at construction")
	}

	// Fresh transition View should render without panic at t≈0.
	view0 := tr.View()
	if view0.Content == "" {
		t.Fatal("View() returned empty content at t=0")
	}
	lines := strings.Split(view0.Content, "\n")
	if len(lines) != h {
		t.Errorf("expected %d rows at t=0, got %d", h, len(lines))
	}
	for i, line := range lines {
		if got := ansi.StringWidth(line); got != w {
			t.Errorf("row %d width = %d, want %d", i, got, w)
			break
		}
	}
	// At t=0 the outgoing fills the entire frame, so every row must contain
	// the 'F' fill character from fromSnap.
	if !strings.Contains(view0.Content, "F") {
		t.Error("t=0 view missing outgoing content")
	}

	// Trigger WindowSizeMsg handling — width/height/snapshots should update.
	const (
		w2 = 100
		h2 = 30
	)
	updated, _ := tr.Update(tea.WindowSizeMsg{Width: w2, Height: h2})
	tr2, ok := updated.(transitionModel)
	if !ok {
		t.Fatalf("WindowSizeMsg should not trigger handoff; got %T", updated)
	}
	if tr2.width != w2 || tr2.height != h2 {
		t.Errorf("resize not propagated: got %dx%d want %dx%d", tr2.width, tr2.height, w2, h2)
	}

	// Force the duration elapsed and assert Update hands off to `to`.
	tr2.startedAt = time.Now().Add(-time.Hour)
	handed, _ := tr2.Update(transitionFrameMsg(time.Now()))
	if _, still := handed.(transitionModel); still {
		t.Error("transition did not hand off after duration elapsed")
	}
	if _, isStub := handed.(stubModel); !isStub {
		t.Errorf("handoff returned %T, want stubModel", handed)
	}

	// Any key press during the transition should short-circuit to `to`.
	freshTr := newTransition(from, to, w, h, pin)
	short, _ := freshTr.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if _, still := short.(transitionModel); still {
		t.Error("key press did not short-circuit transition")
	}
}

func TestComposeSlideCollapseProgression(t *testing.T) {
	const (
		w = 40
		h = 10
	)
	fromRow := strings.Repeat("F", w)
	toRow := strings.Repeat("T", w)
	fromSnap := strings.Repeat(fromRow+"\n", h-1) + fromRow
	toSnap := strings.Repeat(toRow+"\n", h-1) + toRow
	pin := renderFeaturePin()

	cases := []struct {
		name string
		t    float64
	}{
		{"t=0", 0.0},
		{"t=0.25", 0.25},
		{"t=0.5", 0.5},
		{"t=0.75", 0.75},
		{"t=1", 1.0},
	}

	const pinY = h - 1
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := composeSlideCollapse(fromSnap, toSnap, pin, w, h, pinY, c.t)
			rows := strings.Split(out, "\n")
			if len(rows) != h {
				t.Fatalf("rows = %d, want %d", len(rows), h)
			}
			for i, row := range rows {
				if got := ansi.StringWidth(row); got != w {
					t.Errorf("row %d width = %d, want %d", i, got, w)
				}
			}
			if c.t >= pinFadeInStart {
				pinRow := rows[pinY]
				head := ansi.Cut(pinRow, 0, ansi.StringWidth(pin))
				if lipgloss.Width(head) == 0 {
					t.Errorf("pin area at t=%v was empty", c.t)
				}
			}
		})
	}
}

func TestComposeSlideCollapseCustomPinRow(t *testing.T) {
	const (
		w    = 40
		h    = 10
		pinY = 6
	)
	fromSnap := strings.Repeat(strings.Repeat("F", w)+"\n", h-1) + strings.Repeat("F", w)
	toSnap := strings.Repeat(strings.Repeat("T", w)+"\n", h-1) + strings.Repeat("T", w)
	pin := renderFeaturePin()

	out := composeSlideCollapse(fromSnap, toSnap, pin, w, h, pinY, 1.0)
	rows := strings.Split(out, "\n")
	if len(rows) != h {
		t.Fatalf("rows = %d, want %d", len(rows), h)
	}
	// At t=1 the pin should be at row pinY, not at row h-1.
	pinRow := rows[pinY]
	head := ansi.Cut(pinRow, 0, ansi.StringWidth(pin))
	if lipgloss.Width(head) == 0 {
		t.Errorf("pin missing at custom pinY=%d", pinY)
	}
	// The bottom rows (below pinY) should be entirely incoming.
	for i := pinY + 1; i < h; i++ {
		if !strings.Contains(rows[i], "T") {
			t.Errorf("row %d below pinY should be incoming 'T', got %q", i, rows[i])
		}
	}
}

func TestEaseOutCubic(t *testing.T) {
	cases := []struct {
		in, want float64
	}{
		{-0.5, 0},
		{0, 0},
		{1, 1},
		{1.5, 1},
	}
	for _, c := range cases {
		if got := easeOutCubic(c.in); got != c.want {
			t.Errorf("easeOutCubic(%v) = %v, want %v", c.in, got, c.want)
		}
	}
	if v := easeOutCubic(0.5); v <= 0.5 {
		t.Errorf("easeOutCubic(0.5) = %v, expected > 0.5 (ease-out)", v)
	}
}
