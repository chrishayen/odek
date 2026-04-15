# Requirement: "an embedded key-value store optimized for read-heavy workloads"

On-disk hash index with memory-mapped data regions and an append-only WAL for durability.

std
  std.fs
    std.fs.open_append
      fn (path: string) -> result[file_handle, string]
      + opens a file for append, creating it when missing
      # filesystem
    std.fs.write_all
      fn (handle: file_handle, data: bytes) -> result[void, string]
      + writes and flushes the buffer
      # filesystem
    std.fs.read_at
      fn (path: string, offset: i64, n: i32) -> result[bytes, string]
      + reads n bytes starting at offset
      - returns error when offset lies past end of file
      # filesystem
    std.fs.file_size
      fn (path: string) -> result[i64, string]
      + returns the file size in bytes
      # filesystem
  std.hash
    std.hash.fnv1a_64
      fn (data: bytes) -> u64
      + computes a 64-bit FNV-1a hash
      # hashing

readykv
  readykv.open
    fn (dir: string) -> result[kv_state, string]
    + opens the store at dir, building the hash index from the data file
    - returns error when the directory is unreadable
    # lifecycle
    -> std.fs.file_size
    -> std.fs.read_at
  readykv.close
    fn (state: kv_state) -> result[void, string]
    + flushes pending writes and releases file handles
    # lifecycle
  readykv.put
    fn (state: kv_state, key: string, value: bytes) -> result[kv_state, string]
    + appends a record and updates the in-memory hash index
    # mutation
    -> std.fs.open_append
    -> std.fs.write_all
    -> std.hash.fnv1a_64
  readykv.get
    fn (state: kv_state, key: string) -> optional[bytes]
    + looks up a key via the hash index and reads the record from disk
    # query
    -> std.fs.read_at
    -> std.hash.fnv1a_64
  readykv.delete
    fn (state: kv_state, key: string) -> result[kv_state, string]
    + appends a tombstone and removes the key from the index
    # mutation
    -> std.fs.open_append
    -> std.fs.write_all
  readykv.compact
    fn (state: kv_state) -> result[kv_state, string]
    + rewrites the data file with live records only and rebuilds the index
    # storage
    -> std.fs.read_at
    -> std.fs.write_all
  readykv.keys
    fn (state: kv_state) -> list[string]
    + returns every live key in the store
    # query
