# Requirement: "a high-level event-driven consumer/producer supporting pluggable message broker dialects"

Dialects are represented as an opaque handle implementing the transport; the project layer provides a uniform pub/sub facade.

std
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      # serialization

commander
  commander.register_dialect
    @ (name: string, impl: dialect_impl) -> void
    + registers a transport under a name (e.g. "kafka", "nats")
    # dialects
  commander.connect
    @ (dialect: string, endpoint: string) -> result[broker_conn, string]
    + opens a connection using the named dialect
    - returns error when the dialect is unregistered
    - returns error when the transport fails
    # connection
  commander.close
    @ (conn: broker_conn) -> result[void, string]
    + closes the underlying transport
    # connection
  commander.publish
    @ (conn: broker_conn, topic: string, headers: map[string, string], payload: bytes) -> result[void, string]
    + serializes headers and payload as an event and publishes it to topic
    # producer
    -> std.json.encode_object
  commander.subscribe
    @ (conn: broker_conn, topic: string, group: string) -> result[subscription, string]
    + joins a consumer group on topic and returns a subscription handle
    # consumer
  commander.poll
    @ (sub: subscription, max_events: i32, timeout_ms: i64) -> result[list[event], string]
    + blocks up to timeout_ms and returns up to max_events buffered events
    - returns an empty list when no events arrive within the timeout
    # consumer
    -> std.json.parse_object
  commander.commit
    @ (sub: subscription, event_ids: list[string]) -> result[void, string]
    + acknowledges processing of the listed event ids
    # consumer
  commander.on_event
    @ (sub: subscription, handler: fn(ev: event) -> bool) -> result[void, string]
    + runs a dispatch loop invoking handler per event; commits when handler returns true
    # consumer
    -> commander.poll
    -> commander.commit
