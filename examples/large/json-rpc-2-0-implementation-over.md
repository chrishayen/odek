# Requirement: "a JSON-RPC 2.0 implementation over WebSockets"

A full bidirectional RPC layer: client and server sides share framing, with a websocket transport and pending-call tracking.

std
  std.websocket
    std.websocket.dial
      @ (url: string) -> result[ws_conn, string]
      + opens a websocket connection to the given URL
      - returns error when the handshake fails
      # websocket
    std.websocket.send_text
      @ (conn: ws_conn, message: string) -> result[void, string]
      + sends a text frame
      - returns error when the connection is closed
      # websocket
    std.websocket.recv_text
      @ (conn: ws_conn) -> result[string, string]
      + receives the next text frame
      - returns error when the connection is closed
      # websocket
    std.websocket.close
      @ (conn: ws_conn) -> void
      + closes the connection
      # websocket
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object
      - returns error on malformed input
      # serialization
  std.sync
    std.sync.new_channel
      @ () -> channel_state
      + creates a single-value channel for awaiting an RPC reply
      # synchronization
    std.sync.send
      @ (ch: channel_state, value: string) -> void
      + delivers the value to the channel receiver
      # synchronization
    std.sync.recv
      @ (ch: channel_state) -> string
      + blocks until a value is delivered
      # synchronization

jsonrpc_ws
  jsonrpc_ws.new_client
    @ (url: string) -> result[client_state, string]
    + dials the websocket and returns a ready client
    - returns error on handshake failure
    # client_construction
    -> std.websocket.dial
  jsonrpc_ws.next_id
    @ (state: client_state) -> tuple[string, client_state]
    + returns a fresh request id
    # id_generation
  jsonrpc_ws.encode_request
    @ (id: string, method: string, params: map[string, string]) -> string
    + returns a JSON-RPC 2.0 request envelope
    # framing
    -> std.json.encode_object
  jsonrpc_ws.encode_response
    @ (id: string, result_value: map[string, string]) -> string
    + returns a JSON-RPC 2.0 success response envelope
    # framing
    -> std.json.encode_object
  jsonrpc_ws.encode_error
    @ (id: string, code: i32, message: string) -> string
    + returns a JSON-RPC 2.0 error response envelope
    # framing
    -> std.json.encode_object
  jsonrpc_ws.parse_message
    @ (raw: string) -> result[map[string, string], string]
    + parses an incoming request or response envelope
    - returns error when the envelope is missing jsonrpc version
    # framing
    -> std.json.parse_object
  jsonrpc_ws.call
    @ (state: client_state, method: string, params: map[string, string]) -> result[map[string, string], string]
    + sends a request, waits for the matching response, and returns the result
    - returns error when the server replies with an error
    # rpc_call
    -> std.websocket.send_text
    -> std.sync.new_channel
    -> std.sync.recv
  jsonrpc_ws.dispatch_loop
    @ (state: client_state) -> void
    + reads frames from the socket and routes responses to waiting callers
    # dispatcher
    -> std.websocket.recv_text
    -> std.sync.send
  jsonrpc_ws.register_handler
    @ (state: server_state, method: string, handler_id: string) -> server_state
    + registers a method handler on the server side
    # server_registration
  jsonrpc_ws.serve
    @ (state: server_state, conn: ws_conn) -> void
    + reads frames, invokes registered handlers, and writes responses back
    - encodes an error envelope when a method is unknown
    # server_loop
    -> std.websocket.recv_text
    -> std.websocket.send_text
  jsonrpc_ws.close
    @ (state: client_state) -> void
    + closes the underlying connection
    # teardown
    -> std.websocket.close
