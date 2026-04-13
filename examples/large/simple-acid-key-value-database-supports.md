# Requirement: "an ACID key-value store with transactions"

Durability is achieved via an append-only write-ahead log; atomicity via transaction batches that commit as a single log record.

std
  std.fs
    std.fs.open_append
      @ (path: string) -> result[file_handle, string]
      + opens a file for appending, creating it if missing
      - returns error when the containing directory does not exist
      # filesystem
    std.fs.write
      @ (handle: file_handle, data: bytes) -> result[void, string]
      + appends bytes to the file
      # filesystem
    std.fs.fsync
      @ (handle: file_handle) -> result[void, string]
      + flushes file contents to disk
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file into memory
      - returns error when the file does not exist
      # filesystem
  std.encoding
    std.encoding.u32_be_encode
      @ (n: u32) -> bytes
      + encodes a 32-bit unsigned integer as big-endian
      # encoding
    std.encoding.u32_be_decode
      @ (data: bytes, offset: i32) -> result[u32, string]
      + decodes a big-endian u32 at the given offset
      - returns error when offset + 4 exceeds length
      # encoding
  std.hash
    std.hash.crc32
      @ (data: bytes) -> u32
      + computes the CRC-32 checksum of data
      # hashing

kv_store
  kv_store.open
    @ (path: string) -> result[kv_handle, string]
    + opens or creates the store at path, replaying the write-ahead log into memory
    - returns error when a record has an invalid checksum
    # lifecycle
    -> std.fs.read_all
    -> std.fs.open_append
    -> std.encoding.u32_be_decode
    -> std.hash.crc32
  kv_store.close
    @ (handle: kv_handle) -> result[void, string]
    + flushes pending writes and closes the underlying file
    # lifecycle
    -> std.fs.fsync
  kv_store.get
    @ (handle: kv_handle, key: string) -> optional[bytes]
    + returns the value for the given key or none
    # read
  kv_store.begin
    @ (handle: kv_handle) -> txn_state
    + starts a new transaction with an empty write set
    # transactions
  kv_store.txn_put
    @ (txn: txn_state, key: string, value: bytes) -> txn_state
    + stages a key-value pair in the transaction's write set
    # transactions
  kv_store.txn_delete
    @ (txn: txn_state, key: string) -> txn_state
    + stages a deletion in the transaction's write set
    # transactions
  kv_store.commit
    @ (handle: kv_handle, txn: txn_state) -> result[void, string]
    + atomically appends the transaction as a single log record and applies it to the in-memory map
    - returns error when the underlying write fails, leaving state unchanged
    # transactions
    -> std.encoding.u32_be_encode
    -> std.hash.crc32
    -> std.fs.write
    -> std.fs.fsync
  kv_store.abort
    @ (txn: txn_state) -> void
    + discards the transaction's staged writes without touching the store
    # transactions
  kv_store.keys
    @ (handle: kv_handle) -> list[string]
    + returns all live keys in insertion order
    # read
