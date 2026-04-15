# Requirement: "a client for reliable queues backed by a stream-based store"

Publish and consume messages with at-least-once delivery via explicit acknowledgement.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.new_ulid
      fn () -> string
      + returns a lexicographically sortable unique identifier
      # identifiers

reliable_queue
  reliable_queue.connect
    fn (endpoint: string) -> result[queue_state, string]
    + opens a session against the given store endpoint
    - returns error when the endpoint is unreachable
    # connection
  reliable_queue.publish
    fn (state: queue_state, topic: string, payload: bytes) -> result[string, string]
    + appends a message to the topic stream and returns its id
    - returns error when the topic cannot accept writes
    # publishing
    -> std.id.new_ulid
    -> std.time.now_millis
  reliable_queue.subscribe
    fn (state: queue_state, topic: string, group: string) -> result[subscription_state, string]
    + creates or resumes a consumer group on the topic
    - returns error when the group is locked by another consumer
    # subscription
  reliable_queue.next
    fn (sub: subscription_state) -> result[optional[queue_message], string]
    + returns the next unacknowledged message for the group
    + returns none when the stream is drained
    - returns error when the connection is closed
    # consumption
    -> std.time.now_millis
  reliable_queue.ack
    fn (sub: subscription_state, message_id: string) -> result[void, string]
    + marks the message as successfully processed
    - returns error when the message is not owned by this consumer
    # acknowledgement
