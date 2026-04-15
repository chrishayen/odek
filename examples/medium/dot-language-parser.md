# Requirement: "a Graphviz DOT language parser"

Returns a graph AST with nodes, edges, attributes, and subgraphs. std supplies generic string scanning.

std
  std.strings
    std.strings.byte_at
      fn (s: string, i: i32) -> u8
      + returns the byte at index i
      # strings
    std.strings.slice
      fn (s: string, start: i32, end: i32) -> string
      + returns the substring between byte offsets
      # strings

dot
  dot.tokenize
    fn (source: string) -> result[list[dot_token], string]
    + emits idents, numerals, strings, punctuation, and edge operators
    - returns error on unterminated strings
    # tokenization
    -> std.strings.byte_at
    -> std.strings.slice
  dot.parse
    fn (source: string) -> result[dot_graph, string]
    + returns an AST with graph type, id, and statement list
    - returns error on missing braces or unknown statements
    # parsing
  dot.is_directed
    fn (g: dot_graph) -> bool
    + returns true for "digraph", false for "graph"
    # query
  dot.nodes
    fn (g: dot_graph) -> list[dot_node]
    + returns every declared node with its attributes
    # query
  dot.edges
    fn (g: dot_graph) -> list[dot_edge]
    + returns every edge with source, target, and attributes
    # query
  dot.subgraphs
    fn (g: dot_graph) -> list[dot_graph]
    + returns every nested subgraph
    # query
  dot.attribute
    fn (target: dot_node, key: string) -> optional[string]
    + returns the attribute value for key when set
    # query
