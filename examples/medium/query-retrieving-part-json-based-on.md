# Requirement: "a JSONPath query library"

Compiles a JSONPath expression, then evaluates it against a parsed JSON value. Two entry points in the project package.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses JSON text into a generic json_value tree
      - returns error on malformed input
      # parsing

jsonpath
  jsonpath.compile
    @ (expression: string) -> result[jsonpath_program, string]
    + returns a compiled program for expressions using $, ., .., [*], [index], and [?(filter)]
    - returns error on unbalanced brackets or unknown operators
    # compilation
  jsonpath.evaluate
    @ (program: jsonpath_program, root: json_value) -> list[json_value]
    + returns all matching json values in document order
    + returns an empty list when nothing matches
    # evaluation
  jsonpath.query
    @ (expression: string, raw_json: string) -> result[list[json_value], string]
    + convenience that compiles, parses, and evaluates in one call
    - returns error when either the expression or the JSON is invalid
    # convenience
    -> std.json.parse
