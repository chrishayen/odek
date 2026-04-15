# Requirement: "a batch file rename library"

The project layer plans renames from a pattern and executes them; filesystem and regex primitives live in std.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns the names of entries directly inside path
      - returns error when path is not a directory
      # filesystem
    std.fs.rename
      fn (from: string, to: string) -> result[void, string]
      + renames a file or directory
      - returns error when the source does not exist or destination already exists
      # filesystem
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex_handle, string]
      + compiles a regular expression
      - returns error on invalid syntax
      # regex
    std.regex.replace_all
      fn (re: regex_handle, input: string, replacement: string) -> string
      + returns input with every match replaced, supporting $1..$n backreferences
      # regex

batch_rename
  batch_rename.plan
    fn (dir: string, pattern: string, replacement: string) -> result[list[rename_op], string]
    + returns the list of (from, to) pairs that would be applied
    - returns error when pattern fails to compile
    - skips entries whose name does not match the pattern
    # planning
    -> std.fs.list_dir
    -> std.regex.compile
    -> std.regex.replace_all
  batch_rename.validate
    fn (ops: list[rename_op]) -> result[void, string]
    + returns ok when no two ops target the same destination
    - returns error describing the first collision detected
    # validation
  batch_rename.apply
    fn (ops: list[rename_op]) -> result[i32, string]
    + performs each rename in order and returns the count applied
    - stops and returns error on the first failing rename
    # execution
    -> std.fs.rename
  batch_rename.dry_run
    fn (ops: list[rename_op]) -> list[string]
    + returns human-readable "from -> to" lines, one per op
    # reporting
