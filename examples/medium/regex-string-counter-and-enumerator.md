# Requirement: "a library that counts and enumerates all strings matching a regular expression"

Parse a regex into an AST, then fold the AST into a match count and an iterator-style expansion.

std: (all units exist)

regen
  regen.parse
    fn (pattern: string) -> result[regex_node, string]
    + parses literals, character classes, concatenation, alternation, and bounded repetition
    - returns error on unbalanced parentheses or unbounded repetition (* or +)
    ? unbounded quantifiers are rejected because they produce infinite languages
    # parsing
  regen.count
    fn (node: regex_node) -> result[u64, string]
    + returns the exact number of strings matching the regex
    - returns error when the regex matches an infinite language
    # counting
  regen.expand_all
    fn (node: regex_node) -> list[string]
    + returns every matching string in lexicographic order
    - returns error when the count exceeds an implementation-defined cap
    # enumeration
  regen.nth
    fn (node: regex_node, index: u64) -> result[string, string]
    + returns the nth matching string (0-indexed) without materializing the full set
    - returns error when index is out of range
    # indexed_access
