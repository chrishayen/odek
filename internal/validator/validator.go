package validator

import (
	"strings"

	"github.com/chrishayen/odek/internal/llm"
)

// Validator runs LLM-based validation on decomposition and hydration output.
type Validator struct {
	client     *llm.Client
	maxRetries int
}

// Result holds the outcome of a validation check.
type Result struct {
	Passed bool     `json:"passed"`
	Issues []string `json:"issues,omitempty"`
}

// New creates a Validator.
func New(client *llm.Client, maxRetries int) *Validator {
	return &Validator{client: client, maxRetries: maxRetries}
}

// MaxRetries returns the configured retry limit.
func (v *Validator) MaxRetries() int {
	return v.maxRetries
}

// parseResult parses an LLM validation response into a Result.
// Expected format: lines of issues, ending with "RESULT: PASS" or "RESULT: FAIL".
func parseResult(output string) *Result {
	if strings.Contains(output, "RESULT: PASS") {
		return &Result{Passed: true}
	}
	var issues []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "RESULT:") {
			continue
		}
		issues = append(issues, line)
	}
	return &Result{Passed: false, Issues: issues}
}
