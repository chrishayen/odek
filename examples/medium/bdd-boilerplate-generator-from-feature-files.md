# Requirement: "a boilerplate generator for behavior-driven tests from feature files"

Parses Gherkin-style feature text into a structured form, then renders a test skeleton.

std
  std.strings
    std.strings.split_lines
      fn (s: string) -> list[string]
      + splits on \n, dropping the trailing empty segment if present
      # strings
    std.strings.trim
      fn (s: string) -> string
      + returns s with leading and trailing whitespace removed
      # strings

bdd_gen
  bdd_gen.parse_feature
    fn (source: string) -> result[feature_doc, string]
    + returns a feature with its title, scenarios, and steps
    - returns error when no Feature: header is found
    - returns error on a step outside any scenario
    # parsing
    -> std.strings.split_lines
    -> std.strings.trim
  bdd_gen.step_signature
    fn (step: feature_step) -> string
    + returns a canonical function-name-safe identifier for a step
    ? punctuation is stripped; spaces become underscores
    # naming
  bdd_gen.render_skeleton
    fn (feature: feature_doc) -> string
    + returns source text with an empty test stub per scenario and per unique step
    # generation
  bdd_gen.extract_steps
    fn (feature: feature_doc) -> list[feature_step]
    + returns the deduplicated list of steps across all scenarios
    # analysis
