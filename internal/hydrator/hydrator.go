package hydrator

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chrishayen/odek/framework"
	"github.com/chrishayen/odek/internal/claude"
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
	client   *claude.Client
	language string
}

func New(store *runepkg.Store, client *claude.Client, language string) *Hydrator {
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

	if err := extractFiles(codeDir, output); err != nil {
		return nil, fmt.Errorf("extracting files: %w", err)
	}

	coverage, testsRan := runTests(codeDir, h.language)

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

	if err := extractFiles(codeDir, output); err != nil {
		return nil, fmt.Errorf("extracting files: %w", err)
	}

	coverage, testsRan := runTests(codeDir, h.language)

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

	// Build children map for topological sort.
	pathSet := make(map[string]bool)
	for _, rn := range runes {
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
			cd := computeDepth(c)
			if cd > maxChild {
				maxChild = cd
			}
		}
		depths[p] = maxChild + 1
		return depths[p]
	}
	for p := range pathSet {
		computeDepth(p)
	}

	type levelTarget struct {
		depth int
		rune  runepkg.Rune
	}
	var lts []levelTarget
	for _, t := range targets {
		lts = append(lts, levelTarget{depth: depths[t.Name], rune: t})
	}
	sort.Slice(lts, func(i, j int) bool {
		if lts[i].depth != lts[j].depth {
			return lts[i].depth < lts[j].depth
		}
		return lts[i].rune.Name < lts[j].rune.Name
	})

	type specLevel struct {
		depth int
		runes []runepkg.Rune
	}
	var levels []specLevel
	for _, lt := range lts {
		if len(levels) == 0 || levels[len(levels)-1].depth != lt.depth {
			levels = append(levels, specLevel{depth: lt.depth})
		}
		levels[len(levels)-1].runes = append(levels[len(levels)-1].runes, lt.rune)
	}

	result := &HydrateAllResult{}

	for _, level := range levels {
		if logOut != nil {
			fmt.Fprintf(logOut, "Level %d: hydrating %d runes\n", level.depth, len(level.runes))
		}

		results := h.parallelHydrate(ctx, level.runes, concurrency, logOut)
		for _, hr := range results {
			if hr.err != nil {
				if logOut != nil {
					fmt.Fprintf(logOut, "  FAIL %s: %v\n", hr.name, hr.err)
				}
				result.Failed++
			} else {
				if logOut != nil {
					fmt.Fprintf(logOut, "  OK   %s\n", hr.name)
				}
				result.Hydrated++
			}
		}
	}

	return result, nil
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

	sb.WriteString(`
Instructions:
1. This component must be isolated from other runes.
   - Do NOT import or call any other runes directly.
   - All inter-rune communication goes through the dispatcher via serializable types.
2. Implement the component as described above, covering all specified behavior.
3. Write tests that verify every positive and negative test case listed above.
   Each test case should be its own test function with a clear name.
4. Output each file using this format exactly:

=== FILE: <filename> ===
<file contents>
=== END FILE ===

Keep the implementation minimal and focused on the described behavior.
Do not include explanations outside of file blocks.`)

	return sb.String()
}

func extractFiles(dir, output string) error {
	re := regexp.MustCompile(`(?s)=== FILE: (.+?) ===\n(.+?)=== END FILE ===`)
	matches := re.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		return os.WriteFile(filepath.Join(dir, "main.go"), []byte(output), 0644)
	}
	for _, m := range matches {
		filename := strings.TrimSpace(m[1])
		content := m[2]
		path := filepath.Join(dir, filename)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func runTests(dir, language string) (coverage float64, ran bool) {
	cmd, args := testCommand(language)
	if cmd == "" {
		return -1, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c := exec.CommandContext(ctx, cmd, args...)
	c.Dir = dir
	out, _ := c.CombinedOutput()

	coverage = parseCoverage(string(out))
	return coverage, true
}

func testCommand(language string) (string, []string) {
	switch language {
	case "go":
		return "go", []string{"test", "-cover", "."}
	case "ts":
		return "node", []string{"--test"}
	case "py":
		return "python", []string{"-m", "pytest", "-q"}
	default:
		return "", nil
	}
}

var coverageRe = regexp.MustCompile(`coverage:\s+([\d.]+)%`)

func parseCoverage(output string) float64 {
	m := coverageRe.FindStringSubmatch(output)
	if len(m) < 2 {
		return -1
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return -1
	}
	return v
}
