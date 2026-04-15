# Requirement: "a library to read a JSON value via a dotted path expression"

Lets callers pull a single leaf out of a JSON document without parsing it into a typed tree. Supports object keys, array indices, and simple wildcards.

std: (all units exist)

jsonpath
  jsonpath.tokenize_path
    fn (path: string) -> result[list[string], string]
    + splits a dotted path into segments, handling quoted keys and bracketed indices
    - returns error on unbalanced brackets or quotes
    ? "users.0.name" and "users[0].name" both yield ["users", "0", "name"]
    # path_parsing
  jsonpath.scan_value
    fn (raw: string, offset: i32) -> result[tuple[string, i32], string]
    + reads a JSON value starting at offset and returns its raw slice and the offset after it
    - returns error on malformed JSON at that offset
    ? handles nested objects and arrays by tracking bracket depth
    # scanning
  jsonpath.field
    fn (raw: string, key: string) -> result[string, string]
    + returns the raw JSON for the immediate field `key` of the object at the start of raw
    - returns error when raw is not an object or the key is absent
    # object_access
    -> jsonpath.scan_value
  jsonpath.index
    fn (raw: string, idx: i32) -> result[string, string]
    + returns the raw JSON for the element at `idx` of the array at the start of raw
    - returns error when raw is not an array or idx is out of bounds
    # array_access
    -> jsonpath.scan_value
  jsonpath.get
    fn (raw: string, path: string) -> result[string, string]
    + walks the dotted path and returns the raw JSON at that location
    - returns error when any segment cannot be resolved
    # entry
    -> jsonpath.tokenize_path
    -> jsonpath.field
    -> jsonpath.index
  jsonpath.unquote_string
    fn (raw: string) -> result[string, string]
    + interprets a raw JSON string literal, decoding escape sequences
    - returns error when raw is not a quoted string
    # string_decoding
  jsonpath.get_string
    fn (raw: string, path: string) -> result[string, string]
    + like get but also unquotes the result as a JSON string
    - returns error when the value at path is not a string
    # convenience
    -> jsonpath.get
    -> jsonpath.unquote_string
