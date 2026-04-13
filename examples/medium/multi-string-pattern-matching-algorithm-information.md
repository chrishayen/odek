# Requirement: "a multi-string pattern matching library"

Builds an Aho-Corasick automaton from a set of patterns and finds all occurrences in a text in a single pass.

std: (all units exist)

multimatch
  multimatch.automaton_build
    @ (patterns: list[string]) -> ac_automaton
    + builds the goto, failure, and output links for the pattern set
    - returns an automaton that matches nothing when patterns is empty
    # construction
  multimatch.patterns
    @ (automaton: ac_automaton) -> list[string]
    + returns the patterns in the order they were added
    # inspection
  multimatch.find_all
    @ (automaton: ac_automaton, text: string) -> list[pattern_match]
    + returns all matches as (pattern_index, start_offset, end_offset)
    + returns overlapping matches for overlapping patterns
    # matching
  multimatch.find_first
    @ (automaton: ac_automaton, text: string) -> optional[pattern_match]
    + returns the first match by end offset
    - returns none when no pattern occurs
    # matching
  multimatch.replace_all
    @ (automaton: ac_automaton, text: string, replacements: list[string]) -> string
    + replaces each match with the corresponding replacement by pattern index
    ? overlapping matches resolve to the leftmost-longest
    # replacement
