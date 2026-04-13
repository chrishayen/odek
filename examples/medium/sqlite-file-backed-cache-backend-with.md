# Requirement: "a persistent key-value cache backed by a local database and file blobs"

Small values live in a local database; large values spill to files. Lookups are durable across process restarts.

std
  std.fs
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + atomically writes data to path
      - returns error when the parent directory does not exist
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when the path does not exist
      # filesystem
    std.fs.remove
      @ (path: string) -> result[void, string]
      + deletes the file at path
      - returns error when the path does not exist
      # filesystem
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

diskcache
  diskcache.open
    @ (dir: string, inline_limit_bytes: i32) -> result[cache_state, string]
    + opens or creates a cache rooted at dir; values up to inline_limit_bytes are stored inline
    - returns error when dir cannot be created
    # construction
  diskcache.set
    @ (c: cache_state, key: string, value: bytes, ttl_seconds: i64) -> result[void, string]
    + stores value under key with an expiry of now + ttl_seconds
    + spills to a blob file when value exceeds inline_limit_bytes
    - returns error when the underlying store rejects the write
    # write
    -> std.fs.write_all
    -> std.time.now_seconds
  diskcache.get
    @ (c: cache_state, key: string) -> result[optional[bytes], string]
    + returns the value when present and not expired
    + returns none when the key is absent
    - returns none when the stored entry has expired
    # read
    -> std.fs.read_all
    -> std.time.now_seconds
  diskcache.delete
    @ (c: cache_state, key: string) -> result[void, string]
    + removes the entry and any spill file for key
    ? delete on a missing key is not an error
    # write
    -> std.fs.remove
  diskcache.evict_expired
    @ (c: cache_state) -> result[i32, string]
    + removes all entries whose expiry is in the past and returns the number removed
    # maintenance
    -> std.time.now_seconds
