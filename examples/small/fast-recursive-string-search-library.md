# Requirement: "a fast recursive string search library"

Walks a directory tree and returns files and line positions that contain a literal substring.

std
  std.fs
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns every regular file under root, recursively
      - returns error when root is not a directory
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns file contents as a string
      - returns error when the file cannot be read
      # filesystem

stringsearch
  stringsearch.search_text
    @ (text: string, needle: string) -> list[match]
    + returns one match per occurrence with line number and byte offset
    - returns an empty list when needle does not occur
    ? line numbers are 1-based
    # matching
  stringsearch.search_file
    @ (path: string, needle: string) -> result[list[match], string]
    + returns matches within the file
    - returns error when the file cannot be read
    # matching
    -> std.fs.read_all
    -> stringsearch.search_text
  stringsearch.search_tree
    @ (root: string, needle: string) -> result[map[string, list[match]], string]
    + returns a map from relative path to matches for every file containing needle
    - returns error when root is not a directory
    + silently skips files that cannot be read
    # orchestration
    -> std.fs.walk
    -> stringsearch.search_file
