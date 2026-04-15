# Requirement: "verifying that JSON payloads match an expected shape"

A JSON assertion library: given an actual document and an expected document, report structural and value differences. The project layer walks parsed trees; std supplies the JSON parser.

std
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document into a tagged tree
      - returns error on malformed JSON
      # serialization
    std.json.kind
      fn (value: json_value) -> string
      + returns one of "null", "bool", "number", "string", "array", "object"
      # introspection
    std.json.as_object
      fn (value: json_value) -> result[map[string, json_value], string]
      + returns the object entries when value is an object
      - returns error when value is not an object
      # introspection
    std.json.as_array
      fn (value: json_value) -> result[list[json_value], string]
      + returns the array elements when value is an array
      - returns error when value is not an array
      # introspection

jsonassert
  jsonassert.compare
    fn (expected: string, actual: string) -> result[list[json_diff], string]
    + returns an empty diff list when the documents are structurally equal
    + returns per-path mismatches when types or values differ
    - returns error when either input is not valid JSON
    # comparison
    -> std.json.parse
  jsonassert.compare_values
    fn (expected: json_value, actual: json_value, path: string) -> list[json_diff]
    + returns diffs whose path is rooted at the given path
    + recursively descends into matching objects and arrays
    # comparison
    -> std.json.kind
    -> std.json.as_object
    -> std.json.as_array
  jsonassert.format_diff
    fn (diffs: list[json_diff]) -> string
    + returns a human-readable report with one line per diff
    + returns "OK" when the list is empty
    # presentation
  jsonassert.assert_equal
    fn (expected: string, actual: string) -> result[void, string]
    + returns void when the documents match
    - returns error whose message is the formatted diff when they differ
    # assertion
