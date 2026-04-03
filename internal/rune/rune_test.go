package rune

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSemverString(t *testing.T) {
	tests := []struct {
		v    Semver
		want string
	}{
		{Semver{1, 0, 0}, "1.0.0"},
		{Semver{2, 3, 4}, "2.3.4"},
		{Semver{0, 0, 0}, "0.0.0"},
	}
	for _, tt := range tests {
		if got := tt.v.String(); got != tt.want {
			t.Errorf("Semver%v.String() = %q, want %q", tt.v, got, tt.want)
		}
	}
}

func TestSemverIsZero(t *testing.T) {
	zero := Semver{}
	if !zero.IsZero() {
		t.Error("zero Semver should be zero")
	}
	nonzero := Semver{1, 0, 0}
	if nonzero.IsZero() {
		t.Error("1.0.0 should not be zero")
	}
}

func TestParseSemver(t *testing.T) {
	tests := []struct {
		input string
		want  Semver
	}{
		{"1.0.0", Semver{1, 0, 0}},
		{"2.3.4", Semver{2, 3, 4}},
		{"0.0.0", Semver{0, 0, 0}},
		{"garbage", Semver{}},
		{"", Semver{}},
	}
	for _, tt := range tests {
		got := ParseSemver(tt.input)
		if got != tt.want {
			t.Errorf("ParseSemver(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestSemverLess(t *testing.T) {
	tests := []struct {
		a, b Semver
		want bool
	}{
		{Semver{1, 0, 0}, Semver{2, 0, 0}, true},
		{Semver{1, 0, 0}, Semver{1, 1, 0}, true},
		{Semver{1, 0, 0}, Semver{1, 0, 1}, true},
		{Semver{2, 0, 0}, Semver{1, 0, 0}, false},
		{Semver{1, 0, 0}, Semver{1, 0, 0}, false},
	}
	for _, tt := range tests {
		if got := tt.a.Less(tt.b); got != tt.want {
			t.Errorf("%v.Less(%v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSemverBump(t *testing.T) {
	v := Semver{1, 2, 3}
	if got := v.BumpMajor(); got != (Semver{2, 0, 0}) {
		t.Errorf("BumpMajor() = %v", got)
	}
	if got := v.BumpMinor(); got != (Semver{1, 3, 0}) {
		t.Errorf("BumpMinor() = %v", got)
	}
	if got := v.BumpPatch(); got != (Semver{1, 2, 4}) {
		t.Errorf("BumpPatch() = %v", got)
	}
}

func TestSemverJSON(t *testing.T) {
	v := Semver{1, 2, 3}
	data, err := v.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"1.2.3"` {
		t.Errorf("MarshalJSON() = %s", data)
	}

	var v2 Semver
	if err := v2.UnmarshalJSON(data); err != nil {
		t.Fatal(err)
	}
	if v2 != v {
		t.Errorf("UnmarshalJSON round-trip: got %v, want %v", v2, v)
	}
}

func TestIsSemverFilename(t *testing.T) {
	tests := []struct {
		name string
		ok   bool
		ver  Semver
	}{
		{"1.0.0.md", true, Semver{1, 0, 0}},
		{"2.3.4.md", true, Semver{2, 3, 4}},
		{"feature.md", false, Semver{}},
		{"app.md", false, Semver{}},
		{"1.0.md", false, Semver{}},
		{"not-a-version.md", false, Semver{}},
	}
	for _, tt := range tests {
		v, ok := IsSemverFilename(tt.name)
		if ok != tt.ok {
			t.Errorf("IsSemverFilename(%q) ok = %v, want %v", tt.name, ok, tt.ok)
		}
		if ok && v != tt.ver {
			t.Errorf("IsSemverFilename(%q) = %v, want %v", tt.name, v, tt.ver)
		}
	}
}

func TestLatestVersion(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "1.0.0.md"), []byte("---\nversion: 1.0.0\n---\n"), 0644)
	os.WriteFile(filepath.Join(dir, "1.1.0.md"), []byte("---\nversion: 1.1.0\n---\n"), 0644)
	os.WriteFile(filepath.Join(dir, "2.0.0.md"), []byte("---\nversion: 2.0.0\n---\n"), 0644)
	os.WriteFile(filepath.Join(dir, "feature.md"), []byte("not a version"), 0644)

	v, found := LatestVersion(dir)
	if !found {
		t.Fatal("expected to find a version")
	}
	if v != (Semver{2, 0, 0}) {
		t.Errorf("LatestVersion() = %v, want 2.0.0", v)
	}
}

func TestLatestVersionEmpty(t *testing.T) {
	dir := t.TempDir()
	_, found := LatestVersion(dir)
	if found {
		t.Error("expected no version in empty dir")
	}
}

func TestIsDotPath(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"std.auth.validate_email", true},
		{"auth.login", true},
		{"a.b.c.d", true},
		{"a1.b2", true},
		{"", false},
		{".leading.dot", false},
		{"trailing.dot.", false},
		{"has spaces", false},
		{"HAS_UPPER", false},
		{"has/slash", false},
	}
	for _, tt := range tests {
		if got := IsDotPath(tt.input); got != tt.want {
			t.Errorf("IsDotPath(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseTree(t *testing.T) {
	input := `std
  std.auth
    @ (email: string) -> bool
    + validates email format
    - rejects empty string
    -> std.util
    ? assumes ASCII only
  std.auth.hash
    @ (password: string) -> string
    + hashes password
`
	nodes := ParseTree(input)
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}

	// std
	if nodes[0].Path != "std" {
		t.Errorf("node 0 path = %q", nodes[0].Path)
	}

	// std.auth
	auth := nodes[1]
	if auth.Path != "std.auth" {
		t.Errorf("node 1 path = %q", auth.Path)
	}
	if auth.Signature != "(email: string) -> bool" {
		t.Errorf("signature = %q", auth.Signature)
	}
	if len(auth.Pos) != 1 || auth.Pos[0] != "validates email format" {
		t.Errorf("pos = %v", auth.Pos)
	}
	if len(auth.Neg) != 1 || auth.Neg[0] != "rejects empty string" {
		t.Errorf("neg = %v", auth.Neg)
	}
	if len(auth.Refs) != 1 || auth.Refs[0] != "std.util" {
		t.Errorf("refs = %v", auth.Refs)
	}
	if len(auth.Assumptions) != 1 || auth.Assumptions[0] != "assumes ASCII only" {
		t.Errorf("assumptions = %v", auth.Assumptions)
	}

	// std.auth.hash
	if nodes[2].Signature != "(password: string) -> string" {
		t.Errorf("node 2 signature = %q", nodes[2].Signature)
	}
}

func TestParseTreeExtend(t *testing.T) {
	input := `~> std.auth
  @ (email: string) -> bool
  + new test case`
	nodes := ParseTree(input)
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	if !nodes[0].Extend {
		t.Error("expected Extend = true")
	}
	if nodes[0].Path != "std.auth" {
		t.Errorf("path = %q", nodes[0].Path)
	}
}

func TestBuildChildrenMap(t *testing.T) {
	paths := []string{"std", "std.auth", "std.auth.validate", "std.cli"}
	children := BuildChildrenMap(paths)

	if len(children["std"]) != 2 {
		t.Errorf("std children = %v", children["std"])
	}
	if len(children["std.auth"]) != 1 || children["std.auth"][0] != "std.auth.validate" {
		t.Errorf("std.auth children = %v", children["std.auth"])
	}
}

func TestParseRef(t *testing.T) {
	tests := []struct {
		input string
		path  string
		major int
	}{
		{"std.auth@1", "std.auth", 1},
		{"std.cli@2", "std.cli", 2},
		{"no-at-sign", "", 0},
	}
	for _, tt := range tests {
		path, major := ParseRef(tt.input)
		if path != tt.path || major != tt.major {
			t.Errorf("ParseRef(%q) = (%q, %d), want (%q, %d)", tt.input, path, major, tt.path, tt.major)
		}
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		r    Rune
		ok   bool
	}{
		{"valid", Rune{Name: "auth.login", Description: "desc", Signature: "sig"}, true},
		{"no name", Rune{Description: "desc", Signature: "sig"}, false},
		{"no namespace", Rune{Name: "login", Description: "desc", Signature: "sig"}, false},
		{"no desc", Rune{Name: "auth.login", Signature: "sig"}, false},
		{"no sig", Rune{Name: "auth.login", Description: "desc"}, false},
		{"bad path", Rune{Name: "HAS CAPS", Description: "desc", Signature: "sig"}, false},
	}
	for _, tt := range tests {
		err := validate(tt.r)
		if (err == nil) != tt.ok {
			t.Errorf("%s: validate() error = %v, wantOK = %v", tt.name, err, tt.ok)
		}
	}
}

func TestStoreCreateGetDeleteRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir, filepath.Join(dir, "src"))

	r := Rune{
		Name:          "test.hello",
		Description:   "says hello",
		Signature:     "(name: string) -> string",
		PositiveTests: []string{"returns greeting"},
		NegativeTests: []string{"rejects empty"},
	}

	if err := s.Create(r); err != nil {
		t.Fatal(err)
	}

	got, err := s.Get("test.hello")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "test.hello" {
		t.Errorf("Name = %q", got.Name)
	}
	if got.Description != "says hello" {
		t.Errorf("Description = %q", got.Description)
	}
	if got.Version != (Semver{1, 0, 0}) {
		t.Errorf("Version = %v", got.Version)
	}
	if len(got.PositiveTests) != 1 {
		t.Errorf("PositiveTests = %v", got.PositiveTests)
	}

	// Duplicate
	if err := s.Create(r); err == nil {
		t.Error("expected error creating duplicate")
	}

	// List
	runes, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(runes) != 1 {
		t.Errorf("List() len = %d", len(runes))
	}

	// Update
	got.Description = "updated"
	got.Version = got.Version.BumpMinor()
	if err := s.Update(*got); err != nil {
		t.Fatal(err)
	}
	updated, _ := s.Get("test.hello")
	if updated.Description != "updated" {
		t.Errorf("after update: Description = %q", updated.Description)
	}

	// Delete
	if err := s.Delete("test.hello"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("test.hello"); err == nil {
		t.Error("expected error after delete")
	}
}

func TestStoreGetNotFound(t *testing.T) {
	s := NewStore(t.TempDir(), "")
	_, err := s.Get("no.such.rune")
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestParseRuneContent(t *testing.T) {
	content := `---
version: 1.2.0
hydrated: true
coverage: 85.5
signature: '(x: int) -> bool'
dependencies:
  - std.util@1
---

# test.my_rune

A test rune

## Signature

(x: int) -> bool

## Behavior

Does things

## Tests

+ positive case
- negative case

## Assumptions

? assumes something
`
	r, err := parse(content)
	if err != nil {
		t.Fatal(err)
	}
	if r.Version != (Semver{1, 2, 0}) {
		t.Errorf("Version = %v", r.Version)
	}
	if !r.Hydrated {
		t.Error("expected Hydrated = true")
	}
	if r.Coverage != 85.5 {
		t.Errorf("Coverage = %v", r.Coverage)
	}
	if r.Name != "test.my_rune" {
		t.Errorf("Name = %q", r.Name)
	}
	if r.Description != "A test rune" {
		t.Errorf("Description = %q", r.Description)
	}
	if r.Behavior != "Does things" {
		t.Errorf("Behavior = %q", r.Behavior)
	}
	if len(r.PositiveTests) != 1 || r.PositiveTests[0] != "positive case" {
		t.Errorf("PositiveTests = %v", r.PositiveTests)
	}
	if len(r.NegativeTests) != 1 || r.NegativeTests[0] != "negative case" {
		t.Errorf("NegativeTests = %v", r.NegativeTests)
	}
	if len(r.Assumptions) != 1 || r.Assumptions[0] != "assumes something" {
		t.Errorf("Assumptions = %v", r.Assumptions)
	}
	if len(r.Dependencies) != 1 || r.Dependencies[0] != "std.util@1" {
		t.Errorf("Dependencies = %v", r.Dependencies)
	}
}
