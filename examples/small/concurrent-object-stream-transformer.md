# Requirement: "a library for transforming object streams concurrently"

An object-at-a-time stream transformer with bounded parallelism. Output order matches input order.

std: (all units exist)

concurrent_stream
  concurrent_stream.new
    fn (concurrency: i32, transform: fn(any) -> any) -> stream_state
    + returns a transformer configured with the worker count and mapping function
    ? concurrency must be >= 1; callers pass 1 for sequential semantics
    # construction
  concurrent_stream.push
    fn (state: stream_state, item: any) -> stream_state
    + enqueues an input item for processing
    # input
  concurrent_stream.drain
    fn (state: stream_state) -> list[any]
    + runs the transform over all pending items and returns the outputs in input order
    + preserves ordering even when later items finish before earlier ones
    # execution
  concurrent_stream.close
    fn (state: stream_state) -> void
    + releases any worker resources; further push/drain calls are no-ops
    # lifecycle
