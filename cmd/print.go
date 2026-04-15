package main

import (
	"fmt"
	"strings"

	"shotgun.dev/odek/internal/decomposer"
)

func printBanner() {
	fmt.Println("=== Auto-Recursive Rune Decomposition Engine ===")
}

func printInitialDecomposition(resp *decomposer.DecompositionResponse) {
	fmt.Printf("\n🌳 INITIAL DECOMPOSITION:\n")
	if len(resp.ProjectPackage.Runes) > 0 {
		fmt.Printf("   📦 %s\n", resp.ProjectPackage.Name)
		printRunesIndented(resp.ProjectPackage.Runes, 1)
	}
	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		fmt.Printf("   📚 %s\n", resp.StdPackage.Name)
		printRunesIndented(resp.StdPackage.Runes, 1)
	}
}

func printCompleteTree(allDecompositions []*decomposer.AutoDecomposition, path string, depth int, isRoot bool) {
	var decomposition *decomposer.AutoDecomposition
	for _, d := range allDecompositions {
		if d.Path == path {
			decomposition = d
			break
		}
	}

	if decomposition == nil || decomposition.Response == nil {
		return
	}

	resp := decomposition.Response

	if isRoot {
		fmt.Printf("🌳 ROOT DECOMPOSITION: %s\n", path)
	} else {
		indent := strings.Repeat("   ", depth)
		fmt.Printf("%s🔸 EXPANDED: %s\n", indent, path)
	}

	if len(resp.ProjectPackage.Runes) > 0 {
		pkgHeader := fmt.Sprintf("   📦 %s", resp.ProjectPackage.Name)
		if !isRoot {
			pkgHeader = strings.Repeat("   ", depth) + pkgHeader
		}
		fmt.Printf("%s\n", pkgHeader)
		printRunesIndented(resp.ProjectPackage.Runes, depth+1)
	}

	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		pkgHeader := fmt.Sprintf("   📚 %s", resp.StdPackage.Name)
		if !isRoot {
			pkgHeader = strings.Repeat("   ", depth) + pkgHeader
		}
		fmt.Printf("%s\n", pkgHeader)
		printRunesIndented(resp.StdPackage.Runes, depth+1)
	}

	for _, childPath := range decomposition.ChildPaths {
		printCompleteTree(allDecompositions, childPath, depth+1, false)
	}
}

func printRunesIndented(runes map[string]decomposer.Rune, indentLevel int) {
	if len(runes) == 0 {
		return
	}

	indent := strings.Repeat("   ", indentLevel)

	for name, r := range runes {
		fmt.Printf("%s├─ %s\n", indent, name)
		if r.Description != "" {
			descIndent := strings.Repeat("   ", indentLevel+1)
			wrappedDesc := wrapText(r.Description, 70-len(descIndent))
			fmt.Printf("%s│  └─ %s\n", descIndent, wrappedDesc)
		}
		if sig := decomposer.NormalizeFunctionSig(r.FunctionSig); sig != "" {
			sigIndent := strings.Repeat("   ", indentLevel+1)
			fmt.Printf("%s│     fn: %s\n", sigIndent, sig)
		}
	}
}

func wrapText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}
	return text[:maxWidth-3] + "..."
}

// printExpansionEvent renders a streaming ExpansionEvent to stdout. Designed
// to produce output roughly equivalent to the pre-refactor in-place prints
// from expandRecursively/expandOne, so CLI smoke-runs look familiar.
func printExpansionEvent(evt decomposer.ExpansionEvent) {
	switch e := evt.(type) {
	case decomposer.EventLevelStarted:
		fmt.Printf("\n📤 Dispatching %d expansions (depth %d)...\n", e.Count, e.Depth)
	case decomposer.EventRuneStarted:
		// Intentionally quiet: rune-level completion lines below are the
		// natural progress signal. Uncomment if per-rune dispatch lines are
		// ever desired.
	case decomposer.EventRuneExpanded:
		dur := fmt.Sprintf("%dms", e.ElapsedMs)
		if e.ChildCount == 0 {
			fmt.Printf("   ✓ %s: leaf (%s)\n", e.Path, dur)
		} else {
			fmt.Printf("   ➜ %s: %d sub-runes (%s)\n", e.Path, e.ChildCount, dur)
		}
	case decomposer.EventRuneError:
		fmt.Printf("   ⚠️  %s: %s (%dms)\n", e.Path, e.Err, e.ElapsedMs)
	case decomposer.EventLevelCompleted:
		factor := 0.0
		if e.WallClockMs > 0 {
			factor = float64(e.SumRequestMs) / float64(e.WallClockMs)
		}
		fmt.Printf("   ⏱️  level %d wall-clock: %dms, sum of requests: %dms (parallelism factor: %.1fx)\n",
			e.Depth, e.WallClockMs, e.SumRequestMs, factor)
	case decomposer.EventReadExample:
		fmt.Printf("🔎 read_example (%d handle%s)\n", len(e.Paths), plural(len(e.Paths)))
		for _, p := range e.Paths {
			fmt.Printf("   → %s\n", p)
		}
	case decomposer.EventCapReached:
		fmt.Printf("\n⚠️  Max total runes (%d) reached. Stopping expansion.\n", e.Cap)
	case decomposer.EventCancelled:
		fmt.Printf("\n⚠️  Expansion cancelled.\n")
	case decomposer.EventDone:
		// Summary is printed by the main after the loop so this is a no-op.
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
