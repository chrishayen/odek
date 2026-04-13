# Requirement: "a realtime distributed messaging platform"

A pub/sub broker with topics, channels, persistent-on-disk buffers, and consumer acknowledgements. The project layer owns topic/channel bookkeeping and delivery; std provides time, disk, and identifier primitives.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.new_ulid
      @ () -> string
      + returns a newly generated ulid string
      # identifiers
  std.fs
    std.fs.append_bytes
      @ (path: string, data: bytes) -> result[void, string]
      + appends data to the file at path, creating it if needed
      - returns error when the parent directory does not exist
      # filesystem
    std.fs.read_range
      @ (path: string, offset: i64, length: i64) -> result[bytes, string]
      + reads length bytes from path starting at offset
      - returns error on out-of-range reads
      # filesystem

messaging
  messaging.new_broker
    @ (data_dir: string) -> broker_state
    + creates a broker that persists buffered messages under data_dir
    # construction
  messaging.create_topic
    @ (broker: broker_state, topic: string) -> broker_state
    + registers a new topic with an empty channel set
    ? duplicate create is a no-op
    # topic_management
  messaging.create_channel
    @ (broker: broker_state, topic: string, channel: string) -> result[broker_state, string]
    + registers a new channel on the topic with its own consumer cursor
    - returns error when the topic does not exist
    # channel_management
  messaging.publish
    @ (broker: broker_state, topic: string, body: bytes) -> result[tuple[message_id, broker_state], string]
    + assigns a message id, appends the body to the topic log, fanouts to channels
    - returns error when the topic does not exist
    # publishing
    -> std.id.new_ulid
    -> std.time.now_millis
    -> std.fs.append_bytes
  messaging.subscribe
    @ (broker: broker_state, topic: string, channel: string, consumer: consumer_id) -> result[broker_state, string]
    + attaches a consumer to the channel's delivery queue
    - returns error when the channel does not exist
    # consumption
  messaging.fetch_next
    @ (broker: broker_state, consumer: consumer_id) -> result[optional[message], string]
    + returns the next unacknowledged message for the consumer (or none)
    + marks the message as in-flight with a redelivery deadline
    # consumption
    -> std.fs.read_range
    -> std.time.now_millis
  messaging.ack
    @ (broker: broker_state, consumer: consumer_id, id: message_id) -> result[broker_state, string]
    + acknowledges a message, removing it from in-flight
    - returns error when the message is not in-flight for the consumer
    # consumption
  messaging.requeue_expired
    @ (broker: broker_state) -> broker_state
    + moves messages whose in-flight deadline passed back to pending
    # delivery_guarantees
    -> std.time.now_millis
  messaging.topic_depth
    @ (broker: broker_state, topic: string, channel: string) -> i64
    + returns the number of pending messages in the channel
    # query
  messaging.delete_channel
    @ (broker: broker_state, topic: string, channel: string) -> result[broker_state, string]
    + removes the channel, discarding its pending and in-flight messages
    - returns error when the channel does not exist
    # channel_management
