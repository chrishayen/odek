package decompose

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rune represents a single hierarchical unit in the decomposition
type Rune struct {
	Path         string   `json:"path"`
	Version      string   `json:"version"`
	Signature    string   `json:"signature"`
	Dependencies []string `json:"dependencies,omitempty"` // Both std and sibling rune paths
	Tests        []Test   `json:"tests,omitempty"`        // Test cases for this rune
	Assumptions  []string `json:"assumptions,omitempty"`
	Description  string   `json:"description,omitempty"`
	Children     []*Rune  `json:"children,omitempty"` // Child runes (nested hierarchy)
}

// Test represents a test case for a rune
type Test struct {
	Name       string      `json:"name"`                 // Test name/description
	Input      interface{} `json:"input,omitempty"`      // Input data/parameters
	Expected   interface{} `json:"expected,omitempty"`   // Expected output/result
	Conditions []string    `json:"conditions,omitempty"` // Pre-conditions or setup requirements
}

// Decomposition represents the complete hierarchical decomposition of a feature
type Decomposition struct {
	FeatureName string `json:"feature_name"`          // Name of the feature being decomposed
	Description string `json:"description,omitempty"` // High-level description of the feature
	RuneTree    *Rune  `json:"rune_tree"`             // Root of the hierarchical rune structure
}

// stripMarkdownJSON removes markdown code fences from JSON strings
func stripMarkdownJSON(input string) string {
	input = strings.TrimSpace(input)

	if strings.HasPrefix(input, "```") {
		openEnd := strings.IndexByte(input, '\n')
		if openEnd == -1 {
			return input
		}
		content := strings.TrimSpace(input[openEnd+1:])
		if strings.HasSuffix(content, "```") {
			content = strings.TrimSpace(strings.TrimSuffix(content, "```"))
		}
		return content
	}

	return input
}

// ParseDecomposition parses an LLM response into a structured Decomposition
func ParseDecomposition(responseText string) (*Decomposition, error) {
	var decomposition Decomposition

	// Strip markdown code fences if present
	cleaned := stripMarkdownJSON(responseText)

	err := json.Unmarshal([]byte(cleaned), &decomposition)
	if err != nil {
		return nil, fmt.Errorf("failed to parse decomposition JSON: %w", err)
	}

	return &decomposition, nil
}

// Validate recursively validates the rune hierarchy
func (r *Rune) Validate(path string) error {
	if r == nil {
		return fmt.Errorf("rune at path '%s' is nil", path)
	}
	if r.Path == "" {
		return fmt.Errorf("rune at path '%s' has empty Path field", path)
	}
	if r.Version == "" {
		return fmt.Errorf("rune '%s' has empty Version field", r.Path)
	}
	if r.Signature == "" {
		return fmt.Errorf("rune '%s' has empty Signature field", r.Path)
	}

	for i, child := range r.Children {
		childPath := fmt.Sprintf("%s.child[%d]", path, i)
		if err := child.Validate(childPath); err != nil {
			return err
		}
	}

	return nil
}

// FormatJSON returns the decomposition as formatted JSON string
func (d *Decomposition) FormatJSON() (string, error) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal decomposition: %w", err)
	}
	return string(data), nil
}
