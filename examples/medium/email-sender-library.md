# Requirement: "an email sending library"

A client that composes messages and delivers them via SMTP. The std layer provides the transport and encoding primitives.

std
  std.net
    std.net.tcp_connect
      fn (host: string, port: i32) -> result[tcp_conn, string]
      + opens a TCP connection to host:port
      - returns error when the host cannot be resolved
      # networking
    std.net.tcp_write
      fn (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      - returns error when the connection is closed
      # networking
    std.net.tcp_read_line
      fn (conn: tcp_conn) -> result[string, string]
      + reads up to the next CRLF
      - returns error on read failure
      # networking
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes to standard base64
      # encoding
  std.crypto
    std.crypto.tls_upgrade
      fn (conn: tcp_conn, host: string) -> result[tcp_conn, string]
      + upgrades a plain connection to TLS
      - returns error on handshake failure
      # cryptography

email
  email.message
    fn (from: string, to: list[string], subject: string, body: string) -> email_msg
    + builds a message with the given headers and body
    # composition
  email.connect
    fn (host: string, port: i32, use_tls: bool) -> result[smtp_session, string]
    + dials the SMTP server and performs the EHLO handshake
    - returns error when the greeting is not 220
    # transport
    -> std.net.tcp_connect
    -> std.net.tcp_read_line
    -> std.crypto.tls_upgrade
  email.authenticate
    fn (session: smtp_session, user: string, password: string) -> result[smtp_session, string]
    + authenticates with AUTH PLAIN
    - returns error when credentials are rejected
    # authentication
    -> std.encoding.base64_encode
    -> std.net.tcp_write
  email.send
    fn (session: smtp_session, msg: email_msg) -> result[void, string]
    + issues MAIL FROM, RCPT TO, DATA and the message body
    - returns error when any recipient is rejected
    # delivery
    -> std.net.tcp_write
    -> std.net.tcp_read_line
  email.close
    fn (session: smtp_session) -> void
    + sends QUIT and releases the connection
    # transport
    -> std.net.tcp_write
