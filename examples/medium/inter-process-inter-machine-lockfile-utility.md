# Requirement: "a cross-process and cross-machine lockfile utility"

Coordinates mutual exclusion by creating and maintaining a lockfile with an owner token and freshness timestamp.

std
  std.fs
    std.fs.create_exclusive
      @ (path: string, data: bytes) -> result[void, string]
      + creates the file atomically, failing if it already exists
      - returns error when the file exists
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when missing
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the path, replacing existing contents
      # filesystem
    std.fs.remove
      @ (path: string) -> result[void, string]
      + removes the file at path
      - returns error when absent
      # filesystem
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.random
    std.random.hex_token
      @ (n_bytes: i32) -> string
      + returns a cryptographically random hex token of the given byte length
      # random

lockfile
  lockfile.acquire
    @ (path: string, stale_after_seconds: i64) -> result[lock_handle, string]
    + creates the lockfile with a fresh owner token and returns a handle
    + breaks and reclaims a lock whose heartbeat is older than stale_after_seconds
    - returns error when a live lock is already held
    # acquire
    -> std.fs.create_exclusive
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.time.now_seconds
    -> std.random.hex_token
  lockfile.refresh
    @ (handle: lock_handle) -> result[lock_handle, string]
    + updates the heartbeat timestamp for the held lock
    - returns error when the owner token no longer matches
    # heartbeat
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.time.now_seconds
  lockfile.release
    @ (handle: lock_handle) -> result[void, string]
    + removes the lockfile when the owner token still matches
    - returns error when another owner has taken over
    # release
    -> std.fs.read_all
    -> std.fs.remove
  lockfile.is_locked
    @ (path: string, stale_after_seconds: i64) -> bool
    + returns true when a fresh lockfile exists at the path
    + returns false when the file is absent or stale
    # query
    -> std.fs.read_all
    -> std.time.now_seconds
