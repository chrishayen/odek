# Requirement: "a key/value database using an lsm tree with a b+tree index"

A write-ahead log feeds a mutable memtable, which flushes into sorted on-disk segments; reads consult the memtable first and then fall through the segments via a b+tree index.

std
  std.fs
    std.fs.append
      @ (path: string, data: bytes) -> result[void, string]
      + appends data to path, creating it if missing
      - returns error when the directory is not writable
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file fully
      # filesystem
    std.fs.rename
      @ (from: string, to: string) -> result[void, string]
      + atomically renames a file
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + lists directory entries
      # filesystem
  std.hash
    std.hash.crc32
      @ (data: bytes) -> u32
      + returns the crc32 checksum of data
      # hashing

kvdb
  kvdb.open
    @ (dir: string) -> result[kvdb_state, string]
    + opens or creates a database at dir, replaying the write-ahead log
    - returns error when the directory cannot be read
    # lifecycle
    -> std.fs.list_dir
    -> std.fs.read_all
  kvdb.put
    @ (db: kvdb_state, key: bytes, value: bytes) -> result[kvdb_state, string]
    + appends the write to the log and updates the memtable
    # write_path
    -> std.fs.append
    -> std.hash.crc32
  kvdb.delete
    @ (db: kvdb_state, key: bytes) -> result[kvdb_state, string]
    + records a tombstone in the log and memtable
    # write_path
    -> std.fs.append
  kvdb.get
    @ (db: kvdb_state, key: bytes) -> optional[bytes]
    + returns the newest value across the memtable and segments, honoring tombstones
    - returns none for a missing or deleted key
    # read_path
  kvdb.flush_memtable
    @ (db: kvdb_state) -> result[kvdb_state, string]
    + writes the memtable to a new sorted segment and truncates the log
    # flush
    -> std.fs.rename
  kvdb.compact
    @ (db: kvdb_state) -> result[kvdb_state, string]
    + merges overlapping segments, dropping tombstones and obsolete versions
    # compaction
  kvdb.build_btree_index
    @ (segment_path: string) -> result[btree_index, string]
    + builds an in-memory b+tree over a segment's keys for lookups
    - returns error when the segment is corrupt
    # indexing
    -> std.fs.read_all
    -> std.hash.crc32
  kvdb.range
    @ (db: kvdb_state, start: bytes, end: bytes) -> list[tuple[bytes, bytes]]
    + returns all live key/value pairs whose keys fall in [start, end)
    # scan
  kvdb.close
    @ (db: kvdb_state) -> result[void, string]
    + flushes the memtable and releases handles
    # lifecycle
