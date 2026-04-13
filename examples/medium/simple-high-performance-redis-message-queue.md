# Requirement: "a high-performance message queue backed by a pluggable key-value store"

Producer/consumer queue with acknowledgements, retries, and delayed delivery. The storage backend is abstracted.

std
  std.kv
    std.kv.lpush
      @ (store: kv_handle, key: string, value: bytes) -> result[void, string]
      + prepends a value to the list at key
      - returns error when the key exists and holds a non-list type
      # storage
    std.kv.rpop
      @ (store: kv_handle, key: string) -> result[optional[bytes], string]
      + removes and returns the tail element, or none when the list is empty
      # storage
    std.kv.zadd
      @ (store: kv_handle, key: string, score: f64, member: bytes) -> result[void, string]
      + adds a scored member to a sorted set
      # storage
    std.kv.zrange_by_score
      @ (store: kv_handle, key: string, max_score: f64, limit: i32) -> result[list[bytes], string]
      + returns members with score <= max_score, up to limit
      # storage
  std.time
    std.time.now_millis
      @ () -> i64
      + returns the current unix time in milliseconds
      # time
  std.id
    std.id.new_v4
      @ () -> string
      + returns a random 128-bit identifier as a canonical string
      # identifiers

mq
  mq.new_queue
    @ (store: kv_handle, name: string) -> queue_handle
    + binds a queue handle to a backing store and logical name
    # construction
  mq.publish
    @ (q: queue_handle, payload: bytes, delay_ms: i64) -> result[string, string]
    + enqueues a message and returns its id
    + when delay_ms > 0, the message is hidden until that time elapses
    - returns error when the store rejects the write
    # producer
    -> std.id.new_v4
    -> std.time.now_millis
    -> std.kv.lpush
    -> std.kv.zadd
  mq.consume
    @ (q: queue_handle) -> result[optional[queue_message], string]
    + returns the next due message and marks it in-flight
    + returns none when no ready message exists
    # consumer
    -> std.time.now_millis
    -> std.kv.rpop
    -> std.kv.zrange_by_score
  mq.ack
    @ (q: queue_handle, message_id: string) -> result[void, string]
    + removes the in-flight message
    - returns error when the message id is unknown
    # consumer
  mq.nack
    @ (q: queue_handle, message_id: string, retry_delay_ms: i64) -> result[void, string]
    + requeues an in-flight message after retry_delay_ms
    - returns error when the message id is unknown
    # consumer
    -> std.kv.zadd
  mq.reclaim_expired
    @ (q: queue_handle, visibility_ms: i64) -> result[i32, string]
    + requeues in-flight messages whose visibility window expired; returns the count reclaimed
    # reliability
    -> std.time.now_millis
