package tui

import (
	"context"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	openai "shotgun.dev/odek/openai"
)

func makeSplit(t *testing.T, width, height int) splitPaneModel {
	t.Helper()
	left := newCreateFeatureModel(context.Background(), (*openai.Client)(nil), nil, width, height)
	leftW, rightW := splitWidths(width)
	left.resize(leftW, height)
	right := newFeatureDecompModel(rightW, height, nil, left.state)
	return newSplitPaneModel(left, right, width, height)
}

func TestSplitPaneView(t *testing.T) {
	const (
		w = 160
		h = 40
	)
	split := makeSplit(t, w, h)

	view := split.View()
	if view.Content == "" {
		t.Fatal("split view rendered empty")
	}
	lines := strings.Split(view.Content, "\n")
	if len(lines) == 0 {
		t.Fatal("split view has no lines")
	}
	first := lines[0]
	got := ansi.StringWidth(first)
	// JoinHorizontal of two full-width panes + 3-col separator should equal total.
	if got != w {
		t.Errorf("first row width = %d, want %d", got, w)
	}
}

func TestSplitPaneTabTogglesFocus(t *testing.T) {
	split := makeSplit(t, 160, 40)
	if split.focus != 0 {
		t.Fatalf("initial focus = %d, want 0", split.focus)
	}

	next, _ := split.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	after, ok := next.(splitPaneModel)
	if !ok {
		t.Fatalf("tab returned %T, want splitPaneModel", next)
	}
	if after.focus != 1 {
		t.Errorf("focus after tab = %d, want 1", after.focus)
	}

	next2, _ := after.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	after2 := next2.(splitPaneModel)
	if after2.focus != 0 {
		t.Errorf("focus after second tab = %d, want 0", after2.focus)
	}
}

func TestSplitPaneNarrowResizeCollapses(t *testing.T) {
	split := makeSplit(t, 160, 40)
	next, _ := split.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if _, stillSplit := next.(splitPaneModel); stillSplit {
		t.Fatal("narrow resize did not collapse split")
	}
	if _, isLeft := next.(createFeatureModel); !isLeft {
		t.Errorf("narrow resize returned %T, want createFeatureModel (focus=0)", next)
	}
}

func TestSplitPaneNarrowResizeCollapsesToFocused(t *testing.T) {
	split := makeSplit(t, 160, 40)
	// Focus the right pane before shrinking.
	tabbed, _ := split.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	split = tabbed.(splitPaneModel)

	next, _ := split.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if _, isRight := next.(featureDecompModel); !isRight {
		t.Errorf("narrow resize with focus=1 returned %T, want featureDecompModel", next)
	}
}

func TestSplitPaneWideResizeStaysSplit(t *testing.T) {
	split := makeSplit(t, 160, 40)
	next, _ := split.Update(tea.WindowSizeMsg{Width: 180, Height: 42})
	after, ok := next.(splitPaneModel)
	if !ok {
		t.Fatalf("wide resize returned %T, want splitPaneModel", next)
	}
	if after.width != 180 || after.height != 42 {
		t.Errorf("split dimensions after resize = %dx%d, want 180x42", after.width, after.height)
	}
	leftW, rightW := splitWidths(180)
	if after.left.width != leftW {
		t.Errorf("left pane width = %d, want %d", after.left.width, leftW)
	}
	if after.right.width != rightW {
		t.Errorf("right pane width = %d, want %d", after.right.width, rightW)
	}
}

func TestSplitPaneCtrlEnterNoOp(t *testing.T) {
	split := makeSplit(t, 160, 40)
	next, _ := split.Update(tea.KeyPressMsg{Code: tea.KeyEnter, Mod: tea.ModCtrl})
	if _, stillSplit := next.(splitPaneModel); !stillSplit {
		t.Errorf("ctrl+enter in split should no-op, got %T", next)
	}
}

func TestSplitPaneShiftTabTogglesFocus(t *testing.T) {
	split := makeSplit(t, 160, 40)
	if split.focus != 0 {
		t.Fatalf("initial focus = %d, want 0", split.focus)
	}
	next, _ := split.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	after := next.(splitPaneModel)
	if after.focus != 1 {
		t.Errorf("shift+tab focus = %d, want 1", after.focus)
	}
}

func TestSplitPaneResizeBrackets(t *testing.T) {
	split := makeSplit(t, 160, 40)
	initialLeftW := split.leftW

	// ctrl+] grows the left pane
	next, _ := split.Update(tea.KeyPressMsg{Code: ']', Mod: tea.ModCtrl})
	after := next.(splitPaneModel)
	if after.leftW != initialLeftW+splitResizeStep {
		t.Errorf("ctrl+] leftW = %d, want %d", after.leftW, initialLeftW+splitResizeStep)
	}
	if after.left.width != after.leftW {
		t.Errorf("left pane width not propagated: got %d want %d", after.left.width, after.leftW)
	}

	// ctrl+[ shrinks the left pane
	next2, _ := after.Update(tea.KeyPressMsg{Code: '[', Mod: tea.ModCtrl})
	after2 := next2.(splitPaneModel)
	if after2.leftW != after.leftW-splitResizeStep {
		t.Errorf("ctrl+[ leftW = %d, want %d", after2.leftW, after.leftW-splitResizeStep)
	}
}

func TestSplitPaneResizeClamps(t *testing.T) {
	split := makeSplit(t, 160, 40)
	// Drive ctrl+[ enough times to hit the minimum clamp.
	for range 100 {
		next, _ := split.Update(tea.KeyPressMsg{Code: '[', Mod: tea.ModCtrl})
		split = next.(splitPaneModel)
	}
	if split.leftW < splitMinSubW {
		t.Errorf("leftW = %d below min %d", split.leftW, splitMinSubW)
	}
	if split.leftW != splitMinSubW {
		t.Errorf("leftW should clamp to %d after many ctrl+[, got %d", splitMinSubW, split.leftW)
	}

	// Drive ctrl+] enough times to hit the maximum clamp.
	for range 100 {
		next, _ := split.Update(tea.KeyPressMsg{Code: ']', Mod: tea.ModCtrl})
		split = next.(splitPaneModel)
	}
	avail := 160 - splitSepCols
	wantMax := avail - splitMinSubW
	if split.leftW != wantMax {
		t.Errorf("leftW should clamp to %d after many ctrl+], got %d", wantMax, split.leftW)
	}
}

func TestSplitWidths(t *testing.T) {
	cases := []struct {
		total, wantL, wantR int
	}{
		{150, 48, 98},
		{160, 52, 104},
		{200, 65, 131},
		{151, 49, 98},
	}
	for _, c := range cases {
		l, r := splitWidths(c.total)
		if l != c.wantL || r != c.wantR {
			t.Errorf("splitWidths(%d) = (%d, %d), want (%d, %d)", c.total, l, r, c.wantL, c.wantR)
		}
		if l+r+splitSepCols != c.total {
			t.Errorf("splitWidths(%d) sum mismatch: %d + %d + %d != %d", c.total, l, r, splitSepCols, c.total)
		}
		if r <= l {
			t.Errorf("splitWidths(%d): right pane (%d) should be wider than left (%d)", c.total, r, l)
		}
	}
}
