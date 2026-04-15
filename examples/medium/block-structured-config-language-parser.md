# Requirement: "a parser for a block-structured configuration language"

Tokenizer and parser in the project; only a generic string scanner lives in std.

std
  std.strings
    std.strings.is_digit
      fn (c: string) -> bool
      + returns true when c is a single ASCII digit
      # strings
    std.strings.is_alpha
      fn (c: string) -> bool
      + returns true when c is a single ASCII letter or underscore
      # strings

config_lang
  config_lang.tokenize
    fn (source: string) -> result[list[config_token], string]
    + produces tokens for identifiers, strings, numbers, braces, equals, and newlines
    - returns error with line and column on unterminated strings
    - returns error on unexpected characters
    # lexing
    -> std.strings.is_digit
    -> std.strings.is_alpha
  config_lang.parse
    fn (tokens: list[config_token]) -> result[config_node, string]
    + returns a tree of blocks and assignments
    - returns error on missing closing brace
    - returns error on assignment without a value
    # parsing
  config_lang.load
    fn (source: string) -> result[config_node, string]
    + convenience entry point: tokenize then parse
    - returns the first error encountered
    # loading
  config_lang.get_string
    fn (node: config_node, path: string) -> result[string, string]
    + returns the string value at a dotted path like "server.host"
    - returns error when the path is missing or not a string
    # lookup
  config_lang.get_int
    fn (node: config_node, path: string) -> result[i64, string]
    + returns the integer value at a dotted path
    - returns error when the path is missing or not an integer
    # lookup
