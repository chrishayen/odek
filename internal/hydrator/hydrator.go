package hydrator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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

// Hydrator runs a sandbox agent to generate code for a rune.
type Hydrator struct {
	store *runepkg.Store
}

func New(store *runepkg.Store) *Hydrator {
	return &Hydrator{store: store}
}

// Hydrate generates code for the named rune using the given runner, stores it, runs tests.
func (h *Hydrator) Hydrate(ctx context.Context, name string, r runner.Runner) (*Result, error) {
	rune, err := h.store.Get(name)
	if err != nil {
		return nil, err
	}

	// Build prompt: instruct the agent to generate code + tests
	prompt := buildPrompt(rune.Name, rune.Description)

	// Create code directory
	codeDir := h.store.CodeDir(name)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return nil, fmt.Errorf("creating code dir: %w", err)
	}

	// Run the sandbox agent
	output, err := r.Run(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("sandbox run failed: %w", err)
	}

	// Extract and store code files from agent output
	if err := extractFiles(codeDir, output); err != nil {
		return nil, fmt.Errorf("extracting files: %w", err)
	}

	// Run tests and get coverage
	coverage, testsRan := runTests(codeDir)

	// Update rune record
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

func buildPrompt(name, description string) string {
	return fmt.Sprintf(`You are implementing a software component called "%s".

Description: %s

Instructions:
1. Implement the component as described above.
2. Write behavior tests that verify the described functionality.
3. Output each file using this format exactly:

=== FILE: <filename> ===
<file contents>
=== END FILE ===

Keep the implementation minimal and focused on the described behavior.
Do not include explanations outside of file blocks.`, name, description)
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
func runTests(dir string) (coverage float64, ran bool) {
	cmd, args := detectTestCommand(dir)
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

func detectTestCommand(dir string) (string, []string) {
	// Go
	if hasFile(dir, "go.mod") || hasGlob(dir, "*.go") {
		return "go", []string{"test", "-cover", "."}
	}
	// Python
	if hasGlob(dir, "*.py") {
		return "python", []string{"-m", "pytest", "--tb=short", dir}
	}
	// Node/TypeScript
	if hasFile(dir, "package.json") {
		return "npm", []string{"test", "--prefix", dir}
	}
	return "", nil
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

func hasFile(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

func hasGlob(dir, pattern string) bool {
	matches, _ := filepath.Glob(filepath.Join(dir, pattern))
	return len(matches) > 0
}
