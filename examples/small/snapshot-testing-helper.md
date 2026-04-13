# Requirement: "a snapshot testing helper"

Compares a value's canonical form against a stored snapshot file, writing the file on first run.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads the entire file as text
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes text, creating parent directories when needed
      # filesystem
    std.fs.exists
      @ (path: string) -> bool
      + returns true when a file exists at the path
      # filesystem

snapshot
  snapshot.canonicalize
    @ (value: string) -> string
    + normalizes line endings to lf and strips trailing whitespace per line
    # normalization
  snapshot.match
    @ (snapshot_path: string, actual: string) -> result[bool, string]
    + returns true when the stored snapshot equals the canonicalized value
    + creates the snapshot file with the value on first run and returns true
    - returns false with a non-empty diff path when stored snapshot differs
    # matching
    -> std.fs.exists
    -> std.fs.read_all
    -> std.fs.write_all
  snapshot.update
    @ (snapshot_path: string, actual: string) -> result[void, string]
    + overwrites the snapshot with the canonicalized value
    # maintenance
    -> std.fs.write_all
