# Requirement: "a caching library with multiple pluggable backend stores"

A single cache surface backed by interchangeable stores: in-memory, on-disk, and a generic remote key-value store.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes (and overwrites) the file
      # filesystem
    std.fs.remove
      @ (path: string) -> result[void, string]
      + deletes the file
      # filesystem
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

cache
  cache.new_memory_store
    @ () -> cache_store
    + creates an in-process map-backed store
    # backend
  cache.new_file_store
    @ (directory: string) -> cache_store
    + creates a store that serializes entries to files under directory
    ? key becomes the filename; value is the file contents
    # backend
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.fs.remove
  cache.new_remote_store
    @ (get: remote_get_fn, set: remote_set_fn, del: remote_del_fn) -> cache_store
    + creates a store that forwards operations to injected transport functions
    ? allows adapting any remote key-value service
    # backend
  cache.set
    @ (store: cache_store, key: string, value: bytes, ttl_ms: i64) -> result[cache_store, string]
    + stores the entry with an expiry ttl_ms from now
    # mutation
    -> std.time.now_millis
  cache.get
    @ (store: cache_store, key: string) -> result[optional[bytes], string]
    + returns the value when present and not expired
    - returns none when the entry is missing or expired
    # lookup
    -> std.time.now_millis
  cache.delete
    @ (store: cache_store, key: string) -> result[cache_store, string]
    + removes the entry
    # mutation
