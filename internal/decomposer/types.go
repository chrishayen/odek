package decomposer

import (
	"sort"
	"sync"

	openai "shotgun.dev/odek/openai"
)

// Rune is a single function specification as returned by the `decompose`
// tool call. Field shape must match the tool's JSON schema in decomposer.go.
type Rune struct {
	Description   string   `json:"description"`
	FunctionSig   string   `json:"function_signature"`
	PositiveTests []string `json:"positive_tests"`
	NegativeTests []string `json:"negative_tests"`
	Assumptions   []string `json:"assumptions"`
}

// PackageNode is a flat map of rune name → rune, with a package name. Both
// the project package and (optionally) a stdlib package are returned per
// decomposition.
type PackageNode struct {
	Name  string          `json:"name"`
	Runes map[string]Rune `json:"runes"`
}

// DecompositionResponse is the parsed JSON arguments of a single `decompose`
// tool call. Summary is a 1-2 sentence narrative the model writes alongside
// the tree: on a fresh pass it introduces what the feature is; on a
// refinement pass it describes what changed based on the user's latest
// feedback. The chat surfaces it verbatim; the right-hand pane ignores it.
type DecompositionResponse struct {
	Summary        string       `json:"summary"`
	ProjectPackage PackageNode  `json:"project_package"`
	StdPackage     *PackageNode `json:"std_package,omitempty"`
}

// ClarificationRequest is returned instead of a DecompositionResponse when
// the model replies in plain text (did not call the decompose tool).
type ClarificationRequest struct {
	Message string
}

// ClarificationNeeded is returned as an error from NewSession when the
// model replied with a clarification question instead of a decomposition.
// This is NOT a hard failure — callers should surface Message to the user
// as a normal assistant turn and let them answer before retrying the
// /decompose command. The user's prior session (if any) must be left
// untouched.
type ClarificationNeeded struct {
	Message string
}

func (e *ClarificationNeeded) Error() string {
	return "clarification needed: " + e.Message
}

// AutoDecomposition is one node in the recursive expansion tree. The root
// node holds the initial decomposition; each child holds a later expansion
// of a single rune at a deeper level.
type AutoDecomposition struct {
	Path       string
	Depth      int
	Response   *DecompositionResponse
	ParentPath string
	ChildPaths []string
}

// RuneExpansionInfo is the expansion queue entry for a single rune.
type RuneExpansionInfo struct {
	FullPath            string
	Depth               int
	ParentDecomposition *AutoDecomposition
}

// RuneStatus tracks where each rune sits in the expansion lifecycle.
type RuneStatus int

const (
	StatusPending RuneStatus = iota
	StatusInFlight
	StatusDone
	StatusLeaf
	StatusError
)

// Snapshot is a copy of session state suitable for rendering without
// holding the session's mutex. Keys in TopLevelNames are sorted
// lexicographically for stable rendering.
type Snapshot struct {
	HasSession      bool
	PackageName     string
	TopLevelNames   []string
	RunesByName     map[string]Rune
	StatusByName    map[string]RuneStatus
	ChildrenByName  map[string][]string
	Requirement     string
	Summary         string
	TotalRunes      int
	MaxDepthReached int
	Expanding       bool
	InFlightCount   int
	ErrorCount      int

	// PackagePaths lists the column-0 entries under the synthetic "root"
	// parent. Ordered deterministically: std first (when it has runes),
	// then the project package.
	PackagePaths []string
	// RuneByPath holds every rune spec discovered across the full
	// expansion tree, keyed by fully-qualified path. Unlike RunesByName
	// it includes sub-runes produced by recursive expansion, not just
	// the root-level decomposition.
	RuneByPath map[string]Rune
	// DisplayNameByPath maps a fully-qualified path to the short leaf
	// label shown in a column (final dot-segment). Package roots map to
	// themselves.
	DisplayNameByPath map[string]string
}

// Session is the mutable tree-plus-channel shared across chat, the
// decomposition page, and the expansion goroutine. It lives on the heap and
// is passed around as *Session so Bubble Tea's copy-on-Update model sharing
// doesn't create stale state.
type Session struct {
	Requirement  string
	EffortLevel  int
	EffortReason string
	Root         *AutoDecomposition
	BaseMessages []openai.ChatMessage

	mu         sync.Mutex
	tree       map[string]*AutoDecomposition
	treeOrder  []string
	status     map[string]RuneStatus
	totalRunes int
	maxDepth   int
	expanding  bool

	Events <-chan ExpansionEvent
	Cancel func()
}

func newSession(req string, level int, reason string, root *AutoDecomposition, baseMsgs []openai.ChatMessage) *Session {
	s := &Session{
		Requirement:  req,
		EffortLevel:  level,
		EffortReason: reason,
		Root:         root,
		BaseMessages: baseMsgs,
		tree:         map[string]*AutoDecomposition{root.Path: root},
		treeOrder:    []string{root.Path},
		status:       map[string]RuneStatus{},
	}
	if root.Response != nil {
		s.totalRunes = countTotalRunes(root.Response)
		for _, path := range initialTopLevelPaths(root.Response) {
			s.status[path] = StatusPending
		}
	}
	return s
}

// Apply mutates session state in response to a single expansion event.
// Safe to call from any goroutine. Idempotent for terminal status events
// at the same path (later events overwrite earlier status, but do not
// corrupt structure).
func (s *Session) Apply(evt ExpansionEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch e := evt.(type) {
	case EventLevelStarted:
		s.expanding = true
	case EventRuneStarted:
		s.status[e.Path] = StatusInFlight
	case EventRuneExpanded:
		if e.ChildCount == 0 {
			s.status[e.Path] = StatusLeaf
		} else {
			s.status[e.Path] = StatusDone
		}
		if _, exists := s.tree[e.Path]; !exists {
			child := &AutoDecomposition{
				Path:       e.Path,
				Depth:      e.Depth,
				Response:   e.Response,
				ParentPath: e.ParentPath,
				ChildPaths: []string{},
			}
			s.tree[e.Path] = child
			s.treeOrder = append(s.treeOrder, e.Path)
			if parent, ok := s.tree[e.ParentPath]; ok && e.ParentPath != "" {
				parent.ChildPaths = append(parent.ChildPaths, e.Path)
			}
		}
		s.totalRunes += e.ChildCount
		if e.Depth > s.maxDepth {
			s.maxDepth = e.Depth
		}
		if e.Response != nil {
			for _, path := range initialTopLevelPaths(e.Response) {
				if _, ok := s.status[path]; !ok {
					s.status[path] = StatusPending
				}
			}
		}
	case EventRuneError:
		s.status[e.Path] = StatusError
	case EventCancelled, EventDone, EventCapReached:
		s.expanding = false
	}
}

// Snapshot returns a stable view of the session. Safe to call from the
// render path.
func (s *Session) Snapshot() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	snap := Snapshot{
		HasSession:        s.Root != nil && s.Root.Response != nil,
		Requirement:       s.Requirement,
		TotalRunes:        s.totalRunes,
		MaxDepthReached:   s.maxDepth,
		Expanding:         s.expanding,
		RunesByName:       map[string]Rune{},
		StatusByName:      map[string]RuneStatus{},
		ChildrenByName:    map[string][]string{},
		RuneByPath:        map[string]Rune{},
		DisplayNameByPath: map[string]string{},
	}

	if !snap.HasSession {
		return snap
	}

	root := s.Root.Response
	snap.PackageName = root.ProjectPackage.Name
	snap.Summary = root.Summary

	names := make([]string, 0, len(root.ProjectPackage.Runes))
	for name := range root.ProjectPackage.Runes {
		names = append(names, name)
	}
	sort.Strings(names)
	snap.TopLevelNames = names

	for name, r := range root.ProjectPackage.Runes {
		snap.RunesByName[name] = r
	}
	if root.StdPackage != nil {
		for name, r := range root.StdPackage.Runes {
			snap.RunesByName[name] = r
		}
	}

	for path, st := range s.status {
		snap.StatusByName[path] = st
	}

	// Walk the full decomposition tree so sub-runes produced by
	// recursive expansion (and std runes introduced at any depth) show
	// up in RuneByPath / DisplayNameByPath.
	projPkg := root.ProjectPackage.Name
	stdPkg := ""
	hasAnyStd := false
	for _, path := range s.treeOrder {
		d := s.tree[path]
		if d == nil || d.Response == nil {
			continue
		}
		pp := d.Response.ProjectPackage
		for name, r := range pp.Runes {
			full := qualify(pp.Name, name)
			snap.RuneByPath[full] = r
			snap.DisplayNameByPath[full] = lastSegment(full)
		}
		if sp := d.Response.StdPackage; sp != nil && sp.Name != "" {
			if stdPkg == "" {
				stdPkg = sp.Name
			}
			for name, r := range sp.Runes {
				full := qualify(sp.Name, name)
				snap.RuneByPath[full] = r
				snap.DisplayNameByPath[full] = lastSegment(full)
				hasAnyStd = true
			}
		}
	}

	// PackagePaths is the display order for column-0 sections: std
	// first (when present), then the project package. The TUI renders
	// this as headers with runes listed directly underneath — not as
	// selectable rows.
	if stdPkg != "" && hasAnyStd {
		snap.PackagePaths = append(snap.PackagePaths, stdPkg)
		snap.DisplayNameByPath[stdPkg] = stdPkg
	}
	if projPkg != "" {
		snap.PackagePaths = append(snap.PackagePaths, projPkg)
		snap.DisplayNameByPath[projPkg] = projPkg
	}

	// Build each package's top-level rune list. std aggregates every
	// std rune encountered across all decomposition nodes so primitives
	// introduced during deeper expansions show up in column 0 too.
	if projPkg != "" {
		for name := range root.ProjectPackage.Runes {
			full := qualify(projPkg, name)
			snap.ChildrenByName[projPkg] = append(snap.ChildrenByName[projPkg], full)
		}
		sort.Strings(snap.ChildrenByName[projPkg])
	}
	if stdPkg != "" {
		for _, path := range s.treeOrder {
			d := s.tree[path]
			if d == nil || d.Response == nil || d.Response.StdPackage == nil {
				continue
			}
			sp := d.Response.StdPackage
			for name := range sp.Runes {
				full := qualify(sp.Name, name)
				snap.ChildrenByName[stdPkg] = append(snap.ChildrenByName[stdPkg], full)
			}
		}
		sort.Strings(snap.ChildrenByName[stdPkg])
		snap.ChildrenByName[stdPkg] = dedupeSorted(snap.ChildrenByName[stdPkg])
	}

	// Sub-rune children (expanded decomposition nodes) attach to their
	// parent rune path — expand.go sets ParentPath to the parent's full
	// path. Skip the root-level decomposition whose ParentPath is "" or
	// "root"; those are handled separately below.
	for _, child := range s.tree {
		if child == nil || child.ParentPath == "" || child.ParentPath == "root" {
			continue
		}
		snap.ChildrenByName[child.ParentPath] = append(snap.ChildrenByName[child.ParentPath], child.Path)
	}
	// Sort + dedupe sub-rune children. Skip the keys we've already
	// ordered deliberately above (per-package lists, and "root" which
	// must stay std-first not lexicographic).
	for parent, children := range snap.ChildrenByName {
		if parent == stdPkg || parent == projPkg {
			continue
		}
		sort.Strings(children)
		snap.ChildrenByName[parent] = dedupeSorted(children)
	}

	// Navigation feed: flat ordered list of every top-level rune across
	// packages, in display order (std section first, then project).
	// Up/down in column 0 moves through this list; section headers are
	// pure presentation and don't take a nav slot.
	var rootKids []string
	for _, pkg := range snap.PackagePaths {
		rootKids = append(rootKids, snap.ChildrenByName[pkg]...)
	}
	snap.ChildrenByName["root"] = rootKids

	for _, st := range s.status {
		switch st {
		case StatusInFlight:
			snap.InFlightCount++
		case StatusError:
			snap.ErrorCount++
		}
	}

	return snap
}

// lastSegment returns the substring after the final '.' in path, or the
// whole path if there is no dot. Used to turn fully-qualified rune paths
// into the short labels shown in a column.
func lastSegment(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return path
}

// dedupeSorted collapses consecutive duplicates in a sorted slice. The
// caller passes an already-sorted slice.
func dedupeSorted(in []string) []string {
	if len(in) < 2 {
		return in
	}
	out := in[:1]
	for _, s := range in[1:] {
		if s != out[len(out)-1] {
			out = append(out, s)
		}
	}
	return out
}

// TopLevelPaths returns the fully-qualified paths of the top-level runes
// from the root decomposition, for the CLI queue driver.
func (s *Session) TopLevelPaths() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Root == nil || s.Root.Response == nil {
		return nil
	}
	return initialTopLevelPaths(s.Root.Response)
}

// AllDecompositions returns the session's internal tree in insertion order,
// for CLI printing of the complete tree.
func (s *Session) AllDecompositions() []*AutoDecomposition {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*AutoDecomposition, 0, len(s.treeOrder))
	for _, path := range s.treeOrder {
		if d, ok := s.tree[path]; ok {
			out = append(out, d)
		}
	}
	return out
}

// setEvents + clearEvents bracket a live ExpansionEvent channel. Separate
// from Session.mu so the pump goroutine never waits on the session's tree
// mutex while draining.
func (s *Session) setEvents(ch <-chan ExpansionEvent, cancel func()) {
	s.mu.Lock()
	s.Events = ch
	s.Cancel = cancel
	s.mu.Unlock()
}

func (s *Session) clearEvents() {
	s.mu.Lock()
	s.Events = nil
	s.Cancel = nil
	s.mu.Unlock()
}

// initialTopLevelPaths returns a stable sorted list of fully-qualified paths
// for the runes in a DecompositionResponse's project package (and std).
func initialTopLevelPaths(resp *DecompositionResponse) []string {
	if resp == nil {
		return nil
	}
	paths := make([]string, 0)
	if resp.ProjectPackage.Name != "" {
		for name := range resp.ProjectPackage.Runes {
			paths = append(paths, qualify(resp.ProjectPackage.Name, name))
		}
	}
	if resp.StdPackage != nil && resp.StdPackage.Name != "" {
		for name := range resp.StdPackage.Runes {
			paths = append(paths, qualify(resp.StdPackage.Name, name))
		}
	}
	sort.Strings(paths)
	return paths
}

func qualify(pkgName, runeName string) string {
	if len(runeName) > len(pkgName)+1 && runeName[:len(pkgName)] == pkgName && runeName[len(pkgName)] == '.' {
		return runeName
	}
	return pkgName + "." + runeName
}
