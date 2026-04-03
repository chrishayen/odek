package composer

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/chrishayen/odek/framework"
	"github.com/chrishayen/odek/internal/claude"
	"github.com/chrishayen/odek/internal/codegen"
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

	if err := codegen.ExtractFiles(codeDir, output); err != nil {
		return nil, fmt.Errorf("extracting files: %w", err)
	}

	coverage, testsRan := codegen.RunTests(codeDir, c.language)

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

