# Requirement: "a file type identification library based on content signatures"

Looks up file type by matching magic bytes at the start of a buffer against a signature table.

std
  std.fs
    std.fs.read_head
      fn (path: string, max_bytes: i32) -> result[bytes, string]
      + reads up to max_bytes from the beginning of a file
      - returns error when the file does not exist
      # filesystem

magic
  magic.new
    fn () -> magic_state
    + creates a detector preloaded with a built-in signature table
    # construction
  magic.load_signatures
    fn (state: magic_state, entries: list[signature_entry]) -> magic_state
    + returns a new state with additional signatures appended
    ? entries are matched in insertion order, first match wins
    # registry
  magic.identify_bytes
    fn (state: magic_state, data: bytes) -> optional[file_type]
    + returns the mime type and description when a signature matches
    - returns none when no signature matches
    # detection
  magic.identify_file
    fn (state: magic_state, path: string) -> result[optional[file_type], string]
    + returns the mime type for the file at the given path
    - returns error when the file cannot be read
    # detection
    -> std.fs.read_head
