# Requirement: "a library for pattern-matching for filtering messages and events"

Patterns are compiled once into an automaton; event matching walks the automaton over the event's field values.

std
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses JSON into a generic value tree
      - returns error on malformed JSON
      # serialization

matcher
  matcher.compile_pattern
    @ (pattern: string) -> result[compiled_pattern, string]
    + compiles a pattern describing required fields and allowed values
    - returns error when the pattern is not valid JSON
    - returns error when a field uses an unknown matcher keyword
    # pattern_compilation
    -> std.json.parse_value
  matcher.add_pattern
    @ (state: matcher_state, name: string, pattern: compiled_pattern) -> matcher_state
    + registers a named pattern in the automaton
    ? patterns sharing prefixes reuse automaton nodes
    # automaton_build
  matcher.match_event
    @ (state: matcher_state, event: string) -> result[list[string], string]
    + returns the names of every pattern that matches the event
    + returns an empty list when nothing matches
    - returns error when the event is not valid JSON
    # matching
    -> std.json.parse_value
  matcher.new
    @ () -> matcher_state
    + creates an empty matcher with no patterns
    # construction
