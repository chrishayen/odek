# Requirement: "a minimal MQTT broker that runs over any byte-stream transport"

Parses MQTT control packets, maintains client sessions and subscriptions, and routes PUBLISH messages. The transport is abstracted as a byte-stream reader/writer supplied by the caller.

std
  std.io
    std.io.read_exact
      @ (stream: io_stream, n: i32) -> result[bytes, string]
      + reads exactly n bytes or returns an error
      - returns error on premature EOF
      # io
    std.io.write_all
      @ (stream: io_stream, data: bytes) -> result[void, string]
      + writes all bytes to the stream
      # io
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

mqtt
  mqtt.decode_varint
    @ (data: bytes, offset: i32) -> result[tuple[i32, i32], string]
    + returns (value, bytes_consumed) for an MQTT remaining-length varint
    - returns error when the varint exceeds four bytes
    # framing
  mqtt.encode_varint
    @ (value: i32) -> bytes
    + returns the MQTT varint encoding of value
    # framing
  mqtt.read_packet
    @ (stream: io_stream) -> result[packet, string]
    + returns the next MQTT control packet from the stream
    - returns error on malformed fixed header
    # framing
    -> std.io.read_exact
  mqtt.encode_packet
    @ (pkt: packet) -> bytes
    + returns the wire bytes for a control packet
    # framing
  mqtt.write_packet
    @ (stream: io_stream, pkt: packet) -> result[void, string]
    + encodes and writes a packet to the stream
    # framing
    -> std.io.write_all
  mqtt.new_broker
    @ () -> broker_state
    + returns a broker with no clients or subscriptions
    # construction
  mqtt.handle_connect
    @ (broker: broker_state, pkt: packet, stream: io_stream) -> result[tuple[broker_state, client_id], string]
    + registers a new session and returns the accepted client id
    - returns error with a CONNACK reason code when the protocol version is unsupported
    # session
    -> std.time.now_millis
  mqtt.handle_subscribe
    @ (broker: broker_state, client: client_id, pkt: packet) -> broker_state
    + records the client's subscriptions and their granted QoS levels
    # subscriptions
  mqtt.handle_unsubscribe
    @ (broker: broker_state, client: client_id, pkt: packet) -> broker_state
    + removes the client's matching subscriptions
    # subscriptions
  mqtt.topic_matches
    @ (filter: string, topic: string) -> bool
    + returns true when a publish topic matches a subscription filter with + and # wildcards
    # routing
  mqtt.handle_publish
    @ (broker: broker_state, pkt: packet) -> list[delivery]
    + returns the deliveries that must be written to subscribed clients
    + honors retained flag for late-joining subscribers
    # routing
  mqtt.handle_disconnect
    @ (broker: broker_state, client: client_id) -> broker_state
    + removes the session, preserving subscriptions only when clean_session is false
    # session
  mqtt.tick_keepalive
    @ (broker: broker_state) -> tuple[broker_state, list[client_id]]
    + returns clients whose keepalive deadline has elapsed so the caller may close them
    # session
    -> std.time.now_millis
