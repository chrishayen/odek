# Requirement: "a distributed key-value store based on an LSM tree"

A log-structured merge tree key-value store with writes going to a memtable and write-ahead log, flushed to sorted level-0 files, and compacted into deeper levels. Replication sends committed writes to peers.

std
  std.fs
    std.fs.append
      fn (path: string, data: bytes) -> result[void, string]
      + appends bytes to a file, creating it if needed
      - returns error on write failure
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the full file contents
      - returns error when the path does not exist
      # filesystem
    std.fs.write_atomic
      fn (path: string, data: bytes) -> result[void, string]
      + writes to a temp file and renames into place
      - returns error when the target directory is not writable
      # filesystem
  std.hash
    std.hash.crc32
      fn (data: bytes) -> u32
      + returns the CRC32 of data
      # hashing
  std.net
    std.net.send_message
      fn (peer_addr: string, payload: bytes) -> result[void, string]
      + sends a framed message to a peer
      - returns error when the peer is unreachable
      # networking

lsmkv
  lsmkv.open
    fn (data_dir: string) -> result[store_state, string]
    + opens or creates a store at the directory, replaying the write-ahead log
    - returns error when data_dir cannot be created
    # construction
    -> std.fs.read_all
  lsmkv.put
    fn (store: store_state, key: bytes, value: bytes) -> result[store_state, string]
    + inserts into the memtable and appends to the write-ahead log
    - returns error on log write failure
    # writes
    -> std.fs.append
    -> std.hash.crc32
  lsmkv.get
    fn (store: store_state, key: bytes) -> optional[bytes]
    + searches memtable then level-0 files then deeper levels
    + returns none when the key is absent or tombstoned
    # reads
  lsmkv.delete
    fn (store: store_state, key: bytes) -> result[store_state, string]
    + writes a tombstone entry
    - returns error on log write failure
    # writes
    -> std.fs.append
  lsmkv.flush_memtable
    fn (store: store_state) -> result[store_state, string]
    + sorts the memtable and writes a new level-0 SSTable
    - returns error on file write failure
    # flushing
    -> std.fs.write_atomic
  lsmkv.compact_level
    fn (store: store_state, level: i32) -> result[store_state, string]
    + merges files at the given level into the next level, dropping overwritten keys
    - returns error on file write failure
    # compaction
    -> std.fs.write_atomic
  lsmkv.sstable_search
    fn (path: string, key: bytes) -> result[optional[bytes], string]
    + binary-searches an SSTable for the given key
    - returns error on a corrupted table
    # reads
    -> std.fs.read_all
  lsmkv.iterate_range
    fn (store: store_state, start: bytes, end: bytes) -> list[tuple[bytes, bytes]]
    + returns key-value pairs within [start, end) in sorted order
    # reads
  lsmkv.add_peer
    fn (store: store_state, peer_addr: string) -> store_state
    + registers a replication peer
    # replication
  lsmkv.replicate_write
    fn (store: store_state, key: bytes, value: bytes) -> list[result[void, string]]
    + ships the write to every peer and returns per-peer results
    # replication
    -> std.net.send_message
  lsmkv.apply_remote_write
    fn (store: store_state, key: bytes, value: bytes, peer_seq: u64) -> result[store_state, string]
    + applies a write received from a peer if peer_seq is newer than the last applied sequence
    - returns error on out-of-order sequence
    # replication
