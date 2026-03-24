package analyzer

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	runepkg "github.com/chrishayen/valkyrie/internal/rune"
	"github.com/chrishayen/valkyrie/internal/runner"
)

//go:embed rune-agent.md
var Instructions string

// ProposedRune is a rune the agent wants to create.
type ProposedRune struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Signature     string   `json:"signature"`
	Behavior      string   `json:"behavior"`
	PositiveTests []string `json:"positive_tests"`
	NegativeTests []string `json:"negative_tests"`
}

// ExistingMatch is a rune already in the registry that covers part of the requirements.
type ExistingMatch struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Covers      string `json:"covers"`
}

// Result is the structured output of an analysis.
type Result struct {
	NewRunes      []ProposedRune  `json:"new_runes"`
	ExistingRunes []ExistingMatch `json:"existing_runes"`
}

// Analyzer decomposes requirements into runes via a sandbox agent.
type Analyzer struct {
	store *runepkg.Store
}

// New creates an Analyzer backed by the given store.
func New(store *runepkg.Store) *Analyzer {
	return &Analyzer{store: store}
}

// Analyze sends requirements to a runner and returns the decomposition.
// If logOut is non-nil, sandbox output is streamed to it in real time.
func (a *Analyzer) Analyze(ctx context.Context, requirements string, r runner.Runner, logOut io.Writer) (*Result, error) {
	existing, err := a.store.List()
	if err != nil {
		return nil, fmt.Errorf("listing existing runes: %w", err)
	}

	prompt := buildAnalyzePrompt(requirements, existing)

	output, err := r.Run(ctx, prompt, logOut)
	if err != nil {
		return nil, fmt.Errorf("sandbox run failed: %w", err)
	}

	return parseResult(output)
}

// ToRune converts a ProposedRune to a Rune for storage.
func (p ProposedRune) ToRune() runepkg.Rune {
	return runepkg.Rune{
		Name:          p.Name,
		Description:   p.Description,
		Signature:     p.Signature,
		Behavior:      p.Behavior,
		PositiveTests: p.PositiveTests,
		NegativeTests: p.NegativeTests,
	}
}

func buildAnalyzePrompt(requirements string, existing []runepkg.Rune) string {
	var b strings.Builder

	b.WriteString(Instructions)
	b.WriteString("\n\n---\n\n")

	// Existing runes context
	if len(existing) > 0 {
		b.WriteString("## Existing runes in the registry\n\n")
		for _, r := range existing {
			fmt.Fprintf(&b, "- **%s**: %s\n", r.Name, r.Description)
		}
		b.WriteString("\n---\n\n")
	}

	// User requirements
	b.WriteString("## Requirements to analyze\n\n")
	b.WriteString(requirements)
	b.WriteString("\n\n---\n\n")

	// Output format instructions
	b.WriteString(`## Output format

Respond with ONLY a JSON object matching this exact structure. No markdown, no explanation, no code fences — just the JSON:

{
  "new_runes": [
    {
      "name": "feature/verb-noun-slug",
      "description": "One or two sentences: what the function does, accepts, and returns.",
      "signature": "(param: type, param: type) -> return_type",
      "behavior": "- Input: description of inputs and types\n- Output: description of output and type\n- Edge case details\n- Constraint details",
      "positive_tests": ["Given X, returns Y", "Given A, returns B"],
      "negative_tests": ["Given invalid X, throws error", "Given null, returns error"]
    }
  ],
  "existing_runes": [
    {
      "name": "existing-rune-name",
      "description": "Its current description",
      "covers": "Which requirement it satisfies"
    }
  ]
}

CRITICAL RULES for field separation:
- "description" is ONLY a 1-2 sentence summary. Nothing else.
- "signature" is the function signature with precise types. Use result[T, E] for functions that can fail.
- "behavior" is ONLY inputs, outputs, edge cases, and constraints. DO NOT put test cases in behavior.
- "positive_tests" is a separate array of passing test cases. NEVER put these in behavior.
- "negative_tests" is a separate array of failure/error test cases. NEVER put these in behavior.
- Use newlines (\n) within behavior for readability, with each point on its own line starting with "- ".`)

	return b.String()
}

func parseResult(output string) (*Result, error) {
	// Try to find JSON in the output — the agent might wrap it in text
	trimmed := strings.TrimSpace(output)

	// Strip markdown code fences if present
	if strings.HasPrefix(trimmed, "```") {
		lines := strings.Split(trimmed, "\n")
		// Remove first line (```json or ```) and last line (```)
		if len(lines) >= 3 {
			end := len(lines) - 1
			for end > 0 && strings.TrimSpace(lines[end]) != "```" {
				end--
			}
			trimmed = strings.Join(lines[1:end], "\n")
		}
	}

	// Find JSON object boundaries
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("no JSON object found in agent output:\n%s", output)
	}
	jsonStr := trimmed[start : end+1]

	var result Result
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("parsing agent JSON: %w\nraw output:\n%s", err, output)
	}
	return &result, nil
}
