# Requirement: "a decentralized blob storage library with content addressing, erasure coding, and peer replication"

Stores blobs across a pluggable set of peer nodes with erasure-coded redundancy and cryptographic addressing.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns 32 bytes of SHA-256 digest
      # cryptography
    std.crypto.ed25519_sign
      fn (key: bytes, data: bytes) -> bytes
      + returns a 64-byte Ed25519 signature
      # cryptography
    std.crypto.ed25519_verify
      fn (pubkey: bytes, data: bytes, sig: bytes) -> bool
      + returns true when the signature is valid
      # cryptography
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + returns lowercase hex representation
      # encoding
    std.encoding.hex_decode
      fn (s: string) -> result[bytes, string]
      + returns decoded bytes
      - returns error on non-hex characters
      # encoding
  std.math
    std.math.gf256_mul
      fn (a: u8, b: u8) -> u8
      + returns Galois-field multiplication of a and b in GF(256)
      # math

dstore
  dstore.new
    fn (replication_k: i32, replication_m: i32) -> store_state
    + returns a store configured for (k, m) Reed-Solomon erasure coding
    - returns error-marker state when k + m exceeds 255 or k < 1
    # construction
  dstore.add_peer
    fn (state: store_state, peer_id: string, pubkey: bytes) -> store_state
    + registers a peer by id and public key
    # peers
  dstore.remove_peer
    fn (state: store_state, peer_id: string) -> store_state
    + removes a peer from the active set
    # peers
  dstore.content_id
    fn (data: bytes) -> string
    + returns the hex-encoded SHA-256 of data
    # addressing
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  dstore.encode_shards
    fn (state: store_state, data: bytes) -> list[bytes]
    + returns k+m shards under the configured Reed-Solomon parameters
    # erasure_coding
    -> std.math.gf256_mul
  dstore.decode_shards
    fn (state: store_state, shards: list[optional[bytes]]) -> result[bytes, string]
    + returns reconstructed data when at least k shards are present
    - returns error when fewer than k shards are available
    # erasure_coding
    -> std.math.gf256_mul
  dstore.assign_shards
    fn (state: store_state, shards: list[bytes]) -> map[string, bytes]
    + returns a peer-id to shard map using rendezvous hashing
    # placement
    -> std.crypto.sha256
  dstore.put
    fn (state: store_state, data: bytes) -> result[string, string]
    + returns the content id after encoding and staging shards for peers
    - returns error when there are fewer than k+m peers
    # operations
    -> dstore.content_id
    -> dstore.encode_shards
    -> dstore.assign_shards
  dstore.get
    fn (state: store_state, content_id: string, shards: map[string, optional[bytes]]) -> result[bytes, string]
    + returns the original bytes after verifying the hash of the decoded result
    - returns error when decoding fails or the hash does not match
    # operations
    -> dstore.decode_shards
    -> dstore.content_id
  dstore.sign_manifest
    fn (privkey: bytes, content_id: string, shard_map: map[string, bytes]) -> bytes
    + returns an Ed25519 signature over the canonical manifest bytes
    # manifests
    -> std.crypto.ed25519_sign
  dstore.verify_manifest
    fn (pubkey: bytes, content_id: string, shard_map: map[string, bytes], sig: bytes) -> bool
    + returns true when the manifest signature matches
    # manifests
    -> std.crypto.ed25519_verify
