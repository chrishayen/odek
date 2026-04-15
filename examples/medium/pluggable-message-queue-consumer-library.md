# Requirement: "a consumer library for a pluggable message queue service"

Polls a queue, dispatches messages to a handler, and manages acknowledgement and retry.

std
  std.time
    std.time.sleep_millis
      fn (millis: i64) -> void
      + blocks the current fiber for the given duration
      # time

consumer
  consumer.new
    fn (queue_url: string, handler: callback) -> consumer_state
    + creates a consumer bound to a queue url and user handler
    # construction
  consumer.set_batch_size
    fn (state: consumer_state, size: i32) -> consumer_state
    + sets how many messages to request per poll (1..10)
    # configuration
  consumer.poll_once
    fn (state: consumer_state, fetch: callback) -> result[i32, string]
    + fetches one batch via fetch and dispatches each message to the handler
    + returns the number of successfully handled messages
    - records a failure when the handler raises
    # polling
  consumer.ack
    fn (state: consumer_state, receipt: string, ack_fn: callback) -> result[void, string]
    + acknowledges a single message by receipt via ack_fn
    - returns error when ack_fn reports a transport failure
    # acknowledgement
  consumer.run_forever
    fn (state: consumer_state, fetch: callback, poll_interval_ms: i64) -> void
    + loops polling and sleeping between empty batches
    # lifecycle
    -> std.time.sleep_millis
