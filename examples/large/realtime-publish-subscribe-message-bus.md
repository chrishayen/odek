# Requirement: "a real-time client-server publish/subscribe message bus"

Clients handshake with a server, subscribe to channels, and receive published messages in real time. The library provides both sides as pure state transitions over a transport.

std
  std.ws
    std.ws.connect
      @ (url: string) -> result[ws_conn, string]
      + opens a websocket connection
      - returns error when the handshake fails
      # websocket
    std.ws.receive
      @ (conn: ws_conn) -> result[string, string]
      + returns the next text message
      - returns error when the connection is closed
      # websocket
    std.ws.send
      @ (conn: ws_conn, message: string) -> result[void, string]
      + sends a text message
      - returns error when the connection is closed
      # websocket
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses JSON text into a value tree
      - returns error on invalid JSON
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + serializes a JSON value to a string
      # serialization
  std.rand
    std.rand.token_hex
      @ (length: i32) -> string
      + returns a hex-encoded random token of the given byte length
      # randomness

msgbus
  msgbus.server_new
    @ () -> server_state
    + creates an empty server with no clients or channels
    # construction
  msgbus.handshake
    @ (server: server_state) -> tuple[server_state, string]
    + returns a new client id and a server updated with the client registered
    # handshake
    -> std.rand.token_hex
  msgbus.subscribe
    @ (server: server_state, client_id: string, channel: string) -> result[server_state, string]
    + returns a server with client_id added as a subscriber of channel
    - returns error when the client id is unknown
    # subscription
  msgbus.unsubscribe
    @ (server: server_state, client_id: string, channel: string) -> server_state
    + returns a server with the client removed from the channel's subscribers
    # subscription
  msgbus.publish
    @ (server: server_state, channel: string, payload: json_value) -> tuple[server_state, list[delivery]]
    + returns a list of deliveries (one per subscribed client) for the transport to send
    # publication
  msgbus.encode_message
    @ (msg: bus_message) -> string
    + returns the wire form of a server or client message
    # serialization
    -> std.json.encode
  msgbus.decode_message
    @ (raw: string) -> result[bus_message, string]
    + parses a wire message
    - returns error when the channel field is missing
    - returns error when the type field is not recognized
    # serialization
    -> std.json.parse
  msgbus.client_connect
    @ (url: string) -> result[client_state, string]
    + opens a connection and performs the handshake exchange
    - returns error when the handshake response is not received
    # client
    -> std.ws.connect
    -> std.ws.send
    -> std.ws.receive
  msgbus.client_subscribe
    @ (client: client_state, channel: string) -> result[client_state, string]
    + sends a subscribe message and records the channel in the client state
    - returns error when the connection is closed
    # client
    -> std.ws.send
  msgbus.client_publish
    @ (client: client_state, channel: string, payload: json_value) -> result[void, string]
    + sends a publish message
    - returns error when the connection is closed
    # client
    -> std.ws.send
  msgbus.client_next
    @ (client: client_state) -> result[bus_message, string]
    + reads and decodes the next message from the server
    - returns error when the connection is closed
    # client
    -> std.ws.receive
