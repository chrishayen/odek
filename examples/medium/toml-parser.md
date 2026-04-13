# Requirement: "a TOML parser"

Tokenize once, then parse into a flat key-path map. Real TOML has more types; this decomposition keeps the surface small.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads an entire file as a UTF-8 string
      - returns error when the file does not exist
      # io

toml
  toml.tokenize
    @ (source: string) -> result[list[toml_token], string]
    + splits source into section headers, key-value pairs, and values
    + ignores comments starting with #
    - returns error on unterminated strings
    # lexing
  toml.parse
    @ (source: string) -> result[map[string, toml_value], string]
    + returns a flat map keyed by dotted paths like "server.port"
    + supports string, integer, float, boolean, and array values
    - returns error when a key appears twice in the same table
    # parsing
  toml.parse_file
    @ (path: string) -> result[map[string, toml_value], string]
    + reads and parses a TOML file in one step
    - returns error when the file cannot be read
    # parsing
    -> std.fs.read_all
  toml.get_string
    @ (doc: map[string, toml_value], key: string) -> optional[string]
    + returns the string at the given dotted key
    - returns none when the key is missing or not a string
    # access
  toml.get_int
    @ (doc: map[string, toml_value], key: string) -> optional[i64]
    + returns the integer at the given dotted key
    - returns none when the key is missing or not an integer
    # access
