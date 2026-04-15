# Requirement: "decode an HTTP request's query string, form body, and headers into a named field map using a declared schema"

The library takes a pre-parsed request shape and a schema that binds fields to sources, then returns one resolved map.

std
  std.encoding
    std.encoding.url_decode
      fn (raw: string) -> result[string, string]
      + percent-decodes a URL-encoded token
      - returns error on invalid percent sequences
      # encoding
  std.http
    std.http.parse_query
      fn (raw: string) -> result[map[string, list[string]], string]
      + parses a query string into keys and decoded values
      - returns error on malformed pairs
      # http
    std.http.parse_form
      fn (body: string) -> result[map[string, list[string]], string]
      + parses an application/x-www-form-urlencoded body into keys and decoded values
      - returns error on malformed pairs
      # http

http_decode
  http_decode.new_schema
    fn () -> schema_state
    + creates an empty schema
    # construction
  http_decode.from_query
    fn (schema: schema_state, field: string, key: string, required: bool) -> schema_state
    + binds a field to a query-string key
    # schema
  http_decode.from_form
    fn (schema: schema_state, field: string, key: string, required: bool) -> schema_state
    + binds a field to a form-body key
    # schema
  http_decode.from_header
    fn (schema: schema_state, field: string, header: string, required: bool) -> schema_state
    + binds a field to a request header
    # schema
  http_decode.decode
    fn (schema: schema_state, query: string, body: string, headers: map[string, string]) -> result[map[string, string], string]
    + returns one resolved map with each schema field set from its source
    - returns error when a required field is missing from its source
    - returns error when a source cannot be parsed
    # decoding
    -> std.http.parse_query
    -> std.http.parse_form
    -> std.encoding.url_decode
