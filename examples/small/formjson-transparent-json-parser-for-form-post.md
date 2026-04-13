# Requirement: "a library that transparently handles json input as a standard form post"

Normalize incoming HTTP request bodies so a handler can read fields the same way regardless of whether the client sent JSON or url-encoded form data.

std
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.url
    std.url.parse_form
      @ (raw: string) -> map[string, string]
      + parses application/x-www-form-urlencoded bodies into a map
      + decodes percent-escapes and plus-as-space
      # url

formjson
  formjson.normalize_body
    @ (content_type: string, body: string) -> result[map[string, string], string]
    + returns a flat field map regardless of content type
    + treats content types starting with "application/json" as JSON
    + treats content types starting with "application/x-www-form-urlencoded" as form
    - returns error for unsupported content types
    # normalization
    -> std.json.parse_object
    -> std.url.parse_form
  formjson.get_field
    @ (fields: map[string, string], name: string) -> optional[string]
    + returns the value for the given field name
    + returns none when the field is missing
    # access
