package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/chrishayen/valkyrie/internal/adaptor"
	"github.com/chrishayen/valkyrie/internal/server/store"
)

// RuneProposal is a single rune returned by the pipeline.
type RuneProposal struct {
	FQN      string `json:"fqn"`
	Existing bool   `json:"existing"`
	Spec     *store.Rune `json:"spec,omitempty"`
}

// PipelineResult is the output of a completed requirements job.
type PipelineResult struct {
	Proposals []RuneProposal `json:"proposals"`
}

// Pipeline runs the 3-stage requirements decomposition.
type Pipeline struct {
	adaptor adaptor.Adaptor
	store   *store.RuneStore
}

func NewPipeline(a adaptor.Adaptor, s *store.RuneStore) *Pipeline {
	return &Pipeline{adaptor: a, store: s}
}

// Run executes the full pipeline: decompose → classify → design.
func (p *Pipeline) Run(ctx context.Context, project, requirements string) (*PipelineResult, error) {
	// Stage 1: Decompose requirements into candidate functions
	candidates, err := p.decompose(ctx, project, requirements)
	if err != nil {
		return nil, fmt.Errorf("decompose: %w", err)
	}

	// Stage 2: Classify each candidate and check for existing matches
	classified, err := p.classify(ctx, project, candidates)
	if err != nil {
		return nil, fmt.Errorf("classify: %w", err)
	}

	// Stage 3: Design specs for new runes (parallel)
	proposals, err := p.design(ctx, project, classified)
	if err != nil {
		return nil, fmt.Errorf("design: %w", err)
	}

	// Persist new runes as drafts
	for i := range proposals {
		if !proposals[i].Existing && proposals[i].Spec != nil {
			if err := p.store.Create(*proposals[i].Spec); err != nil {
				log.Printf("warning: failed to persist rune %s: %v", proposals[i].FQN, err)
			}
		}
	}

	return &PipelineResult{Proposals: proposals}, nil
}

// --- Stage 1: Decompose ---

type candidateFunc struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Pure        bool   `json:"pure"`
}

func (p *Pipeline) decompose(ctx context.Context, project, requirements string) ([]candidateFunc, error) {
	existing, _ := p.store.List("", "")
	var existingList string
	for _, r := range existing {
		existingList += fmt.Sprintf("- %s: %s\n", r.FQN, r.Description)
	}
	if existingList == "" {
		existingList = "(none)\n"
	}

	prompt := fmt.Sprintf(`You are decomposing software requirements into atomic, pure functions (runes).

## Requirements
%s

## Project name
%s

## Existing runes on the server
%s
## Instructions

Break the requirements into the smallest possible atomic pure functions. Each function should:
- Do exactly one thing
- Have no side effects (pure computation)
- Be testable in isolation

For each function, provide:
- name: a short descriptive name using snake_case
- description: one sentence explaining what it does
- pure: true if it's a pure computation, false if it requires I/O or side effects

Respond with ONLY a JSON array, no other text:
[{"name": "validate_email", "description": "Validates email format per RFC 5322", "pure": true}, ...]`, requirements, project, existingList)

	resp, err := p.adaptor.Run(ctx, prompt)
	if err != nil {
		return nil, err
	}

	resp = extractJSON(resp)

	var candidates []candidateFunc
	if err := json.Unmarshal([]byte(resp), &candidates); err != nil {
		return nil, fmt.Errorf("parsing decomposition response: %w\nraw: %s", err, resp)
	}
	return candidates, nil
}

// --- Stage 2: Classify ---

type classifiedFunc struct {
	candidateFunc
	FQN      string `json:"fqn"`
	Existing bool   `json:"existing"`
}

func (p *Pipeline) classify(ctx context.Context, project string, candidates []candidateFunc) ([]classifiedFunc, error) {
	existing, _ := p.store.List("", "")
	var existingList string
	for _, r := range existing {
		existingList += fmt.Sprintf("- %s: %s\n", r.FQN, r.Description)
	}
	if existingList == "" {
		existingList = "(none)\n"
	}

	var candidateJSON string
	for _, c := range candidates {
		candidateJSON += fmt.Sprintf("- %s: %s (pure: %v)\n", c.Name, c.Description, c.Pure)
	}

	prompt := fmt.Sprintf(`You are classifying functions into a standard-library-style namespace using dot notation.

## Project name
%s

## Functions to classify
%s
## Existing runes on the server
%s
## Classification rules

- Generic reusable functions use stdlib-style namespaces: net.http, crypto.hash, text.validate, math, time, data, io
- Project-specific functions use the project name as the top-level namespace: %s.auth, %s.payment, etc.
- Nesting depth should be sensible — usually 2-4 levels
- If a function matches an existing rune (same purpose), mark it as existing and use that FQN
- The final segment of the FQN is the function name in snake_case

For each function, provide:
- fqn: the fully-qualified dot-notation name
- existing: true if it matches an existing rune, false if new

Respond with ONLY a JSON array, no other text:
[{"name": "validate_email", "description": "...", "pure": true, "fqn": "text.validate.email", "existing": false}, ...]`, project, candidateJSON, existingList, project, project)

	resp, err := p.adaptor.Run(ctx, prompt)
	if err != nil {
		return nil, err
	}

	resp = extractJSON(resp)

	var classified []classifiedFunc
	if err := json.Unmarshal([]byte(resp), &classified); err != nil {
		return nil, fmt.Errorf("parsing classification response: %w\nraw: %s", err, resp)
	}
	return classified, nil
}

// --- Stage 3: Design ---

func (p *Pipeline) design(ctx context.Context, project string, classified []classifiedFunc) ([]RuneProposal, error) {
	proposals := make([]RuneProposal, len(classified))
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error

	for i, c := range classified {
		if c.Existing {
			proposals[i] = RuneProposal{FQN: c.FQN, Existing: true}
			continue
		}

		wg.Add(1)
		go func(idx int, cf classifiedFunc) {
			defer wg.Done()

			spec, err := p.designOne(ctx, project, cf)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("designing %s: %w", cf.FQN, err)
				}
				proposals[idx] = RuneProposal{FQN: cf.FQN, Existing: false}
				return
			}
			proposals[idx] = RuneProposal{FQN: cf.FQN, Existing: false, Spec: spec}
		}(i, c)
	}
	wg.Wait()

	if firstErr != nil {
		return proposals, firstErr
	}
	return proposals, nil
}

func (p *Pipeline) designOne(ctx context.Context, project string, cf classifiedFunc) (*store.Rune, error) {
	prompt := fmt.Sprintf(`You are designing a rune (pure function specification).

## Function
- FQN: %s
- Name: %s
- Description: %s
- Pure: %v

## Instructions

Write a complete rune specification. Include:
1. A precise typed signature using these types: string, bool, bytes, i8/i16/i32/i64, u8/u16/u32/u64, f32/f64, list[T], map[K,V], optional[T], result[T,E]
2. Behavior description covering inputs, outputs, edge cases, and constraints
3. At least 2 positive test cases (expected successful behavior)
4. At least 2 negative test cases (expected error/failure behavior)

Respond with ONLY JSON, no other text:
{
  "fqn": "%s",
  "description": "...",
  "signature": "(param: type, ...) -> return_type",
  "behavior": "- ...\n- ...",
  "positive_tests": ["test case 1", "test case 2"],
  "negative_tests": ["test case 1", "test case 2"]
}`, cf.FQN, cf.Name, cf.Description, cf.Pure, cf.FQN)

	resp, err := p.adaptor.Run(ctx, prompt)
	if err != nil {
		return nil, err
	}

	resp = extractJSON(resp)

	var spec store.Rune
	if err := json.Unmarshal([]byte(resp), &spec); err != nil {
		return nil, fmt.Errorf("parsing design response: %w\nraw: %s", err, resp)
	}

	// Ensure correct metadata
	spec.FQN = cf.FQN
	spec.Status = "draft"
	spec.Version = "0.1.0"
	if isProjectRune(cf.FQN, project) {
		spec.Project = project
	}

	return &spec, nil
}

// extractJSON strips markdown code fences and surrounding text from a JSON response.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)

	// Try to extract from code fences
	if idx := strings.Index(s, "```json"); idx != -1 {
		s = s[idx+7:]
		if end := strings.Index(s, "```"); end != -1 {
			s = s[:end]
		}
	} else if idx := strings.Index(s, "```"); idx != -1 {
		s = s[idx+3:]
		if end := strings.Index(s, "```"); end != -1 {
			s = s[:end]
		}
	}

	s = strings.TrimSpace(s)

	// If it doesn't start with [ or {, try to find the first JSON structure
	if len(s) > 0 && s[0] != '[' && s[0] != '{' {
		if idx := strings.IndexAny(s, "[{"); idx != -1 {
			s = s[idx:]
		}
	}

	// Trim trailing non-JSON
	if len(s) > 0 {
		if s[0] == '[' {
			if idx := strings.LastIndex(s, "]"); idx != -1 {
				s = s[:idx+1]
			}
		} else if s[0] == '{' {
			if idx := strings.LastIndex(s, "}"); idx != -1 {
				s = s[:idx+1]
			}
		}
	}

	return s
}

func isProjectRune(fqn, project string) bool {
	return project != "" && strings.HasPrefix(fqn, project+".")
}
