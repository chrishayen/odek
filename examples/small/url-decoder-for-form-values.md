# Requirement: "decode URL-encoded form values into typed field values"

std
  std.encoding
    std.encoding.url_decode
      fn (raw: string) -> result[string, string]
      + percent-decodes a single URL-encoded token
      - returns error on invalid percent sequences
      # encoding

form_decode
  form_decode.parse_form
    fn (body: string) -> result[map[string, list[string]], string]
    + splits an application/x-www-form-urlencoded body into keys and decoded values
    + groups repeated keys into a list preserving order
    - returns error on malformed key=value pairs
    # parsing
    -> std.encoding.url_decode
  form_decode.get_string
    fn (values: map[string, list[string]], key: string) -> result[string, string]
    + returns the first value for a key
    - returns error when the key is missing
    # access
  form_decode.get_int
    fn (values: map[string, list[string]], key: string) -> result[i64, string]
    + returns the first value for a key parsed as a signed integer
    - returns error when the key is missing or not an integer
    # access
  form_decode.get_bool
    fn (values: map[string, list[string]], key: string) -> result[bool, string]
    + returns the first value for a key parsed as a boolean ("1"/"0", "true"/"false")
    - returns error when the key is missing or unrecognized
    # access
