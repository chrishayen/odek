# Requirement: "a library to process HTTP file uploads"

Parses multipart/form-data request bodies into fields and files.

std
  std.bytes
    std.bytes.index_of
      @ (haystack: bytes, needle: bytes, from: i32) -> i32
      + returns the first index of needle at or after from, -1 when absent
      # bytes
    std.bytes.slice
      @ (data: bytes, start: i32, end: i32) -> bytes
      + returns a subrange without copying semantics guaranteed
      # bytes
  std.http
    std.http.parse_header_params
      @ (header_value: string) -> map[string, string]
      + parses a header like 'form-data; name="x"; filename="y"' into a param map
      # http

multipart
  multipart.boundary_from_content_type
    @ (content_type: string) -> result[string, string]
    + extracts the boundary parameter from a Content-Type header
    - returns error when the header is not multipart/form-data
    - returns error when boundary is missing
    # header_parsing
    -> std.http.parse_header_params
  multipart.split_parts
    @ (body: bytes, boundary: string) -> result[list[bytes], string]
    + returns the raw bytes of each part between boundary delimiters
    - returns error on truncated body or missing closing boundary
    # framing
    -> std.bytes.index_of
    -> std.bytes.slice
  multipart.parse_part
    @ (raw_part: bytes) -> result[part, string]
    + splits header block from body at the blank line
    + parses each header line into (name, value)
    - returns error when headers are malformed
    # part_parsing
    -> std.bytes.index_of
    -> std.http.parse_header_params
  multipart.parse_body
    @ (body: bytes, content_type: string) -> result[parsed_upload, string]
    + returns a struct separating text fields from uploaded files
    # aggregation
    -> std.bytes.index_of
