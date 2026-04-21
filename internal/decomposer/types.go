package decomposer

import (
	"sort"
	"sync"

	openai "shotgun.dev/odek/openai"
)

// Rune is a single function specification as returned by the `decompose`
// tool call. The tree is expressed via nested Children — each key is the
// next path segment under this rune. Leaf runes have empty Children.
//
// Field shape must match the tool's JSON schema in decomposer.go.
type Rune struct {
	Description   string          `json:"description"`
	FunctionSig   string          `json:"function_signature"`
	PositiveTests []string        `json:"positive_tests"`
	NegativeTests []string        `json:"negative_tests"`
	Assumptions   []string        `json:"assumptions"`
	Dependencies  []string        `json:"dependencies,omitempty"`
	Children      map[string]Rune `json:"children,omitempty"`
}

// PackageNode is the top-level container: a package name plus its
// first-level Runes. Deeper levels live inside each Rune.Children.
type PackageNode struct {
	Name  string          `json:"name"`
	Runes map[string]Rune `json:"runes"`
}

// DecompositionResponse is the parsed JSON arguments of a single `decompose`
// tool call — the output of pass 2. Summary is a 1-2 sentence narrative
// the model writes alongside the tree.
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
type ClarificationNeeded struct {
	Message string
}

func (e *ClarificationNeeded) Error() string {
	return "clarification needed: " + e.Message
}

// Snapshot is a copy of session state suitable for rendering without
// holding the session's mutex.
type Snapshot struct {
	HasSession  bool
	Requirement string

	// Contract is the pass-1 output text, accumulated as chunks arrive.
	Contract string

	// Phase is one of "", PhaseContract, PhaseExtraction, "done", "error".
	Phase string

	// ExtractionBytes is the running count of pass-2 tool-arg bytes seen.
	// Only meaningful while Phase == PhaseExtraction.
	ExtractionBytes int

	// ErrorMsg is set when Phase == "error".
	ErrorMsg string

	// Response is non-nil once pass 2 has completed successfully.
	Response *DecompositionResponse

	// PackageName / Summary convenience accessors for the project package.
	PackageName string
	Summary     string

	// PackagePaths lists the column-0 entries under the synthetic "root"
	// parent: std first (when populated), then project.
	PackagePaths []string

	// ChildrenByName maps a parent path to its immediate child paths.
	// The key "root" holds the flat list of every top-level rune across
	// packages, used by the TUI's column-0 navigator.
	ChildrenByName map[string][]string

	// RuneByPath holds every rune spec in the tree, keyed by its
	// fully-qualified dotted path.
	RuneByPath map[string]Rune

	// DisplayNameByPath maps a fully-qualified path to the short label
	// shown in a column (final dot-segment). Package roots map to themselves.
	DisplayNameByPath map[string]string

	// TotalRunes is the count of all runes in the tree, leaves + parents.
	TotalRunes int

	// MaxDepth is the deepest dot-depth reached by any path (0 = just
	// package roots, 1 = first level of runes, etc.).
	MaxDepth int
}

// Session is the mutable state shared across the chat pane, the
// decomposition page, and the two-pass producer goroutine. Passed around
// as *Session so Bubble Tea's copy-on-Update model sharing doesn't create
// stale state.
type Session struct {
	Requirement  string
	EffortLevel  int
	EffortReason string
	BaseMessages []openai.ChatMessage

	mu              sync.Mutex
	contract        string
	phase           string
	extractionBytes int
	errorMsg        string
	response        *DecompositionResponse

	Events <-chan DecompositionEvent
	Cancel func()
}

func newSession(req string, level int, reason string, baseMsgs []openai.ChatMessage) *Session {
	return &Session{
		Requirement:  req,
		EffortLevel:  level,
		EffortReason: reason,
		BaseMessages: baseMsgs,
	}
}

// Apply mutates session state in response to a single decomposition event.
// Safe to call from any goroutine.
func (s *Session) Apply(evt DecompositionEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch e := evt.(type) {
	case EventPhaseStarted:
		s.phase = e.Phase
	case EventContractChunk:
		s.contract += e.Text
	case EventContractComplete:
		s.contract = e.Full
	case EventExtractionProgress:
		s.extractionBytes = e.Bytes
	case EventRunesComplete:
		s.response = e.Response
		s.phase = "done"
	case EventError:
		s.errorMsg = e.Err
		s.phase = "error"
	case EventCancelled:
		if s.phase != "done" && s.phase != "error" {
			s.phase = "error"
			if s.errorMsg == "" {
				s.errorMsg = "cancelled"
			}
		}
	}
}

// Snapshot returns a stable view of the session. Safe to call from the
// render path.
func (s *Session) Snapshot() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	snap := Snapshot{
		HasSession:        s.response != nil || s.contract != "",
		Requirement:       s.Requirement,
		Contract:          s.contract,
		Phase:             s.phase,
		ExtractionBytes:   s.extractionBytes,
		ErrorMsg:          s.errorMsg,
		Response:          s.response,
		ChildrenByName:    map[string][]string{},
		RuneByPath:        map[string]Rune{},
		DisplayNameByPath: map[string]string{},
	}

	if s.response == nil {
		return snap
	}

	snap.PackageName = s.response.ProjectPackage.Name
	snap.Summary = s.response.Summary

	var stdName string
	var stdHasRunes bool
	if sp := s.response.StdPackage; sp != nil {
		stdName = sp.Name
		stdHasRunes = len(sp.Runes) > 0
	}
	projName := s.response.ProjectPackage.Name

	if stdName != "" && stdHasRunes {
		snap.PackagePaths = append(snap.PackagePaths, stdName)
		snap.DisplayNameByPath[stdName] = stdName
	}
	if projName != "" {
		snap.PackagePaths = append(snap.PackagePaths, projName)
		snap.DisplayNameByPath[projName] = projName
	}

	if stdName != "" && stdHasRunes {
		walkPackage(stdName, s.response.StdPackage.Runes, &snap)
	}
	if projName != "" {
		walkPackage(projName, s.response.ProjectPackage.Runes, &snap)
	}

	// Sort each parent's children alphabetically.
	for parent, kids := range snap.ChildrenByName {
		sort.Strings(kids)
		snap.ChildrenByName[parent] = kids
	}

	// "root" is the flat top-level nav feed: std section first, then project.
	var rootKids []string
	for _, pkg := range snap.PackagePaths {
		rootKids = append(rootKids, snap.ChildrenByName[pkg]...)
	}
	snap.ChildrenByName["root"] = rootKids

	snap.TotalRunes = len(snap.RuneByPath)
	for p := range snap.RuneByPath {
		d := dotDepth(p)
		if d > snap.MaxDepth {
			snap.MaxDepth = d
		}
	}

	return snap
}

// walkPackage recurses into a package's rune map, populating RuneByPath,
// DisplayNameByPath, and ChildrenByName in snap. parentPath is the
// fully-qualified path to attach children under (the package name for the
// top call, then each rune path as we descend).
func walkPackage(parentPath string, runes map[string]Rune, snap *Snapshot) {
	for name, r := range runes {
		full := qualify(parentPath, name)
		snap.RuneByPath[full] = r
		snap.DisplayNameByPath[full] = lastSegment(full)
		snap.ChildrenByName[parentPath] = append(snap.ChildrenByName[parentPath], full)
		if len(r.Children) > 0 {
			walkPackage(full, r.Children, snap)
		}
	}
}

// lastSegment returns the substring after the final '.' in path, or the
// whole path if there is no dot.
func lastSegment(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return path
}

// dotDepth counts the number of '.' characters in path.
func dotDepth(path string) int {
	n := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '.' {
			n++
		}
	}
	return n
}

// TopLevelPaths returns the fully-qualified paths of the top-level runes
// from the completed response, for the CLI queue driver. Empty if pass 2
// has not completed.
func (s *Session) TopLevelPaths() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.response == nil {
		return nil
	}
	paths := make([]string, 0)
	if s.response.ProjectPackage.Name != "" {
		for name := range s.response.ProjectPackage.Runes {
			paths = append(paths, qualify(s.response.ProjectPackage.Name, name))
		}
	}
	if s.response.StdPackage != nil && s.response.StdPackage.Name != "" {
		for name := range s.response.StdPackage.Runes {
			paths = append(paths, qualify(s.response.StdPackage.Name, name))
		}
	}
	sort.Strings(paths)
	return paths
}

// Response returns the completed decomposition response, or nil if pass 2
// has not finished.
func (s *Session) Response() *DecompositionResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.response
}

// setEvents + clearEvents bracket a live DecompositionEvent channel.
func (s *Session) setEvents(ch <-chan DecompositionEvent, cancel func()) {
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

func qualify(pkgName, runeName string) string {
	if len(runeName) > len(pkgName)+1 && runeName[:len(pkgName)] == pkgName && runeName[len(pkgName)] == '.' {
		return runeName
	}
	return pkgName + "." + runeName
}
