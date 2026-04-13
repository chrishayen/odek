# Requirement: "an email message parser that reads an RFC 5322 message and emits a structured representation"

Parses headers, body, and simple multipart structure from a raw email byte stream.

std
  std.io
    std.io.read_all_stdin
      @ () -> result[bytes, string]
      + reads standard input until EOF
      - returns error on read failure
      # input
  std.encoding
    std.encoding.decode_quoted_printable
      @ (data: string) -> result[string, string]
      + decodes a quoted-printable encoded body
      - returns error on an invalid escape sequence
      # encoding
    std.encoding.decode_base64
      @ (data: string) -> result[bytes, string]
      + decodes a base64 encoded body
      - returns error on characters outside the alphabet
      # encoding

email_parser
  email_parser.parse
    @ (raw: bytes) -> result[email_message, string]
    + parses a complete RFC 5322 message into headers and body parts
    - returns error when headers are missing the blank-line terminator
    # parsing
    -> std.encoding.decode_quoted_printable
    -> std.encoding.decode_base64
  email_parser.split_headers
    @ (raw: bytes) -> result[tuple[list[header_line], bytes], string]
    + splits the raw message into its header lines and body bytes
    + unfolds continuation lines per RFC 5322
    - returns error when no blank line separator is found
    # parsing
  email_parser.parse_header_line
    @ (line: string) -> result[tuple[string, string], string]
    + returns a name and value pair from a header line
    - returns error when the line has no colon separator
    # parsing
  email_parser.parse_address_list
    @ (value: string) -> list[email_address]
    + parses a header value into zero or more name+address pairs
    # parsing
  email_parser.split_multipart
    @ (body: bytes, boundary: string) -> list[bytes]
    + returns the raw parts of a multipart body split by the boundary
    # parsing
