# Requirement: "a minimal embedded key-value database"

A small database with a page-based storage engine, B+ tree index, and write-ahead log.

std
  std.fs
    std.fs.read_at
      @ (path: string, offset: i64, length: i32) -> result[bytes, string]
      + reads a slice of a file at an offset
      - returns error on short read
      # filesystem
    std.fs.write_at
      @ (path: string, offset: i64, data: bytes) -> result[void, string]
      + writes bytes at an absolute file offset
      - returns error on partial write
      # filesystem
    std.fs.fsync
      @ (path: string) -> result[void, string]
      + flushes pending writes to stable storage
      # filesystem
  std.hash
    std.hash.crc32
      @ (data: bytes) -> u32
      + computes CRC32 of the input
      # hashing
  std.bytes
    std.bytes.read_u32_le
      @ (data: bytes, offset: i32) -> result[u32, string]
      + reads a little-endian u32 at offset
      - returns error on out-of-range offset
      # byte_reading
    std.bytes.write_u32_le
      @ (buf: bytes, offset: i32, value: u32) -> result[bytes, string]
      + writes a little-endian u32 at offset
      - returns error on out-of-range offset
      # byte_writing

nosqldb
  nosqldb.open
    @ (path: string) -> result[db_state, string]
    + opens or creates a database file and replays the write-ahead log
    - returns error on corrupted header
    # lifecycle
    -> std.fs.read_at
  nosqldb.close
    @ (state: db_state) -> result[void, string]
    + flushes the log and closes file handles
    # lifecycle
    -> std.fs.fsync
  nosqldb.get
    @ (state: db_state, key: bytes) -> result[optional[bytes], string]
    + returns the value for a key or none when missing
    # read
  nosqldb.put
    @ (state: db_state, key: bytes, value: bytes) -> result[db_state, string]
    + inserts or replaces the value for a key
    - returns error when the key exceeds max size
    # write
  nosqldb.delete
    @ (state: db_state, key: bytes) -> result[db_state, string]
    + removes a key from the tree
    # write
  nosqldb.alloc_page
    @ (state: db_state) -> tuple[db_state, page_id]
    + reserves a fresh page id from the free list or extends the file
    # pager
  nosqldb.read_page
    @ (state: db_state, id: page_id) -> result[bytes, string]
    + reads the raw bytes for a page
    - returns error when the page id is out of range
    # pager
    -> std.fs.read_at
  nosqldb.write_page
    @ (state: db_state, id: page_id, data: bytes) -> result[db_state, string]
    + writes a page and updates the cache
    # pager
    -> std.fs.write_at
  nosqldb.btree_search
    @ (state: db_state, root: page_id, key: bytes) -> result[optional[bytes], string]
    + walks the B+ tree from root looking for the key
    # index
  nosqldb.btree_insert
    @ (state: db_state, root: page_id, key: bytes, value: bytes) -> result[tuple[db_state, page_id], string]
    + inserts into the B+ tree and returns the new root after splits
    # index
  nosqldb.wal_append
    @ (state: db_state, record: bytes) -> result[db_state, string]
    + appends a framed record to the write-ahead log with a checksum
    # wal
    -> std.hash.crc32
    -> std.fs.write_at
  nosqldb.wal_replay
    @ (state: db_state) -> result[db_state, string]
    + replays the log into the tree at open time
    - returns error when a record checksum mismatches
    # wal
    -> std.hash.crc32
