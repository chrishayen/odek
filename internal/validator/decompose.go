package validator

import (
	"fmt"
	"strings"

	runepkg "github.com/chrishayen/odek/internal/rune"
)

const decomposeSemanticPrompt = `You validate the semantic quality of a decomposition tree. The structural format has already been verified — do NOT check signatures, test case counts, naming format, or package structure. Focus ONLY on:

1. No two sibling runes should have overlapping or duplicate responsibilities. Siblings are runes that share the same parent (e.g. std.text.wrap and std.text.pad_right are siblings).
2. Each leaf rune should represent a single, focused responsibility — flag any that try to do too much.
3. References (->) should be semantically reasonable (the referenced rune makes sense for the referencing context). References may point to runes in an external registry not shown in the tree.

If all checks pass, respond with exactly:
RESULT: PASS

If any checks fail, list each issue on its own line, then end with:
RESULT: FAIL

No other output.`

// ValidateDecomposition checks the raw tree output from decomposition.
// Structural rules are checked deterministically; semantic rules use the LLM.
func (v *Validator) ValidateDecomposition(treeOutput string) (*Result, error) {
	// Deterministic structural checks.
	if r := validateStructure(treeOutput); !r.Passed {
		return r, nil
	}

	// Semantic checks via LLM.
	resp, err := v.client.Call(decomposeSemanticPrompt, treeOutput)
	if err != nil {
		// Structural validation passed; don't fail on LLM error.
		return &Result{Passed: true}, nil
	}
	return parseResult(resp), nil
}

// validateStructure runs deterministic checks on the tree output.
func validateStructure(treeOutput string) *Result {
	nodes := runepkg.ParseTree(treeOutput)
	if len(nodes) == 0 {
		return &Result{Issues: []string{"no runes found in tree output"}}
	}

	allNames := make([]string, len(nodes))
	for i, n := range nodes {
		allNames[i] = n.Path
	}

	var issues []string
	for _, n := range nodes {
		leaf := runepkg.IsLeaf(n.Path, allNames)

		if leaf {
			if n.Signature == "" {
				issues = append(issues, fmt.Sprintf("%s: leaf rune missing signature (@ line)", n.Path))
			}
			if len(n.Pos) == 0 {
				issues = append(issues, fmt.Sprintf("%s: leaf rune missing positive test case (+ line)", n.Path))
			}
			if len(n.Neg) == 0 {
				issues = append(issues, fmt.Sprintf("%s: leaf rune missing negative test case (- line)", n.Path))
			}
		} else if n.Signature != "" {
			issues = append(issues, fmt.Sprintf("%s: package rune should not have a signature", n.Path))
		}
	}

	if len(issues) == 0 {
		return &Result{Passed: true}
	}
	return &Result{Issues: issues}
}

// FormatDecompositionFeedback formats validation issues as refinement instructions
// for the decomposer's prevDecomposition path.
func FormatDecompositionFeedback(issues []string) string {
	return fmt.Sprintf(
		"The validator found the following issues:\n%s\n\nFix these issues and re-output the complete trees.",
		strings.Join(issues, "\n"),
	)
}
