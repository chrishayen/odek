package composer

import (
	"context"
	_ "embed"
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
	"github.com/chrishayen/valkyrie/internal/feature"
	runepkg "github.com/chrishayen/valkyrie/internal/rune"
	"github.com/chrishayen/valkyrie/internal/runner"
)

//go:embed compose-agent.md
var instructions string

// Result holds the outcome of composing a feature.
type Result struct {
	FeatureName string  `json:"feature_name"`
	Output      string  `json:"output"`
	Coverage    float64 `json:"coverage"`
	TestsRan    bool    `json:"tests_ran"`
}

// Composer generates dispatcher and wiring code for a feature.
type Composer struct {
	featureStore *feature.Store
	runeStore    *runepkg.Store
	language     string
}

func New(featureStore *feature.Store, runeStore *runepkg.Store, language string) *Composer {
	return &Composer{featureStore: featureStore, runeStore: runeStore, language: language}
}

// Compose generates wiring code for the named feature.
// Reads the raw feature.md and passes it to the sandbox agent
// along with all rune signatures. The agent uses the prebuilt
// dispatch framework to wire runes together.
func (c *Composer) Compose(ctx context.Context, name string, r runner.Runner, logOut io.Writer) (*Result, error) {
	// Read raw feature file — the agent reads this as a document
	raw, err := c.featureStore.ReadRaw(name)
	if err != nil {
		return nil, err
	}

	// List all runes for context
	runes, err := c.runeStore.List()
	if err != nil {
		return nil, fmt.Errorf("listing runes: %w", err)
	}

	prompt := buildPrompt(raw, runes, c.language)

	// Ensure dispatch framework exists
	if err := framework.EnsureDispatch(c.featureStore.OutputPath()); err != nil {
		return nil, fmt.Errorf("ensuring dispatch framework: %w", err)
	}

	// Create code directory
	codeDir := c.featureStore.CodeDir(name)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return nil, fmt.Errorf("creating code dir: %w", err)
	}

	// Run the sandbox agent
	output, err := r.Run(ctx, prompt, logOut)
	if err != nil {
		return nil, fmt.Errorf("sandbox run failed: %w", err)
	}

	// Extract generated files
	if err := extractFiles(codeDir, output); err != nil {
		return nil, fmt.Errorf("extracting files: %w", err)
	}

	// Run tests
	coverage, testsRan := runTests(codeDir, c.language)

	// Update feature frontmatter
	feat, err := c.featureStore.Get(name)
	if err != nil {
		return nil, fmt.Errorf("reading feature for update: %w", err)
	}
	feat.Hydrated = true
	feat.Coverage = coverage
	if err := c.featureStore.Update(*feat); err != nil {
		return nil, fmt.Errorf("updating feature: %w", err)
	}

	return &Result{
		FeatureName: name,
		Output:      output,
		Coverage:    coverage,
		TestsRan:    testsRan,
	}, nil
}

func buildPrompt(rawFeature string, runes []runepkg.Rune, language string) string {
	var b strings.Builder

	b.WriteString(instructions)
	fmt.Fprintf(&b, "\n\nWrite all code in %s.\n", language)
	b.WriteString("\n---\n\n")

	// Raw feature spec — the agent reads it as a document
	b.WriteString("## Feature spec\n\n")
	b.WriteString(rawFeature)
	b.WriteString("\n\n---\n\n")

	// All rune signatures for context
	if len(runes) > 0 {
		b.WriteString("## Available runes\n\n")
		for _, r := range runes {
			fmt.Fprintf(&b, "- **%s**: `%s` — %s\n", r.Name, r.Signature, r.Description)
		}
	}

	return b.String()
}

// extractFiles parses agent output for FILE blocks and writes them to disk.
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

// runTests attempts to run tests in the code dir and parse coverage.
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

