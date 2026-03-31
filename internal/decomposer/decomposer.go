package decomposer

import (
	"context"
	_ "embed"
	"encoding/json"
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
	FeatureName   string          `json:"feature_name,omitempty"`
	Summary       string          `json:"summary,omitempty"`
	FlowDiagram   string          `json:"flow_diagram,omitempty"`
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

const metaSystemPrompt = `You name features. Given a requirement, respond with exactly this JSON and nothing else:
{"name":"snake_case_slug","summary":"One sentence summary."}
The name is a short slug (e.g. auth, payment, http_serve). The summary describes what the feature does.`

const flowSystemPrompt = `You draw flow diagrams for software features using box-drawing characters. Given a requirement, show how the components connect to deliver the feature's functionality. Use arrows (──>, <──) to show data/control flow between boxes drawn with ┌─┐│└─┘ characters. Label arrows with what flows between components. Keep it compact — fit within 80 columns. Show the happy path top-to-bottom. No prose, just the diagram.`

// Decompose sends requirements to Claude and returns the decomposition.
func (d *Decomposer) Decompose(_ context.Context, requirements string) (*Result, error) {
	userPrompt, err := d.buildPrompt(requirements)
	if err != nil {
		return nil, err
	}

	type treeOut struct {
		output string
		err    error
	}
	type metaOut struct {
		name    string
		summary string
		err     error
	}
	type flowOut struct {
		diagram string
		err     error
	}

	treeCh := make(chan treeOut, 1)
	metaCh := make(chan metaOut, 1)
	flowCh := make(chan flowOut, 1)

	go func() {
		output, err := d.client.Call(systemPrompt, userPrompt)
		treeCh <- treeOut{output, err}
	}()
	go func() {
		name, summary, err := d.generateMeta(requirements)
		metaCh <- metaOut{name, summary, err}
	}()
	go func() {
		diagram, err := d.client.Call(flowSystemPrompt, requirements)
		flowCh <- flowOut{strings.TrimSpace(diagram), err}
	}()

	tr := <-treeCh
	if tr.err != nil {
		return nil, fmt.Errorf("claude call failed: %w", tr.err)
	}

	result, err := d.parseResult(tr.output)
	if err != nil {
		return nil, err
	}

	mr := <-metaCh
	if mr.err == nil {
		result.FeatureName = mr.name
		result.Summary = mr.summary
	}

	fr := <-flowCh
	if fr.err == nil {
		result.FlowDiagram = fr.diagram
	}

	return result, nil
}

func (d *Decomposer) generateMeta(requirements string) (string, string, error) {
	output, err := d.client.Call(metaSystemPrompt, requirements)
	if err != nil {
		return "", "", err
	}
	var meta struct {
		Name    string `json:"name"`
		Summary string `json:"summary"`
	}
	if err := json.Unmarshal([]byte(claude.StripCodeFences(output)), &meta); err != nil {
		return "", "", fmt.Errorf("meta parse: %w", err)
	}
	return meta.Name, meta.Summary, nil
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
