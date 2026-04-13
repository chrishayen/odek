# Requirement: "a low-overhead publish-subscribe network protocol library"

Encodes and decodes protocol frames and tracks subscription state for a lightweight pub/sub transport.

std
  std.net
    std.net.send
      @ (socket: socket_handle, data: bytes) -> result[i32, string]
      + sends bytes on a connected socket and returns the number written
      - returns error on transport failure
      # networking
    std.net.recv
      @ (socket: socket_handle, max: i32) -> result[bytes, string]
      + reads up to max bytes from a socket
      - returns error on transport failure
      # networking

zproto
  zproto.encode_publish
    @ (topic: string, payload: bytes) -> bytes
    + encodes a publish frame with topic and payload
    # framing
  zproto.encode_subscribe
    @ (subscription_id: i64, pattern: string) -> bytes
    + encodes a subscription request for a topic pattern
    # framing
  zproto.encode_unsubscribe
    @ (subscription_id: i64) -> bytes
    + encodes an unsubscribe request
    # framing
  zproto.decode_frame
    @ (buffer: bytes) -> result[tuple[frame, bytes], string]
    + decodes one frame and returns it with the remaining buffer
    - returns error on truncated or malformed input
    # framing
  zproto.new_router
    @ () -> router_state
    + creates an empty subscription router
    # routing
  zproto.add_subscription
    @ (router: router_state, subscription_id: i64, pattern: string) -> router_state
    + registers a topic pattern under the given subscription id
    # routing
  zproto.remove_subscription
    @ (router: router_state, subscription_id: i64) -> router_state
    + removes the subscription with the given id
    # routing
  zproto.match_subscribers
    @ (router: router_state, topic: string) -> list[i64]
    + returns the subscription ids whose patterns match the topic
    # routing
  zproto.send_frame
    @ (socket: socket_handle, frame_bytes: bytes) -> result[void, string]
    + writes a frame to a socket in full
    - returns error on transport failure
    # transport
    -> std.net.send
  zproto.read_frame
    @ (socket: socket_handle) -> result[frame, string]
    + reads one frame from a socket, buffering as needed
    - returns error on transport or decode failure
    # transport
    -> std.net.recv
    -> zproto.decode_frame
