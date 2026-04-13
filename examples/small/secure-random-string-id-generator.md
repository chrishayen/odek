# Requirement: "a cryptographically secure random string id generator"

Produces a fresh id of a requested length using a secure random source and a configurable alphabet.

std
  std.crypto
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically secure random bytes
      ? n must be non-negative
      # cryptography

secure_id
  secure_id.generate
    @ (length: i32, alphabet: string) -> string
    + returns a fresh id of the given length drawn uniformly from alphabet
    - returns "" when length is 0 or alphabet is empty
    ? uses rejection sampling over the byte source to avoid modulo bias
    # generation
    -> std.crypto.random_bytes
