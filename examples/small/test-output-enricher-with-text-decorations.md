# Requirement: "enrich test runner output with text decorations"

Parses a line-oriented test runner's output and returns a decorated rendering with color markers around statuses.

std
  std.ansi
    std.ansi.colorize
      @ (text: string, color: string) -> string
      + wraps text in the ANSI escape for the named color and a reset
      + returns text unchanged when color is empty
      # terminal

test_output_enricher
  test_output_enricher.classify_line
    @ (line: string) -> line_kind
    + returns one of pass, fail, skip, run, summary, other based on prefix
    # classification
  test_output_enricher.decorate_line
    @ (line: string, kind: line_kind) -> string
    + returns the line with a color applied to the status keyword
    + returns the line unchanged for kind other
    # decoration
    -> std.ansi.colorize
  test_output_enricher.enrich
    @ (raw: string) -> string
    + returns the full output with every recognized line decorated
    # pipeline
