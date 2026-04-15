# Requirement: "a YAML 1.2 parser and emitter"

A full YAML 1.2 library: scanner, parser, composer, and emitter. Project runes consume std primitives for IO and string handling.

std
  std.io
    std.io.read_all
      fn (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when path is not readable
      # io
    std.io.write_all
      fn (path: string, data: string) -> result[void, string]
      + writes data to path, creating or truncating
      - returns error when path is not writable
      # io
  std.text
    std.text.is_whitespace
      fn (c: i32) -> bool
      + returns true for space, tab, newline, carriage return
      # text
    std.text.is_digit
      fn (c: i32) -> bool
      + returns true for code points 0-9
      # text

yaml
  yaml.scan
    fn (source: string) -> result[list[yaml_token], string]
    + emits tokens for block sequences, block mappings, flow sequences, flow mappings, scalars, anchors, aliases, tags, and stream markers
    + tracks indentation to emit block-end tokens
    - returns error with line and column on invalid indentation
    # scanning
    -> std.text.is_whitespace
  yaml.parse_events
    fn (tokens: list[yaml_token]) -> result[list[yaml_event], string]
    + turns a token stream into stream-start, document-start, mapping, sequence, scalar, and end events
    - returns error on unbalanced flow collections
    # parsing
  yaml.compose
    fn (events: list[yaml_event]) -> result[yaml_node, string]
    + builds a node tree from an event stream, resolving anchors and aliases
    - returns error on dangling alias references
    # composition
  yaml.resolve_scalar
    fn (raw: string, style: yaml_style) -> yaml_node
    + resolves plain scalars to null, bool, int, float, or string per the YAML 1.2 core schema
    + leaves quoted scalars as strings
    # scalar_resolution
    -> std.text.is_digit
  yaml.parse
    fn (source: string) -> result[yaml_node, string]
    + returns the root node of the first document in the source
    - returns error when the source contains no documents
    # top_level_parse
  yaml.parse_all
    fn (source: string) -> result[list[yaml_node], string]
    + returns every document in a multi-document stream
    # multi_document_parse
  yaml.emit
    fn (node: yaml_node) -> string
    + serializes a node tree as a YAML 1.2 document
    + uses block style for collections and plain style for safe scalars
    # emission
  yaml.emit_all
    fn (nodes: list[yaml_node]) -> string
    + emits multiple documents separated by "---"
    # multi_document_emission
  yaml.load_file
    fn (path: string) -> result[yaml_node, string]
    + reads a file and parses its first document
    - returns error when file is unreadable or invalid YAML
    # file_load
    -> std.io.read_all
  yaml.dump_file
    fn (path: string, node: yaml_node) -> result[void, string]
    + serializes a node and writes it to the given path
    # file_dump
    -> std.io.write_all
