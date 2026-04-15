# Requirement: "a helper library for working with ordered sequences"

Convenience operations over ordered generic sequences: mapping, filtering, uniqueness, chunking.

std: (all units exist)

sequence
  sequence.map
    fn (items: list[string], transform: string) -> list[string]
    + applies the named transform to each element and returns a new list
    + returns empty list for empty input
    # transformation
  sequence.filter
    fn (items: list[string], predicate: string) -> list[string]
    + keeps only elements for which the named predicate returns true
    + returns empty list when nothing matches
    # transformation
  sequence.unique
    fn (items: list[string]) -> list[string]
    + returns elements in original order with duplicates removed
    # deduplication
  sequence.chunk
    fn (items: list[string], size: i32) -> result[list[list[string]], string]
    + splits the list into sublists of the given size
    + returns a shorter final sublist when the length is not divisible
    - returns error when size is less than 1
    # chunking
