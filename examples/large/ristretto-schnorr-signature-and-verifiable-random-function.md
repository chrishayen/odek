# Requirement: "schnorr signatures and verifiable random functions on the ristretto group"

Schnorr signatures and VRFs built from ristretto group primitives. The project layer is a thin policy layer; all real cryptography lives in std.

std
  std.rand
    std.rand.bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
  std.hash
    std.hash.sha512
      fn (data: bytes) -> bytes
      + returns the 64-byte sha512 digest of data
      # hashing
  std.crypto
    std.crypto.ristretto_generator
      fn () -> ristretto_point
      + returns the canonical ristretto base point
      # group
    std.crypto.ristretto_scalar_from_bytes
      fn (raw: bytes) -> result[ristretto_scalar, string]
      + reduces 64 bytes to a ristretto scalar
      - returns error on wrong-length input
      # group
    std.crypto.ristretto_scalar_mul_base
      fn (s: ristretto_scalar) -> ristretto_point
      + returns s * G where G is the base point
      # group
    std.crypto.ristretto_scalar_mul
      fn (p: ristretto_point, s: ristretto_scalar) -> ristretto_point
      + returns s * p
      # group
    std.crypto.ristretto_point_add
      fn (a: ristretto_point, b: ristretto_point) -> ristretto_point
      + returns the group sum of two points
      # group
    std.crypto.ristretto_point_encode
      fn (p: ristretto_point) -> bytes
      + returns the 32-byte canonical encoding of a point
      # group
    std.crypto.ristretto_point_decode
      fn (raw: bytes) -> result[ristretto_point, string]
      + parses a 32-byte encoding into a point
      - returns error when the encoding is not canonical
      # group
    std.crypto.ristretto_scalar_add
      fn (a: ristretto_scalar, b: ristretto_scalar) -> ristretto_scalar
      + returns a + b mod group order
      # group
    std.crypto.ristretto_scalar_mul_scalar
      fn (a: ristretto_scalar, b: ristretto_scalar) -> ristretto_scalar
      + returns a * b mod group order
      # group
    std.crypto.hash_to_ristretto
      fn (data: bytes) -> ristretto_point
      + maps arbitrary bytes to a ristretto point using a ro hash-to-curve
      # group

schnorr
  schnorr.keygen
    fn () -> tuple[bytes, bytes]
    + returns (secret_key, public_key) where the public key is secret * G
    # key_generation
    -> std.rand.bytes
    -> std.crypto.ristretto_scalar_from_bytes
    -> std.crypto.ristretto_scalar_mul_base
    -> std.crypto.ristretto_point_encode
  schnorr.sign
    fn (secret_key: bytes, message: bytes) -> result[bytes, string]
    + returns a 64-byte schnorr signature
    - returns error when the secret key is not 32 bytes
    # signing
    -> std.rand.bytes
    -> std.hash.sha512
    -> std.crypto.ristretto_scalar_from_bytes
    -> std.crypto.ristretto_scalar_mul_base
    -> std.crypto.ristretto_scalar_add
    -> std.crypto.ristretto_scalar_mul_scalar
  schnorr.verify
    fn (public_key: bytes, message: bytes, signature: bytes) -> result[bool, string]
    + returns true when the signature verifies against the public key and message
    - returns error on malformed public key or signature
    # verification
    -> std.crypto.ristretto_point_decode
    -> std.crypto.ristretto_scalar_mul_base
    -> std.crypto.ristretto_scalar_mul
    -> std.crypto.ristretto_point_add
    -> std.hash.sha512
  schnorr.vrf_prove
    fn (secret_key: bytes, input: bytes) -> result[tuple[bytes, bytes], string]
    + returns (output, proof): a 32-byte pseudo-random output and a proof that can be verified
    - returns error when the secret key is invalid
    # vrf
    -> std.crypto.hash_to_ristretto
    -> std.crypto.ristretto_scalar_mul
    -> std.crypto.ristretto_point_encode
    -> std.hash.sha512
  schnorr.vrf_verify
    fn (public_key: bytes, input: bytes, output: bytes, proof: bytes) -> result[bool, string]
    + returns true when the proof shows output was derived from input under public_key
    - returns error when public_key or proof is malformed
    # vrf
    -> std.crypto.ristretto_point_decode
    -> std.crypto.hash_to_ristretto
    -> std.crypto.ristretto_scalar_mul
    -> std.hash.sha512
