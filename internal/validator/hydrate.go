package validator

import (
	"fmt"
	"strings"

	"github.com/chrishayen/odek/internal/llm"
	runepkg "github.com/chrishayen/odek/internal/rune"
)

const hydrateValidationPrompt = `You validate generated code for an isolated library component. Given the rune spec and the raw LLM output containing === FILE: === blocks, check ALL of the following rules:

1. No config files: output must NOT contain package.json, tsconfig.json, vitest.config.ts, jest.config.ts, jest.config.js, go.mod, go.sum, setup.py, pyproject.toml, or similar project config files.
2. No subdirectories: all filenames in === FILE: === blocks must be plain filenames with no path separators (no src/, no nested folders).
3. No cross-rune imports: the code must not import or reference other runes directly.
4. No main() functions or CLI scaffolding: this is a library component, not an executable.
5. Correct file naming: files should use the short name provided in the spec.
6. Test coverage: every positive and negative test case from the spec should have a corresponding test function.
7. Files use the correct === FILE: <filename> === / === END FILE === format.

If all checks pass, respond with exactly:
RESULT: PASS

If any checks fail, list each issue on its own line, then end with:
RESULT: FAIL

No other output.`

// ValidateHydration checks raw LLM hydration output against the rune spec.
func (v *Validator) ValidateHydration(rn *runepkg.Rune, output, language string) (*Result, error) {
	shortName := runepkg.ShortName(rn.Name)
	var sb strings.Builder
	fmt.Fprintf(&sb, "Rune: %s (short name: %s)\nLanguage: %s\n", rn.Name, shortName, language)
	fmt.Fprintf(&sb, "Signature: %s\n", rn.Signature)
	if len(rn.PositiveTests) > 0 {
		sb.WriteString("\nPositive tests:\n")
		for _, t := range rn.PositiveTests {
			fmt.Fprintf(&sb, "+ %s\n", t)
		}
	}
	if len(rn.NegativeTests) > 0 {
		sb.WriteString("\nNegative tests:\n")
		for _, t := range rn.NegativeTests {
			fmt.Fprintf(&sb, "- %s\n", t)
		}
	}
	fmt.Fprintf(&sb, "\n--- Generated Output ---\n%s", output)

	resp, err := v.client.Call(hydrateValidationPrompt, sb.String())
	if err != nil {
		return nil, fmt.Errorf("validation call failed: %w", err)
	}
	return parseResult(resp), nil
}

// BuildRetryMessages constructs a multi-turn conversation that feeds validation
// issues back to the hydration LLM for correction.
func (v *Validator) BuildRetryMessages(originalPrompt, previousOutput string, issues []string) []llm.ChatMessage {
	return []llm.ChatMessage{
		{Role: "user", Content: originalPrompt},
		{Role: "assistant", Content: previousOutput},
		{Role: "user", Content: fmt.Sprintf(
			"The validator found the following issues with your output:\n%s\n\nPlease regenerate the complete output, fixing all issues listed above. Output all files using the === FILE: === / === END FILE === format.",
			strings.Join(issues, "\n"),
		)},
	}
}
