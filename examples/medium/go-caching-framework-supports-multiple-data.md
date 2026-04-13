# Requirement: "a caching framework that supports pluggable backend drivers"

A cache facade whose backend is chosen at construction time; drivers implement get/set/delete against any store.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.serialization
    std.serialization.encode
      @ (value: bytes) -> bytes
      + returns a length-prefixed encoding of value
      # serialization
    std.serialization.decode
      @ (raw: bytes) -> result[bytes, string]
      + returns the payload from a length-prefixed encoding
      - returns error when the prefix does not match the payload length
      # serialization

cache
  cache.register_driver
    @ (name: string, driver: cache_driver) -> void
    + associates a driver implementation with a name
    # registration
  cache.new
    @ (driver_name: string) -> result[cache_state, string]
    + constructs a cache backed by the named driver
    - returns error when the driver name is not registered
    # construction
  cache.set
    @ (state: cache_state, key: string, value: bytes, ttl_seconds: i64) -> result[void, string]
    + stores value under key with the given ttl
    - returns error when the driver rejects the write
    # storage
    -> std.time.now_seconds
    -> std.serialization.encode
  cache.get
    @ (state: cache_state, key: string) -> result[optional[bytes], string]
    + returns the value for key when present and unexpired
    + returns none when the entry has expired
    - returns error when the driver fails to read
    # retrieval
    -> std.time.now_seconds
    -> std.serialization.decode
  cache.delete
    @ (state: cache_state, key: string) -> result[void, string]
    + removes the entry for key
    - returns error when the driver fails to delete
    # eviction
