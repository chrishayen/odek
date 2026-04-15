package tui

import (
	"fmt"
	"sort"
	"strings"

	"shotgun.dev/odek/internal/decomposer"
)

// renderDecompositionSummary returns the short prose summary shown in the
// chat after a decompose finishes. Two branches:
//
//   - Fresh (prior == nil): one or two sentences introducing the package
//     and listing up to 5 top-level rune names, plus the std count if any.
//   - Refinement (prior != nil): a name-level diff against the prior
//     decomposition, rendered as "Updated **x**: added …; removed …"
//     prose so the user can see what their latest message changed.
//
// Full detail (signatures, tests, assumptions, dependencies) still lives
// on the right-side decomposition pane; this summary only belongs in the
// chat where the user is iterating.
func renderDecompositionSummary(sess *decomposer.Session, prior *decomposer.DecompositionResponse) string {
	if sess == nil || sess.Root == nil || sess.Root.Response == nil {
		return "(no decomposition available)"
	}

	resp := sess.Root.Response
	newProject := resp.ProjectPackage
	var newStd decomposer.PackageNode
	if resp.StdPackage != nil {
		newStd = *resp.StdPackage
	}

	if prior == nil {
		return renderFreshSummary(newProject, newStd)
	}

	var priorProject decomposer.PackageNode
	var priorStd decomposer.PackageNode
	priorProject = prior.ProjectPackage
	if prior.StdPackage != nil {
		priorStd = *prior.StdPackage
	}
	return renderRefinementSummary(newProject, newStd, priorProject, priorStd)
}

// renderFreshSummary handles the no-prior case: a first-pass decomposition
// with nothing to diff against. Introduces the package and lists a few
// rune names.
func renderFreshSummary(project, std decomposer.PackageNode) string {
	var b strings.Builder
	total := len(project.Runes)
	fmt.Fprintf(&b, "Decomposed into **%s** with %d top-level %s.", project.Name, total, pluralRune(total))
	if total > 0 {
		fmt.Fprintf(&b, " Includes %s.", joinNames(sortedRuneNames(project.Runes), 5))
	}
	if len(std.Runes) > 0 {
		fmt.Fprintf(&b, " Plus %d shared %s in `%s`: %s.",
			len(std.Runes), pluralRune(len(std.Runes)), std.Name,
			joinNames(sortedRuneNames(std.Runes), 3))
	}
	return b.String()
}

// renderRefinementSummary handles the refinement case: compute a name-
// level diff between the new and prior decompositions and describe the
// change in prose so the chat shows what the user's latest message
// actually changed.
func renderRefinementSummary(newP, newS, oldP, oldS decomposer.PackageNode) string {
	addedP, removedP := diffRuneNames(newP.Runes, oldP.Runes)
	addedS, removedS := diffRuneNames(newS.Runes, oldS.Runes)

	total := len(newP.Runes)
	if len(addedP) == 0 && len(removedP) == 0 && len(addedS) == 0 && len(removedS) == 0 {
		return fmt.Sprintf("No structural changes to **%s**; same %d %s.", newP.Name, total, pluralRune(total))
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Updated **%s**", newP.Name)

	var clauses []string
	if len(addedP) > 0 {
		clauses = append(clauses, "added "+joinNames(addedP, 5))
	}
	if len(removedP) > 0 {
		clauses = append(clauses, "removed "+joinNames(removedP, 5))
	}
	if len(clauses) > 0 {
		b.WriteString(": ")
		b.WriteString(strings.Join(clauses, "; "))
	}
	b.WriteString(".")

	if len(addedS) > 0 || len(removedS) > 0 {
		var stdClauses []string
		if len(addedS) > 0 {
			stdClauses = append(stdClauses, "added "+joinNames(addedS, 5))
		}
		if len(removedS) > 0 {
			stdClauses = append(stdClauses, "removed "+joinNames(removedS, 5))
		}
		fmt.Fprintf(&b, " std: %s.", strings.Join(stdClauses, "; "))
	}

	fmt.Fprintf(&b, " Now %d %s total.", total, pluralRune(total))
	return b.String()
}

// diffRuneNames returns sorted lists of names that appear in newRunes but
// not in oldRunes (added) and vice-versa (removed). Exact name match only
// — no fuzzy rename detection.
func diffRuneNames(newRunes, oldRunes map[string]decomposer.Rune) (added, removed []string) {
	for name := range newRunes {
		if _, ok := oldRunes[name]; !ok {
			added = append(added, name)
		}
	}
	for name := range oldRunes {
		if _, ok := newRunes[name]; !ok {
			removed = append(removed, name)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	return
}

// sortedRuneNames returns the map keys of a rune map in sorted order.
func sortedRuneNames(runes map[string]decomposer.Rune) []string {
	names := make([]string, 0, len(runes))
	for name := range runes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// joinNames renders names as a comma-separated list of inline-code spans
// (`name`), truncated at max with an "and N more" tail.
func joinNames(names []string, max int) string {
	if len(names) == 0 {
		return ""
	}
	shown := names
	moreCount := 0
	if len(names) > max {
		shown = names[:max]
		moreCount = len(names) - max
	}
	quoted := make([]string, len(shown))
	for i, n := range shown {
		quoted[i] = "`" + n + "`"
	}
	result := strings.Join(quoted, ", ")
	if moreCount > 0 {
		result += fmt.Sprintf(", and %d more", moreCount)
	}
	return result
}

func pluralRune(n int) string {
	if n == 1 {
		return "rune"
	}
	return "runes"
}
