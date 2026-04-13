# Requirement: "a library that pretty-prints HTTP requests and responses for debugging"

Takes structured request/response records and produces a human-readable multi-line dump.

std: (all units exist)

httpdump
  httpdump.format_request
    @ (method: string, url: string, headers: map[string, string], body: bytes) -> string
    + returns a request-line, sorted headers, a blank line, and a decoded body
    ? headers are sorted by name for stable output
    # formatting
  httpdump.format_response
    @ (status: i32, reason: string, headers: map[string, string], body: bytes) -> string
    + returns a status line, sorted headers, a blank line, and a decoded body
    # formatting
  httpdump.redact_headers
    @ (headers: map[string, string], sensitive: list[string]) -> map[string, string]
    + returns a copy with the named headers' values replaced by "<redacted>"
    ? header name matching is case-insensitive
    # redaction
