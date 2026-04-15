# Requirement: "a minimal json parser and transformer for unstructured json"

Walks tokens without building a full tree, letting the caller rewrite values in place.

std: (all units exist)

ujson
  ujson.walk
    fn (source: string, visitor: json_visitor) -> result[string, string]
    + invokes visitor for each key/value pair and returns the (possibly rewritten) output
    + preserves original whitespace between tokens
    - returns error on malformed json
    # traversal
  ujson.read_value
    fn (source: string, path: list[string]) -> optional[string]
    + returns the raw json text of the value at a dotted path
    - returns none when the path does not exist
    # query
  ujson.replace_value
    fn (source: string, path: list[string], replacement: string) -> result[string, string]
    + returns source with the value at path replaced by replacement
    - returns error when the path does not exist
    # transform
