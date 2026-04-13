# Requirement: "a fuzzy text library that sanitizes, normalizes, and compares strings"

Strips noise, folds to a canonical form, and reports a similarity ratio between two canonicalized strings.

std
  std.unicode
    std.unicode.nfkd
      @ (input: string) -> string
      + returns the NFKD-normalized form
      # unicode
    std.unicode.strip_marks
      @ (input: string) -> string
      + removes combining marks such as accents
      # unicode
    std.unicode.to_lower
      @ (input: string) -> string
      + returns the lowercase form
      # unicode

fuzzy_text
  fuzzy_text.sanitize
    @ (input: string) -> string
    + collapses runs of whitespace into a single space and trims ends
    + removes control characters
    # sanitize
  fuzzy_text.normalize
    @ (input: string) -> string
    + returns a canonical form: sanitized, lowercased, stripped of accents
    # normalize
    -> std.unicode.nfkd
    -> std.unicode.strip_marks
    -> std.unicode.to_lower
  fuzzy_text.edit_distance
    @ (a: string, b: string) -> i32
    + returns the Levenshtein distance between two strings
    # distance
  fuzzy_text.similarity
    @ (a: string, b: string) -> f64
    + returns a ratio in [0.0, 1.0] based on edit distance after normalization
    + returns 1.0 when both normalized strings are identical
    + returns 0.0 when one normalized string is empty and the other is not
    # compare
