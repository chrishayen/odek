# Requirement: "an expressive general-purpose client for an RPC framework"

A client that loads service schemas, discovers methods, constructs request messages dynamically, and sends them over a framed transport. Message encoding and transport framing are primitives in std.

std
  std.encoding
    std.encoding.varint_encode
      fn (value: u64) -> bytes
      + encodes an unsigned integer using variable-length encoding
      # encoding
    std.encoding.varint_decode
      fn (data: bytes, offset: i64) -> result[tuple[u64, i64], string]
      + decodes a varint at offset and returns the value plus new offset
      - returns error on truncated input
      # encoding
  std.net
    std.net.dial_tcp
      fn (host: string, port: i32) -> result[socket, string]
      + opens a TCP connection
      - returns error on refused or unresolvable host
      # network
    std.net.send_bytes
      fn (s: socket, data: bytes) -> result[void, string]
      + writes all bytes to the socket
      - returns error on broken connection
      # network
    std.net.recv_bytes
      fn (s: socket, n: i64) -> result[bytes, string]
      + reads exactly n bytes or returns error
      - returns error on early close
      # network
  std.text
    std.text.split_lines
      fn (source: string) -> list[string]
      + splits text on line separators
      # text

rpc_client
  rpc_client.load_schema
    fn (source: string) -> result[schema, string]
    + parses a schema file describing services, methods, and message fields
    - returns error on malformed schema
    -> std.text.split_lines
    # schema
  rpc_client.list_services
    fn (s: schema) -> list[string]
    + returns the names of services declared in the schema
    # reflection
  rpc_client.list_methods
    fn (s: schema, service: string) -> result[list[method_info], string]
    + returns methods for the named service
    - returns error when the service is unknown
    # reflection
  rpc_client.describe_method
    fn (s: schema, service: string, method: string) -> result[method_info, string]
    + returns the full description including input and output field types
    - returns error when the method is unknown
    # reflection
  rpc_client.build_request
    fn (info: method_info, fields: map[string, string]) -> result[bytes, string]
    + encodes a request message for the method from field name/value pairs
    - returns error on missing required field or type mismatch
    -> std.encoding.varint_encode
    # request
  rpc_client.decode_response
    fn (info: method_info, payload: bytes) -> result[map[string, string], string]
    + decodes a response message into a field map
    - returns error when bytes do not match the method output schema
    -> std.encoding.varint_decode
    # response
  rpc_client.connect
    fn (host: string, port: i32) -> result[client_state, string]
    + opens a connection and returns client state
    -> std.net.dial_tcp
    # connection
  rpc_client.call
    fn (state: client_state, service: string, method: string, fields: map[string, string], s: schema) -> result[map[string, string], string]
    + sends a unary request and returns the decoded response
    - returns error on transport failure
    - returns error when the service returns a non-success status
    -> std.net.send_bytes
    -> std.net.recv_bytes
    # invocation
  rpc_client.close
    fn (state: client_state) -> void
    + releases the underlying connection
    # teardown
  rpc_client.format_message
    fn (fields: map[string, string]) -> string
    + renders a decoded message as readable text
    # rendering
