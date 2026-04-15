# Requirement: "a client and server implementation of the GraphQL-over-WebSocket subprotocol"

The subprotocol multiplexes subscription operations over a single WebSocket connection using a small set of typed messages.

std
  std.websocket
    std.websocket.accept
      fn (request_headers: map[string,string]) -> result[ws_socket, string]
      + completes the server handshake and returns a bound socket
      - returns error on missing or malformed Sec-WebSocket-Key
      # websocket
    std.websocket.dial
      fn (url: string, headers: map[string,string]) -> result[ws_socket, string]
      + completes the client handshake and returns a bound socket
      # websocket
    std.websocket.send_text
      fn (socket: ws_socket, frame: string) -> result[void, string]
      + writes one text frame to the socket
      # websocket
    std.websocket.recv_text
      fn (socket: ws_socket) -> result[string, string]
      + reads one text frame, returning error on close
      # websocket
    std.websocket.close
      fn (socket: ws_socket, code: i32, reason: string) -> result[void, string]
      + sends a close frame and tears down the socket
      # websocket
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string,string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

gqlws
  gqlws.encode_message
    fn (type: string, id: string, payload: string) -> string
    + produces a protocol frame as a JSON text envelope
    # wire
    -> std.json.encode_object
  gqlws.decode_message
    fn (frame: string) -> result[gqlws_msg, string]
    + parses a frame into its type, id, and payload
    - returns error on unknown message type
    # wire
    -> std.json.parse_object
  gqlws.server_new
    fn () -> gqlws_server
    + creates a new server state with no active connections
    # server
  gqlws.server_handle_connection
    fn (server: gqlws_server, headers: map[string,string]) -> result[gqlws_conn, string]
    + accepts a handshake and waits for connection_init before marking the conn active
    - returns error when the client does not send connection_init
    # server
    -> std.websocket.accept
    -> std.websocket.send_text
    -> std.websocket.recv_text
  gqlws.server_dispatch
    fn (conn: gqlws_conn, msg: gqlws_msg) -> result[list[gqlws_msg], string]
    + handles one incoming message and returns any outbound messages to send
    + supports subscribe, next, complete, and ping
    # server
  gqlws.server_stream_next
    fn (conn: gqlws_conn, op_id: string, data: string) -> result[void, string]
    + pushes a next message for an active subscription
    - returns error when op_id is unknown
    # server
    -> std.websocket.send_text
  gqlws.client_connect
    fn (url: string) -> result[gqlws_client, string]
    + opens a socket and completes the connection_init / connection_ack exchange
    - returns error when the server does not acknowledge
    # client
    -> std.websocket.dial
    -> std.websocket.send_text
    -> std.websocket.recv_text
  gqlws.client_subscribe
    fn (client: gqlws_client, query: string) -> result[string, string]
    + sends a subscribe message and returns the operation id assigned
    # client
    -> std.websocket.send_text
  gqlws.client_next
    fn (client: gqlws_client) -> result[gqlws_msg, string]
    + reads the next inbound message
    - returns error when the connection is closed
    # client
    -> std.websocket.recv_text
  gqlws.client_complete
    fn (client: gqlws_client, op_id: string) -> result[void, string]
    + sends a complete message for the given operation
    # client
    -> std.websocket.send_text
  gqlws.close
    fn (conn: gqlws_conn, code: i32, reason: string) -> result[void, string]
    + cancels all active operations and closes the underlying socket
    # lifecycle
    -> std.websocket.close
