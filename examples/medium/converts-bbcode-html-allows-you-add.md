# Requirement: "a BBCode to HTML converter with custom tag support"

Tokenize bracketed tags, walk them into a tree, and render each node via a registry that callers can extend.

std
  std.text
    std.text.find_char
      @ (s: string, c: u8, start: i32) -> i32
      + returns the index of the first occurrence of c at or after start
      - returns -1 when not found
      # text
  std.html
    std.html.escape_text
      @ (s: string) -> string
      + escapes &, <, >, and quotes for safe HTML text content
      # serialization

bbcode
  bbcode.tokenize
    @ (input: string) -> list[token]
    + splits the input into open tags, close tags, and text runs
    + unmatched brackets become literal text
    # lexing
    -> std.text.find_char
  bbcode.parse
    @ (tokens: list[token]) -> node
    + builds a tree of tag and text nodes
    + auto-closes dangling open tags at end of input
    # parsing
  bbcode.new_renderer
    @ () -> renderer
    + creates a renderer with built-in handlers for b, i, u, url, img, quote
    # construction
  bbcode.register_tag
    @ (r: renderer, name: string, handler: tag_handler) -> renderer
    + registers or overrides a handler for a named tag
    # extensibility
  bbcode.render
    @ (r: renderer, tree: node) -> string
    + walks the tree and emits HTML using the registered handlers
    + text nodes are HTML-escaped
    # rendering
    -> std.html.escape_text
