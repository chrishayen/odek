# Requirement: "byte counters that wrap readers, writers, and response writers"

Transparent wrappers that observe byte volume flowing through a stream without changing semantics.

std: (all units exist)

byte_counter
  byte_counter.wrap_reader
    fn (inner: reader_handle) -> counted_reader_state
    + returns a reader that delegates to inner and tracks total bytes read
    # wrapping
  byte_counter.wrap_writer
    fn (inner: writer_handle) -> counted_writer_state
    + returns a writer that delegates to inner and tracks total bytes written
    # wrapping
  byte_counter.wrap_response
    fn (inner: response_writer_handle) -> counted_response_state
    + returns a response writer that tracks bytes written to the response body
    # wrapping
  byte_counter.total
    fn (state: counter_handle) -> i64
    + returns the total bytes observed by any counter wrapper
    ? the same query works across reader, writer, and response wrappers
    # reporting
