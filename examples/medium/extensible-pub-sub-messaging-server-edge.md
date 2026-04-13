# Requirement: "an extensible pub/sub messaging library for edge applications"

Topic-based broker with pluggable filters. Subscribers receive published messages that pass the topic's filter chain.

std
  std.ids
    std.ids.new_id
      @ () -> string
      + returns a fresh opaque identifier unique within the process
      # identifiers
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

pubsub
  pubsub.new_broker
    @ () -> broker_state
    + creates a broker with no topics and no subscribers
    # construction
  pubsub.create_topic
    @ (state: broker_state, name: string) -> result[broker_state, string]
    + registers a new topic with an empty filter chain
    - returns error when the topic name already exists
    # topics
  pubsub.subscribe
    @ (state: broker_state, topic: string, handler: fn(bytes) -> void) -> result[tuple[string, broker_state], string]
    + returns (subscription_id, new_state) when the topic exists
    - returns error when the topic is unknown
    # subscription
    -> std.ids.new_id
  pubsub.unsubscribe
    @ (state: broker_state, subscription_id: string) -> broker_state
    + removes the subscription; idempotent when the id is unknown
    # subscription
  pubsub.add_filter
    @ (state: broker_state, topic: string, filter: fn(bytes) -> bool) -> result[broker_state, string]
    + appends a filter to the topic's filter chain
    - returns error when the topic is unknown
    # filters
  pubsub.publish
    @ (state: broker_state, topic: string, message: bytes) -> result[i32, string]
    + dispatches the message to every subscriber whose filter chain accepts it; returns the count delivered
    - returns error when the topic is unknown
    # publishing
    -> std.time.now_millis
