package decompose

import (
	"encoding/json"
	"fmt"
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

// ParseDecomposition parses an LLM response into a structured Decomposition
func ParseDecomposition(responseText string) (*Decomposition, error) {
	var decomposition Decomposition

	err := json.Unmarshal([]byte(responseText), &decomposition)
	if err != nil {
		return nil, fmt.Errorf("failed to parse decomposition JSON: %w", err)
	}

	return &decomposition, nil
}

// Validate recursively validates the rune hierarchy
func (r *Rune) Validate(path string) error {
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
