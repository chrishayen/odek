# Requirement: "a disk-backed key-value store"

A simple store that sharded keys across files in a directory, one file per key.

std
  std.fs
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes the entire file atomically via a temp file and rename
      - returns error on I/O failure
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file into memory
      - returns error when the file does not exist
      # filesystem
    std.fs.remove_file
      @ (path: string) -> result[void, string]
      + deletes a file
      - returns error when the file does not exist
      # filesystem
    std.fs.make_dir_all
      @ (path: string) -> result[void, string]
      + creates a directory and all missing parents
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns the names of entries in a directory
      - returns error when the directory does not exist
      # filesystem
  std.crypto
    std.crypto.sha1_hex
      @ (data: bytes) -> string
      + returns the lowercase hex SHA-1 of the input
      # cryptography

kv_disk
  kv_disk.open
    @ (base_dir: string) -> result[store_state, string]
    + opens or creates a store rooted at base_dir
    - returns error when the directory cannot be created
    # construction
    -> std.fs.make_dir_all
  kv_disk.key_path
    @ (store: store_state, key: string) -> string
    + returns the on-disk path for a key using a 2-character shard prefix
    ? sharding uses the first two hex chars of SHA-1 of the key
    # layout
    -> std.crypto.sha1_hex
  kv_disk.put
    @ (store: store_state, key: string, value: bytes) -> result[void, string]
    + writes a value for a key atomically
    - returns error on I/O failure
    # writes
    -> std.fs.write_all
  kv_disk.get
    @ (store: store_state, key: string) -> result[optional[bytes], string]
    + returns the value for a key
    - returns none when the key does not exist
    - returns error on I/O failure other than not-found
    # reads
    -> std.fs.read_all
  kv_disk.delete
    @ (store: store_state, key: string) -> result[void, string]
    + removes a key's file
    - returns error when the key does not exist
    # writes
    -> std.fs.remove_file
  kv_disk.keys
    @ (store: store_state) -> result[list[string], string]
    + returns all keys by scanning every shard directory
    ? scans are O(n); suitable for maintenance, not hot paths
    # iteration
    -> std.fs.list_dir
