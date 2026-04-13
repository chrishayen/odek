# Requirement: "a linter that detects literal values that could be replaced by standard-library constants"

Given source code and a catalogue of known standard-library constants, reports each literal that matches one.

std
  std.source
    std.source.parse
      @ (src: string) -> result[ast_node, string]
      + parses source text into an abstract syntax tree
      - returns error on syntax failure
      # parsing
    std.source.walk
      @ (root: ast_node) -> list[ast_node]
      + returns every node in pre-order
      # traversal

stdlibvars
  stdlibvars.load_catalogue
    @ (entries: list[tuple[string,string]]) -> catalogue_state
    + builds a catalogue from (constant_name, literal_value) pairs
    # construction
  stdlibvars.analyze
    @ (src: string, cat: catalogue_state) -> result[list[suggestion], string]
    + returns a suggestion for every literal matching a catalogued value
    - returns error when source fails to parse
    # analysis
    -> std.source.parse
    -> std.source.walk
  stdlibvars.literal_value
    @ (node: ast_node) -> optional[string]
    + returns the literal's textual value if node is a string or numeric literal, else none
    # extraction
  stdlibvars.format_suggestion
    @ (s: suggestion) -> string
    + formats as "file:line: replace \"value\" with const_name"
    # reporting
