# Requirement: "a Merkle root and inclusion proof library"

Streaming root computation without retaining the whole tree, plus proof build and verify.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + computes the SHA-256 digest of data
      + returns 32 bytes
      # cryptography

merkle
  merkle.new_builder
    fn () -> merkle_builder
    + creates a streaming root builder holding at most O(log N) partial hashes
    ? uses a stack of per-level hashes; pushes collapse upward when paired
    # construction
  merkle.push_leaf
    fn (b: merkle_builder, leaf: bytes) -> void
    + hashes the leaf and folds it into the running stack
    # streaming
    -> std.crypto.sha256
  merkle.finalize_root
    fn (b: merkle_builder) -> bytes
    + returns the Merkle root, duplicating odd nodes as needed
    + returns a 32-byte zero digest when no leaves were pushed
    # finalization
    -> std.crypto.sha256
  merkle.build_proof
    fn (leaves: list[bytes], index: i32) -> result[list[bytes], string]
    + returns the sibling hashes needed to prove inclusion of the leaf at index
    - returns error when the index is out of range
    # proof_construction
    -> std.crypto.sha256
  merkle.verify_proof
    fn (leaf: bytes, index: i32, proof: list[bytes], root: bytes) -> bool
    + returns true when hashing the leaf with the proof rebuilds the given root
    - returns false on any mismatch
    # proof_verification
    -> std.crypto.sha256
