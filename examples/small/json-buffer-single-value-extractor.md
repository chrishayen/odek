# Requirement: "a library for extracting a single value from a JSON buffer without parsing the whole document"

Scan the raw bytes for a top-level key and return its value substring. No full parse.

std: (all units exist)

jsonpeek
  jsonpeek.find_key
    fn (buffer: bytes, key: string) -> result[i32, string]
    + returns the byte offset just after the closing quote of the key in the top-level object
    - returns error when the key is not present at the top level
    - returns error when the buffer is not a JSON object
    # scanning
  jsonpeek.extract_value
    fn (buffer: bytes, key: string) -> result[bytes, string]
    + returns the raw bytes of the value associated with key at the top level
    - returns error when the key is absent
    - returns error when the value is not a complete JSON token (unterminated string, bracket, or brace)
    # extraction
    -> jsonpeek.find_key
  jsonpeek.extract_string
    fn (buffer: bytes, key: string) -> result[string, string]
    + returns the decoded string value for key
    - returns error when the value is not a JSON string
    # extraction
    -> jsonpeek.extract_value
