# Requirement: "a library that batches small inserts and flushes them as large bulk requests to a database backend"

Accumulates rows in memory and flushes when a size or time threshold is reached.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

bulkwriter
  bulkwriter.new
    @ (max_rows: i32, max_bytes: i64, flush_interval_ms: i64) -> bulkwriter_state
    + creates a writer with the given thresholds
    # construction
  bulkwriter.add
    @ (state: bulkwriter_state, row: string) -> tuple[bool, bulkwriter_state]
    + appends a row; first element is true when a flush is due
    + triggers a due flush when row count reaches max_rows
    + triggers a due flush when accumulated bytes reach max_bytes
    # accumulation
    -> std.time.now_millis
  bulkwriter.take_pending
    @ (state: bulkwriter_state) -> tuple[list[string], bulkwriter_state]
    + returns the pending batch and a reset writer
    # flush
  bulkwriter.should_flush
    @ (state: bulkwriter_state) -> bool
    + returns true when the flush interval has elapsed since the first pending row
    # flush
    -> std.time.now_millis
