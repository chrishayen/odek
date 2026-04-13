# Requirement: "a distributed locking library backed by a shared key-value store"

Implements a leased-lock protocol on top of an abstract atomic SETNX-with-ttl primitive. The backend is pluggable so the locking logic can be tested against any store that exposes conditional set and delete.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.rand
    std.rand.bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding

dlock
  dlock.new_token
    @ () -> string
    + returns a fresh 16-byte random lease token as hex
    # tokens
    -> std.rand.bytes
    -> std.encoding.hex_encode
  dlock.try_acquire
    @ (backend: kv_backend, key: string, ttl_millis: i64) -> result[lock_handle, string]
    + returns a handle when the backend accepts a conditional set of key with a new token
    - returns error when the key already exists and has not expired
    # acquisition
    -> dlock.new_token
    -> std.time.now_millis
  dlock.release
    @ (backend: kv_backend, handle: lock_handle) -> result[void, string]
    + deletes the key only when its current value matches the lease token
    - returns error when the token no longer matches (lease was lost)
    # release
  dlock.refresh
    @ (backend: kv_backend, handle: lock_handle, ttl_millis: i64) -> result[lock_handle, string]
    + extends the ttl when the lease is still held
    - returns error when the token no longer matches
    # renewal
    -> std.time.now_millis
