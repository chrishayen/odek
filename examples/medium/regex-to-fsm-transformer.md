# Requirement: "a library for transforming regular expressions into finite state machines"

Parse a regex into an AST, convert to an NFA via Thompson's construction, then to a DFA via subset construction, with a minimization pass.

std: (all units exist)

regex_fsm
  regex_fsm.parse
    fn (pattern: string) -> result[regex_ast, string]
    + parses concatenation, alternation, star, plus, question, character classes, and grouping
    - returns error with position on unbalanced parentheses
    - returns error on unexpected metacharacters
    # parsing
  regex_fsm.to_nfa
    fn (ast: regex_ast) -> nfa
    + returns an NFA built via Thompson's construction with a single accept state
    # nfa
  regex_fsm.to_dfa
    fn (machine: nfa) -> dfa
    + returns a DFA via subset construction
    + resolves epsilon closures during conversion
    # dfa
  regex_fsm.minimize
    fn (machine: dfa) -> dfa
    + returns an equivalent DFA with the minimum number of states using Hopcroft's algorithm
    # optimization
  regex_fsm.compile
    fn (pattern: string) -> result[dfa, string]
    + parses, converts to NFA, converts to DFA, and minimizes in one call
    - returns error when parsing fails
    # pipeline
  regex_fsm.match
    fn (machine: dfa, input: string) -> bool
    + returns true when the input is accepted by the DFA from its start state
    - returns false when the input contains characters with no transition
    # matching
