# Requirement: "a case-insensitive string comparison and search library"

Wraps common string operations so that uppercase and lowercase variants are treated as equal. Case folding goes through a thin std primitive so it stays Unicode-correct.

std
  std.unicode
    std.unicode.fold_case
      @ (s: string) -> string
      + returns the full Unicode case-folded form of s
      ? uses the full mapping (e.g. sharp-s to "ss")
      # unicode

cistr
  cistr.equals
    @ (a: string, b: string) -> bool
    + returns true when a and b are equal after case folding
    + returns true for ("Hello", "HELLO")
    - returns false for ("cat", "dog")
    # comparison
    -> std.unicode.fold_case
  cistr.contains
    @ (haystack: string, needle: string) -> bool
    + returns true when needle appears in haystack ignoring case
    # search
    -> std.unicode.fold_case
  cistr.index_of
    @ (haystack: string, needle: string) -> i32
    + returns the byte offset of the first case-insensitive match, or -1
    # search
    -> std.unicode.fold_case
  cistr.has_prefix
    @ (s: string, prefix: string) -> bool
    + case-insensitive prefix check
    # comparison
    -> std.unicode.fold_case
  cistr.has_suffix
    @ (s: string, suffix: string) -> bool
    + case-insensitive suffix check
    # comparison
    -> std.unicode.fold_case
