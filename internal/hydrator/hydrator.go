package hydrator

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chrishayen/valkyrie/framework"
	runepkg "github.com/chrishayen/valkyrie/internal/rune"
	"github.com/chrishayen/valkyrie/internal/runner"
)

// Result holds the outcome of hydrating a rune.
type Result struct {
	RuneName string  `json:"rune_name"`
	Output   string  `json:"output"`    // raw agent output
	Coverage float64 `json:"coverage"`  // test coverage %, -1 if unavailable
	TestsRan bool    `json:"tests_ran"`
}

// HydrationSpec contains everything a sub-agent needs to hydrate a rune.
type HydrationSpec struct {
	RuneName string `json:"rune_name"`
	Prompt   string `json:"prompt"` // enriched prompt with behavior, tests, isolation instructions
}

// Hydrator runs a sandbox agent to generate code for a rune.
type Hydrator struct {
	store    *runepkg.Store
	language string
}

func New(store *runepkg.Store, language string) *Hydrator {
	return &Hydrator{store: store, language: language}
}

// GetHydrationSpec returns the prompt for a rune.
// Used in non-sandbox mode: the calling agent spawns a sub-agent with this prompt,
// collects the FILE-block output, and passes it to FinalizeHydration.
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

// FinalizeHydration extracts files from the sub-agent's output, runs tests,
// and updates the rune record. The output must contain === FILE: ... === blocks.
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

// Hydrate generates code for the named rune using a sandbox runner, stores it, runs tests.
// Used in sandbox mode only. If logOut is non-nil, sandbox output is streamed to it in real time.
func (h *Hydrator) Hydrate(ctx context.Context, name string, r runner.Runner, logOut io.Writer) (*Result, error) {
	rune, err := h.store.Get(name)
	if err != nil {
		return nil, err
	}

	prompt := buildPrompt(rune, h.language)

	if err := framework.EnsureDispatch(h.store.OutputPath()); err != nil {
		return nil, fmt.Errorf("ensuring dispatch framework: %w", err)
	}

	codeDir := h.store.CodeDir(name)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return nil, fmt.Errorf("creating code dir: %w", err)
	}

	output, err := r.Run(ctx, prompt, logOut)
	if err != nil {
		return nil, fmt.Errorf("sandbox run failed: %w", err)
	}

	if err := extractFiles(codeDir, output); err != nil {
		return nil, fmt.Errorf("extracting files: %w", err)
	}

	coverage, testsRan := runTests(codeDir, h.language)

	rune.Hydrated = true
	rune.Coverage = coverage
	if err := h.store.Update(*rune); err != nil {
		return nil, fmt.Errorf("updating rune: %w", err)
	}

	return &Result{
		RuneName: name,
		Output:   output,
		Coverage: coverage,
		TestsRan: testsRan,
	}, nil
}

func buildPrompt(r *runepkg.Rune, language string) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, `You are implementing a single, isolated software component called "%s".
Write all code in %s.

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

// extractFiles parses agent output for FILE blocks and writes them to disk.
func extractFiles(dir, output string) error {
	re := regexp.MustCompile(`(?s)=== FILE: (.+?) ===\n(.+?)=== END FILE ===`)
	matches := re.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		// fallback: write raw output as main.go
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

// runTests attempts to run tests in the code dir and parse coverage.
// Returns coverage % and whether tests ran successfully.
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

