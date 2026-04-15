# Requirement: "a URL shortener service"

Two project entry points (shorten and resolve) over an opaque state value. Std holds the URL validation and random-string generation — both generic primitives.

std
  std.random
    std.random.alphanumeric_string
      fn (length: u32) -> string
      + returns a cryptographically random alphanumeric string of the given length
      # randomness
  std.url
    std.url.validate
      fn (raw: string) -> result[void, string]
      + returns ok when the string is a syntactically valid http or https URL
      - returns error on missing scheme
      - returns error on malformed authority component
      # validation

url_shortener
  url_shortener.shorten
    fn (state: shortener_state, long_url: string) -> result[tuple[string, shortener_state], string]
    + generates a short code and stores the mapping
    + returns (short_code, new_state)
    - returns error when long_url is not a valid URL
    ? short codes are 7 characters; collisions are retried up to 5 times
    # creation
    -> std.url.validate
    -> std.random.alphanumeric_string
  url_shortener.resolve
    fn (state: shortener_state, short_code: string) -> optional[string]
    + returns the long URL for a known short code
    + returns none when the short code does not exist
    # lookup
