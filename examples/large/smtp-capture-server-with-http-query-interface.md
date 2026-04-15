# Requirement: "an SMTP capture server with an HTTP query interface"

Captures SMTP traffic into a store and exposes read-only HTTP handlers that serve the captured messages.

std
  std.net
    std.net.tcp_listen
      fn (host: string, port: u16) -> result[listener_state, string]
      + binds and listens on the address
      # networking
    std.net.tcp_accept
      fn (listener: listener_state) -> result[conn_state, string]
      + blocks until a new connection arrives
      # networking
    std.net.read_line
      fn (conn: conn_state) -> result[string, string]
      + reads until CRLF
      # networking
    std.net.write
      fn (conn: conn_state, data: bytes) -> result[void, string]
      + writes data to the connection
      # networking
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body from a raw HTTP/1.1 request
      - returns error on malformed request line
      # http
    std.http.format_response
      fn (status: u16, headers: map[string, string], body: bytes) -> bytes
      + serializes a response with content-length
      # http
  std.json
    std.json.encode_value
      fn (value: json_value) -> string
      + encodes a generic JSON value to text
      # serialization
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

mail_capture
  mail_capture.new_store
    fn () -> capture_store
    + creates an empty mail store
    # construction
  mail_capture.start_smtp
    fn (store: capture_store, host: string, port: u16) -> result[server_state, string]
    + begins accepting SMTP connections
    # lifecycle
    -> std.net.tcp_listen
  mail_capture.handle_smtp
    fn (store: capture_store, conn: conn_state) -> result[void, string]
    + drives the SMTP dialogue and stores each completed message
    # protocol
    -> std.net.read_line
    -> std.net.write
    -> std.time.now_millis
  mail_capture.parse_message
    fn (raw: bytes) -> result[captured_message, string]
    + extracts headers, body text, html part, and attachments
    - returns error on unterminated header block
    # mime
  mail_capture.list_messages
    fn (store: capture_store, offset: i32, limit: i32) -> list[captured_message]
    + returns a page of captured messages newest first
    # query
  mail_capture.get_message
    fn (store: capture_store, id: string) -> optional[captured_message]
    + retrieves a single message by id
    # query
  mail_capture.delete_message
    fn (store: capture_store, id: string) -> bool
    + removes the named message, returning whether it existed
    # mutation
  mail_capture.clear
    fn (store: capture_store) -> void
    + removes every captured message
    # mutation
  mail_capture.start_http
    fn (store: capture_store, host: string, port: u16) -> result[server_state, string]
    + begins serving the read and delete HTTP endpoints
    # lifecycle
    -> std.net.tcp_listen
  mail_capture.handle_http
    fn (store: capture_store, conn: conn_state) -> result[void, string]
    + routes GET /messages, GET /messages/{id}, DELETE /messages/{id}, DELETE /messages
    + writes JSON bodies with appropriate status codes
    - returns 404 for unknown message ids
    # http
    -> std.http.parse_request
    -> std.http.format_response
    -> std.json.encode_value
