# Requirement: "a fault-tolerant message streaming library layered on a pub/sub transport"

Durable, replicated message streams with append, subscribe-from-offset, and leader-election semantics on top of a pluggable pub/sub transport.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.fs
    std.fs.append_bytes
      fn (path: string, data: bytes) -> result[void, string]
      + appends bytes to a file, creating it if missing
      - returns error when the parent directory does not exist
      # filesystem
    std.fs.read_range
      fn (path: string, offset: i64, length: i64) -> result[bytes, string]
      + reads a contiguous range of bytes from a file
      - returns error on out-of-range offset
      # filesystem
  std.encoding
    std.encoding.varint_encode
      fn (value: u64) -> bytes
      + encodes an unsigned integer as a variable-length byte sequence
      # encoding
    std.encoding.varint_decode
      fn (data: bytes, pos: i64) -> result[tuple[u64, i64], string]
      + returns (value, next_position)
      - returns error on truncated input
      # encoding

message_stream
  message_stream.create_stream
    fn (name: string, partitions: i32, replicas: i32) -> result[stream_state, string]
    + creates a new stream with the given partition and replica counts
    - returns error when partitions or replicas are non-positive
    # stream_lifecycle
  message_stream.append
    fn (state: stream_state, partition: i32, payload: bytes) -> result[i64, string]
    + appends payload and returns the assigned monotonic offset
    - returns error when the partition id is unknown
    # append
    -> std.fs.append_bytes
    -> std.encoding.varint_encode
    -> std.time.now_millis
  message_stream.read_from
    fn (state: stream_state, partition: i32, offset: i64, max_messages: i32) -> result[list[bytes], string]
    + returns up to max_messages payloads starting at offset
    - returns error when offset is beyond the current tail
    # read
    -> std.fs.read_range
    -> std.encoding.varint_decode
  message_stream.subscribe
    fn (state: stream_state, partition: i32, start_offset: i64) -> subscription_state
    + creates a subscription cursor positioned at start_offset
    # subscription
  message_stream.poll
    fn (sub: subscription_state, max_messages: i32) -> tuple[list[bytes], subscription_state]
    + returns next messages and advances the cursor
    + returns empty list when no new messages are available
    # subscription
  message_stream.elect_leader
    fn (state: stream_state, partition: i32, candidate_id: string, term: i64) -> result[bool, string]
    + returns true when the candidate becomes leader for the partition
    - returns false when another candidate holds a higher term
    # leadership
  message_stream.replicate
    fn (state: stream_state, partition: i32, follower_id: string, up_to_offset: i64) -> result[i64, string]
    + returns the offset through which the follower has been replicated
    - returns error when the partition has no leader
    # replication
  message_stream.acknowledge
    fn (state: stream_state, partition: i32, offset: i64, replicas: i32) -> bool
    + returns true when the offset is acknowledged by the required replica count
    # durability
  message_stream.commit_offset
    fn (sub: subscription_state, offset: i64) -> subscription_state
    + persists a consumer-committed offset for the subscription
    # consumer_offsets
    -> std.fs.append_bytes
  message_stream.truncate_before
    fn (state: stream_state, partition: i32, offset: i64) -> result[void, string]
    + removes all messages before the given offset
    - returns error on unknown partition
    # retention
