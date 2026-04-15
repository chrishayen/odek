# Requirement: "a client library for S3-compatible object storage"

The library builds signed requests for object operations. Actual transport is the caller's concern: every operation returns a prepared request struct and consumes a response struct.

std
  std.crypto
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256 of data under key
      + returns 32 bytes
      # cryptography
    std.crypto.sha256_hex
      fn (data: bytes) -> string
      + returns the lowercase hex SHA-256 digest
      # cryptography
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + returns lowercase hex
      # encoding
    std.encoding.url_encode
      fn (s: string) -> string
      + percent-encodes reserved characters
      # encoding
  std.time
    std.time.now_unix
      fn () -> i64
      + returns current unix time in seconds
      # time
    std.time.format_iso8601_basic
      fn (unix_sec: i64) -> string
      + formats as YYYYMMDDTHHMMSSZ
      # time
  std.xml
    std.xml.parse_element
      fn (raw: string) -> result[xml_element, string]
      + parses a single XML document root
      - returns error on malformed XML
      # serialization
    std.xml.child_text
      fn (elem: xml_element, tag: string) -> optional[string]
      + returns the text of the first child with the given tag
      # serialization

object_store
  object_store.new_client
    fn (endpoint: string, access_key: string, secret_key: string, region: string) -> client_state
    + returns a client bound to the given credentials
    # construction
  object_store.build_put_request
    fn (client: client_state, bucket: string, key: string, body: bytes) -> signed_request
    + returns a signed PUT request for uploading body to bucket/key
    # requests
    -> std.crypto.sha256_hex
    -> std.crypto.hmac_sha256
    -> std.encoding.hex_encode
    -> std.time.now_unix
    -> std.time.format_iso8601_basic
  object_store.build_get_request
    fn (client: client_state, bucket: string, key: string) -> signed_request
    + returns a signed GET request for bucket/key
    # requests
    -> std.crypto.hmac_sha256
  object_store.build_delete_request
    fn (client: client_state, bucket: string, key: string) -> signed_request
    + returns a signed DELETE request
    # requests
    -> std.crypto.hmac_sha256
  object_store.build_list_request
    fn (client: client_state, bucket: string, prefix: string) -> signed_request
    + returns a signed LIST request
    # requests
    -> std.encoding.url_encode
    -> std.crypto.hmac_sha256
  object_store.parse_list_response
    fn (raw: string) -> result[list[string], string]
    + returns the list of object keys in a ListObjects response body
    - returns error on malformed XML
    # parsing
    -> std.xml.parse_element
    -> std.xml.child_text
  object_store.parse_error_response
    fn (status: i32, raw: string) -> string
    + returns a human-readable error string from a non-2xx response body
    # parsing
    -> std.xml.parse_element
    -> std.xml.child_text
