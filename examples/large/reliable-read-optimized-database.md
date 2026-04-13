# Requirement: "a reliable database optimised for read operations"

A log-structured key-value store: writes go to an append-only log, reads hit an in-memory index, and compaction rewrites live records into a fresh file.

std
  std.fs
    std.fs.open_append
      @ (path: string) -> result[file_handle, string]
      + opens a file for appending, creating it if absent
      - returns error when the parent directory is missing
      # filesystem
    std.fs.append_bytes
      @ (handle: file_handle, data: bytes) -> result[i64, string]
      + appends bytes and returns the offset where they were written
      - returns error on io failure
      # filesystem
    std.fs.read_at
      @ (path: string, offset: i64, length: i64) -> result[bytes, string]
      + reads a range of bytes from a file
      - returns error when the range is out of bounds
      # filesystem
    std.fs.rename
      @ (from: string, to: string) -> result[void, string]
      + atomically renames a file
      - returns error on io failure
      # filesystem
  std.hash
    std.hash.crc32
      @ (data: bytes) -> u32
      + returns the CRC32 of the data
      # hashing

readdb
  readdb.open
    @ (path: string) -> result[db_state, string]
    + opens a database at the given path, replaying the log to rebuild the index
    - returns error when the log is corrupt past the last valid record
    # lifecycle
    -> std.fs.open_append
    -> std.fs.read_at
    -> std.hash.crc32
  readdb.put
    @ (state: db_state, key: bytes, value: bytes) -> result[db_state, string]
    + appends a write record and updates the in-memory index
    - returns error on io failure
    # write
    -> std.fs.append_bytes
    -> std.hash.crc32
  readdb.get
    @ (state: db_state, key: bytes) -> result[optional[bytes], string]
    + returns the current value for a key, or none if absent
    - returns error on io failure or checksum mismatch
    # read
    -> std.fs.read_at
    -> std.hash.crc32
  readdb.delete
    @ (state: db_state, key: bytes) -> result[db_state, string]
    + appends a tombstone and removes the key from the index
    # write
    -> std.fs.append_bytes
    -> std.hash.crc32
  readdb.keys
    @ (state: db_state) -> list[bytes]
    + returns every live key in insertion order
    # read
  readdb.compact
    @ (state: db_state) -> result[db_state, string]
    + rewrites live records to a new log file and swaps it into place
    - returns error on io failure during rewrite
    # compaction
    -> std.fs.open_append
    -> std.fs.append_bytes
    -> std.fs.rename
  readdb.snapshot
    @ (state: db_state) -> db_state
    + returns a read-only view that ignores subsequent writes
    # read
  readdb.close
    @ (state: db_state) -> result[void, string]
    + flushes pending writes and releases the file handle
    # lifecycle
