# Requirement: "a library that encodes and decodes short URL-safe ids from lists of non-negative integers"

Reversibly maps a list of numbers to and from a compact id string using a configurable alphabet.

std: (all units exist)

sqids
  sqids.new
    fn (alphabet: string, min_length: i32) -> result[sqids_config, string]
    + creates a config from an alphabet and minimum output length
    - returns error when alphabet has fewer than 5 characters or contains duplicates
    # construction
  sqids.encode
    fn (config: sqids_config, numbers: list[u64]) -> string
    + encodes a list of non-negative integers into a compact id string of at least min_length
    ? shuffles the alphabet between numbers so identical inputs don't produce visible repetition
    # encoding
  sqids.decode
    fn (config: sqids_config, id: string) -> list[u64]
    + decodes a previously encoded id back into its list of numbers
    - returns an empty list when id contains characters outside the alphabet
    # decoding
