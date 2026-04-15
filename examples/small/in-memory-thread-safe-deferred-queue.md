# Requirement: "an in-memory, thread-safe deferred queue"

A priority queue keyed by ready-at time. Items become visible once their deferral has elapsed.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

dqueue
  dqueue.new
    fn () -> dqueue_state
    + returns an empty deferred queue
    # construction
  dqueue.push
    fn (q: dqueue_state, item: bytes, delay_ms: i64) -> dqueue_state
    + inserts an item that becomes ready after delay_ms
    ? a delay of 0 means immediately available
    # producer
    -> std.time.now_millis
  dqueue.pop_ready
    fn (q: dqueue_state) -> tuple[optional[bytes], dqueue_state]
    + returns the earliest ready item and removes it
    - returns none when no item is due yet
    # consumer
    -> std.time.now_millis
  dqueue.size
    fn (q: dqueue_state) -> i32
    + returns the total number of queued items, ready or not
    # introspection
