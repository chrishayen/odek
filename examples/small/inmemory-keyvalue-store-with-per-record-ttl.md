# Requirement: "an in-memory key-value store with per-record ttl"

Each record carries its own ttl; expired records are invisible to reads.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

kv
  kv.new
    fn () -> kv_state
    + creates an empty store
    # construction
  kv.set_with_ttl
    fn (k: kv_state, key: string, value: bytes, ttl_millis: i64) -> kv_state
    + stores the record with an absolute expiry of now + ttl_millis
    # write
    -> std.time.now_millis
  kv.get
    fn (k: kv_state, key: string) -> optional[bytes]
    + returns the value when present and unexpired
    - returns none when the record has expired
    # read
    -> std.time.now_millis
  kv.remove
    fn (k: kv_state, key: string) -> kv_state
    + drops the record regardless of expiry
    # write
  kv.keys
    fn (k: kv_state) -> list[string]
    + returns the keys of all currently unexpired records
    # query
    -> std.time.now_millis
