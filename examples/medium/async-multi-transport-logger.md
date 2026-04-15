# Requirement: "a multi-transport async logging library"

A logger that fans out structured records to one or more transports asynchronously via a bounded queue.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.sync
    std.sync.channel_new
      fn (capacity: i32) -> channel_state
      + creates a bounded channel with the given capacity
      # concurrency
    std.sync.channel_send
      fn (ch: channel_state, payload: bytes) -> result[void, string]
      + enqueues a payload, blocking when the channel is full
      - returns error when the channel is closed
      # concurrency
    std.sync.channel_recv
      fn (ch: channel_state) -> result[bytes, string]
      + dequeues the next payload, blocking until one is available
      - returns error when the channel is closed and drained
      # concurrency

logger
  logger.new
    fn (queue_capacity: i32) -> logger_state
    + creates a logger with an empty transport list and a bounded queue
    # construction
    -> std.sync.channel_new
  logger.add_transport
    fn (state: logger_state, name: string, sink: transport_sink) -> logger_state
    + registers a transport under the given name
    ? transport_sink is an opaque callable that accepts an encoded record
    # configuration
  logger.log
    fn (state: logger_state, level: string, message: string, fields: map[string,string]) -> result[void, string]
    + enqueues a record stamped with the current time for async delivery
    - returns error when the queue is closed
    # logging
    -> std.time.now_millis
    -> std.sync.channel_send
  logger.run_worker
    fn (state: logger_state) -> result[void, string]
    + drains the queue and dispatches each record to every registered transport
    + isolates transport failures so one slow sink cannot block the others
    - returns error when the queue is closed and drained
    # dispatch
    -> std.sync.channel_recv
  logger.close
    fn (state: logger_state) -> result[void, string]
    + closes the queue so the worker exits after draining
    # lifecycle
