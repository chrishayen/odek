package decomposer

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/chrishayen/odek/internal/claude"
	runepkg "github.com/chrishayen/odek/internal/rune"
)

func logProgress(w io.Writer, format string, args ...any) {
	if w != nil {
		fmt.Fprintf(w, format+"\n", args...)
	}
}

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
	Assumptions   []string `json:"assumptions,omitempty"`
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
	FeatureName      string            `json:"feature_name,omitempty"`
	Summary          string            `json:"summary,omitempty"`
	FlowDiagram      string            `json:"flow_diagram,omitempty"`
	NewRunes         []ProposedRune    `json:"new_runes"`
	ExistingRunes    []ExistingMatch   `json:"existing_runes"`
	TreeOutput       string            `json:"tree_output,omitempty"`
	NewCount         int               `json:"new_count"`
	UpdatedCount     int               `json:"updated_count"`
	PackageSummaries map[string]string `json:"package_summaries,omitempty"`
}

// Decomposer decomposes requirements into runes.
type Decomposer struct {
	store   *runepkg.Store
	client  *claude.Client
	project string
}

// New creates a Decomposer backed by the given store and client.
func New(store *runepkg.Store, client *claude.Client, project string) *Decomposer {
	return &Decomposer{store: store, client: client, project: project}
}

const metaSystemPrompt = `You name features. Given a requirement, respond with exactly this JSON and nothing else:
{"name":"snake_case_slug","summary":"One sentence summary."}
The name is a short slug (e.g. auth, payment, http_serve). The summary describes what the feature does.`

const flowSystemPrompt = `You draw flow diagrams for software features using box-drawing characters. Given a requirement, show how the components connect to deliver the feature's functionality. Use arrows (──>, <──) to show data/control flow between boxes drawn with ┌─┐│└─┘ characters. Label arrows with what flows between components. Keep it compact — fit within 80 columns. Show the happy path top-to-bottom. No prose, just the diagram.`

// Decompose sends requirements to Claude and returns the decomposition.
// If prevDecomposition is non-empty, it is included as context so the LLM can iterate.
func (d *Decomposer) Decompose(_ context.Context, requirements, prevDecomposition string, logOut io.Writer) (*Result, error) {
	userPrompt, err := d.buildPrompt(requirements, prevDecomposition)
	if err != nil {
		return nil, err
	}

	logProgress(logOut, "Analyzing requirements...")

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

	prompt := strings.ReplaceAll(systemPrompt, "project_name", d.project)
	go func() {
		output, err := d.client.Call(prompt, userPrompt)
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
	logProgress(logOut, "Parsing composition tree...")

	result, err := d.parseResult(tr.output)
	if err != nil {
		return nil, err
	}
	logProgress(logOut, "Found %d runes", len(result.NewRunes))

	d.generatePackageSummaries(result, logOut)

	mr := <-metaCh
	if mr.err == nil {
		result.FeatureName = mr.name
		result.Summary = mr.summary
	}

	fr := <-flowCh
	if fr.err == nil {
		result.FlowDiagram = fr.diagram
	}

	logProgress(logOut, "Decomposition complete")
	return result, nil
}

const pkgSummaryPrompt = `You summarize software packages in one sentence. Given a package name and its contained runes, write a single concise sentence describing what the package provides. No quotes, no markdown, no prefix — just the sentence.`

// generatePackageSummaries identifies packages in the result and generates
// a one-sentence summary for each via Claude. Calls run in parallel.
func (d *Decomposer) generatePackageSummaries(result *Result, logOut io.Writer) {
	// Find packages: any rune whose name is a prefix of another rune's name.
	packages := map[string][]ProposedRune{}
	for _, r := range result.NewRunes {
		prefix := r.Name + "."
		for _, other := range result.NewRunes {
			if strings.HasPrefix(other.Name, prefix) {
				packages[r.Name] = append(packages[r.Name], other)
				break
			}
		}
	}
	if len(packages) == 0 {
		return
	}
	logProgress(logOut, "Generating %d package summaries...", len(packages))

	type summaryResult struct {
		name    string
		summary string
	}
	ch := make(chan summaryResult, len(packages))

	for name, children := range packages {
		go func(name string, children []ProposedRune) {
			var b strings.Builder
			fmt.Fprintf(&b, "Package: %s\nContains:\n", name)
			for _, c := range children {
				fmt.Fprintf(&b, "  %s — %s\n", c.Name, c.Description)
			}
			resp, err := d.client.Call(pkgSummaryPrompt, b.String())
			if err != nil {
				ch <- summaryResult{name, ""}
				return
			}
			ch <- summaryResult{name, strings.TrimSpace(resp)}
		}(name, children)
	}

	result.PackageSummaries = make(map[string]string, len(packages))
	for range packages {
		sr := <-ch
		if sr.summary != "" {
			result.PackageSummaries[sr.name] = sr.summary
		}
	}
}

const askSystemPrompt = `You answer questions about software decompositions concisely and directly. No markdown headers.`

// Ask answers a question about a decomposition given context.
func (d *Decomposer) Ask(_ context.Context, question, decompContext string) (string, error) {
	return d.client.Call(askSystemPrompt, decompContext+"\n\nQuestion: "+question)
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
		Assumptions:   p.Assumptions,
		Dependencies:  p.Refs,
	}
}

func (d *Decomposer) buildPrompt(requirements, prevDecomposition string) (string, error) {
	var b strings.Builder

	existingCtx, err := d.store.FormatExistingContext()
	if err != nil {
		return "", fmt.Errorf("formatting existing context: %w", err)
	}
	if existingCtx != "" {
		b.WriteString(existingCtx)
		b.WriteString("\n")
	}

	if prevDecomposition != "" {
		b.WriteString("Your previous output (to be refined — output the COMPLETE updated trees including all std and project units, not just changes):\n")
		b.WriteString(prevDecomposition)
		b.WriteString("\n\nThe user wants to refine the above. Apply this change and re-output both complete trees:\n")
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
			Name:        n.Path,
			Signature:   n.Signature,
			Assumptions: n.Assumptions,
			Refs:        n.Refs,
			Extend:      n.Extend,
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
