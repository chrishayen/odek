# Requirement: "a parallelized formatter for BDD feature files"

Formats feature files (scenarios, steps, tables) with consistent indentation across many files concurrently. The project layer handles parsing, formatting, and a work dispatcher; std supplies file IO and a parallel map.

std
  std.fs
    std.fs.read_text
      @ (path: string) -> result[string, string]
      + returns the file contents as UTF-8 text
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_text
      @ (path: string, data: string) -> result[void, string]
      + writes data to path, truncating any existing file
      - returns error when the path cannot be written
      # filesystem
  std.concurrency
    std.concurrency.parallel_map
      @ (items: list[string], worker_count: i32, fn: fn(string) -> result[string, string]) -> list[result[string, string]]
      + applies fn to each item across worker_count workers, preserving input order
      + uses at most one worker per item when items are fewer than worker_count
      # parallelism

featfmt
  featfmt.parse
    @ (source: string) -> result[feature_ast, string]
    + returns an AST with the feature, background, and scenarios populated
    - returns error when a step keyword is unknown
    # parsing
  featfmt.format_ast
    @ (ast: feature_ast) -> string
    + returns the canonical pretty-printed form with 2-space indentation
    + aligns table columns to the widest cell
    # formatting
  featfmt.format_source
    @ (source: string) -> result[string, string]
    + returns the formatted source
    - returns error when parsing fails
    # formatting
  featfmt.format_file
    @ (path: string) -> result[void, string]
    + reads, formats, and rewrites the file in place
    - returns error when reading, parsing, or writing fails
    # file_pipeline
    -> std.fs.read_text
    -> std.fs.write_text
  featfmt.format_files
    @ (paths: list[string], worker_count: i32) -> list[result[string, string]]
    + returns per-path results in the same order as paths, running formatting in parallel
    # orchestration
    -> std.concurrency.parallel_map
