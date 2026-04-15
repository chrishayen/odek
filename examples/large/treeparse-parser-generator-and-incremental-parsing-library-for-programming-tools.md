# Requirement: "a parser generator and incremental parsing library for programming tools"

The project has two faces: building parse tables from a grammar and running a parser over source text with incremental re-parsing on edits.

std
  std.collections
    std.collections.map_get
      fn (m: map[string, string], key: string) -> optional[string]
      + returns the value for key when present
      # collections
    std.collections.list_append
      fn (xs: list[i32], x: i32) -> list[i32]
      + returns xs with x appended
      # collections
  std.strings
    std.strings.byte_at
      fn (s: string, i: i32) -> u8
      + returns the byte at index i
      # strings

treeparse
  treeparse.grammar_new
    fn () -> grammar
    + returns an empty grammar ready for rule additions
    # grammar
  treeparse.grammar_add_rule
    fn (g: grammar, name: string, body: rule_body) -> grammar
    + registers a named production
    # grammar
  treeparse.grammar_compile
    fn (g: grammar) -> result[parse_table, string]
    + builds a parse table from the grammar rules
    - returns error on left-recursion or unresolved references
    # compilation
    -> std.collections.map_get
  treeparse.parser_new
    fn (table: parse_table) -> parser
    + constructs a parser bound to a compiled table
    # parsing
  treeparse.parse
    fn (p: parser, source: string) -> result[syntax_tree, parse_error]
    + returns a full syntax tree for source
    - returns a parse_error pointing at the first unexpected token
    # parsing
    -> std.strings.byte_at
  treeparse.reparse_incremental
    fn (p: parser, old_tree: syntax_tree, edits: list[text_edit], new_source: string) -> result[syntax_tree, parse_error]
    + reuses unchanged subtrees from old_tree when they overlap unedited regions
    + falls back to a full parse when edits invalidate the root
    # incremental
    -> std.collections.list_append
  treeparse.tree_root
    fn (t: syntax_tree) -> syntax_node
    + returns the root node of the syntax tree
    # tree
  treeparse.node_children
    fn (n: syntax_node) -> list[syntax_node]
    + returns the immediate children of n
    # tree
  treeparse.node_kind
    fn (n: syntax_node) -> string
    + returns the production name that produced this node
    # tree
  treeparse.node_range
    fn (n: syntax_node) -> tuple[i32, i32]
    + returns (start_byte, end_byte) of the node's extent in the source
    # tree
  treeparse.walk
    fn (t: syntax_tree, visit: node_visitor) -> void
    + performs a pre-order traversal invoking visit on each node
    # traversal
  treeparse.query
    fn (t: syntax_tree, pattern: string) -> result[list[syntax_node], string]
    + returns nodes matching the s-expression pattern
    - returns error on malformed pattern
    # query
