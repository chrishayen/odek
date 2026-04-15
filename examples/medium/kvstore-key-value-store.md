# Requirement: "a key/value store"

Values live in an in-memory map, backed by an append-only file so state survives restarts.

std
  std.fs
    std.fs.open_append
      fn (path: string) -> result[file_handle, string]
      + opens a file for appending, creating it if absent
      - returns error when the parent directory is missing
      # filesystem
    std.fs.append_bytes
      fn (handle: file_handle, data: bytes) -> result[i64, string]
      + appends bytes and returns the offset written
      - returns error on io failure
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire contents of a file
      - returns error when the file does not exist
      # filesystem

kvstore
  kvstore.open
    fn (path: string) -> result[kv_state, string]
    + opens a store, replaying the log file to rebuild the in-memory map
    - returns error when the log is truncated mid-record
    # lifecycle
    -> std.fs.read_all
    -> std.fs.open_append
  kvstore.set
    fn (state: kv_state, key: string, value: bytes) -> result[kv_state, string]
    + updates the map and appends a set record to the log
    - returns error on io failure
    # write
    -> std.fs.append_bytes
  kvstore.get
    fn (state: kv_state, key: string) -> optional[bytes]
    + returns the current value or none
    # read
  kvstore.delete
    fn (state: kv_state, key: string) -> result[kv_state, string]
    + removes the key and appends a delete record
    - returns error on io failure
    # write
    -> std.fs.append_bytes
  kvstore.keys
    fn (state: kv_state) -> list[string]
    + returns every live key
    # read
  kvstore.close
    fn (state: kv_state) -> result[void, string]
    + releases the file handle
    # lifecycle
