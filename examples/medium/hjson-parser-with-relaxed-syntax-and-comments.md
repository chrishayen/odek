# Requirement: "a parser for a human-friendly JSON variant with relaxed syntax and comments"

The format permits unquoted keys, trailing commas, and line/block comments. Tokenization and parsing are split so the parser consumes a flat stream.

std
  std.text
    std.text.is_digit
      fn (c: u8) -> bool
      + returns true for ASCII 0-9
      # text
    std.text.is_ident_start
      fn (c: u8) -> bool
      + returns true for letters and underscore
      # text

hjson
  hjson.tokenize
    fn (source: string) -> result[list[token], string]
    + produces tokens for braces, brackets, commas, colons, strings, numbers, and identifiers
    + skips // line comments and /* block comments
    - returns error on an unterminated string or block comment
    # tokenization
    -> std.text.is_digit
    -> std.text.is_ident_start
  hjson.parse
    fn (source: string) -> result[hjson_value, string]
    + parses a document into a tagged value tree
    + accepts unquoted object keys and trailing commas
    - returns error on a missing closing brace or bracket
    # parsing
  hjson.get_string
    fn (value: hjson_value, path: list[string]) -> optional[string]
    + returns the string at the given object path
    - returns none when the path is missing or the leaf is not a string
    # access
  hjson.get_number
    fn (value: hjson_value, path: list[string]) -> optional[f64]
    + returns the numeric value at the given object path
    - returns none when the path is missing or the leaf is not numeric
    # access
  hjson.to_strict_json
    fn (value: hjson_value) -> string
    + emits the value as standard JSON with quoted keys and no comments
    # serialization
