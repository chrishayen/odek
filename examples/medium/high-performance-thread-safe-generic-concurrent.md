# Requirement: "a thread-safe concurrent hash map using open addressing with striped locks"

Shards the key space across fixed-size striped locks so independent keys don't contend.

std
  std.hash
    std.hash.hash_bytes
      @ (data: bytes) -> u64
      + returns a stable 64-bit hash of the input
      # hashing
  std.sync
    std.sync.mutex_new
      @ () -> mutex_handle
      + creates an unlocked mutex
      # concurrency
    std.sync.mutex_with
      @ (m: mutex_handle, body: fn() -> void) -> void
      + runs body while holding the mutex
      # concurrency

concurrent_map
  concurrent_map.new
    @ (num_shards: i32, initial_capacity: i32) -> cmap_state
    + allocates num_shards open-addressed tables, each with its own mutex
    # construction
    -> std.sync.mutex_new
  concurrent_map.shard_for
    @ (m: cmap_state, key: bytes) -> i32
    + returns the shard index for key using the hash low bits
    # sharding
    -> std.hash.hash_bytes
  concurrent_map.put
    @ (m: cmap_state, key: bytes, value: bytes) -> void
    + inserts or replaces under the shard lock
    + resizes the shard's table when load factor exceeds a threshold
    # mutation
    -> concurrent_map.shard_for
    -> std.sync.mutex_with
  concurrent_map.get
    @ (m: cmap_state, key: bytes) -> optional[bytes]
    + returns the value for key under the shard lock
    - returns none when key is absent
    # lookup
    -> concurrent_map.shard_for
    -> std.sync.mutex_with
  concurrent_map.remove
    @ (m: cmap_state, key: bytes) -> bool
    + deletes the entry and returns true when it was present
    # mutation
    -> concurrent_map.shard_for
    -> std.sync.mutex_with
  concurrent_map.len
    @ (m: cmap_state) -> i64
    + returns the total element count across all shards
    # introspection
