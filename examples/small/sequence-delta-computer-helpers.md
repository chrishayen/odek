# Requirement: "helpers for computing deltas between sequences"

Compute longest-common-subsequences and produce line-level diffs in a few common formats.

std: (all units exist)

difflib
  difflib.longest_common_subsequence
    @ (a: list[string], b: list[string]) -> list[string]
    + returns one longest common subsequence of two token lists
    + returns an empty list when there is no shared element
    # algorithm
  difflib.ratio
    @ (a: string, b: string) -> f64
    + returns a similarity score in [0, 1] based on matching characters
    + returns 1.0 for identical inputs
    + returns 0.0 when inputs share no characters
    # similarity
  difflib.unified_diff
    @ (a: list[string], b: list[string], context: i32) -> list[string]
    + returns a unified diff of the two line lists with the given context size
    + returns an empty list when inputs are equal
    # diff
  difflib.ndiff
    @ (a: list[string], b: list[string]) -> list[string]
    + returns a line-by-line diff prefixing each line with ' ', '-', '+', or '?'
    # diff
