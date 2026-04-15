# Requirement: "a message streaming bridge between a range of protocols"

Reads messages from a pluggable source, applies a pipeline of transformations, and writes to a pluggable sink. Backpressure, batching, and acknowledgement are first-class.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current time in milliseconds
      # time
  std.json
    std.json.parse_object
      fn (raw: bytes) -> result[map[string,string], string]
      + parses a JSON object
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string,string]) -> bytes
      + encodes an object to JSON bytes
      # serialization

msgbridge
  msgbridge.new_bridge
    fn (source: source_handle, sink: sink_handle) -> bridge_state
    + creates a bridge binding one source to one sink
    # construction
  msgbridge.register_source
    fn (name: string, reader: reader_fn, acker: acker_fn) -> source_handle
    + returns a source handle using the given read and ack callbacks
    # sources
  msgbridge.register_sink
    fn (name: string, writer: writer_fn) -> sink_handle
    + returns a sink handle using the given write callback
    # sinks
  msgbridge.add_transform
    fn (bridge: bridge_state, transform: transform_fn) -> bridge_state
    + appends a per-message transform to the pipeline
    ? transforms return a new message or drop the message
    # pipeline
  msgbridge.add_filter
    fn (bridge: bridge_state, predicate: predicate_fn) -> bridge_state
    + appends a predicate that drops messages returning false
    # pipeline
  msgbridge.apply_pipeline
    fn (bridge: bridge_state, msg: bytes) -> optional[bytes]
    + runs the message through all transforms and filters in order
    + returns none when any stage drops the message
    # pipeline
  msgbridge.step
    fn (bridge: bridge_state, batch_size: i32) -> result[i32, string]
    + reads up to batch_size messages, runs the pipeline, writes survivors to the sink, and acks originals
    + returns the number of messages processed
    - returns error when the source read fails
    - returns error when the sink write fails
    # execution
    -> std.time.now_millis
  msgbridge.run_until_empty
    fn (bridge: bridge_state) -> result[i64, string]
    + repeatedly steps the bridge until the source yields zero messages
    + returns total messages processed
    # execution
  msgbridge.encode_as_json
    fn (fields: map[string,string]) -> bytes
    + convenience transform that encodes a field map as JSON
    # codec
    -> std.json.encode_object
  msgbridge.decode_as_json
    fn (raw: bytes) -> result[map[string,string], string]
    + convenience transform that parses a JSON message into fields
    # codec
    -> std.json.parse_object
  msgbridge.stats
    fn (bridge: bridge_state) -> map[string,i64]
    + returns counters for messages_in, messages_out, dropped, and errors
    # observability
