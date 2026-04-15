# Requirement: "an OAuth 1.0 request-signing library"

Produces the Authorization header for an OAuth 1.0 signed HTTP request using HMAC-SHA1. No network IO; the caller sends the request.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.random
    std.random.bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes as standard base64 with padding
      # encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding
    std.encoding.percent_encode
      fn (raw: string) -> string
      + percent-encodes per RFC 3986 unreserved character set
      + encodes space as %20, not +
      # url_encoding
  std.crypto
    std.crypto.hmac_sha1
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA1 of data under key
      + returns 20 bytes
      # cryptography
  std.url
    std.url.parse
      fn (raw: string) -> result[map[string, string], string]
      + returns a map with keys "scheme", "host", "path", "query"
      - returns error on malformed URL
      # url_parsing

oauth
  oauth.nonce
    fn () -> string
    + returns a fresh random nonce suitable as a single-use token
    ? nonce is hex-encoded random bytes
    # nonce
    -> std.random.bytes
    -> std.encoding.hex_encode
  oauth.timestamp
    fn () -> string
    + returns the current unix time in seconds as a decimal string
    # timestamp
    -> std.time.now_seconds
  oauth.base_string
    fn (method: string, url: string, params: map[string, string]) -> result[string, string]
    + builds the RFC 5849 signature base string from method, normalized URL, and sorted params
    - returns error when url is not parseable
    ? params are sorted lexicographically by encoded key then encoded value
    # signature_base
    -> std.url.parse
    -> std.encoding.percent_encode
  oauth.signing_key
    fn (consumer_secret: string, token_secret: string) -> string
    + returns "encoded_consumer_secret&encoded_token_secret" per spec
    ? token_secret may be empty; the trailing ampersand is always present
    # signing_key
    -> std.encoding.percent_encode
  oauth.sign_hmac_sha1
    fn (base_string: string, signing_key: string) -> string
    + returns the base64-encoded HMAC-SHA1 signature
    # signing
    -> std.crypto.hmac_sha1
    -> std.encoding.base64_encode
  oauth.authorization_header
    fn (params: map[string, string]) -> string
    + returns the `OAuth k1="v1", k2="v2"` header value with values percent-encoded and quoted
    + includes only oauth_ prefixed params
    # header_formatting
    -> std.encoding.percent_encode
  oauth.sign_request
    fn (method: string, url: string, body_params: map[string, string], consumer_key: string, consumer_secret: string, token: string, token_secret: string) -> result[string, string]
    + returns the full Authorization header value for the request
    - returns error on malformed url
    ? adds oauth_consumer_key, oauth_nonce, oauth_signature_method ("HMAC-SHA1"), oauth_timestamp, oauth_token, oauth_version ("1.0")
    # end_to_end
