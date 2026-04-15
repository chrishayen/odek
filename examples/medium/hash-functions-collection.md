# Requirement: "a collection of cryptographic hash functions"

A library exposing several hash primitives through a common streaming interface and convenience one-shot wrappers.

std: (all units exist)

hashes
  hashes.new_sha256
    fn () -> hash_state
    + creates a SHA-256 hasher
    # construction
  hashes.new_sha512
    fn () -> hash_state
    + creates a SHA-512 hasher
    # construction
  hashes.new_blake2b
    fn (digest_size: i32) -> result[hash_state, string]
    + creates a BLAKE2b hasher producing the requested number of output bytes
    - returns error when digest_size is not between 1 and 64
    # construction
  hashes.update
    fn (state: hash_state, data: bytes) -> hash_state
    + absorbs more input into the hasher
    # streaming
  hashes.finalize
    fn (state: hash_state) -> bytes
    + returns the final digest
    ? the hasher must not be updated after finalize
    # finalization
  hashes.sha256
    fn (data: bytes) -> bytes
    + one-shot: returns the 32-byte SHA-256 digest of the input
    # one_shot
  hashes.sha512
    fn (data: bytes) -> bytes
    + one-shot: returns the 64-byte SHA-512 digest of the input
    # one_shot
  hashes.blake2b
    fn (data: bytes, digest_size: i32) -> result[bytes, string]
    + one-shot: returns the BLAKE2b digest of the input
    - returns error when digest_size is not between 1 and 64
    # one_shot
