# Requirement: "an encoder and decoder for a TOML-like configuration format"

A full subsystem: tokenize, parse, build a typed value tree, and emit it back as text. Values can be strings, integers, floats, booleans, arrays, or nested tables.

std: (all units exist)

toml
  toml.tokenize
    @ (source: string) -> result[list[toml_token], string]
    + splits the input into keys, equals, strings, numbers, booleans, brackets, and newlines
    - returns error on an unterminated string literal
    - returns error on an invalid number literal
    # lexing
  toml.parse
    @ (tokens: list[toml_token]) -> result[toml_value, string]
    + returns the root table value containing every top-level key
    + supports nested tables via "[section.sub]" headers
    + supports inline arrays of homogeneous type
    - returns error on duplicate keys within the same table
    - returns error on unclosed section headers
    # parsing
  toml.decode
    @ (source: string) -> result[toml_value, string]
    + tokenizes then parses, returning the root table
    - propagates tokenize and parse errors
    # decoding
    -> toml.tokenize
    -> toml.parse
  toml.get_string
    @ (root: toml_value, path: string) -> optional[string]
    + returns the string value at a dotted path like "server.host"
    - returns none when any segment is missing or the leaf is not a string
    # lookup
  toml.get_int
    @ (root: toml_value, path: string) -> optional[i64]
    + returns the integer value at a dotted path
    - returns none when missing or not an integer
    # lookup
  toml.get_bool
    @ (root: toml_value, path: string) -> optional[bool]
    + returns the boolean value at a dotted path
    - returns none when missing or not a boolean
    # lookup
  toml.get_array
    @ (root: toml_value, path: string) -> optional[list[toml_value]]
    + returns the array at a dotted path
    - returns none when missing or not an array
    # lookup
  toml.encode_value
    @ (value: toml_value) -> string
    + renders a scalar value or inline array in TOML literal syntax
    ? used internally by encode; exposed for callers who build values directly
    # encoding
  toml.encode
    @ (root: toml_value) -> result[string, string]
    + returns a TOML document whose decode round-trips to the same value tree
    + emits nested tables as "[section.sub]" headers ordered by appearance
    - returns error when the root is not a table
    # encoding
    -> toml.encode_value
