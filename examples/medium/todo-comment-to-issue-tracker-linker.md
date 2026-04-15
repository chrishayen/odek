# Requirement: "a static analyzer that links TODO comments in source to issues in an issue tracker"

Scans source files for TODO comments that reference issue ids, looks up their status through a pluggable tracker, and flags TODOs whose issue is closed or missing.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the path does not exist
      # filesystem
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns all file paths under root in depth-first order
      - returns error when root is not a directory
      # filesystem

todo_linker
  todo_linker.scan_file
    fn (path: string, content: string) -> list[todo_item]
    + returns a todo item per comment starting with TODO, capturing line number and trailing text
    + recognizes line comments and the first line of block comments
    # scanning
  todo_linker.extract_issue_id
    fn (text: string, patterns: list[string]) -> optional[string]
    + returns the first issue id matching any pattern
    ? patterns use a simple syntax: a literal prefix followed by digits
    # parsing
  todo_linker.scan_tree
    fn (root: string, patterns: list[string]) -> result[list[todo_item], string]
    + returns all todos across files under root, each with its optional issue id filled in
    - returns error when a file cannot be read
    # orchestration
    -> std.fs.walk
    -> std.fs.read_all
    -> todo_linker.scan_file
    -> todo_linker.extract_issue_id
  todo_linker.resolve_statuses
    fn (items: list[todo_item], fetch: fn(string) -> result[issue_status, string]) -> list[resolved_todo]
    + returns each todo paired with the status returned by fetch, or "unknown" when fetch errors
    + caches fetch results so each issue id is queried once
    # resolution
  todo_linker.select_violations
    fn (resolved: list[resolved_todo]) -> list[resolved_todo]
    + returns todos whose issue is closed or whose issue was not found
    + returns todos with no issue id when that policy is enabled
    ? missing ids are considered a violation by default
    # policy
  todo_linker.format_report
    fn (violations: list[resolved_todo]) -> string
    + returns a newline-separated report with file, line, issue id and status
    + returns "no violations" when the list is empty
    # reporting
