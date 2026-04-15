# Requirement: "a spec-compliant HTML parser"

A tokenizer and tree constructor exposed through a single parse entry point, producing a DOM-like node tree.

std
  std.strings
    std.strings.to_lower
      fn (s: string) -> string
      + returns the lowercase form of an ASCII string
      # strings
    std.strings.starts_with
      fn (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings
  std.unicode
    std.unicode.decode_utf8
      fn (data: bytes) -> result[list[i32], string]
      + decodes a utf-8 byte stream into codepoints
      - returns error on invalid sequences
      # unicode

html_parser
  html_parser.tokenize
    fn (source: string) -> list[html_token]
    + splits source into start-tag, end-tag, text, comment, and doctype tokens
    + handles attributes with single, double, and unquoted values
    - emits a parse-error token for unterminated tags
    # tokenization
    -> std.strings.to_lower
  html_parser.parse
    fn (source: string) -> html_document
    + builds a document tree from the token stream using the insertion-mode state machine
    + inserts implicit html, head, and body elements when missing
    + reparents misnested elements following the spec's adoption agency rules
    # tree_construction
  html_parser.serialize
    fn (doc: html_document) -> string
    + round-trips a parsed document back to canonical html text
    # serialization
  html_parser.query_tag
    fn (doc: html_document, tag: string) -> list[html_node]
    + returns every element node whose tag name matches
    # traversal
    -> std.strings.to_lower
