# Requirement: "distributed generation of super short, unique, non-sequential, url-friendly ids"

A generator parameterized by a worker id so separate processes never collide. Output is a short base-57 encoding of time + worker + counter.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + current unix time in milliseconds
      # time
  std.random
    std.random.u32
      @ () -> u32
      + returns a uniformly random u32
      # random

shortid
  shortid.new
    @ (worker_id: u16, alphabet: string) -> result[shortid_state, string]
    + returns a generator for this worker using the given base-n alphabet
    - returns error when alphabet has fewer than 32 characters
    - returns error when alphabet contains duplicates
    # construction
  shortid.next
    @ (state: shortid_state) -> tuple[string, shortid_state]
    + returns the next id and the updated state
    + packs timestamp_ms, worker_id, and per-millisecond counter, then base-n encodes
    ? counter resets when the timestamp advances; within the same millisecond the counter increments
    # generation
    -> std.time.now_millis
  shortid.encode_base
    @ (value: u64, alphabet: string) -> string
    + encodes a u64 using the given alphabet as the radix
    # encoding
  shortid.shuffle_alphabet
    @ (alphabet: string, seed: u32) -> string
    + returns a deterministic permutation of alphabet so that ids appear non-sequential
    # obfuscation
    -> std.random.u32
