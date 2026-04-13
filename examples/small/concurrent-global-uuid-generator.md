# Requirement: "a globally unique identifier generator safe for concurrent server use"

Generates compact 12-byte identifiers composed of timestamp, machine id, process id, and a counter.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns the current unix time in seconds
      # time
  std.random
    std.random.bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness

xid
  xid.new_generator
    @ (machine_id: bytes, process_id: i32) -> xid_state
    + creates a generator seeded with a random counter
    ? machine_id is three bytes; process_id occupies two bytes in the id
    # construction
    -> std.random.bytes
  xid.generate
    @ (state: xid_state) -> tuple[bytes, xid_state]
    + returns a new 12-byte id and the advanced generator state
    + the counter increments monotonically across calls
    # id_generation
    -> std.time.now_seconds
  xid.encode
    @ (id: bytes) -> string
    + returns the 20-character base32 representation of a 12-byte id
    # formatting
  xid.decode
    @ (encoded: string) -> result[bytes, string]
    + parses a 20-character base32 string back to 12 bytes
    - returns error on invalid length or characters
    # parsing
