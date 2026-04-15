# Requirement: "a stream library with filter, map, flat_map, peek, sorted, for_each, reduce"

A lazy stream pipeline over a sequence of items. Each operator returns a new stream; terminal operators trigger evaluation.

std: (all units exist)

stream
  stream.from_list
    fn (items: list[i64]) -> stream_state
    + wraps a list as a stream source
    # construction
  stream.filter
    fn (s: stream_state, pred: fn(i64) -> bool) -> stream_state
    + returns a new stream that yields only items where pred is true
    # intermediate
  stream.map
    fn (s: stream_state, f: fn(i64) -> i64) -> stream_state
    + returns a new stream with f applied to each item
    # intermediate
  stream.flat_map
    fn (s: stream_state, f: fn(i64) -> list[i64]) -> stream_state
    + returns a new stream where each item is expanded into zero or more items
    # intermediate
  stream.peek
    fn (s: stream_state, observer: fn(i64) -> void) -> stream_state
    + returns a stream that calls observer on each item as it flows past
    ? peek is non-consuming; pipeline still requires a terminal op to run
    # intermediate
  stream.sorted
    fn (s: stream_state, less: fn(i64, i64) -> bool) -> stream_state
    + returns a new stream whose items are ordered by less
    ? sorted is a buffering operator; it materializes upstream before yielding
    # intermediate
  stream.for_each
    fn (s: stream_state, action: fn(i64) -> void) -> void
    + runs the pipeline and invokes action on every emitted item
    # terminal
  stream.reduce
    fn (s: stream_state, seed: i64, combine: fn(i64, i64) -> i64) -> i64
    + folds the stream left-to-right starting from seed
    + returns seed when the stream is empty
    # terminal
