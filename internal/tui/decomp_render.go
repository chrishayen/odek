package tui

import (
	"strings"

	"shotgun.dev/odek/internal/decomposer"
)

// renderDecompositionSummary returns the prose shown in the chat after a
// decompose finishes. The decomposer tool schema requires the model to
// provide a `summary` field alongside the rune tree — a 1-2 sentence
// narrative that introduces what the feature is (fresh) or describes what
// changed (refinement). We surface that string verbatim here.
//
// Structural detail (names, signatures, tests, assumptions) lives on the
// right-side decomposition pane; the chat gets the narrative only.
func renderDecompositionSummary(sess *decomposer.Session) string {
	if sess == nil || sess.Root == nil || sess.Root.Response == nil {
		return "(no decomposition available)"
	}
	if s := strings.TrimSpace(sess.Root.Response.Summary); s != "" {
		return s
	}
	return "Decomposed."
}
