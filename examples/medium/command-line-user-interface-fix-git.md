# Requirement: "a library for resolving merge conflicts in text files"

Parses conflict markers out of a file and lets callers pick a side per conflict, then renders the resolved file.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full contents of a file as text
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + writes content to path, replacing any existing file
      # filesystem

conflicts
  conflicts.parse
    @ (text: string) -> result[list[conflict_block], string]
    + returns each conflict as (before_lines, ours_lines, theirs_lines, after_lines)
    + returns a single non-conflict block when no markers are present
    - returns error on an unterminated conflict marker
    # parsing
  conflicts.load
    @ (path: string) -> result[list[conflict_block], string]
    + reads the file and parses its conflict blocks
    # loading
    -> std.fs.read_all
    -> conflicts.parse
  conflicts.resolve
    @ (blocks: list[conflict_block], choices: list[conflict_choice]) -> result[string, string]
    + renders the file using each choice (ours, theirs, or union) per conflict
    - returns error when choices length does not match the number of conflicts
    # resolution
  conflicts.save
    @ (path: string, resolved: string) -> result[void, string]
    + writes the resolved text back to disk
    # persistence
    -> std.fs.write_all
