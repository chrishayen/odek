# Requirement: "a library for a lightweight smtp client providing an http api"

Accepts mail submission requests through a small HTTP handler and delivers messages by speaking SMTP over a TCP connection.

std
  std.net
    std.net.tcp_connect
      @ (host: string, port: i32) -> result[conn, string]
      + opens a TCP connection to host:port
      - returns error on dns failure or refused connection
      # networking
    std.net.tcp_write_line
      @ (c: conn, line: string) -> result[void, string]
      + writes line followed by CRLF
      - returns error on write failure
      # networking
    std.net.tcp_read_line
      @ (c: conn) -> result[string, string]
      + reads a single CRLF-terminated line
      - returns error on closed connection
      # networking

smtp_http
  smtp_http.send_mail
    @ (server: string, port: i32, from: string, to: list[string], subject: string, body: string) -> result[void, string]
    + runs HELO, MAIL FROM, RCPT TO, DATA, and QUIT with CRLF line endings
    - returns error when any SMTP reply code is outside 2xx or 3xx
    # smtp_dialog
    -> std.net.tcp_connect
    -> std.net.tcp_write_line
    -> std.net.tcp_read_line
  smtp_http.read_reply
    @ (c: conn) -> result[tuple[i32, string], string]
    + parses the numeric code and text from one or more continuation lines
    - returns error when no digit code can be read
    # smtp_dialog
    -> std.net.tcp_read_line
  smtp_http.build_data_payload
    @ (from: string, to: list[string], subject: string, body: string) -> string
    + returns an RFC 5322 message ending with a lone dot line
    # message_building
  smtp_http.parse_submit_request
    @ (method: string, path: string, body: string) -> result[mail_request, string]
    + decodes a JSON body into from, to, subject, and body fields
    - returns error when method is not POST or path is not "/submit"
    - returns error on invalid JSON or missing fields
    # http_handler
  smtp_http.handle
    @ (method: string, path: string, body: string, server: string, port: i32) -> result[string, string]
    + parses the request, sends the message, and returns a "queued" response body
    - returns error when parsing or delivery fails
    # http_handler
