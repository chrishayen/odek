# Requirement: "a cloud services client library"

A thin client for a generic cloud provider API: credential loading, request signing, HTTP transport, and typed wrappers for a handful of representative services.

std
  std.env
    std.env.get
      fn (name: string) -> optional[string]
      + returns the value for the given environment variable
      - returns none when the variable is not set
      # environment
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file's contents as a string
      - returns error when the file does not exist
      # filesystem
  std.time
    std.time.now_iso8601
      fn () -> string
      + returns the current UTC time formatted as ISO 8601
      # time
  std.crypto
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256
      # cryptography
    std.crypto.sha256_hex
      fn (data: bytes) -> string
      + returns the lowercase hex SHA-256 digest
      # cryptography
  std.http
    std.http.send
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an HTTP request and returns status, headers, and body
      - returns error on transport failure
      # networking
  std.encoding
    std.encoding.json_parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document
      - returns error on malformed input
      # serialization
    std.encoding.json_encode
      fn (value: json_value) -> string
      + encodes a JSON value to a string
      # serialization
    std.encoding.ini_parse
      fn (raw: string) -> result[map[string, map[string, string]], string]
      + parses an INI document into sections of key-value pairs
      - returns error on malformed input
      # serialization

cloud
  cloud.credentials_from_env
    fn () -> result[credentials, string]
    + returns credentials from standard environment variables
    - returns error when any required variable is missing
    # credentials
    -> std.env.get
  cloud.credentials_from_file
    fn (path: string, profile: string) -> result[credentials, string]
    + returns credentials for the named profile in a credentials file
    - returns error when the profile is absent
    # credentials
    -> std.fs.read_all
    -> std.encoding.ini_parse
  cloud.new_client
    fn (creds: credentials, region: string) -> client_state
    + creates a client bound to the given credentials and region
    # construction
  cloud.sign_request
    fn (state: client_state, method: string, url: string, headers: map[string, string], body: bytes) -> map[string, string]
    + returns headers including authorization and timestamp
    ? follows a canonical-request + signing-key construction
    # request_signing
    -> std.time.now_iso8601
    -> std.crypto.hmac_sha256
    -> std.crypto.sha256_hex
  cloud.send_request
    fn (state: client_state, method: string, service: string, path: string, body: bytes) -> result[http_response, string]
    + signs and dispatches an HTTP request to the named service endpoint
    - returns error on transport failure
    - returns error when response status indicates a client or server error
    # transport
    -> std.http.send
  cloud.object_store_put
    fn (state: client_state, bucket: string, key: string, data: bytes) -> result[void, string]
    + uploads an object to the given bucket and key
    - returns error when the bucket does not exist
    # object_store
  cloud.object_store_get
    fn (state: client_state, bucket: string, key: string) -> result[bytes, string]
    + downloads an object by bucket and key
    - returns error when the object is missing
    # object_store
  cloud.queue_send
    fn (state: client_state, queue: string, message: string) -> result[string, string]
    + enqueues a message and returns its assigned id
    - returns error when the queue does not exist
    # queue
  cloud.queue_receive
    fn (state: client_state, queue: string, max_messages: i32) -> result[list[queue_message], string]
    + returns up to max_messages pending messages
    - returns error when the queue does not exist
    # queue
  cloud.table_put_item
    fn (state: client_state, table: string, item: map[string, json_value]) -> result[void, string]
    + writes an item to a key-value table
    - returns error when the primary key is missing from the item
    # table
    -> std.encoding.json_encode
  cloud.table_get_item
    fn (state: client_state, table: string, key: map[string, json_value]) -> result[optional[map[string, json_value]], string]
    + returns the item for the given key, or none if absent
    - returns error when the table does not exist
    # table
    -> std.encoding.json_parse
