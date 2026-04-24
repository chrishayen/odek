package decompose

import "testing"

func TestParseDecompositionStripsMarkdownFence(t *testing.T) {
	input := "```json\n" + `{
  "feature_name": "demo",
  "rune_tree": {
    "path": "demo.root",
    "version": "1.0.0",
    "signature": "() -> string"
  }
}` + "\n```"

	got, err := ParseDecomposition(input)
	if err != nil {
		t.Fatalf("ParseDecomposition: %v", err)
	}
	if got.FeatureName != "demo" {
		t.Fatalf("FeatureName = %q, want demo", got.FeatureName)
	}
	if got.RuneTree == nil || got.RuneTree.Path != "demo.root" {
		t.Fatalf("RuneTree = %#v, want path demo.root", got.RuneTree)
	}
}

func TestRuneValidateRejectsNilRune(t *testing.T) {
	var r *Rune
	if err := r.Validate("root"); err == nil {
		t.Fatal("Validate(nil) returned nil error")
	}
}
