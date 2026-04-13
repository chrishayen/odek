# Requirement: "an in-memory email capture server for development"

Accepts SMTP connections, parses incoming mail, and stores it in memory for later inspection.

std
  std.net
    std.net.tcp_listen
      @ (host: string, port: u16) -> result[listener_state, string]
      + binds and listens on the address
      # networking
    std.net.tcp_accept
      @ (listener: listener_state) -> result[conn_state, string]
      + blocks for the next connection
      # networking
    std.net.read_line
      @ (conn: conn_state) -> result[string, string]
      + reads a CRLF-terminated line
      # networking
    std.net.write
      @ (conn: conn_state, data: bytes) -> result[void, string]
      + writes raw bytes to the connection
      # networking
  std.time
    std.time.now_millis
      @ () -> i64
      + returns the current unix time in milliseconds
      # time

mail_dev
  mail_dev.new_store
    @ () -> mail_store
    + creates an empty in-memory mail store
    # construction
  mail_dev.start
    @ (store: mail_store, host: string, port: u16) -> result[server_state, string]
    + begins accepting SMTP traffic
    # lifecycle
    -> std.net.tcp_listen
  mail_dev.handle_connection
    @ (store: mail_store, conn: conn_state) -> result[void, string]
    + walks the SMTP state machine and appends completed messages to the store
    - returns error on commands sent out of order
    # protocol
    -> std.net.read_line
    -> std.net.write
    -> std.time.now_millis
  mail_dev.parse_message
    @ (raw: bytes) -> result[captured_message, string]
    + extracts from, to, subject, body, and attachments
    - returns error on malformed headers
    # mime
  mail_dev.list
    @ (store: mail_store) -> list[captured_message]
    + returns captured messages in receipt order
    # query
  mail_dev.get
    @ (store: mail_store, id: string) -> optional[captured_message]
    + retrieves a single captured message
    # query
  mail_dev.clear
    @ (store: mail_store) -> void
    + empties the store
    # mutation
