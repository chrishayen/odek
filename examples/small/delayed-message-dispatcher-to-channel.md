# Requirement: "a library for scheduling when to dispatch a message to a channel"

A deterministic delay queue. Callers drive time by calling tick; the library returns messages whose dispatch time has arrived.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

delay_queue
  delay_queue.new
    fn () -> delay_queue_state
    + creates an empty delay queue
    # construction
  delay_queue.schedule
    fn (state: delay_queue_state, payload: bytes, dispatch_at_ms: i64) -> delay_queue_state
    + inserts a message into the queue ordered by dispatch time
    # scheduling
  delay_queue.schedule_in
    fn (state: delay_queue_state, payload: bytes, delay_ms: i64) -> delay_queue_state
    + schedules a message relative to the current time
    # scheduling
    -> std.time.now_millis
  delay_queue.drain_due
    fn (state: delay_queue_state, now_ms: i64) -> tuple[list[bytes], delay_queue_state]
    + returns all messages whose dispatch time is at or before now_ms, in order
    + leaves future messages in the queue
    # dispatch
  delay_queue.next_due_at
    fn (state: delay_queue_state) -> optional[i64]
    + returns the dispatch time of the earliest pending message
    - returns none when the queue is empty
    # introspection
