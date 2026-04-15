# Requirement: "a multi-model NoSQL database"

A storage engine that supports key-value, hash, and list models against a single keyspace, with a write-ahead log for durability.

std
  std.fs
    std.fs.append_file
      fn (path: string, data: bytes) -> result[void, string]
      + appends bytes to a file, creating it if necessary
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.fsync
      fn (path: string) -> result[void, string]
      + flushes the file to durable storage
      # filesystem
  std.encoding
    std.encoding.varint_encode_u64
      fn (value: u64) -> bytes
      + encodes a u64 as little-endian varint
      # encoding
    std.encoding.varint_decode_u64
      fn (data: bytes, offset: i32) -> result[tuple[u64, i32], string]
      + returns (value, bytes_consumed)
      - returns error on truncated input
      # encoding

db
  db.engine_new
    fn (wal_path: string) -> result[engine_state, string]
    + creates an engine and replays the write-ahead log
    - returns error when the log is corrupt
    # construction
    -> std.fs.read_all
  db.kv_set
    fn (state: engine_state, key: string, value: bytes) -> result[engine_state, string]
    + sets a value under the key-value model
    # kv
    -> std.fs.append_file
    -> std.fs.fsync
  db.kv_get
    fn (state: engine_state, key: string) -> optional[bytes]
    + returns the stored value or none
    # kv
  db.hash_set_field
    fn (state: engine_state, key: string, field: string, value: bytes) -> result[engine_state, string]
    + sets a field on a hash-model key
    - returns error when the key is already bound to a different model
    # hash_model
    -> std.fs.append_file
  db.hash_get_field
    fn (state: engine_state, key: string, field: string) -> optional[bytes]
    + returns the field value or none
    # hash_model
  db.list_push
    fn (state: engine_state, key: string, value: bytes) -> result[engine_state, string]
    + appends a value to a list-model key
    - returns error when the key is bound to a non-list model
    # list_model
    -> std.fs.append_file
  db.list_range
    fn (state: engine_state, key: string, start: i32, stop: i32) -> list[bytes]
    + returns values in the inclusive range, negative indices count from the end
    # list_model
  db.delete
    fn (state: engine_state, key: string) -> engine_state
    + removes a key regardless of its model
    # deletion
    -> std.fs.append_file
  db.snapshot
    fn (state: engine_state, path: string) -> result[void, string]
    + writes a compact snapshot of the current state to disk
    # durability
    -> std.encoding.varint_encode_u64
    -> std.fs.append_file
    -> std.fs.fsync
  db.load_snapshot
    fn (path: string) -> result[engine_state, string]
    + rebuilds engine state from a snapshot
    - returns error on version mismatch
    # durability
    -> std.fs.read_all
    -> std.encoding.varint_decode_u64
