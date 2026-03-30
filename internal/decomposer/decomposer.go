package decomposer

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/chrishayen/valkyrie/internal/claude"
	runepkg "github.com/chrishayen/valkyrie/internal/rune"
)

//go:embed rune-agent.md
var Instructions string

//go:embed decompose-agent.md
var systemPrompt string

// ProposedRune is a rune the agent wants to create.
type ProposedRune struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Signature     string   `json:"signature"`
	Behavior      string   `json:"behavior"`
	PositiveTests []string `json:"positive_tests"`
	NegativeTests []string `json:"negative_tests"`
	Refs          []string `json:"refs,omitempty"`
	Extend        bool     `json:"extend,omitempty"`
}

// ExistingMatch is a rune already in the registry that covers part of the requirements.
type ExistingMatch struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Covers      string `json:"covers"`
}

// Result is the structured output of a decomposition.
type Result struct {
	NewRunes      []ProposedRune  `json:"new_runes"`
	ExistingRunes []ExistingMatch `json:"existing_runes"`
	TreeOutput    string          `json:"tree_output,omitempty"`
	NewCount      int             `json:"new_count"`
	UpdatedCount  int             `json:"updated_count"`
}

// Decomposer decomposes requirements into runes.
type Decomposer struct {
	store  *runepkg.Store
	client *claude.Client
}

// New creates a Decomposer backed by the given store and client.
func New(store *runepkg.Store, client *claude.Client) *Decomposer {
	return &Decomposer{store: store, client: client}
}

// Decompose sends requirements to Claude and returns the decomposition.
func (d *Decomposer) Decompose(_ context.Context, requirements string) (*Result, error) {
	userPrompt, err := d.buildPrompt(requirements)
	if err != nil {
		return nil, err
	}

	output, err := d.client.Call(systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("claude call failed: %w", err)
	}

	return d.parseResult(output)
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
		Dependencies:  p.Refs,
	}
}

func (d *Decomposer) buildPrompt(requirements string) (string, error) {
	var b strings.Builder

	existingCtx, err := d.store.FormatExistingContext()
	if err != nil {
		return "", fmt.Errorf("formatting existing context: %w", err)
	}
	if existingCtx != "" {
		b.WriteString(existingCtx)
		b.WriteString("\n")
	}

	b.WriteString(requirements)

	return b.String(), nil
}

func (d *Decomposer) parseResult(output string) (*Result, error) {
	nodes := runepkg.ParseTree(output)

	result := &Result{
		TreeOutput: output,
	}

	for _, n := range nodes {
		if len(n.Refs) > 0 && n.Signature == "" && len(n.Pos) == 0 && len(n.Neg) == 0 {
			for _, ref := range n.Refs {
				result.ExistingRunes = append(result.ExistingRunes, ExistingMatch{
					Name:   ref,
					Covers: "Referenced by " + n.Path,
				})
			}
			continue
		}

		pr := ProposedRune{
			Name:      n.Path,
			Signature: n.Signature,
			Refs:      n.Refs,
			Extend:    n.Extend,
		}

		if len(n.Pos) > 0 {
			pr.Description = n.Pos[0]
			pr.PositiveTests = n.Pos
		}
		if len(n.Neg) > 0 {
			pr.NegativeTests = n.Neg
		}

		result.NewRunes = append(result.NewRunes, pr)
	}

	return result, nil
}
