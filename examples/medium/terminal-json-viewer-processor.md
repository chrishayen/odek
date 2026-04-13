# Requirement: "a json viewer and query processor"

Parses json, renders a collapsible tree, and evaluates path-style queries against the parsed document.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses any json value
      - returns error on malformed input
      # serialization

json_viewer
  json_viewer.build_tree
    @ (value: json_value) -> tree_node
    + constructs a renderable tree with collapsed/expanded state per node
    # tree_construction
    -> std.json.parse
  json_viewer.toggle_node
    @ (tree: tree_node, path: list[string]) -> tree_node
    + flips the expanded state at path
    ? path is unchanged when no node matches
    # tree_mutation
  json_viewer.render
    @ (tree: tree_node, max_depth: i32) -> string
    + produces an indented text rendering of expanded nodes up to max_depth
    + collapsed objects render as {...}
    # rendering
  json_viewer.query
    @ (value: json_value, path: string) -> result[json_value, string]
    + evaluates a dotted path like ".users[0].name" against value
    - returns error when a segment does not exist
    - returns error when indexing into a non-array
    # query
  json_viewer.pretty_print
    @ (value: json_value, indent: i32) -> string
    + returns canonical pretty-printed json with the given indent
    # formatting
