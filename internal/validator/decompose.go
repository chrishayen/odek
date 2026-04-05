package validator

import (
	"fmt"
	"strings"
)

const decomposeValidationPrompt = `You validate a decomposition tree for a software project. The tree uses an indented format where each node is a rune (isolated library function). Check ALL of the following rules:

1. Every leaf rune must have a signature line (@ prefix).
2. Every leaf rune must have at least one positive test case (+ prefix) and at least one negative test case (- prefix).
3. No two sibling runes should have overlapping or duplicate responsibilities.
4. Names must follow snake_case dot-path conventions (e.g. std.auth.validate_email).
5. References (->) must point to runes that exist elsewhere in the tree.
6. Package (non-leaf) runes should NOT have signatures — they are organizational groupings.

If all checks pass, respond with exactly:
RESULT: PASS

If any checks fail, list each issue on its own line, then end with:
RESULT: FAIL

No other output.`

// ValidateDecomposition checks the raw tree output from decomposition.
func (v *Validator) ValidateDecomposition(treeOutput string) (*Result, error) {
	resp, err := v.client.Call(decomposeValidationPrompt, treeOutput)
	if err != nil {
		return nil, fmt.Errorf("validation call failed: %w", err)
	}
	return parseResult(resp), nil
}

// FormatDecompositionFeedback formats validation issues as refinement instructions
// for the decomposer's prevDecomposition path.
func FormatDecompositionFeedback(issues []string) string {
	return fmt.Sprintf(
		"The validator found the following issues:\n%s\n\nFix these issues and re-output the complete trees.",
		strings.Join(issues, "\n"),
	)
}
