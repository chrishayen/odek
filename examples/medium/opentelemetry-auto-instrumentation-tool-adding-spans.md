# Requirement: "a distributed tracing auto-instrumentation library for wrapping functions in spans"

Creates and finishes spans, propagates context, and exports to a pluggable sink.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns current unix time in nanoseconds
      # time
  std.crypto
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n bytes from a cryptographic RNG
      # randomness

tracer
  tracer.new_trace_id
    @ () -> string
    + returns a hex-encoded 16-byte trace id
    # ids
    -> std.crypto.random_bytes
  tracer.new_span_id
    @ () -> string
    + returns a hex-encoded 8-byte span id
    # ids
    -> std.crypto.random_bytes
  tracer.start_span
    @ (parent: optional[span_context], name: string) -> span_handle
    + returns a span stamped with the current time and a fresh span id
    + inherits the parent's trace id when parent is present
    # span_lifecycle
    -> std.time.now_nanos
    -> tracer.new_trace_id
    -> tracer.new_span_id
  tracer.set_attribute
    @ (s: span_handle, key: string, value: string) -> span_handle
    + adds a string attribute to the span
    # span_lifecycle
  tracer.record_error
    @ (s: span_handle, message: string) -> span_handle
    + marks the span as errored and records the message
    # span_lifecycle
  tracer.end_span
    @ (s: span_handle) -> finished_span
    + stamps the end time and returns an immutable finished span
    # span_lifecycle
    -> std.time.now_nanos
  tracer.context_from
    @ (s: span_handle) -> span_context
    + returns the trace/span id pair a child needs
    # propagation
  tracer.export
    @ (spans: list[finished_span], sink: sink_handle) -> result[i32, string]
    + sends all finished spans to the sink and returns the count accepted
    - returns error when the sink rejects a span
    # export
