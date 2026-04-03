package decomposer

import (
	"testing"
)

func TestToRune(t *testing.T) {
	pr := ProposedRune{
		Name:          "test.hello",
		Description:   "says hello",
		Signature:     "(name: string) -> string",
		Behavior:      "greets the user",
		PositiveTests: []string{"returns greeting"},
		NegativeTests: []string{"rejects empty"},
		Assumptions:   []string{"ASCII only"},
		Refs:          []string{"std.util@1"},
	}

	r := pr.ToRune()
	if r.Name != "test.hello" {
		t.Errorf("Name = %q", r.Name)
	}
	if r.Description != "says hello" {
		t.Errorf("Description = %q", r.Description)
	}
	if r.Signature != "(name: string) -> string" {
		t.Errorf("Signature = %q", r.Signature)
	}
	if r.Behavior != "greets the user" {
		t.Errorf("Behavior = %q", r.Behavior)
	}
	if len(r.PositiveTests) != 1 {
		t.Errorf("PositiveTests = %v", r.PositiveTests)
	}
	if len(r.NegativeTests) != 1 {
		t.Errorf("NegativeTests = %v", r.NegativeTests)
	}
	if len(r.Assumptions) != 1 {
		t.Errorf("Assumptions = %v", r.Assumptions)
	}
	if len(r.Dependencies) != 1 || r.Dependencies[0] != "std.util@1" {
		t.Errorf("Dependencies = %v", r.Dependencies)
	}
}

func TestParseResultFromTree(t *testing.T) {
	d := &Decomposer{}
	input := `std
  std.auth
    @ (email: string) -> bool
    + validates email
    - rejects empty
  std.auth.hash
    @ (pwd: string) -> string
    + hashes password
app
  app.login
    -> std.auth
`

	result, err := d.parseResult(input)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.NewRunes) == 0 {
		t.Fatal("expected new runes")
	}

	// Reference-only nodes should become existing matches
	foundRef := false
	for _, e := range result.ExistingRunes {
		if e.Name == "std.auth" {
			foundRef = true
		}
	}
	if !foundRef {
		t.Error("expected std.auth as existing rune from -> reference")
	}

	if result.TreeOutput != input {
		t.Error("TreeOutput should preserve raw output")
	}
}
