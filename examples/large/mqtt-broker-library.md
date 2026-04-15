# Requirement: "an MQTT broker library"

A library that parses MQTT control packets, manages client sessions, and routes publishes to subscribers by topic filter.

std
  std.net
    std.net.listen_tcp
      fn (address: string, port: i32) -> result[tcp_listener, string]
      + binds and listens on the given address
      - returns error when the port is in use
      # networking
    std.net.accept
      fn (listener: tcp_listener) -> result[tcp_connection, string]
      + blocks for the next incoming connection
      # networking
    std.net.read_bytes
      fn (conn: tcp_connection, max: i32) -> result[bytes, string]
      + reads up to max bytes
      # networking
    std.net.write_bytes
      fn (conn: tcp_connection, data: bytes) -> result[void, string]
      + writes all bytes
      # networking
  std.encoding
    std.encoding.read_varint
      fn (data: bytes, offset: i32) -> result[tuple[i32, i32], string]
      + returns (value, bytes_consumed) for MQTT remaining-length encoding
      - returns error when more than four continuation bytes are seen
      # encoding
    std.encoding.write_varint
      fn (value: i32) -> bytes
      + encodes a non-negative integer in MQTT remaining-length form
      # encoding

mqtt
  mqtt.parse_packet
    fn (raw: bytes) -> result[mqtt_packet, string]
    + parses the fixed header, remaining length, and variable header
    - returns error on truncated input
    - returns error on unknown packet type
    # protocol_parsing
    -> std.encoding.read_varint
  mqtt.encode_packet
    fn (packet: mqtt_packet) -> bytes
    + serializes a packet to wire bytes
    # protocol_encoding
    -> std.encoding.write_varint
  mqtt.topic_matches
    fn (filter: string, topic: string) -> bool
    + honors '+' single-level and '#' multi-level wildcards
    - returns false when '#' is not the last segment of the filter
    # topic_matching
  mqtt.broker_new
    fn () -> broker_state
    + creates an empty broker with no clients or subscriptions
    # construction
  mqtt.broker_on_connect
    fn (state: broker_state, client_id: string, clean_session: bool) -> broker_state
    + registers a client session, optionally resuming prior state
    # session_management
  mqtt.broker_on_disconnect
    fn (state: broker_state, client_id: string) -> broker_state
    + marks the client inactive and drops volatile subscriptions
    # session_management
  mqtt.broker_on_subscribe
    fn (state: broker_state, client_id: string, filter: string, qos: i32) -> broker_state
    + adds a topic filter subscription for the client
    # subscriptions
  mqtt.broker_on_unsubscribe
    fn (state: broker_state, client_id: string, filter: string) -> broker_state
    + removes the matching subscription
    # subscriptions
  mqtt.broker_on_publish
    fn (state: broker_state, topic: string, payload: bytes, qos: i32) -> tuple[broker_state, list[mqtt_delivery]]
    + returns the updated state and the list of deliveries to dispatch
    - no deliveries when no subscription matches
    # routing
    -> mqtt.topic_matches
  mqtt.broker_retained_for
    fn (state: broker_state, filter: string) -> list[mqtt_delivery]
    + returns retained messages matching the filter for new subscribers
    # retained_messages
    -> mqtt.topic_matches
