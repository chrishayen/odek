# Requirement: "an SMTP capture server for developer email testing"

Accepts SMTP traffic, parses MIME messages, and exposes them for inspection via an in-process query API. No outbound delivery.

std
  std.net
    std.net.tcp_listen
      fn (host: string, port: u16) -> result[listener_state, string]
      + binds and starts listening on the address
      - returns error when the port is already in use
      # networking
    std.net.tcp_accept
      fn (listener: listener_state) -> result[conn_state, string]
      + blocks until a new connection arrives
      # networking
    std.net.read_line
      fn (conn: conn_state) -> result[string, string]
      + reads up to and including a CRLF from the connection
      # networking
    std.net.write
      fn (conn: conn_state, data: bytes) -> result[void, string]
      + writes data to the connection
      # networking
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits on each separator occurrence
      # strings
    std.strings.to_lower
      fn (s: string) -> string
      + lowercases ASCII letters
      # strings
  std.encoding
    std.encoding.base64_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes standard base64
      # encoding
    std.encoding.quoted_printable_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes quoted-printable content
      # encoding
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

smtp_capture
  smtp_capture.new_store
    fn () -> capture_store
    + creates an empty store for captured messages
    # construction
  smtp_capture.start
    fn (store: capture_store, host: string, port: u16) -> result[server_state, string]
    + begins accepting SMTP connections on the given address
    # lifecycle
    -> std.net.tcp_listen
  smtp_capture.handle_connection
    fn (store: capture_store, conn: conn_state) -> result[void, string]
    + walks the SMTP command dialogue: HELO, MAIL FROM, RCPT TO, DATA, QUIT
    + enqueues the parsed message into the store on successful DATA termination
    - returns error when a client issues commands out of order
    # protocol
    -> std.net.read_line
    -> std.net.write
    -> std.time.now_millis
  smtp_capture.parse_command
    fn (line: string) -> result[smtp_command, string]
    + recognizes HELO, EHLO, MAIL, RCPT, DATA, RSET, QUIT with their arguments
    - returns error on unknown verbs
    # parsing
    -> std.strings.split
    -> std.strings.to_lower
  smtp_capture.parse_message
    fn (raw: bytes) -> result[captured_message, string]
    + extracts headers, subject, from, to, body parts, and attachments
    + decodes quoted-printable and base64 transfer encodings
    - returns error on malformed header lines
    # mime
    -> std.encoding.base64_decode
    -> std.encoding.quoted_printable_decode
  smtp_capture.list_messages
    fn (store: capture_store) -> list[captured_message]
    + returns all captured messages in receipt order
    # query
  smtp_capture.get_message
    fn (store: capture_store, id: string) -> optional[captured_message]
    + retrieves a single captured message by id
    # query
  smtp_capture.search
    fn (store: capture_store, query: string) -> list[captured_message]
    + returns messages whose subject, from, or to contains the query substring
    # query
  smtp_capture.delete_message
    fn (store: capture_store, id: string) -> bool
    + removes a message from the store and returns whether it existed
    # mutation
  smtp_capture.clear
    fn (store: capture_store) -> void
    + removes every captured message
    # mutation
