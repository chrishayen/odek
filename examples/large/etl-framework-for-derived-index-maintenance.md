# Requirement: "an ETL framework for building and keeping a derived index fresh"

Sources emit records, transforms rewrite them, and sinks apply updates to a target index. The framework tracks source checkpoints so reruns only process new data.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.hashing
    std.hashing.fnv64
      @ (data: bytes) -> u64
      + returns a 64-bit FNV-1a hash of the input
      # hashing

etlkit
  etlkit.source_spec
    @ (name: string, read_fn: string, checkpoint_key: string) -> source_spec
    + declares a named source with a read callback reference and checkpoint identifier
    # declaration
  etlkit.transform_spec
    @ (name: string, transform_fn: string) -> transform_spec
    + declares a named row-level transform
    # declaration
  etlkit.sink_spec
    @ (name: string, apply_fn: string) -> sink_spec
    + declares a named sink with an apply callback reference
    # declaration
  etlkit.pipeline
    @ (source: source_spec, transforms: list[transform_spec], sink: sink_spec) -> pipeline
    + composes a source, ordered transforms, and a sink into a runnable pipeline
    # composition
  etlkit.load_checkpoint
    @ (state: etl_state, key: string) -> optional[i64]
    + returns the last processed offset for the given checkpoint key
    - returns none when no prior run exists
    # checkpointing
  etlkit.save_checkpoint
    @ (state: etl_state, key: string, offset: i64) -> etl_state
    + records the offset as the latest processed value
    -> std.time.now_millis
    # checkpointing
  etlkit.fingerprint_record
    @ (record: map[string, string]) -> u64
    + returns a stable fingerprint for idempotent sink application
    -> std.hashing.fnv64
    # hashing
  etlkit.apply_transforms
    @ (transforms: list[transform_spec], record: map[string, string]) -> optional[map[string, string]]
    + runs transforms in order, stopping when any transform drops the record
    - returns none when a transform signals filter
    # transformation
  etlkit.run_once
    @ (state: etl_state, pipe: pipeline) -> result[run_stats, string]
    + reads from source since last checkpoint, applies transforms, writes to sink, advances checkpoint
    - returns error when source read or sink apply fails
    # execution
  etlkit.run_stats
    @ (read: i64, written: i64, dropped: i64, started_millis: i64, ended_millis: i64) -> run_stats
    + builds a run statistics record
    # reporting
  etlkit.diff_index
    @ (target_count: i64, source_count: i64) -> i64
    + returns source_count - target_count as a freshness gauge
    # reporting
