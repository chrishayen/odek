# Requirement: "a template engine that compiles templates into executable form"

Templates are compiled once into an opcode list; rendering executes the opcodes against a context map.

std: (all units exist)

template
  template.tokenize
    fn (source: string) -> result[list[token], string]
    + splits template source into literal and expression tokens
    - returns error on unterminated expression delimiters
    # lexing
  template.parse
    fn (tokens: list[token]) -> result[ast_node, string]
    + builds an AST with literal, substitution, if, and for nodes
    - returns error on mismatched block openers and closers
    # parsing
  template.compile
    fn (ast: ast_node) -> compiled_template
    + lowers the AST into a flat opcode list
    ? compilation resolves jump targets for if and for blocks
    # compilation
  template.render
    fn (tpl: compiled_template, context: map[string,string]) -> result[string, string]
    + executes the opcodes against the context and returns the output
    - returns error when a substituted name is missing from the context
    # rendering
  template.render_safe
    fn (tpl: compiled_template, context: map[string,string]) -> string
    + renders, substituting empty strings for missing names instead of failing
    # rendering
