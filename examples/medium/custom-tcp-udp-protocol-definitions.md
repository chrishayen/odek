# Requirement: "a library for defining custom tcp and udp wire protocols"

A schema-driven codec: describe a message in terms of fields, then encode/decode wire bytes. Socket IO is thin std.

std
  std.net
    std.net.tcp_dial
      @ (host: string, port: u16) -> result[tcp_conn, string]
      + opens a TCP connection to host:port
      - returns error when the host is unreachable
      # transport
    std.net.tcp_send
      @ (conn: tcp_conn, data: bytes) -> result[u32, string]
      + writes data and returns the number of bytes written
      # transport
    std.net.tcp_recv
      @ (conn: tcp_conn, max: u32) -> result[bytes, string]
      + reads up to max bytes from the connection
      - returns error when the peer has closed
      # transport
    std.net.udp_bind
      @ (port: u16) -> result[udp_socket, string]
      + binds a UDP socket on the given port
      # transport
    std.net.udp_send_to
      @ (sock: udp_socket, host: string, port: u16, data: bytes) -> result[u32, string]
      + sends a datagram to the given address
      # transport
    std.net.udp_recv_from
      @ (sock: udp_socket) -> result[tuple[string, u16, bytes], string]
      + blocks for a datagram, returns (host, port, payload)
      # transport

protocol
  protocol.field_u8
    @ (name: string) -> field_spec
    + describes a single-byte unsigned field
    # schema
  protocol.field_u32_be
    @ (name: string) -> field_spec
    + describes a big-endian 4-byte unsigned field
    # schema
  protocol.field_length_prefixed_bytes
    @ (name: string, length_field: string) -> field_spec
    + describes a byte field whose length is read from an earlier field in the same message
    # schema
  protocol.define_message
    @ (name: string, fields: list[field_spec]) -> message_schema
    + creates a message schema from an ordered list of fields
    - returns a schema even if empty; encoding a mismatched value fails at encode time
    # schema
  protocol.encode
    @ (schema: message_schema, values: map[string, bytes]) -> result[bytes, string]
    + serializes a message according to the schema
    - returns error when a required field is missing
    - returns error when a length-prefixed field's length does not fit its length field
    # encoding
  protocol.decode
    @ (schema: message_schema, data: bytes) -> result[map[string, bytes], string]
    + parses bytes into a field map according to the schema
    - returns error when data is shorter than the schema requires
    # decoding
  protocol.send_tcp
    @ (conn: tcp_conn, schema: message_schema, values: map[string, bytes]) -> result[void, string]
    + encodes and writes a message over a TCP connection
    # transport
    -> std.net.tcp_send
  protocol.recv_tcp
    @ (conn: tcp_conn, schema: message_schema) -> result[map[string, bytes], string]
    + reads and decodes a message from a TCP connection
    - returns error when the connection closes mid-message
    # transport
    -> std.net.tcp_recv
