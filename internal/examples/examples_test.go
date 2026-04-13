package examples

import (
	"path/filepath"
	"strings"
	"testing"
)

func examplesRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..", "examples")
}

func TestLoadFromDir_LoadsAllTiers(t *testing.T) {
	idx, err := LoadFromDir(examplesRoot(t))
	if err != nil {
		t.Fatalf("LoadFromDir: %v", err)
	}
	if idx.Len() < 100 {
		t.Fatalf("expected at least 100 entries, got %d", idx.Len())
	}
	counts := map[string]int{}
	for _, e := range idx.Entries {
		counts[e.Tier]++
	}
	for _, tier := range tiers {
		if counts[tier] == 0 {
			t.Errorf("tier %q produced 0 entries", tier)
		}
	}
}

func TestLoadFromDir_ParsesRequirement(t *testing.T) {
	idx, err := LoadFromDir(examplesRoot(t))
	if err != nil {
		t.Fatalf("LoadFromDir: %v", err)
	}
	for _, e := range idx.Entries {
		if e.Requirement == "" {
			t.Errorf("%s: empty Requirement", e.Path)
		}
		if e.Content == "" {
			t.Errorf("%s: empty Content", e.Path)
		}
	}
}

func TestHandle_ShapeIsTierSlashSlug(t *testing.T) {
	e := &Entry{Tier: "medium", Slug: "csv-reader"}
	if got := e.Handle(); got != "medium/csv-reader" {
		t.Errorf("Handle = %q, want %q", got, "medium/csv-reader")
	}
}

func TestLookup_AcceptsSeveralForms(t *testing.T) {
	idx, err := LoadFromDir(examplesRoot(t))
	if err != nil {
		t.Fatalf("LoadFromDir: %v", err)
	}
	if idx.Len() == 0 {
		t.Skip("corpus is empty")
	}
	target := idx.Entries[0]
	handle := target.Handle()

	cases := []string{
		handle,                       // "tier/slug"
		handle + ".md",               // "tier/slug.md"
		"examples/" + handle,         // "examples/tier/slug"
		"examples/" + handle + ".md", // "examples/tier/slug.md"
		"  " + handle + "  ",         // whitespace
	}
	for _, c := range cases {
		res := idx.Lookup(c)
		if res.Kind != LookupHit {
			t.Errorf("Lookup(%q): kind = %v, want LookupHit", c, res.Kind)
			continue
		}
		if res.Entry.Handle() != handle {
			t.Errorf("Lookup(%q) = %s, want %s", c, res.Entry.Handle(), handle)
		}
	}

	miss := idx.Lookup("no-such-tier/no-such-slug")
	if miss.Kind != LookupMiss {
		t.Errorf("Lookup of nonexistent handle: kind = %v, want LookupMiss", miss.Kind)
	}
}

func TestLookup_TierCorrectsWhenSlugUnique(t *testing.T) {
	idx, err := LoadFromDir(examplesRoot(t))
	if err != nil {
		t.Fatalf("LoadFromDir: %v", err)
	}
	if idx.Len() == 0 {
		t.Skip("corpus is empty")
	}
	// Find a slug that lives in exactly one tier, then ask for it with the
	// wrong tier prefix.
	var unique *Entry
	for i := range idx.Entries {
		e := &idx.Entries[i]
		if len(idx.bySlug[e.Slug]) == 1 {
			unique = e
			break
		}
	}
	if unique == nil {
		t.Skip("no slug with a single tier found")
	}
	// Pick a different tier than the real one.
	wrongTier := "medium"
	if unique.Tier == "medium" {
		wrongTier = "large"
	}
	res := idx.Lookup(wrongTier + "/" + unique.Slug)
	if res.Kind != LookupTierCorrected {
		t.Errorf("expected LookupTierCorrected for wrong-tier query, got kind=%v", res.Kind)
	}
	if res.Entry != nil && res.Entry.Handle() != unique.Handle() {
		t.Errorf("tier correction pointed at %s, want %s", res.Entry.Handle(), unique.Handle())
	}
}

func TestLookup_MissReturnsSuggestions(t *testing.T) {
	idx, err := LoadFromDir(examplesRoot(t))
	if err != nil {
		t.Fatalf("LoadFromDir: %v", err)
	}
	if idx.Len() == 0 {
		t.Skip("corpus is empty")
	}
	// Use a handle with a plausible-looking slug that likely doesn't exist.
	res := idx.Lookup("medium/http-server-framework-with-routing-and-middleware")
	if res.Kind != LookupMiss {
		// It might exist; if so skip. If not, verify we got suggestions.
		t.Skipf("handle unexpectedly resolved (kind=%v); cannot test suggestions", res.Kind)
	}
	if len(res.Suggestions) == 0 {
		t.Error("expected at least one suggestion for an http-ish missing handle")
	}
}

func TestManifest_ContainsAllEntriesGroupedByTier(t *testing.T) {
	idx, err := LoadFromDir(examplesRoot(t))
	if err != nil {
		t.Fatalf("LoadFromDir: %v", err)
	}
	m := idx.Manifest()
	if m == "" {
		t.Fatal("Manifest is empty")
	}
	// Every entry handle should appear somewhere in the manifest.
	for _, e := range idx.Entries {
		if !strings.Contains(m, "- "+e.Handle()) {
			t.Errorf("handle %q missing from manifest", e.Handle())
			break // one failure is enough
		}
	}
	// Tier headings should appear in the expected order, when present.
	var last int
	for _, tier := range tiers {
		heading := "## " + tier
		pos := strings.Index(m, heading)
		if pos == -1 {
			continue
		}
		if pos < last {
			t.Errorf("tier %q heading appears before previous tier", tier)
		}
		last = pos
	}
}
