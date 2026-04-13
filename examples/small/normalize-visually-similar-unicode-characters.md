# Requirement: "a library that normalizes visually similar unicode characters to a canonical form"

Maps homoglyphs (Cyrillic a, Greek omicron, fullwidth digits, etc.) to their closest ASCII-ish equivalents.

std: (all units exist)

unhomoglyph
  unhomoglyph.normalize
    @ (input: string) -> string
    + replaces each homoglyph codepoint with its canonical form
    + leaves non-homoglyph codepoints unchanged
    - returns empty string for empty input
    # normalization
  unhomoglyph.is_homoglyph
    @ (codepoint: i32) -> bool
    + returns true when the codepoint has a defined canonical form
    # lookup
  unhomoglyph.canonical_for
    @ (codepoint: i32) -> optional[string]
    + returns the canonical replacement for a codepoint, if any
    # lookup
