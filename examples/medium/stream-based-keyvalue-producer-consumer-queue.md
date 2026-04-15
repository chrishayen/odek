# Requirement: "a producer/consumer queue backed by a stream-based key-value store"

Queue-over-streams needs a producer that appends entries, a consumer that reads a group, and ack/retry handling. The stream backend is abstracted behind a pluggable connection.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.uuid
    std.uuid.new_v4
      fn () -> string
      + returns a random 128-bit identifier in canonical form
      # identifiers

queue
  queue.new_producer
    fn (conn: stream_conn, stream: string) -> producer_state
    + creates a producer bound to the given stream name
    # construction
  queue.enqueue
    fn (state: producer_state, payload: map[string,string]) -> result[string, string]
    + appends an entry and returns its stream id
    - returns error when the connection rejects the write
    # production
    -> std.uuid.new_v4
  queue.new_consumer
    fn (conn: stream_conn, stream: string, group: string, name: string) -> consumer_state
    + creates or reuses a named consumer within a group
    # construction
  queue.poll
    fn (state: consumer_state, max: i32, wait_millis: i32) -> result[list[queue_message], string]
    + blocks up to wait_millis and returns at most max pending messages for this consumer
    - returns error when the group does not exist
    # consumption
    -> std.time.now_millis
  queue.ack
    fn (state: consumer_state, id: string) -> result[void, string]
    + marks the message as successfully processed and removes it from pending
    - returns error when the id is not pending for this consumer
    # acknowledgement
  queue.reclaim_stale
    fn (state: consumer_state, idle_millis: i64) -> result[list[queue_message], string]
    + reclaims pending messages whose last delivery is older than idle_millis
    # retry
    -> std.time.now_millis
