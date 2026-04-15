# Requirement: "an event sourcing database engine"

An append-only log plus a snapshot mechanism and projection replay. Durability is delegated to a file primitive in std.

std
  std.fs
    std.fs.append_bytes
      fn (path: string, data: bytes) -> result[void, string]
      + appends bytes to a file, creating it if missing
      - returns error on permission failure
      # filesystem
    std.fs.read_range
      fn (path: string, offset: i64, length: i64) -> result[bytes, string]
      + reads a byte range from a file
      - returns error on out-of-bounds read
      # filesystem
    std.fs.file_size
      fn (path: string) -> result[i64, string]
      + returns the size of a file in bytes
      - returns error when the file does not exist
      # filesystem
    std.fs.fsync
      fn (path: string) -> result[void, string]
      + flushes file contents to durable storage
      # filesystem
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.hashing
    std.hashing.crc32
      fn (data: bytes) -> u32
      + returns a CRC-32 checksum
      # hashing

eventdb
  eventdb.open
    fn (log_path: string, snapshot_path: string) -> result[db_state, string]
    + opens or creates the log and snapshot files
    - returns error when paths are inaccessible
    -> std.fs.file_size
    # construction
  eventdb.append
    fn (state: db_state, stream: string, payload: bytes) -> result[tuple[db_state, i64], string]
    + appends one event and returns the assigned sequence number
    - returns error when durable write fails
    -> std.fs.append_bytes
    -> std.fs.fsync
    -> std.hashing.crc32
    -> std.time.now_millis
    # ingestion
  eventdb.read_from
    fn (state: db_state, from_seq: i64, limit: i32) -> result[list[event_record], string]
    + reads up to `limit` events starting at the given sequence number
    - returns error on a truncated or checksum-failing record
    -> std.fs.read_range
    -> std.hashing.crc32
    # querying
  eventdb.latest_seq
    fn (state: db_state) -> i64
    + returns the highest committed sequence number
    # querying
  eventdb.snapshot_write
    fn (state: db_state, stream: string, up_to_seq: i64, data: bytes) -> result[db_state, string]
    + stores a snapshot of a stream at the given sequence number
    -> std.fs.append_bytes
    -> std.fs.fsync
    # snapshot
  eventdb.snapshot_read
    fn (state: db_state, stream: string) -> result[snapshot_entry, string]
    + returns the most recent snapshot for the stream
    - returns error when no snapshot exists
    -> std.fs.read_range
    # snapshot
  eventdb.replay
    fn (state: db_state, stream: string, from_seq: i64, fold_tag: string, initial: bytes) -> result[bytes, string]
    + folds events from the given sequence using a named folding function
    - returns error when read fails
    # projection
  eventdb.compact
    fn (state: db_state, before_seq: i64) -> result[db_state, string]
    + rewrites the log dropping events covered by a snapshot below the given sequence
    - returns error when no snapshot covers the requested cutoff
    # maintenance
  eventdb.close
    fn (state: db_state) -> result[void, string]
    + flushes and releases the underlying files
    -> std.fs.fsync
    # teardown
  eventdb.event_record
    fn (seq: i64, stream: string, payload: bytes, timestamp_millis: i64) -> event_record
    + builds an event record
    # record
