# Requirement: "a library to find files with SQL-like queries"

Parses a SQL-like query over filesystem metadata, walks a root directory, and yields matching entries. The std layer provides parsing and filesystem primitives; the project layer wires them into a query engine.

std
  std.fs
    std.fs.walk
      fn (root: string) -> list[string]
      + yields every path under root recursively
      + includes both files and directories
      # filesystem
    std.fs.stat
      fn (path: string) -> result[file_info, string]
      + returns size, modification time, and type for path
      - returns error when path does not exist
      # filesystem
  std.text
    std.text.tokenize
      fn (input: string) -> list[string]
      + splits input into whitespace-separated tokens, preserving quoted strings
      # lexing
    std.text.glob_match
      fn (pattern: string, value: string) -> bool
      + returns true when value matches a glob pattern with * and ?
      # pattern_matching
  std.time
    std.time.parse_duration
      fn (input: string) -> result[i64, string]
      + parses strings like "7d", "2h", "30m" into seconds
      - returns error on unrecognized units
      # time

file_query
  file_query.parse
    fn (query: string) -> result[query_ast, string]
    + parses "select name, size from /dir where size > 1000" into a query_ast
    - returns error on unexpected tokens
    - returns error when from clause is missing
    # parsing
    -> std.text.tokenize
  file_query.parse_where
    fn (tokens: list[string]) -> result[where_clause, string]
    + parses comparison expressions joined by and/or
    - returns error on unbalanced parentheses
    # parsing
  file_query.parse_select
    fn (tokens: list[string]) -> result[list[string], string]
    + returns the list of requested columns
    - returns error when no columns are listed
    # parsing
  file_query.compile
    fn (ast: query_ast) -> compiled_query
    + produces a compiled_query with a root path and predicate
    # compilation
  file_query.evaluate_predicate
    fn (clause: where_clause, info: file_info, path: string) -> bool
    + returns true when file_info matches the clause
    + supports name, size, mtime, type columns
    # evaluation
    -> std.text.glob_match
  file_query.run
    fn (compiled: compiled_query) -> result[list[query_row], string]
    + walks the root and returns rows matching the predicate
    + each row projects only the selected columns
    - returns error when the root directory cannot be read
    # execution
    -> std.fs.walk
    -> std.fs.stat
  file_query.format_row
    fn (row: query_row) -> string
    + renders a row as a tab-separated line for display
    # rendering
  file_query.execute
    fn (query: string) -> result[list[query_row], string]
    + end-to-end parse, compile, and run
    - returns error on parse failure or execution failure
    # facade
    -> std.time.parse_duration
