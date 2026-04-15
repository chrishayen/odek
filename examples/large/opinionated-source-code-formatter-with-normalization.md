# Requirement: "an opinionated source code formatter that parses a file, normalizes its formatting, and returns the reformatted text"

Lex, parse, pretty-print with a canonical style. The project layer is tokens -> AST -> formatted text.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file into a string
      - returns error when the file cannot be opened
      # filesystem
    std.fs.write_all
      fn (path: string, data: string) -> result[void, string]
      + writes a string to a file, creating or truncating it
      - returns error when the path is not writable
      # filesystem

formatter
  formatter.tokenize
    fn (source: string) -> result[list[token], string]
    + produces a stream of tokens (indents, keywords, identifiers, literals, operators, newlines)
    - returns error on unterminated strings
    # lexing
  formatter.parse
    fn (tokens: list[token]) -> result[ast_node, string]
    + builds an abstract syntax tree from tokens
    - returns error on unexpected token
    # parsing
  formatter.canonicalize_strings
    fn (root: ast_node) -> ast_node
    + rewrites string literals to use a single preferred quote style
    # normalization
  formatter.normalize_whitespace
    fn (root: ast_node) -> ast_node
    + drops redundant blank lines and normalizes spacing around operators
    # normalization
  formatter.split_long_lines
    fn (root: ast_node, max_width: i32) -> ast_node
    + inserts line breaks inside collections and calls that exceed max_width
    # wrapping
  formatter.render
    fn (root: ast_node, indent_width: i32) -> string
    + pretty-prints the AST back to source text with canonical indentation
    # rendering
  formatter.format_source
    fn (source: string, max_width: i32, indent_width: i32) -> result[string, string]
    + end-to-end tokenize, parse, normalize, render
    - returns error on lex or parse failure
    # pipeline
    -> formatter.tokenize
    -> formatter.parse
    -> formatter.canonicalize_strings
    -> formatter.normalize_whitespace
    -> formatter.split_long_lines
    -> formatter.render
  formatter.format_file
    fn (path: string, max_width: i32, indent_width: i32) -> result[bool, string]
    + reads the file, formats it, and writes back only when the text changed
    + returns true when the file was rewritten
    # pipeline
    -> std.fs.read_all
    -> std.fs.write_all
    -> formatter.format_source
