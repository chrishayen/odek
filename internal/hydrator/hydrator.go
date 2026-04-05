package hydrator

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/chrishayen/odek/framework"
	"github.com/chrishayen/odek/internal/llm"
	"github.com/chrishayen/odek/internal/codegen"
	runepkg "github.com/chrishayen/odek/internal/rune"
)

// Result holds the outcome of hydrating a rune.
type Result struct {
	RuneName string  `json:"rune_name"`
	Output   string  `json:"output"`
	Coverage float64 `json:"coverage"`
	TestsRan bool    `json:"tests_ran"`
}

// HydrationSpec contains everything a sub-agent needs to hydrate a rune.
type HydrationSpec struct {
	RuneName string `json:"rune_name"`
	Prompt   string `json:"prompt"`
}

// HydrateAllResult holds the outcome of batch hydration.
type HydrateAllResult struct {
	Hydrated int `json:"hydrated"`
	Verified int `json:"verified"`
	Failed   int `json:"failed"`
}

// Hydrator runs agents to generate code for runes.
type Hydrator struct {
	store    *runepkg.Store
	client   *llm.Client
	language string
}

func New(store *runepkg.Store, client *llm.Client, language string) *Hydrator {
	return &Hydrator{store: store, client: client, language: language}
}

// GetHydrationSpec returns the prompt for a rune.
func (h *Hydrator) GetHydrationSpec(name string) (*HydrationSpec, error) {
	r, err := h.store.Get(name)
	if err != nil {
		return nil, err
	}
	return &HydrationSpec{
		RuneName: name,
		Prompt:   buildPrompt(r, h.language),
	}, nil
}

// FinalizeHydration extracts files from agent output, runs tests, and updates the rune.
func (h *Hydrator) FinalizeHydration(name, output string) (*Result, error) {
	if strings.TrimSpace(output) == "" {
		return nil, fmt.Errorf("output is empty — sub-agent produced no code")
	}

	r, err := h.store.Get(name)
	if err != nil {
		return nil, err
	}

	if err := framework.EnsureDispatch(h.store.OutputPath()); err != nil {
		return nil, fmt.Errorf("ensuring dispatch framework: %w", err)
	}

	codeDir := h.store.CodeDir(name)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return nil, fmt.Errorf("creating code dir: %w", err)
	}

	if err := codegen.ExtractFiles(codeDir, output); err != nil {
		return nil, fmt.Errorf("extracting files: %w", err)
	}

	coverage, testsRan := codegen.RunTests(codeDir, h.language)

	r.Hydrated = true
	r.Coverage = coverage
	if err := h.store.Update(*r); err != nil {
		return nil, fmt.Errorf("updating rune: %w", err)
	}

	return &Result{
		RuneName: name,
		Output:   output,
		Coverage: coverage,
		TestsRan: testsRan,
	}, nil
}

// Hydrate generates code for a single rune.
func (h *Hydrator) Hydrate(_ context.Context, name string) (*Result, error) {
	rn, err := h.store.Get(name)
	if err != nil {
		return nil, err
	}

	prompt := buildPrompt(rn, h.language)

	if err := framework.EnsureDispatch(h.store.OutputPath()); err != nil {
		return nil, fmt.Errorf("ensuring dispatch framework: %w", err)
	}

	codeDir := h.store.CodeDir(name)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return nil, fmt.Errorf("creating code dir: %w", err)
	}

	output, err := h.client.Call("", prompt)
	if err != nil {
		return nil, fmt.Errorf("claude call failed: %w", err)
	}

	if err := codegen.ExtractFiles(codeDir, output); err != nil {
		return nil, fmt.Errorf("extracting files: %w", err)
	}

	coverage, testsRan := codegen.RunTests(codeDir, h.language)

	rn.Hydrated = true
	rn.Coverage = coverage
	if err := h.store.Update(*rn); err != nil {
		return nil, fmt.Errorf("updating rune: %w", err)
	}

	return &Result{
		RuneName: name,
		Output:   output,
		Coverage: coverage,
		TestsRan: testsRan,
	}, nil
}

// HydrateAll orchestrates parallel hydration of all un-hydrated runes.
func (h *Hydrator) HydrateAll(ctx context.Context, concurrency int, verify bool, logOut io.Writer) (*HydrateAllResult, error) {
	runes, err := h.store.List()
	if err != nil {
		return nil, fmt.Errorf("listing runes: %w", err)
	}

	var targets []runepkg.Rune
	for _, rn := range runes {
		if !rn.Hydrated {
			targets = append(targets, rn)
		}
	}
	if len(targets) == 0 {
		return &HydrateAllResult{}, nil
	}

	levels := groupByDepth(runes, targets)
	result := &HydrateAllResult{}

	for _, level := range levels {
		logProgress(logOut, "Level %d: hydrating %d runes\n", level.depth, len(level.runes))

		for _, hr := range h.parallelHydrate(ctx, level.runes, concurrency, logOut) {
			if hr.err != nil {
				logProgress(logOut, "  FAIL %s: %v\n", hr.name, hr.err)
				result.Failed++
				continue
			}
			logProgress(logOut, "  OK   %s\n", hr.name)
			result.Hydrated++
		}
	}

	return result, nil
}

type runeLevel struct {
	depth int
	runes []runepkg.Rune
}

// groupByDepth computes tree depths and groups un-hydrated runes by level
// (leaves first) for bottom-up hydration.
func groupByDepth(all, targets []runepkg.Rune) []runeLevel {
	pathSet := make(map[string]bool, len(all))
	for _, rn := range all {
		pathSet[rn.Name] = true
	}
	children := runepkg.BuildChildrenMap(keys(pathSet))

	depths := make(map[string]int)
	var computeDepth func(string) int
	computeDepth = func(p string) int {
		if d, ok := depths[p]; ok {
			return d
		}
		maxChild := -1
		for _, c := range children[p] {
			if cd := computeDepth(c); cd > maxChild {
				maxChild = cd
			}
		}
		depths[p] = maxChild + 1
		return depths[p]
	}
	for p := range pathSet {
		computeDepth(p)
	}

	type depthTarget struct {
		depth int
		rune  runepkg.Rune
	}
	dts := make([]depthTarget, len(targets))
	for i, t := range targets {
		dts[i] = depthTarget{depth: depths[t.Name], rune: t}
	}
	sort.Slice(dts, func(i, j int) bool {
		if dts[i].depth != dts[j].depth {
			return dts[i].depth < dts[j].depth
		}
		return dts[i].rune.Name < dts[j].rune.Name
	})

	var levels []runeLevel
	for _, dt := range dts {
		if len(levels) == 0 || levels[len(levels)-1].depth != dt.depth {
			levels = append(levels, runeLevel{depth: dt.depth})
		}
		levels[len(levels)-1].runes = append(levels[len(levels)-1].runes, dt.rune)
	}
	return levels
}

func logProgress(w io.Writer, format string, args ...any) {
	if w != nil {
		fmt.Fprintf(w, format, args...)
	}
}

type hydrateResult struct {
	name string
	err  error
}

func (h *Hydrator) parallelHydrate(ctx context.Context, runes []runepkg.Rune, concurrency int, logOut io.Writer) []hydrateResult {
	results := make([]hydrateResult, len(runes))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, rn := range runes {
		wg.Add(1)
		go func(idx int, name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			_, err := h.Hydrate(ctx, name)
			results[idx] = hydrateResult{name: name, err: err}
		}(i, rn.Name)
	}

	wg.Wait()
	return results
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func buildPrompt(r *runepkg.Rune, language string) string {
	var sb strings.Builder
	shortName := runepkg.ShortName(r.Name)

	fmt.Fprintf(&sb, `You are implementing a single, isolated library function called "%s".
Write all code in %s. This is a library component meant to be imported and called by consumers — not an executable entry point. Do not generate main() functions or CLI scaffolding.

Description: %s

Signature: %s
`, r.Name, language, r.Description, r.Signature)

	if r.Behavior != "" {
		fmt.Fprintf(&sb, "\nBehavior:\n%s\n", r.Behavior)
	}

	if len(r.PositiveTests) > 0 {
		sb.WriteString("\nExpected passing test cases:\n")
		for _, t := range r.PositiveTests {
			fmt.Fprintf(&sb, "- %s\n", t)
		}
	}

	if len(r.NegativeTests) > 0 {
		sb.WriteString("\nExpected failing/error test cases:\n")
		for _, t := range r.NegativeTests {
			fmt.Fprintf(&sb, "- %s\n", t)
		}
	}

	fmt.Fprintf(&sb, `
Instructions:
1. This component must be isolated from other runes.
   - Do NOT import or call any other runes directly.
   - All inter-rune communication goes through the dispatcher via serializable types.
2. Implement the component as described above, covering all specified behavior.
3. Write tests that verify every positive and negative test case listed above.
   Each test case should be its own test function with a clear name.
4. Do NOT generate package.json, tsconfig.json, vitest.config.ts, or any project
   configuration files. Only output source (.ts, .js, .go, .py) and test files.
   The build and test infrastructure is managed externally.
5. Name your files using ONLY the short name "%s" — for example "%s.ts" and
   "%s.test.ts". Do NOT create subdirectories or nest files under src/ or any
   other folder. All files must be plain filenames with no path separators.
6. Output each file using this format exactly:

=== FILE: <filename> ===
<file contents>
=== END FILE ===

Keep the implementation minimal and focused on the described behavior.
Do not include explanations outside of file blocks.`, shortName, shortName, shortName)

	return sb.String()
}

