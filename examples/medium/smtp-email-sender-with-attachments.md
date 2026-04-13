# Requirement: "send text/HTML emails with attachments over SMTP"

Builds a MIME message, opens an SMTP session, and hands the message to the server. The project layer is thin; the encoding work lives in std.

std
  std.net
    std.net.dial_tls
      @ (host: string, port: i32) -> result[conn_state, string]
      + opens a TLS-wrapped TCP connection
      - returns error on handshake or network failure
      # networking
    std.net.write_line
      @ (conn: conn_state, line: string) -> result[conn_state, string]
      + writes a CRLF-terminated line to a connection
      # networking
    std.net.read_line
      @ (conn: conn_state) -> result[tuple[string, conn_state], string]
      + reads one CRLF-terminated line from a connection
      # networking
    std.net.close
      @ (conn: conn_state) -> void
      + closes a connection and releases its resources
      # networking
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + encodes bytes as standard base64 with padding
      # encoding
    std.encoding.quoted_printable
      @ (text: string) -> string
      + encodes text per RFC 2045 quoted-printable
      # encoding
  std.mime
    std.mime.guess_type
      @ (filename: string) -> string
      + returns a MIME type guessed from a file name extension, or "application/octet-stream"
      # mime

email
  email.new_message
    @ (from: string, subject: string) -> message_state
    + creates an empty message with the given sender and subject
    # construction
  email.add_recipient
    @ (msg: message_state, kind: string, address: string) -> message_state
    + adds a To, Cc, or Bcc recipient
    # recipients
  email.set_body
    @ (msg: message_state, text: string, html: optional[string]) -> message_state
    + sets the plain-text body and an optional HTML alternative
    # body
  email.attach
    @ (msg: message_state, filename: string, data: bytes) -> message_state
    + adds an attachment; the content type is derived from the filename
    # attachments
    -> std.mime.guess_type
  email.render
    @ (msg: message_state) -> string
    + serializes the message as a multipart MIME document suitable for SMTP DATA
    # encoding
    -> std.encoding.base64_encode
    -> std.encoding.quoted_printable
  email.send
    @ (msg: message_state, host: string, port: i32, user: string, password: string) -> result[void, string]
    + connects over TLS, authenticates, and transmits the rendered message
    - returns error when the server rejects any SMTP command
    # transport
    -> std.net.dial_tls
    -> std.net.write_line
    -> std.net.read_line
    -> std.net.close
