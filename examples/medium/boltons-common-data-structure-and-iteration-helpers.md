# Requirement: "a grab-bag utilities library for common data-structure and iteration helpers"

A small collection of general-purpose helpers for lists, strings, and dictionaries.

std: (all units exist)

boltons
  boltons.chunked
    fn (items: list[string], size: i32) -> list[list[string]]
    + splits items into consecutive chunks of at most size
    - returns an empty list when items is empty
    ? size must be > 0
    # iteration
  boltons.windowed
    fn (items: list[string], size: i32) -> list[list[string]]
    + returns all contiguous sliding windows of the given size
    - returns an empty list when size exceeds the length of items
    # iteration
  boltons.unique
    fn (items: list[string]) -> list[string]
    + returns items in original order with duplicates removed
    # sets
  boltons.group_by
    fn (items: list[string], key_fn: fn_string_to_string) -> map[string, list[string]]
    + groups items by the result of key_fn preserving input order within each group
    # grouping
  boltons.flatten
    fn (nested: list[list[string]]) -> list[string]
    + concatenates the inner lists into a single list
    # iteration
  boltons.merge_maps
    fn (a: map[string,string], b: map[string,string]) -> map[string,string]
    + returns a new map containing all entries from a overlaid by b
    # maps
  boltons.invert_map
    fn (m: map[string,string]) -> map[string,string]
    + returns a map swapping keys and values
    ? if values repeat, the last one wins
    # maps
  boltons.camel_to_snake
    fn (text: string) -> string
    + converts CamelCase to snake_case
    + treats acronym runs as a single group (e.g. "HTTPServer" -> "http_server")
    # strings
  boltons.indent_lines
    fn (text: string, prefix: string) -> string
    + prepends prefix to every line in text, preserving trailing newlines
    # strings
