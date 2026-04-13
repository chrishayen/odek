# Requirement: "a JSON parser that produces helpful error messages"

The parser tracks line/column offsets so failures point at the exact character and return a nearby excerpt.

std
  std.strings
    std.strings.slice
      @ (s: string, start: i32, end: i32) -> string
      + returns the substring between byte offsets
      # strings
    std.strings.byte_at
      @ (s: string, i: i32) -> u8
      + returns the byte at index i
      # strings

parsejson
  parsejson.parse
    @ (source: string) -> result[json_value, parse_error]
    + returns a parsed json_value on success
    - returns a parse_error with line, column, and a pointer excerpt on failure
    # parsing
    -> std.strings.byte_at
  parsejson.error_message
    @ (err: parse_error) -> string
    + formats a multi-line human-readable message with caret under the offending character
    + includes the JSON path where the error occurred
    # diagnostics
    -> std.strings.slice
  parsejson.error_location
    @ (err: parse_error) -> tuple[i32, i32]
    + returns (line, column), both 1-based
    # diagnostics
  parsejson.error_excerpt
    @ (err: parse_error, context_chars: i32) -> string
    + returns the source line containing the error with up to context_chars on each side
    # diagnostics
    -> std.strings.slice
  parsejson.error_path
    @ (err: parse_error) -> string
    + returns a dotted path like "items[3].name" to where the error occurred
    - returns "" when the error is at the document root
    # diagnostics
