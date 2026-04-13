# Requirement: "a message queue library supporting topics, subscriptions, and durable delivery"

Topics, consumers with offsets, ack-based redelivery, and pluggable persistence.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.hash
    std.hash.fnv64
      @ (data: bytes) -> u64
      + computes a 64-bit FNV-1a hash
      # hashing

mq
  mq.new_broker
    @ () -> broker_state
    + creates an empty broker with no topics
    # construction
  mq.create_topic
    @ (broker: broker_state, name: string, partitions: i32) -> result[broker_state, string]
    + creates a topic with the given partition count
    - returns error when name is empty
    - returns error when partitions is non-positive
    - returns error when topic already exists
    # topic_management
  mq.publish
    @ (broker: broker_state, topic: string, key: bytes, payload: bytes) -> result[tuple[i64, broker_state], string]
    + assigns the message to a partition via key hash and returns its offset
    - returns error when topic is unknown
    # publishing
    -> std.hash.fnv64
    -> std.time.now_millis
  mq.subscribe
    @ (broker: broker_state, topic: string, group: string) -> result[tuple[subscription_id, broker_state], string]
    + registers a consumer group on the topic
    - returns error when topic is unknown
    # subscription
  mq.fetch
    @ (broker: broker_state, sub: subscription_id, max: i32) -> result[tuple[list[message], broker_state], string]
    + returns up to max unacked messages and marks them in-flight
    - returns empty list when no messages are pending
    - returns error when subscription is unknown
    # consumption
    -> std.time.now_millis
  mq.ack
    @ (broker: broker_state, sub: subscription_id, offsets: list[i64]) -> result[broker_state, string]
    + removes the given offsets from the in-flight set and advances the committed offset
    - returns error when any offset is not in-flight for this subscription
    # acknowledgement
  mq.nack
    @ (broker: broker_state, sub: subscription_id, offsets: list[i64]) -> result[broker_state, string]
    + returns the given offsets to the pending queue for redelivery
    - returns error when any offset is not in-flight
    # redelivery
  mq.redeliver_expired
    @ (broker: broker_state, ack_timeout_ms: i64) -> broker_state
    + returns in-flight messages older than ack_timeout_ms to the pending queue
    # redelivery
    -> std.time.now_millis
  mq.snapshot
    @ (broker: broker_state) -> bytes
    + serializes the full broker state for durable persistence
    # persistence
  mq.restore
    @ (snapshot: bytes) -> result[broker_state, string]
    + rebuilds a broker from a previous snapshot
    - returns error when snapshot is corrupt
    # persistence
  mq.lag
    @ (broker: broker_state, sub: subscription_id) -> result[i64, string]
    + returns the number of messages between committed offset and latest
    - returns error when subscription is unknown
    # inspection
