package hydrator

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	runepkg "github.com/chrishayen/valkyrie/internal/rune"
)

const VerifySystemPrompt = `You are a code reviewer. You receive a spec (with test descriptions) and an implementation. Your job is to verify the implementation satisfies every + and - test case in the spec.

For each test case in the spec, output one line:
  PASS + <test description> — <brief reason>
  PASS - <test description> — <brief reason>
  FAIL + <test description> — <what's wrong>
  FAIL - <test description> — <what's wrong>

If all pass, end with: RESULT: ALL PASS
If any fail, end with: RESULT: <N> FAILURES

No other output.`

// VerifyResult holds the outcome of verifying a rune.
type VerifyResult struct {
	RuneName string `json:"rune_name"`
	Passed   bool   `json:"passed"`
	Output   string `json:"output"`
}

// VerifyAllResult holds the outcome of batch verification.
type VerifyAllResult struct {
	Passed  int            `json:"passed"`
	Failed  int            `json:"failed"`
	Details []VerifyResult `json:"details"`
}

// Verify checks generated code for a rune against its spec.
func (h *Hydrator) Verify(_ context.Context, name string) (*VerifyResult, error) {
	rn, err := h.store.Get(name)
	if err != nil {
		return nil, err
	}

	codeDir := h.store.CodeDir(name)
	var implContent strings.Builder
	entries, err := os.ReadDir(codeDir)
	if err != nil {
		return nil, fmt.Errorf("reading code dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || strings.HasSuffix(e.Name(), "_test.go") || strings.HasSuffix(e.Name(), ".test.ts") || strings.HasSuffix(e.Name(), "test.py") {
			continue
		}
		data, err := os.ReadFile(codeDir + "/" + e.Name())
		if err != nil {
			continue
		}
		fmt.Fprintf(&implContent, "--- %s ---\n%s\n", e.Name(), string(data))
	}

	var specText strings.Builder
	fmt.Fprintf(&specText, "# %s\n\n%s\n\n## Signature\n\n%s\n", rn.Name, rn.Description, rn.Signature)
	if rn.Behavior != "" {
		fmt.Fprintf(&specText, "\n## Behavior\n\n%s\n", rn.Behavior)
	}
	if len(rn.PositiveTests) > 0 || len(rn.NegativeTests) > 0 {
		specText.WriteString("\n## Tests\n\n")
		for _, t := range rn.PositiveTests {
			fmt.Fprintf(&specText, "+ %s\n", t)
		}
		for _, t := range rn.NegativeTests {
			fmt.Fprintf(&specText, "- %s\n", t)
		}
	}

	userPrompt := fmt.Sprintf("Spec:\n%s\n\nImplementation:\n%s", specText.String(), implContent.String())

	output, err := h.client.Call(VerifySystemPrompt, userPrompt)
	if err != nil {
		return &VerifyResult{RuneName: name, Passed: false, Output: "claude error: " + err.Error()}, nil
	}

	passed := strings.Contains(output, "RESULT: ALL PASS")
	return &VerifyResult{
		RuneName: name,
		Passed:   passed,
		Output:   strings.TrimSpace(output),
	}, nil
}

// VerifyAll verifies all hydrated runes concurrently.
func (h *Hydrator) VerifyAll(ctx context.Context, concurrency int, logOut fmt.Stringer) (*VerifyAllResult, error) {
	runes, err := h.store.List()
	if err != nil {
		return nil, err
	}

	var targets []runepkg.Rune
	for _, rn := range runes {
		if rn.Hydrated {
			targets = append(targets, rn)
		}
	}

	result := &VerifyAllResult{}
	if len(targets) == 0 {
		return result, nil
	}

	type vr struct {
		res VerifyResult
	}
	results := make([]vr, len(targets))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, rn := range targets {
		wg.Add(1)
		go func(idx int, name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			v, err := h.Verify(ctx, name)
			if err != nil {
				results[idx] = vr{VerifyResult{RuneName: name, Passed: false, Output: err.Error()}}
				return
			}
			results[idx] = vr{*v}
		}(i, rn.Name)
	}

	wg.Wait()

	for _, vres := range results {
		result.Details = append(result.Details, vres.res)
		if vres.res.Passed {
			result.Passed++
		} else {
			result.Failed++
		}
	}

	return result, nil
}
