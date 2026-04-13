# Requirement: "an HTTP cookie parser and serializer"

Parse a Cookie header into name-value pairs and serialize a cookie with attributes for Set-Cookie.

std
  std.text
    std.text.trim_spaces
      @ (s: string) -> string
      + removes leading and trailing ASCII whitespace
      # text

cookie
  cookie.parse_header
    @ (header: string) -> map[string, string]
    + parses a Cookie header value into name-to-value pairs
    + ignores empty segments and strips surrounding whitespace
    # parsing
    -> std.text.trim_spaces
  cookie.new
    @ (name: string, value: string) -> cookie_data
    + creates a cookie with the given name and value and empty attributes
    # construction
  cookie.with_attributes
    @ (c: cookie_data, path: optional[string], domain: optional[string], max_age: optional[i64], secure: bool, http_only: bool) -> cookie_data
    + returns a new cookie with the given attributes set
    # attributes
  cookie.serialize
    @ (c: cookie_data) -> string
    + produces a Set-Cookie header value with attributes joined by "; "
    - returns error-marked string when the name is empty or contains invalid characters
    # serialization
