# Requirement: "a placeholder and wildcard text parser for command templates"

Match user input against a template like "say <word> to <target>" and extract the named placeholders.

std: (all units exist)

allot
  allot.compile
    @ (template: string) -> result[pattern_state, string]
    + parses a template with <name> placeholders and returns a matcher
    - returns error when a placeholder is unclosed
    # compilation
  allot.match
    @ (pattern: pattern_state, input: string) -> result[map[string,string], string]
    + returns the placeholder name to captured value map when input matches
    - returns error when the template literals do not match input
    # matching
  allot.placeholder_names
    @ (pattern: pattern_state) -> list[string]
    + returns the placeholder names in declaration order
    # introspection
