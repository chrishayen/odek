# Requirement: "a distributed locking library backed by a key-value store"

Locks are acquired by conditional writes to a shared key-value store with lease expiry and owner tokens. Clock reads go through a thin std utility.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.random_token
      fn (length: i32) -> string
      + returns a random hex token of the given length
      # identity
  std.kv
    std.kv.put_if_absent
      fn (key: string, value: bytes, ttl_ms: i64) -> result[bool, string]
      + returns true when the key was written, false when the key already exists
      - returns error on store failure
      # storage
    std.kv.get
      fn (key: string) -> result[optional[bytes], string]
      + returns the stored value or none when absent
      # storage
    std.kv.delete_if_equal
      fn (key: string, expected: bytes) -> result[bool, string]
      + deletes the key only when the current value matches expected
      - returns false when the value differs
      # storage

distlock
  distlock.acquire
    fn (key: string, lease_ms: i64) -> result[optional[lock_handle], string]
    + returns a handle with an owner token when the lock is acquired
    - returns none when another owner holds the lock and the lease is still valid
    # lock_acquire
    -> std.kv.put_if_absent
    -> std.id.random_token
    -> std.time.now_millis
  distlock.release
    fn (handle: lock_handle) -> result[bool, string]
    + releases the lock only when the stored owner token still matches
    - returns false when the lease already expired and someone else holds it
    # lock_release
    -> std.kv.delete_if_equal
  distlock.refresh
    fn (handle: lock_handle, extend_ms: i64) -> result[lock_handle, string]
    + extends the lease when the current owner still holds the lock
    - returns error when the owner token no longer matches
    # lock_refresh
    -> std.kv.get
    -> std.kv.put_if_absent
    -> std.time.now_millis
