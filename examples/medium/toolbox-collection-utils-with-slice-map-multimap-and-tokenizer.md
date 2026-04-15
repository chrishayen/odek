# Requirement: "a collection utilities toolbox with slice, map, and multimap helpers plus a tokenizer"

A small grab-bag of pure data utilities. Kept narrow: a handful of collection helpers, plus one tokenizer.

std: (all units exist)

toolbox
  toolbox.unique
    fn (xs: list[string]) -> list[string]
    + returns elements in first-seen order with duplicates removed
    + returns an empty list for an empty input
    # collections
  toolbox.group_by_key
    fn (xs: list[key_value]) -> map[string, list[string]]
    + returns a multimap where each key maps to the list of its values in input order
    # collections
  toolbox.map_merge
    fn (a: map[string, string], b: map[string, string]) -> map[string, string]
    + returns a new map with b's entries overriding a's on key collision
    # collections
  toolbox.map_invert
    fn (m: map[string, string]) -> result[map[string, string], string]
    + returns a new map with keys and values swapped
    - returns error when two source keys share a value
    # collections
  toolbox.chunk
    fn (xs: list[string], size: i32) -> result[list[list[string]], string]
    + splits the list into chunks of at most size, in order
    - returns error when size is zero or negative
    # collections
  toolbox.tokenize
    fn (src: string) -> list[string]
    + splits on whitespace and punctuation, dropping empty tokens
    + returns an empty list for an empty or whitespace-only input
    # parsing
