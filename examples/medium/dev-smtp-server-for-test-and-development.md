# Requirement: "an SMTP server for test and development"

Accepts mail, stores it in memory, and exposes an inspection API. Drives a protocol state machine.

std
  std.net
    std.net.listen_tcp
      fn (host: string, port: i32) -> result[listener, string]
      + binds a TCP listener on the given address
      - returns error when the address is already in use
      # networking
    std.net.accept
      fn (l: listener) -> result[tcp_conn, string]
      + blocks until a new connection arrives
      # networking
    std.net.read_line
      fn (c: tcp_conn) -> result[string, string]
      + reads one CRLF-terminated line from the connection
      # networking
    std.net.write_all
      fn (c: tcp_conn, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      # networking
  std.sync
    std.sync.new_mutex
      fn () -> mutex
      + creates an unlocked mutex
      # concurrency

dev_smtp
  dev_smtp.new
    fn (host: string, port: i32) -> result[dev_smtp_state, string]
    + opens the listener and initializes the in-memory mailbox
    - returns error when binding fails
    # construction
    -> std.net.listen_tcp
    -> std.sync.new_mutex
  dev_smtp.serve_once
    fn (state: dev_smtp_state) -> result[void, string]
    + accepts one connection, runs the SMTP dialog to completion, and stores any delivered message
    - returns error when the client aborts mid-session
    # session_loop
    -> std.net.accept
    -> std.net.read_line
    -> std.net.write_all
  dev_smtp.list_messages
    fn (state: dev_smtp_state) -> list[stored_message]
    + returns all messages captured so far, oldest first
    # inspection
  dev_smtp.get_message
    fn (state: dev_smtp_state, id: string) -> optional[stored_message]
    + returns a single message by its assigned id
    # inspection
  dev_smtp.clear
    fn (state: dev_smtp_state) -> void
    + removes all stored messages
    # inspection
