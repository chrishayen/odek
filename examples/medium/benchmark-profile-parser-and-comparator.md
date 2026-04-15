# Requirement: "a library for parsing and comparing benchmark profile results"

Text was an editor extension, but the reusable library is a benchmark parser and comparator.

std
  std.text
    std.text.split_lines
      fn (raw: string) -> list[string]
      + splits on "\n" or "\r\n" and drops the trailing empty line
      # text
    std.text.split_whitespace
      fn (line: string) -> list[string]
      + splits a line on runs of whitespace and drops empty tokens
      # text

benchmark_profile
  benchmark_profile.parse
    fn (raw: string) -> result[list[bench_result], string]
    + parses a benchmark log into a list of (name, iterations, ns_per_op, bytes_per_op, allocs_per_op)
    - returns error when no benchmark lines are found
    # parsing
    -> std.text.split_lines
    -> std.text.split_whitespace
  benchmark_profile.find
    fn (results: list[bench_result], name: string) -> optional[bench_result]
    + returns the result with the given name or none
    # lookup
  benchmark_profile.compare
    fn (before: list[bench_result], after: list[bench_result]) -> list[bench_delta]
    + returns a delta per benchmark present in both sides, with percent change in ns_per_op and bytes_per_op
    # comparison
  benchmark_profile.classify
    fn (delta: bench_delta, threshold_pct: f64) -> string
    + returns "improved", "regressed", or "unchanged" based on the threshold
    # classification
  benchmark_profile.format_table
    fn (deltas: list[bench_delta]) -> string
    + returns a fixed-width text table with name, before, after, and change columns
    # rendering
