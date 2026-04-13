# Requirement: "a utility to compare micro-benchmarks"

Parses two benchmark result sets and reports per-name deltas.

std
  std.text
    std.text.split_lines
      @ (raw: string) -> list[string]
      + splits on newlines and drops a trailing empty line
      # text

benchcmp
  benchcmp.parse_results
    @ (raw: string) -> result[map[string, f64], string]
    + parses "name  duration" lines into a name-to-nanoseconds map
    - returns error when a line has no numeric duration
    # parsing
    -> std.text.split_lines
  benchcmp.compare
    @ (old_results: map[string, f64], new_results: map[string, f64]) -> list[tuple[string, f64, f64, f64]]
    + returns (name, old_ns, new_ns, percent_delta) for every name present in both sets
    + percent_delta is (new - old) / old
    - omits names that appear in only one of the two sets
    # comparison
