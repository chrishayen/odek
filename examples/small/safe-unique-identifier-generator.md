# Requirement: "a safe fast unique identifier generator"

Sortable, time-prefixed IDs encoded in a url-safe alphabet.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.random
    std.random.fill
      fn (dst: bytes) -> bytes
      + fills the given buffer with cryptographically random bytes
      # randomness

uid
  uid.generate
    fn () -> string
    + returns a 26-character sortable id with a millisecond timestamp prefix
    + two ids generated in the same millisecond compare by their random suffix
    # generation
    -> std.time.now_millis
    -> std.random.fill
  uid.parse_timestamp
    fn (id: string) -> result[i64, string]
    + returns the millisecond timestamp embedded in the id
    - returns error when the id is not exactly 26 characters
    - returns error when the timestamp segment contains invalid characters
    # parsing
  uid.is_valid
    fn (id: string) -> bool
    + returns true when the id has the expected length and alphabet
    # validation
