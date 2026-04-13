# Requirement: "a websocket framework for building instant-messaging servers"

The framework maintains connected clients, rooms, and message routing on top of low-level websocket framing primitives.

std
  std.ws
    std.ws.handshake
      @ (raw: bytes) -> result[bytes, string]
      + validates an HTTP upgrade request and returns the handshake response bytes
      - returns error when required headers are missing
      # websocket
    std.ws.encode_frame
      @ (opcode: u8, payload: bytes) -> bytes
      + returns a websocket frame for the given opcode and payload
      # websocket
    std.ws.decode_frame
      @ (raw: bytes) -> result[ws_frame, string]
      + parses a frame, handling fragmentation and masking
      - returns error on truncated input
      # websocket
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string map
      - returns error on non-object input
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string map as JSON
      # serialization

imhub
  imhub.new_hub
    @ () -> hub_state
    + returns an empty hub with no clients or rooms
    # construction
  imhub.accept_client
    @ (hub: hub_state, client_id: string, request: bytes) -> result[tuple[hub_state, bytes], string]
    + validates the upgrade request and registers the client
    - returns error when the handshake is invalid
    # lifecycle
    -> std.ws.handshake
  imhub.disconnect
    @ (hub: hub_state, client_id: string) -> hub_state
    + removes a client from all rooms and drops its registration
    # lifecycle
  imhub.join_room
    @ (hub: hub_state, client_id: string, room: string) -> result[hub_state, string]
    + adds the client to the named room
    - returns error when the client is unknown
    # membership
  imhub.leave_room
    @ (hub: hub_state, client_id: string, room: string) -> hub_state
    + removes the client from the named room
    # membership
  imhub.room_members
    @ (hub: hub_state, room: string) -> list[string]
    + returns the ids of clients in the room
    # membership
  imhub.send_direct
    @ (hub: hub_state, from_id: string, to_id: string, body: string) -> result[bytes, string]
    + returns a framed message payload to deliver to to_id
    - returns error when to_id is not connected
    # messaging
    -> std.json.encode_object
    -> std.ws.encode_frame
  imhub.broadcast_room
    @ (hub: hub_state, from_id: string, room: string, body: string) -> list[tuple[string, bytes]]
    + returns a framed payload per room member (excluding from_id)
    # messaging
    -> std.json.encode_object
    -> std.ws.encode_frame
  imhub.handle_frame
    @ (hub: hub_state, client_id: string, raw: bytes) -> result[hub_command, string]
    + decodes an inbound frame and returns the parsed command
    - returns error when the frame does not decode
    - returns error when the payload is not a JSON command object
    # ingress
    -> std.ws.decode_frame
    -> std.json.parse_object
