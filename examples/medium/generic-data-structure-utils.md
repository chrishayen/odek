# Requirement: "a collection of useful functions for working with generic data structures"

A small toolbox of slice and map operations. Each rune is a single, narrowly-scoped utility.

std: (all units exist)

just
  just.map_slice
    fn (items: list[string], fn_id: string) -> list[string]
    + returns a new slice where each element is the result of applying fn_id
    # slice
  just.filter_slice
    fn (items: list[string], predicate_id: string) -> list[string]
    + returns only the elements for which the predicate returns true
    # slice
  just.reduce_slice
    fn (items: list[string], initial: string, fn_id: string) -> string
    + folds the slice from left to right
    + returns initial when the slice is empty
    # slice
  just.unique
    fn (items: list[string]) -> list[string]
    + returns the input with duplicate elements removed, preserving order
    # slice
  just.group_by
    fn (items: list[string], key_fn_id: string) -> map[string, list[string]]
    + groups elements by the key each produces
    # slice
  just.map_keys
    fn (m: map[string, string]) -> list[string]
    + returns the keys of the map
    # map
  just.map_values
    fn (m: map[string, string]) -> list[string]
    + returns the values of the map
    # map
  just.merge
    fn (left: map[string, string], right: map[string, string]) -> map[string, string]
    + returns a new map containing entries from both, with right winning ties
    # map
