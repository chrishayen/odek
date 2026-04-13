# Requirement: "a library for projecting JSON documents to a caller-specified subset of fields"

Given a JSON document and a field selection expression, return a reshaped JSON containing only the requested fields.

std: (all units exist)

dynjson
  dynjson.parse_selection
    @ (expr: string) -> result[selection, string]
    + parses a dotted and bracketed field expression like "user{id,name},posts{title}"
    - returns error on unbalanced braces
    # parsing
  dynjson.project
    @ (document: string, sel: selection) -> result[string, string]
    + returns a JSON string containing only the selected fields
    + preserves array structure for list-typed selections
    - returns error when the document is not valid JSON
    # projection
  dynjson.project_value
    @ (value: json_value, sel: selection) -> json_value
    + returns a json value with only the selected subtree
    # projection
  dynjson.merge_selections
    @ (a: selection, b: selection) -> selection
    + returns the union of two selections
    # combination
