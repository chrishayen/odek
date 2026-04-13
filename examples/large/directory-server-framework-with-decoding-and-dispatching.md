# Requirement: "a directory-server framework that decodes incoming directory-protocol requests and dispatches them to user-provided handlers"

Parses BER-encoded request messages, classifies them by operation, and invokes registered handlers that return response messages.

std
  std.encoding
    std.encoding.ber_decode
      @ (data: bytes) -> result[ber_node, string]
      + decodes a BER-encoded value into a tag-length-value tree
      - returns error on truncated or malformed input
      # encoding
    std.encoding.ber_encode
      @ (node: ber_node) -> bytes
      + serializes a BER tree back to bytes
      # encoding
  std.net
    std.net.accept
      @ (listener: listener_handle) -> result[conn_handle, string]
      + blocks until a client connects and returns the connection
      - returns error when the listener is closed
      # networking
    std.net.read_bytes
      @ (conn: conn_handle, max: i32) -> result[bytes, string]
      + reads up to max bytes from the connection
      - returns error on read failure
      # networking
    std.net.write_bytes
      @ (conn: conn_handle, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      - returns error on write failure
      # networking
    std.net.listen_tcp
      @ (host: string, port: i32) -> result[listener_handle, string]
      + binds a TCP listener on host and port
      - returns error when the address cannot be bound
      # networking

directory_server
  directory_server.parse_request
    @ (data: bytes) -> result[directory_request, string]
    + decodes the message into an operation with a message id and parameters
    - returns error when the envelope or inner structure is malformed
    # parsing
    -> std.encoding.ber_decode
  directory_server.operation_of
    @ (req: directory_request) -> string
    + returns one of "bind", "unbind", "search", "add", "modify", "delete", "compare", "extended"
    # dispatch
  directory_server.new_router
    @ () -> router_state
    + creates a router with no handlers registered
    # construction
  directory_server.register_handler
    @ (router: router_state, op: string, handler_id: string) -> router_state
    + registers a handler id for a given operation name
    # registration
  directory_server.dispatch
    @ (router: router_state, req: directory_request) -> optional[string]
    + returns the handler id registered for the request's operation
    - returns none when no handler is registered
    # dispatch
  directory_server.build_response
    @ (message_id: i32, code: i32, payload: bytes) -> directory_response
    + packages the fields into a response structure
    # responses
  directory_server.encode_response
    @ (resp: directory_response) -> bytes
    + serializes the response to bytes for transmission
    # responses
    -> std.encoding.ber_encode
  directory_server.handle_connection
    @ (router: router_state, conn: conn_handle, invoker: handler_invoker) -> result[void, string]
    + reads requests from the connection in a loop, routes them, and writes responses
    + stops cleanly on connection close
    - returns error on any unrecoverable protocol error
    # serving
    -> std.net.read_bytes
    -> std.net.write_bytes
  directory_server.serve
    @ (router: router_state, listener: listener_handle, invoker: handler_invoker) -> result[void, string]
    + accepts connections and handles each one
    - returns error when the listener fails
    # serving
    -> std.net.accept
  directory_server.new_listener
    @ (host: string, port: i32) -> result[listener_handle, string]
    + binds a listener for the server
    - returns error when the address cannot be bound
    # setup
    -> std.net.listen_tcp
