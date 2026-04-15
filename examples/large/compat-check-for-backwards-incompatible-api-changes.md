# Requirement: "a library that checks recent changes to a project for backwards-incompatible API changes"

Builds a symbol table from each revision's public interface and reports the diff classified by severity. Works over a source tree accessed via std filesystem primitives.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file
      - returns error when the path does not exist
      # filesystem
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns every file path under root
      - returns error when root is not a directory
      # filesystem
  std.process
    std.process.run
      fn (cmd: string, args: list[string]) -> result[string, string]
      + runs a command and returns stdout
      - returns error on non-zero exit
      # process

compat_check
  compat_check.parse_source
    fn (source: string) -> result[list[symbol], string]
    + returns the public symbols declared in a single source file
    - returns error on syntax error
    # parsing
  compat_check.build_symbol_table
    fn (root: string) -> result[symbol_table, string]
    + returns a symbol table keyed by fully qualified name across every source file under root
    - returns error when any file fails to parse
    # analysis
    -> std.fs.walk
    -> std.fs.read_all
  compat_check.snapshot_at_revision
    fn (repo_path: string, revision: string) -> result[symbol_table, string]
    + returns the symbol table for the working tree at the given revision
    - returns error when the revision is unknown
    # analysis
    -> std.process.run
  compat_check.diff_tables
    fn (before: symbol_table, after: symbol_table) -> list[api_change]
    + returns one change record per added, removed, or modified public symbol
    # diffing
  compat_check.classify_change
    fn (change: api_change) -> change_severity
    + returns "breaking" for removals and signature changes
    + returns "non-breaking" for additions
    + returns "warning" for changes that might be semver-minor
    # classification
  compat_check.filter_breaking
    fn (changes: list[api_change]) -> list[api_change]
    + returns only changes whose severity is breaking
    # filtering
  compat_check.format_report
    fn (changes: list[api_change]) -> string
    + renders a human-readable diff report grouped by severity
    # rendering
  compat_check.symbol_signature
    fn (sym: symbol) -> string
    + returns a canonical signature string used for comparing symbols across revisions
    # normalization
  compat_check.compare_signatures
    fn (a: symbol, b: symbol) -> bool
    + returns true when two symbols have identical canonical signatures
    # comparison
  compat_check.run_against_revisions
    fn (repo_path: string, base: string, head: string) -> result[list[api_change], string]
    + returns the classified diff between two revisions of the repository
    - returns error when either revision cannot be checked out
    # orchestration
