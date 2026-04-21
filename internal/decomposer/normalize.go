package decomposer

import "strings"

// NormalizeFunctionSig strips any leaked marker prefix from a function
// signature string. The tree format in decompose.md and the example corpus
// uses "fn" as a visual marker on signature lines, and the prior convention
// used "@". Both are visual markers for human readers, not part of the
// signature value — but the model occasionally leaks them into the
// function_signature JSON field of the decompose tool call. This helper
// defensively removes either prefix (and loops so that compound leakage
// like "fn @ (...)" still collapses to "(...)").
// normalizePackageSignatures walks every rune in a PackageNode and cleans
// its FunctionSig field in place. Called at decompose tool-parse time so
// the session holds canonical signatures; this prevents leaked "@"/"fn"
// markers from being echoed back to the model on refinement passes (where
// the prior decomposition is marshaled into the next prompt).
func normalizePackageSignatures(pkg *PackageNode) {
	if pkg == nil {
		return
	}
	normalizeRuneMap(pkg.Runes)
}

func normalizeRuneMap(runes map[string]Rune) {
	for name, r := range runes {
		r.FunctionSig = NormalizeFunctionSig(r.FunctionSig)
		if len(r.Children) > 0 {
			normalizeRuneMap(r.Children)
		}
		runes[name] = r
	}
}

func NormalizeFunctionSig(s string) string {
	s = strings.TrimSpace(s)
	for {
		lowered := strings.ToLower(s)
		switch {
		case strings.HasPrefix(lowered, "fn "), strings.HasPrefix(lowered, "fn\t"):
			s = strings.TrimSpace(s[3:])
		case strings.HasPrefix(lowered, "fn("):
			s = strings.TrimSpace(s[2:])
		case lowered == "fn":
			return ""
		case strings.HasPrefix(s, "@ "), strings.HasPrefix(s, "@\t"):
			s = strings.TrimSpace(s[2:])
		case strings.HasPrefix(s, "@"):
			s = strings.TrimSpace(s[1:])
		default:
			return s
		}
	}
}
