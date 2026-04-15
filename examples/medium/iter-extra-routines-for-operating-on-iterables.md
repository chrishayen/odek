# Requirement: "additional routines for operating on iterables"

A small grab-bag of list combinators that are commonly reimplemented ad-hoc: chunking, sliding windows, grouping, interleaving, unique.

std: (all units exist)

iter_extra
  iter_extra.chunked
    fn (items: list[bytes], size: i32) -> list[list[bytes]]
    + splits the list into fixed-size chunks, final chunk may be smaller
    - returns empty list when size <= 0
    # chunking
  iter_extra.windowed
    fn (items: list[bytes], size: i32) -> list[list[bytes]]
    + returns all contiguous sublists of the given size
    - returns empty list when size > len(items) or size <= 0
    # windowing
  iter_extra.unique
    fn (items: list[bytes]) -> list[bytes]
    + returns items with duplicates removed, preserving first occurrence order
    # dedup
  iter_extra.group_by
    fn (items: list[bytes], key: fn[bytes, string]) -> map[string, list[bytes]]
    + groups items by key while preserving input order within each group
    # grouping
  iter_extra.interleave
    fn (a: list[bytes], b: list[bytes]) -> list[bytes]
    + returns items alternating between a and b until both are exhausted
    + extends with the remainder of the longer list
    # combining
  iter_extra.partition
    fn (items: list[bytes], predicate: fn[bytes, bool]) -> tuple[list[bytes], list[bytes]]
    + returns (matching, non_matching) preserving order
    # splitting
  iter_extra.flatten
    fn (items: list[list[bytes]]) -> list[bytes]
    + concatenates all inner lists in order
    # combining
