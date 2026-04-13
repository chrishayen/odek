# Requirement: "a JSONPath query engine with advanced filter expressions"

The path expression is compiled to a sequence of selectors; evaluation walks a JSON tree following those selectors. Filter expressions have their own tiny compiler.

std
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses JSON into a tagged value tree
      - returns error on malformed input
      # serialization
    std.json.get_field
      @ (value: json_value, name: string) -> optional[json_value]
      + returns the named field of an object
      # serialization
    std.json.array_items
      @ (value: json_value) -> list[json_value]
      + returns the elements of an array
      # serialization

jsonslice
  jsonslice.compile_path
    @ (expression: string) -> result[path_program, string]
    + compiles a path expression into a sequence of selectors
    - returns error on syntactically invalid expressions
    # compilation
  jsonslice.compile_filter
    @ (expression: string) -> result[filter_program, string]
    + compiles a bracketed filter expression
    - returns error on malformed comparison operators
    # compilation
  jsonslice.evaluate
    @ (program: path_program, raw: string) -> result[list[json_value], string]
    + parses the input and runs the compiled program, returning matched nodes
    - returns error when the input is not valid JSON
    # evaluation
    -> std.json.parse_value
  jsonslice.step
    @ (selector: path_selector, node: json_value) -> list[json_value]
    + applies a single selector to a node
    # evaluation
    -> std.json.get_field
    -> std.json.array_items
  jsonslice.apply_filter
    @ (program: filter_program, node: json_value) -> bool
    + evaluates a filter expression against a node and returns whether it matched
    # evaluation
  jsonslice.query
    @ (expression: string, raw: string) -> result[list[json_value], string]
    + one-shot compile + evaluate
    - returns error on compile or evaluation failure
    # convenience
