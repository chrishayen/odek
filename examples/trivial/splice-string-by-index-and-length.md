# Requirement: "a library that removes or replaces part of a string by index and length"

A single splice function mirroring array-splice semantics on strings.

std: (all units exist)

splice_string
  splice_string.splice
    fn (s: string, start: i32, delete_count: i32, insert: string) -> string
    + removes delete_count characters starting at start and inserts the replacement
    + a negative start counts from the end of the string
    + clamps start and delete_count to the string's bounds
    ? indexes and lengths count code points, not bytes
    # string_editing
