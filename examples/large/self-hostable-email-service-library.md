# Requirement: "a self-hostable email service library"

Parsing, address validation, MIME composition, and delivery via an SMTP primitive.

std
  std.net
    std.net.tcp_connect
      @ (host: string, port: i32) -> result[conn_handle, string]
      + opens a tcp connection
      - returns error when host is unreachable
      # network
    std.net.tcp_write
      @ (conn: conn_handle, data: bytes) -> result[void, string]
      + writes bytes to the connection
      # network
    std.net.tcp_read_line
      @ (conn: conn_handle) -> result[string, string]
      + reads until CRLF, returning the line without the terminator
      # network
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns 32 bytes
      # cryptography
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + returns standard base64 with padding
      # encoding
    std.encoding.quoted_printable_encode
      @ (data: bytes) -> string
      + returns RFC 2045 quoted-printable
      # encoding

email_service
  email_service.parse_address
    @ (raw: string) -> result[email_address, string]
    + returns local and domain parts of a plain address
    - returns error when "@" is missing or either side is empty
    # parsing
  email_service.validate_domain
    @ (domain: string) -> bool
    + returns true when domain has at least one dot and only valid label characters
    # validation
  email_service.new_message
    @ (from: email_address, to: list[email_address], subject: string) -> message_state
    + returns an empty message with the given envelope
    # composition
  email_service.set_text_body
    @ (m: message_state, body: string) -> message_state
    + sets the plain-text body
    -> std.encoding.quoted_printable_encode
    # composition
  email_service.attach_file
    @ (m: message_state, filename: string, content: bytes, mime_type: string) -> message_state
    + adds an attachment encoded as base64
    -> std.encoding.base64_encode
    # composition
  email_service.render_mime
    @ (m: message_state) -> string
    + returns the full RFC 5322 message with headers and boundaries
    # composition
  email_service.message_id
    @ (m: message_state) -> string
    + returns a deterministic message id derived from the content hash
    -> std.crypto.sha256
    # composition
  email_service.smtp_connect
    @ (host: string, port: i32) -> result[smtp_session, string]
    + opens an smtp session and reads the 220 greeting
    - returns error when the greeting is not 220
    # delivery
    -> std.net.tcp_connect
    -> std.net.tcp_read_line
  email_service.smtp_send
    @ (s: smtp_session, m: message_state) -> result[void, string]
    + executes MAIL FROM, RCPT TO, DATA, and closes the session
    - returns error when any step returns a non-2xx code
    # delivery
    -> std.net.tcp_write
    -> std.net.tcp_read_line
