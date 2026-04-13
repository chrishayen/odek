# Requirement: "a message broker library implementing a publish-subscribe protocol with topic filters, QoS levels, and retained messages"

Encodes and decodes protocol packets, maintains a session and subscription registry, routes publishes to matching subscribers, and remembers retained messages per topic.

std
  std.encoding
    std.encoding.varint_encode
      @ (n: i32) -> bytes
      + encodes a non-negative integer using 7-bit variable-length encoding
      # encoding
    std.encoding.varint_decode
      @ (data: bytes) -> result[tuple[i32, i32], string]
      + returns the decoded integer and bytes consumed
      - returns error when the data is truncated
      # encoding
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits a string by separator
      # strings
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

pubsub_broker
  pubsub_broker.encode_packet
    @ (packet: packet) -> bytes
    + serializes a packet into wire bytes with a fixed header and payload
    # codec
    -> std.encoding.varint_encode
  pubsub_broker.decode_packet
    @ (data: bytes) -> result[packet, string]
    + decodes a single packet from wire bytes
    - returns error on truncated input or unknown packet type
    # codec
    -> std.encoding.varint_decode
  pubsub_broker.new_broker
    @ () -> broker_state
    + creates an empty broker with no sessions or subscriptions
    # construction
  pubsub_broker.connect_client
    @ (state: broker_state, client_id: string, clean: bool) -> tuple[broker_state, bool]
    + registers a session for the client id and returns whether a prior session was present
    + clears any stored session when clean is true
    # sessions
    -> std.time.now_millis
  pubsub_broker.disconnect_client
    @ (state: broker_state, client_id: string) -> broker_state
    + removes the session unless it was created as persistent
    # sessions
  pubsub_broker.subscribe
    @ (state: broker_state, client_id: string, filter: string, qos: i32) -> result[broker_state, string]
    + records the subscription at the given QoS
    - returns error when the client has no session or the filter is invalid
    # subscriptions
    -> std.strings.split
  pubsub_broker.unsubscribe
    @ (state: broker_state, client_id: string, filter: string) -> broker_state
    + removes the matching subscription if present
    # subscriptions
  pubsub_broker.filter_matches_topic
    @ (filter: string, topic: string) -> bool
    + returns true when the topic matches the filter, honoring '+' for one level and '#' for the remainder
    - returns false otherwise
    # routing
    -> std.strings.split
  pubsub_broker.route_publish
    @ (state: broker_state, topic: string, qos: i32) -> list[delivery]
    + returns one delivery per matching subscription with the effective QoS
    # routing
  pubsub_broker.publish
    @ (state: broker_state, topic: string, payload: bytes, qos: i32, retain: bool) -> tuple[broker_state, list[delivery]]
    + routes the message, returning state with any retained update applied
    + clears the retained message for the topic when payload is empty and retain is true
    # publishing
    -> std.time.now_millis
  pubsub_broker.retained_for_filter
    @ (state: broker_state, filter: string) -> list[retained]
    + returns retained messages whose topics match the filter
    # retained
  pubsub_broker.assign_packet_id
    @ (state: broker_state, client_id: string) -> tuple[broker_state, i32]
    + returns a fresh packet id for the client and the updated state
    # sessions
  pubsub_broker.acknowledge
    @ (state: broker_state, client_id: string, packet_id: i32) -> broker_state
    + releases the in-flight message tracked under the id
    # sessions
