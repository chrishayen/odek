# Requirement: "a minimal MQTT broker that runs over any byte-stream transport"

Parses MQTT control packets, maintains client sessions and subscriptions, and routes PUBLISH messages. The transport is abstracted as a byte-stream reader/writer supplied by the caller.

std
  std.io
    std.io.read_exact
      fn (stream: io_stream, n: i32) -> result[bytes, string]
      + reads exactly n bytes or returns an error
      - returns error on premature EOF
      # io
    std.io.write_all
      fn (stream: io_stream, data: bytes) -> result[void, string]
      + writes all bytes to the stream
      # io
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

mqtt
  mqtt.decode_varint
    fn (data: bytes, offset: i32) -> result[tuple[i32, i32], string]
    + returns (value, bytes_consumed) for an MQTT remaining-length varint
    - returns error when the varint exceeds four bytes
    # framing
  mqtt.encode_varint
    fn (value: i32) -> bytes
    + returns the MQTT varint encoding of value
    # framing
  mqtt.read_packet
    fn (stream: io_stream) -> result[packet, string]
    + returns the next MQTT control packet from the stream
    - returns error on malformed fixed header
    # framing
    -> std.io.read_exact
  mqtt.encode_packet
    fn (pkt: packet) -> bytes
    + returns the wire bytes for a control packet
    # framing
  mqtt.write_packet
    fn (stream: io_stream, pkt: packet) -> result[void, string]
    + encodes and writes a packet to the stream
    # framing
    -> std.io.write_all
  mqtt.new_broker
    fn () -> broker_state
    + returns a broker with no clients or subscriptions
    # construction
  mqtt.handle_connect
    fn (broker: broker_state, pkt: packet, stream: io_stream) -> result[tuple[broker_state, client_id], string]
    + registers a new session and returns the accepted client id
    - returns error with a CONNACK reason code when the protocol version is unsupported
    # session
    -> std.time.now_millis
  mqtt.handle_subscribe
    fn (broker: broker_state, client: client_id, pkt: packet) -> broker_state
    + records the client's subscriptions and their granted QoS levels
    # subscriptions
  mqtt.handle_unsubscribe
    fn (broker: broker_state, client: client_id, pkt: packet) -> broker_state
    + removes the client's matching subscriptions
    # subscriptions
  mqtt.topic_matches
    fn (filter: string, topic: string) -> bool
    + returns true when a publish topic matches a subscription filter with + and # wildcards
    # routing
  mqtt.handle_publish
    fn (broker: broker_state, pkt: packet) -> list[delivery]
    + returns the deliveries that must be written to subscribed clients
    + honors retained flag for late-joining subscribers
    # routing
  mqtt.handle_disconnect
    fn (broker: broker_state, client: client_id) -> broker_state
    + removes the session, preserving subscriptions only when clean_session is false
    # session
  mqtt.tick_keepalive
    fn (broker: broker_state) -> tuple[broker_state, list[client_id]]
    + returns clients whose keepalive deadline has elapsed so the caller may close them
    # session
    -> std.time.now_millis
