# Requirement: "an embedded key-value database with multiple value types, backed by an append-only log and an LSM index"

Durable storage via write-ahead log, periodically merged sorted segments, and in-memory indexes for strings, lists, hashes, sets, and sorted sets.

std
  std.fs
    std.fs.open_append
      @ (path: string) -> result[file_handle, string]
      + opens a file for append, creating it when missing
      - returns error on permission failure
      # filesystem
    std.fs.write_all
      @ (handle: file_handle, data: bytes) -> result[void, string]
      + writes the buffer and flushes
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.rename
      @ (from: string, to: string) -> result[void, string]
      + atomically renames a file
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + lists entries in a directory
      # filesystem
  std.hash
    std.hash.crc32
      @ (data: bytes) -> u32
      + returns the CRC32 checksum
      # hashing
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

lsmkv
  lsmkv.open
    @ (dir: string) -> result[lsm_state, string]
    + opens the database at dir, replaying the write-ahead log into memory
    - returns error when the directory is unreadable
    # lifecycle
    -> std.fs.list_dir
    -> std.fs.read_all
  lsmkv.close
    @ (state: lsm_state) -> result[void, string]
    + flushes any buffered writes and releases file handles
    # lifecycle
  lsmkv.put_string
    @ (state: lsm_state, key: string, value: string) -> result[lsm_state, string]
    + writes a string entry to the WAL and memtable
    # strings
    -> std.fs.write_all
    -> std.hash.crc32
  lsmkv.get_string
    @ (state: lsm_state, key: string) -> optional[string]
    + looks up a string value, falling back to on-disk segments on memtable miss
    # strings
  lsmkv.list_push
    @ (state: lsm_state, key: string, value: string, at_head: bool) -> result[lsm_state, string]
    + pushes onto the list at key
    # lists
    -> std.fs.write_all
  lsmkv.list_range
    @ (state: lsm_state, key: string, start: i32, stop: i32) -> list[string]
    + returns a half-open slice of the list at key
    # lists
  lsmkv.hash_set
    @ (state: lsm_state, key: string, field: string, value: string) -> result[lsm_state, string]
    + sets a field inside the hash at key
    # hashes
    -> std.fs.write_all
  lsmkv.hash_get_all
    @ (state: lsm_state, key: string) -> map[string, string]
    + returns all fields and values of the hash at key
    # hashes
  lsmkv.set_add
    @ (state: lsm_state, key: string, member: string) -> result[lsm_state, string]
    + adds a member to the set at key
    # sets
    -> std.fs.write_all
  lsmkv.zset_add
    @ (state: lsm_state, key: string, member: string, score: f64) -> result[lsm_state, string]
    + inserts or updates a scored member in the sorted set
    # sorted_sets
    -> std.fs.write_all
  lsmkv.zset_range
    @ (state: lsm_state, key: string, start: i32, stop: i32) -> list[string]
    + returns members in the given rank range ordered by score
    # sorted_sets
  lsmkv.delete
    @ (state: lsm_state, key: string) -> result[lsm_state, string]
    + writes a tombstone for the key across all types
    # mutation
    -> std.fs.write_all
  lsmkv.flush_memtable
    @ (state: lsm_state) -> result[lsm_state, string]
    + writes the memtable as a sorted on-disk segment and rotates the WAL
    # storage
    -> std.fs.open_append
    -> std.fs.write_all
    -> std.fs.rename
    -> std.time.now_millis
  lsmkv.compact
    @ (state: lsm_state) -> result[lsm_state, string]
    + merges overlapping segments and drops tombstoned keys
    # storage
    -> std.fs.list_dir
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.fs.rename
