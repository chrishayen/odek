package main

import (
	"fmt"
	"sort"
	"strings"

	"shotgun.dev/odek/internal/decomposer"
)

func printBanner() {
	fmt.Println("=== Two-Pass Rune Decomposition Engine ===")
}

// printCompleteTree renders the final decomposition response as an
// indented tree. std package first (when present), then project package.
func printCompleteTree(resp *decomposer.DecompositionResponse) {
	if resp == nil {
		return
	}
	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		fmt.Printf("\n📚 %s\n", resp.StdPackage.Name)
		printRunesIndented(resp.StdPackage.Runes, 1)
	}
	fmt.Printf("\n📦 %s\n", resp.ProjectPackage.Name)
	printRunesIndented(resp.ProjectPackage.Runes, 1)
}

// printRunesIndented recursively renders a rune map with dot-tree indentation.
func printRunesIndented(runes map[string]decomposer.Rune, indentLevel int) {
	if len(runes) == 0 {
		return
	}
	names := make([]string, 0, len(runes))
	for name := range runes {
		names = append(names, name)
	}
	sort.Strings(names)

	indent := strings.Repeat("   ", indentLevel)
	for _, name := range names {
		r := runes[name]
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
		if len(r.Children) > 0 {
			printRunesIndented(r.Children, indentLevel+1)
		}
	}
}

func wrapText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}
	return text[:maxWidth-3] + "..."
}

// printDecompositionEvent renders a streaming DecompositionEvent to stdout.
// Designed for CLI smoke runs of the two-pass pipeline.
func printDecompositionEvent(evt decomposer.DecompositionEvent) {
	switch e := evt.(type) {
	case decomposer.EventPhaseStarted:
		switch e.Phase {
		case decomposer.PhaseContract:
			fmt.Printf("\n📝 Pass 1: designing contract...\n\n")
		case decomposer.PhaseExtraction:
			fmt.Printf("\n\n🔧 Pass 2: extracting runes...\n")
		}
	case decomposer.EventContractChunk:
		// Stream contract text directly to stdout as it arrives.
		fmt.Print(e.Text)
	case decomposer.EventContractComplete:
		fmt.Printf("\n   ⏱️  contract: %dms\n", e.ElapsedMs)
	case decomposer.EventExtractionProgress:
		// Overwrite on the same line — CR resets to column 0.
		fmt.Printf("\r   ...%d bytes received", e.Bytes)
	case decomposer.EventRunesComplete:
		total := countTreeRunes(e.Response)
		fmt.Printf("\n   ⏱️  extraction: %dms · %d runes\n", e.ElapsedMs, total)
	case decomposer.EventReadExample:
		fmt.Printf("\n🔎 read_example (%d handle%s)\n", len(e.Paths), plural(len(e.Paths)))
		for _, p := range e.Paths {
			fmt.Printf("   → %s\n", p)
		}
	case decomposer.EventError:
		fmt.Printf("\n⚠️  [%s] error: %s\n", e.Phase, e.Err)
	case decomposer.EventCancelled:
		fmt.Printf("\n⚠️  Cancelled.\n")
	case decomposer.EventDone:
		// Summary is printed by main after the loop.
	}
}

func countTreeRunes(resp *decomposer.DecompositionResponse) int {
	if resp == nil {
		return 0
	}
	n := countRunesRecursive(resp.ProjectPackage.Runes)
	if resp.StdPackage != nil {
		n += countRunesRecursive(resp.StdPackage.Runes)
	}
	return n
}

func countRunesRecursive(runes map[string]decomposer.Rune) int {
	n := len(runes)
	for _, r := range runes {
		n += countRunesRecursive(r.Children)
	}
	return n
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
