# Requirement: "a library of extra iterator adaptors"

A handful of pure list combinators: chunks, windowed pairs, unique, interleave, and group-by-key.

std: (all units exist)

iterx
  iterx.chunks
    fn (items: list[string], size: i32) -> result[list[list[string]], string]
    + splits items into consecutive chunks of the given size, with a possibly shorter final chunk
    - returns error when size is less than 1
    # chunking
  iterx.windows
    fn (items: list[string], size: i32) -> result[list[list[string]], string]
    + returns all consecutive overlapping windows of the given size
    + returns an empty list when the input is shorter than size
    - returns error when size is less than 1
    # windowing
  iterx.unique
    fn (items: list[string]) -> list[string]
    + returns items in input order with duplicates removed
    # deduplication
  iterx.interleave
    fn (a: list[string], b: list[string]) -> list[string]
    + returns items taken alternately from a and b, continuing with the leftover tail when lengths differ
    # merging
  iterx.group_by
    fn (items: list[string], key_of: fn(string) -> string) -> map[string, list[string]]
    + groups items into buckets keyed by key_of, preserving input order inside each bucket
    # grouping
