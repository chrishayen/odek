# Requirement: "a scalable http and websocket engine that runs across multiple worker processes"

A websocket and http server with a worker pool and pub/sub channels for cross-worker messaging.

std
  std.net
    std.net.listen_tcp
      @ (host: string, port: u16) -> result[listener_handle, string]
      + opens a tcp listener on the given host and port
      - returns error when the port is already in use
      # networking
    std.net.accept
      @ (listener: listener_handle) -> result[conn_handle, string]
      + accepts the next incoming connection
      - returns error when the listener is closed
      # networking
    std.net.read_bytes
      @ (conn: conn_handle, max: i32) -> result[bytes, string]
      + reads up to max bytes from a connection
      # networking
    std.net.write_bytes
      @ (conn: conn_handle, data: bytes) -> result[void, string]
      + writes bytes to a connection
      # networking
  std.crypto
    std.crypto.sha1
      @ (data: bytes) -> bytes
      + returns the 20-byte sha1 digest of data
      # cryptography
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + encodes bytes to base64 with padding
      # encoding
  std.ipc
    std.ipc.publish
      @ (channel: string, payload: bytes) -> result[void, string]
      + publishes a message to a cross-worker channel
      # inter_process
    std.ipc.subscribe
      @ (channel: string) -> result[subscription_handle, string]
      + returns a handle that receives messages published to the channel
      # inter_process

ws_engine
  ws_engine.new
    @ (host: string, port: u16, worker_count: i32) -> result[engine_state, string]
    + creates an engine bound to the address with the requested number of workers
    - returns error when the listener cannot be opened
    # construction
    -> std.net.listen_tcp
  ws_engine.accept_next
    @ (state: engine_state) -> result[tuple[conn_handle, engine_state], string]
    + dispatches a newly accepted connection to the least-loaded worker
    # load_balancing
    -> std.net.accept
  ws_engine.parse_http_request
    @ (raw: bytes) -> result[http_request, string]
    + parses headers and body out of a raw request buffer
    - returns error when the request line is malformed
    # http
  ws_engine.handle_http
    @ (state: engine_state, conn: conn_handle, req: http_request) -> result[void, string]
    + routes an http request to the matching handler and writes the response
    # http_routing
    -> std.net.write_bytes
  ws_engine.upgrade_websocket
    @ (conn: conn_handle, req: http_request) -> result[ws_conn, string]
    + performs the websocket handshake and returns a framed connection
    - returns error when the sec-websocket-key header is missing
    # handshake
    -> std.crypto.sha1
    -> std.encoding.base64_encode
    -> std.net.write_bytes
  ws_engine.read_frame
    @ (ws: ws_conn) -> result[ws_frame, string]
    + reads the next websocket frame, unmasking payload bytes
    - returns error on a malformed frame header
    # framing
    -> std.net.read_bytes
  ws_engine.write_frame
    @ (ws: ws_conn, frame: ws_frame) -> result[void, string]
    + writes a websocket frame with the correct opcode and length encoding
    # framing
    -> std.net.write_bytes
  ws_engine.join_channel
    @ (state: engine_state, ws: ws_conn, channel: string) -> engine_state
    + subscribes a websocket connection to a named channel
    # pub_sub
    -> std.ipc.subscribe
  ws_engine.broadcast
    @ (state: engine_state, channel: string, payload: bytes) -> result[void, string]
    + publishes a message to every connection subscribed to the channel, across all workers
    # pub_sub
    -> std.ipc.publish
  ws_engine.close_conn
    @ (state: engine_state, ws: ws_conn) -> engine_state
    + sends a close frame and removes the connection from all channels
    # lifecycle
