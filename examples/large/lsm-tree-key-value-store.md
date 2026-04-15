# Requirement: "a high-performance NoSQL key-value store backed by an LSM tree"

Multiple data types (strings, hashes, lists, sets, sorted sets) layered over a log-structured storage engine.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file into bytes
      - returns error when the file does not exist
      # filesystem
    std.fs.append
      fn (path: string, data: bytes) -> result[void, string]
      + appends bytes to a file, creating it if absent
      # filesystem
  std.hash
    std.hash.xxh64
      fn (data: bytes) -> u64
      + returns a 64-bit non-cryptographic hash
      # hashing

kvstore
  kvstore.open
    fn (path: string) -> result[kv_state, string]
    + opens or creates an on-disk store rooted at path
    - returns error when the path cannot be created
    # storage
    -> std.fs.read_all
  kvstore.close
    fn (state: kv_state) -> result[void, string]
    + flushes in-memory buffers and closes underlying files
    # storage
  kvstore.set
    fn (state: kv_state, key: string, value: bytes) -> kv_state
    + writes a string-keyed value
    # strings
    -> std.fs.append
  kvstore.get
    fn (state: kv_state, key: string) -> optional[bytes]
    + returns a value for a key if present
    # strings
  kvstore.del
    fn (state: kv_state, key: string) -> kv_state
    + removes a key, no-op if absent
    # strings
  kvstore.hset
    fn (state: kv_state, key: string, field: string, value: bytes) -> kv_state
    + stores a value in a hash under key/field
    # hashes
  kvstore.hget
    fn (state: kv_state, key: string, field: string) -> optional[bytes]
    + returns a field in a hash if present
    # hashes
  kvstore.lpush
    fn (state: kv_state, key: string, value: bytes) -> kv_state
    + prepends a value to a list
    # lists
  kvstore.rpop
    fn (state: kv_state, key: string) -> tuple[optional[bytes], kv_state]
    + removes and returns the last element of a list
    # lists
  kvstore.sadd
    fn (state: kv_state, key: string, member: string) -> kv_state
    + adds a member to a set
    # sets
  kvstore.smembers
    fn (state: kv_state, key: string) -> list[string]
    + returns the members of a set
    # sets
  kvstore.zadd
    fn (state: kv_state, key: string, member: string, score: f64) -> kv_state
    + inserts or updates a scored member in a sorted set
    # sorted_sets
  kvstore.zrange
    fn (state: kv_state, key: string, start: i32, stop: i32) -> list[string]
    + returns members in a score range by index
    # sorted_sets
  kvstore.compact
    fn (state: kv_state) -> result[kv_state, string]
    + merges older segments to reclaim space
    - returns error on partial write
    # compaction
    -> std.hash.xxh64
