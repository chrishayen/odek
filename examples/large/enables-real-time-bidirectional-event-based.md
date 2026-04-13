# Requirement: "a real-time bidirectional event-based communication library"

A framed message protocol with named events, acknowledgements, rooms, and transport-level heartbeats. The project manages sessions and routing; std provides framing, json, and time primitives.

std
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.random
    std.random.alnum
      @ (length: i32) -> string
      + returns a random alphanumeric string of the given length
      # randomness

realtime
  realtime.new_hub
    @ () -> hub_state
    + creates an empty hub with no connections or rooms
    # construction
  realtime.connect
    @ (hub: hub_state) -> tuple[hub_state, session_id]
    + registers a new session and returns its id
    # sessions
    -> std.random.alnum
  realtime.disconnect
    @ (hub: hub_state, sid: session_id) -> hub_state
    + removes the session from all rooms and drops pending acks
    # sessions
  realtime.encode_frame
    @ (event: string, payload: map[string, string], ack_id: optional[i64]) -> string
    + encodes an outgoing frame as a JSON envelope with event, data, and optional ack
    # protocol
    -> std.json.encode_object
  realtime.decode_frame
    @ (raw: string) -> result[frame, string]
    + parses an incoming frame into event, payload, and optional ack id
    - returns error when the envelope is missing the event field
    # protocol
    -> std.json.parse_object
  realtime.emit_to
    @ (hub: hub_state, sid: session_id, event: string, payload: map[string, string]) -> result[tuple[hub_state, list[outbound]], string]
    + returns the frames that should be written to the given session
    - returns error when the session does not exist
    # dispatch
  realtime.broadcast_room
    @ (hub: hub_state, room: string, event: string, payload: map[string, string]) -> tuple[hub_state, list[outbound]]
    + returns frames for every session currently in the room
    # dispatch
  realtime.join_room
    @ (hub: hub_state, sid: session_id, room: string) -> result[hub_state, string]
    + adds the session to the room
    - returns error when the session is not connected
    # rooms
  realtime.leave_room
    @ (hub: hub_state, sid: session_id, room: string) -> hub_state
    + removes the session from the room; no-op if not a member
    # rooms
  realtime.request_ack
    @ (hub: hub_state, sid: session_id, event: string, payload: map[string, string]) -> result[tuple[hub_state, outbound, i64], string]
    + emits a frame with a fresh ack id and records it as pending
    # acknowledgements
    -> std.time.now_millis
  realtime.receive_ack
    @ (hub: hub_state, sid: session_id, ack_id: i64) -> hub_state
    + clears the pending ack so it is no longer considered outstanding
    # acknowledgements
  realtime.tick_heartbeats
    @ (hub: hub_state, interval_ms: i64, now_ms: i64) -> tuple[hub_state, list[outbound], list[session_id]]
    + returns heartbeat frames to send and sessions that have timed out
    # liveness
