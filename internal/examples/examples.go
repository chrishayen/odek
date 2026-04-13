// Package examples loads the odek decomposition example corpus and provides
// direct lookup of entries by their "tier/slug" handle.
//
// The corpus lives at examples/{trivial,small,medium,large}/*.md. Each file
// starts with a `# Requirement: "..."` header followed by a text-DSL
// decomposition. The package parses those into Entry values, builds a
// manifest suitable for inlining into a system prompt, and lets callers
// retrieve full content by handle.
//
// There is no keyword search. The consumer (an LLM) is expected to see the
// full manifest in context and pick handles directly.
package examples

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Entry is one decomposition example file, parsed once at load time.
type Entry struct {
	Path        string // on-disk path, e.g. "examples/medium/csv-reader.md"
	Tier        string // "trivial" | "small" | "medium" | "large"
	Slug        string // filename without .md
	Requirement string // extracted from the `# Requirement: "..."` header
	Content     string // full file text
}

// Handle returns the short identifier the agent uses to refer to this entry,
// e.g. "medium/csv-reader".
func (e *Entry) Handle() string { return e.Tier + "/" + e.Slug }

// Index is the in-memory corpus. Built once via LoadFromDir, then queried.
type Index struct {
	Entries  []Entry
	byHandle map[string]*Entry
	bySlug   map[string][]*Entry
}

// Len returns the number of entries in the index.
func (idx *Index) Len() int { return len(idx.Entries) }

var tiers = []string{"trivial", "small", "medium", "large"}

var requirementRE = regexp.MustCompile(`(?m)^\s*#?\s*Requirement:\s*"([^"]+)"`)

// LoadFromDir walks root/{trivial,small,medium,large}/*.md and builds an Index.
// Files without a parseable requirement header are skipped (but do not error).
func LoadFromDir(root string) (*Index, error) {
	idx := &Index{
		byHandle: map[string]*Entry{},
		bySlug:   map[string][]*Entry{},
	}
	for _, tier := range tiers {
		tierDir := filepath.Join(root, tier)
		info, err := os.Stat(tierDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("stat %s: %w", tierDir, err)
		}
		if !info.IsDir() {
			continue
		}

		entries, err := os.ReadDir(tierDir)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", tierDir, err)
		}
		for _, de := range entries {
			if de.IsDir() || !strings.HasSuffix(de.Name(), ".md") {
				continue
			}
			if de.Name() == "README.md" {
				continue
			}
			path := filepath.Join(tierDir, de.Name())
			entry, ok := parseFile(path, tier)
			if !ok {
				continue
			}
			idx.Entries = append(idx.Entries, entry)
		}
	}
	// Sort so manifest output is deterministic and group-able by tier.
	sort.SliceStable(idx.Entries, func(i, j int) bool {
		ti, tj := tierRank(idx.Entries[i].Tier), tierRank(idx.Entries[j].Tier)
		if ti != tj {
			return ti < tj
		}
		return idx.Entries[i].Slug < idx.Entries[j].Slug
	})
	for i := range idx.Entries {
		e := &idx.Entries[i]
		idx.byHandle[e.Handle()] = e
		idx.bySlug[e.Slug] = append(idx.bySlug[e.Slug], e)
	}
	return idx, nil
}

// LookupResult describes how a caller's reference resolved. Kind tells the
// caller whether we returned the exact handle they asked for, auto-corrected
// the tier because the slug was unambiguous across tiers, or failed entirely
// (in which case Suggestions holds up to 5 plausible alternatives).
type LookupResult struct {
	Kind        LookupKind
	Entry       *Entry
	Suggestions []*Entry // only populated when Kind == LookupMiss
}

// LookupKind classifies Lookup outcomes for callers that want to distinguish
// exact hits from tier-corrected hits and misses.
type LookupKind int

const (
	// LookupHit means the exact tier/slug matched.
	LookupHit LookupKind = iota
	// LookupTierCorrected means the caller's tier was wrong but the slug
	// was unambiguous — we returned the real entry.
	LookupTierCorrected
	// LookupMiss means no match; Suggestions may hold similar candidates.
	LookupMiss
)

// Lookup resolves a handle like "medium/csv-reader", an on-disk path like
// "examples/medium/csv-reader.md", or just a slug like "csv-reader" to an
// Entry. Returns a LookupResult describing the match quality.
//
// Resolution order:
//  1. Exact tier/slug match.
//  2. Slug-only match, if the slug exists in exactly one tier.
//  3. Miss — returns up to 5 suggestions with similar slugs.
func (idx *Index) Lookup(ref string) LookupResult {
	handle := normalizeHandle(ref)

	// 1. Exact handle match.
	if e, ok := idx.byHandle[handle]; ok {
		return LookupResult{Kind: LookupHit, Entry: e}
	}

	// 2. Extract the slug portion and try a slug-only lookup.
	slug := handle
	if _, after, ok := strings.Cut(handle, "/"); ok {
		slug = after
	}
	if bucket, ok := idx.bySlug[slug]; ok && len(bucket) == 1 {
		return LookupResult{Kind: LookupTierCorrected, Entry: bucket[0]}
	}

	// 3. Miss — suggest similar slugs by prefix/substring overlap.
	return LookupResult{Kind: LookupMiss, Suggestions: idx.suggest(slug, 5)}
}

// suggest returns up to n entries whose slugs share a long substring with the
// caller's slug. A simple heuristic: rank by longest common prefix length,
// tiebreak by shorter slug.
func (idx *Index) suggest(slug string, n int) []*Entry {
	if slug == "" || len(idx.Entries) == 0 {
		return nil
	}
	type scored struct {
		entry *Entry
		score int
	}
	scores := make([]scored, 0, 64)
	for i := range idx.Entries {
		e := &idx.Entries[i]
		s := commonPrefixLen(e.Slug, slug)
		// Also credit substring containment — handles where user's slug
		// is a prefix/suffix of the real slug or vice-versa.
		if strings.Contains(e.Slug, slug) || strings.Contains(slug, e.Slug) {
			s += 2
		}
		if s >= 4 {
			scores = append(scores, scored{entry: e, score: s})
		}
	}
	sort.SliceStable(scores, func(i, j int) bool {
		if scores[i].score != scores[j].score {
			return scores[i].score > scores[j].score
		}
		return len(scores[i].entry.Slug) < len(scores[j].entry.Slug)
	})
	if len(scores) > n {
		scores = scores[:n]
	}
	out := make([]*Entry, 0, len(scores))
	for _, s := range scores {
		out = append(out, s.entry)
	}
	return out
}

func commonPrefixLen(a, b string) int {
	n := min(len(a), len(b))
	i := 0
	for i < n && a[i] == b[i] {
		i++
	}
	return i
}

// Manifest returns a formatted, deterministic listing of every entry's handle
// grouped by tier. Intended to be inlined into a system prompt so the LLM sees
// the full corpus index in context.
//
// Shape:
//
//	## trivial (74)
//	- trivial/add-two-integers
//	- trivial/hello-world
//	...
//
//	## small (465)
//	- small/...
//	...
func (idx *Index) Manifest() string {
	counts := map[string]int{}
	groups := map[string][]*Entry{}
	for i := range idx.Entries {
		e := &idx.Entries[i]
		counts[e.Tier]++
		groups[e.Tier] = append(groups[e.Tier], e)
	}
	var b strings.Builder
	for _, tier := range tiers {
		entries := groups[tier]
		if len(entries) == 0 {
			continue
		}
		fmt.Fprintf(&b, "## %s (%d)\n", tier, counts[tier])
		for _, e := range entries {
			fmt.Fprintf(&b, "- %s\n", e.Handle())
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

// normalizeHandle accepts several forms and returns the canonical
// "tier/slug" key used in the byHandle map.
func normalizeHandle(ref string) string {
	ref = strings.TrimSpace(ref)
	ref = strings.TrimPrefix(ref, "./")
	ref = strings.TrimPrefix(ref, "examples/")
	ref = strings.TrimSuffix(ref, ".md")
	return ref
}

func tierRank(tier string) int {
	for i, t := range tiers {
		if t == tier {
			return i
		}
	}
	return len(tiers)
}

// parseFile reads one .md file and returns its Entry. Returns ok=false if the
// requirement header is missing (we skip those files rather than erroring).
func parseFile(path, tier string) (Entry, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, false
	}
	content := string(data)
	m := requirementRE.FindStringSubmatch(content)
	if m == nil {
		return Entry{}, false
	}
	slug := strings.TrimSuffix(filepath.Base(path), ".md")
	requirement := strings.TrimSpace(m[1])
	return Entry{
		Path:        path,
		Tier:        tier,
		Slug:        slug,
		Requirement: requirement,
		Content:     content,
	}, true
}
