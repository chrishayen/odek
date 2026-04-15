# Requirement: "a high-performance asynchronous logging library"

Log calls enqueue records on a bounded channel; a background worker drains the channel, formats records, and writes them to one or more sinks. Formatting and I/O never block the caller.

std
  std.channel
    std.channel.new_bounded
      fn (capacity: i32) -> chan_state
      + creates a multi-producer single-consumer channel with fixed capacity
      # concurrency
    std.channel.try_send
      fn (chan: chan_state, item: log_record) -> bool
      + returns false without blocking when the channel is full
      # concurrency
    std.channel.recv
      fn (chan: chan_state) -> optional[log_record]
      + blocks until an item is available or the channel is closed
      # concurrency
    std.channel.close
      fn (chan: chan_state) -> void
      + marks the channel closed so recv returns none after draining
      # concurrency
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns current unix time in nanoseconds
      # time
  std.fs
    std.fs.append
      fn (path: string, data: bytes) -> result[void, string]
      + appends bytes to a file, creating it if missing
      # filesystem

asynclog
  asynclog.new
    fn (capacity: i32) -> asynclog_state
    + creates a logger with an empty sink list and a bounded record channel
    # construction
    -> std.channel.new_bounded
  asynclog.add_sink
    fn (state: asynclog_state, sink: log_sink) -> asynclog_state
    + registers a sink that will receive every formatted line
    # configuration
  asynclog.log
    fn (state: asynclog_state, level: string, message: string) -> bool
    + returns true if the record was enqueued, false if the channel was full and the record was dropped
    # ingest
    -> std.time.now_nanos
    -> std.channel.try_send
  asynclog.worker_step
    fn (state: asynclog_state) -> bool
    + drains the next record from the channel and writes it to every sink
    + returns false when the channel is closed and drained
    # worker
    -> std.channel.recv
  asynclog.format_record
    fn (record: log_record) -> string
    + renders a record as "<iso_timestamp> <LEVEL> <message>\n"
    # formatting
  asynclog.file_sink
    fn (path: string) -> log_sink
    + returns a sink that appends every formatted line to the given path
    # sink
    -> std.fs.append
  asynclog.shutdown
    fn (state: asynclog_state) -> void
    + closes the channel; subsequent log calls return false
    # lifecycle
    -> std.channel.close
