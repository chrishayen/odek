# Requirement: "a library that flags non-inclusive language in source code"

Scans text for configurable terms and returns findings with locations and suggestions.

std: (all units exist)

woke
  woke.load_rules
    @ (raw: string) -> result[rule_set, string]
    + parses a ruleset describing terms, severity, and suggested alternatives
    - returns error on malformed ruleset
    # rules
  woke.scan_text
    @ (rules: rule_set, source: string) -> list[finding]
    + returns each occurrence of a flagged term with line and column
    + returns empty list when nothing matches
    # scanning
  woke.format_finding
    @ (f: finding) -> string
    + formats a finding as a single line with location, term, and suggestion
    # reporting
