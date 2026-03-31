package composer

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chrishayen/odek/framework"
	"github.com/chrishayen/odek/internal/claude"
	"github.com/chrishayen/odek/internal/feature"
	runepkg "github.com/chrishayen/odek/internal/rune"
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
	client       *claude.Client
	language     string
}

func New(featureStore *feature.Store, runeStore *runepkg.Store, client *claude.Client, language string) *Composer {
	return &Composer{featureStore: featureStore, runeStore: runeStore, client: client, language: language}
}

// Compose generates wiring code for the named feature.
func (c *Composer) Compose(_ context.Context, name string) (*Result, error) {
	raw, err := c.featureStore.ReadRaw(name)
	if err != nil {
		return nil, err
	}

	runes, err := c.runeStore.List()
	if err != nil {
		return nil, fmt.Errorf("listing runes: %w", err)
	}

	prompt := buildPrompt(raw, runes, c.language)

	if err := framework.EnsureDispatch(c.featureStore.OutputPath()); err != nil {
		return nil, fmt.Errorf("ensuring dispatch framework: %w", err)
	}

	codeDir := c.featureStore.CodeDir(name)
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return nil, fmt.Errorf("creating code dir: %w", err)
	}

	output, err := c.client.Call(instructions, prompt)
	if err != nil {
		return nil, fmt.Errorf("claude call failed: %w", err)
	}

	if err := extractFiles(codeDir, output); err != nil {
		return nil, fmt.Errorf("extracting files: %w", err)
	}

	coverage, testsRan := runTests(codeDir, c.language)

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

	fmt.Fprintf(&b, "Write all code in %s.\n", language)
	b.WriteString("\n---\n\n")

	b.WriteString("## Feature spec\n\n")
	b.WriteString(rawFeature)
	b.WriteString("\n\n---\n\n")

	if len(runes) > 0 {
		b.WriteString("## Available runes\n\n")
		for _, r := range runes {
			fmt.Fprintf(&b, "- **%s**: `%s` — %s\n", r.Name, r.Signature, r.Description)
		}
	}

	return b.String()
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
