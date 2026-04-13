# Requirement: "a distributed unique id generator producing time-ordered ids encoded in a compact alphabet"

Ids combine a timestamp, a machine id, and a per-tick sequence. Encoding uses a base58 alphabet.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.encoding
    std.encoding.base58_encode
      @ (value: bytes) -> string
      + encodes bytes to a base58 alphabet (no 0, O, I, l)
      + returns "" for empty input
      # encoding
    std.encoding.base58_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes a base58 string back to bytes
      - returns error on characters outside the alphabet
      # encoding

idgen
  idgen.new
    @ (machine_id: u16, epoch_ms: i64) -> idgen_state
    + creates a generator rooted at the given epoch and machine id
    ? the epoch lets 41 bits of timestamp cover ~69 years
    # construction
  idgen.next
    @ (state: idgen_state) -> tuple[u64, idgen_state]
    + returns a monotonically increasing 64-bit id
    + advances sequence within the same millisecond
    + resets sequence when the millisecond advances
    # id_generation
    -> std.time.now_millis
  idgen.encode
    @ (id: u64) -> string
    + returns the id rendered in the compact alphabet
    # serialization
    -> std.encoding.base58_encode
  idgen.decode
    @ (encoded: string) -> result[u64, string]
    + parses a compact-alphabet string back into a 64-bit id
    - returns error when the decoded byte length is wrong
    # serialization
    -> std.encoding.base58_decode
  idgen.parts
    @ (id: u64) -> tuple[i64, u16, u16]
    + returns (timestamp_ms_since_epoch, machine_id, sequence) extracted from the id
    # introspection
