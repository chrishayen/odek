# Requirement: "a network link speed simulation wrapper for reader and writer interfaces"

Wraps an existing byte stream and throttles reads and writes to a configured bandwidth, blocking the caller when the token budget is exhausted.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns current unix time in nanoseconds
      # time
    std.time.sleep_nanos
      @ (duration: i64) -> void
      + blocks the caller for the given number of nanoseconds
      # time

link_sim
  link_sim.wrap_reader
    @ (inner: reader, bytes_per_second: i64) -> link_reader_state
    + wraps an existing reader with the given throttle rate
    ? bytes_per_second must be positive
    # construction
  link_sim.wrap_writer
    @ (inner: writer, bytes_per_second: i64) -> link_writer_state
    + wraps an existing writer with the given throttle rate
    # construction
  link_sim.read
    @ (state: link_reader_state, max: i32) -> result[tuple[bytes, link_reader_state], string]
    + returns bytes from the underlying reader after sleeping long enough to honor the rate
    - returns error when the underlying reader fails
    # throttling
    -> std.time.now_nanos
    -> std.time.sleep_nanos
  link_sim.write
    @ (state: link_writer_state, data: bytes) -> result[link_writer_state, string]
    + writes data after sleeping long enough to honor the rate
    - returns error when the underlying writer fails
    # throttling
    -> std.time.now_nanos
    -> std.time.sleep_nanos
