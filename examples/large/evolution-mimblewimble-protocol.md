# Requirement: "a privacy-preserving blockchain protocol library"

Implements confidential transaction aggregation: pedersen commitments, range proofs, kernel signatures, and cut-through aggregation. Cryptographic primitives are general and live in std.

std
  std.crypto
    std.crypto.random_scalar
      @ () -> bytes
      + returns a random 32-byte scalar suitable for elliptic curve operations
      # cryptography
    std.crypto.curve_add
      @ (a: bytes, b: bytes) -> bytes
      + returns the point sum of two curve points
      # cryptography
    std.crypto.curve_scalar_mul
      @ (scalar: bytes, point: bytes) -> bytes
      + returns scalar * point on the curve
      # cryptography
    std.crypto.curve_generator
      @ () -> bytes
      + returns the canonical generator point
      # cryptography
    std.crypto.blake2b_256
      @ (data: bytes) -> bytes
      + returns a 32-byte blake2b digest
      # cryptography
    std.crypto.schnorr_sign
      @ (private_key: bytes, message: bytes) -> bytes
      + produces a schnorr signature over message under the given key
      # cryptography
    std.crypto.schnorr_verify
      @ (public_key: bytes, message: bytes, signature: bytes) -> bool
      + returns true when the signature is valid
      - returns false for tampered message or signature
      # cryptography

privchain
  privchain.commit
    @ (value: i64, blinding: bytes) -> bytes
    + returns a pedersen commitment to value under the given blinding factor
    -> std.crypto.curve_generator
    -> std.crypto.curve_scalar_mul
    -> std.crypto.curve_add
    # commitment
  privchain.commit_sum
    @ (commitments: list[bytes]) -> bytes
    + returns the curve-point sum of commitments
    -> std.crypto.curve_add
    # commitment
  privchain.new_blinding
    @ () -> bytes
    + returns a random blinding factor
    -> std.crypto.random_scalar
    # commitment
  privchain.range_proof
    @ (value: i64, blinding: bytes) -> bytes
    + builds a proof that a committed value is non-negative and bounded
    -> std.crypto.blake2b_256
    -> std.crypto.random_scalar
    # range_proof
  privchain.verify_range_proof
    @ (commitment: bytes, proof: bytes) -> bool
    + returns true when the range proof binds to the commitment
    - returns false on tampered proof or mismatched commitment
    -> std.crypto.blake2b_256
    # range_proof
  privchain.build_kernel
    @ (excess: bytes, fee: i64, signature: bytes) -> kernel
    + assembles a transaction kernel carrying the excess, fee, and schnorr signature
    # kernel
  privchain.sign_kernel
    @ (excess_key: bytes, fee: i64) -> bytes
    + signs the fee with the excess key producing a kernel signature
    -> std.crypto.schnorr_sign
    -> std.crypto.blake2b_256
    # kernel
  privchain.verify_kernel
    @ (k: kernel) -> bool
    + returns true when signature, excess, and fee are consistent
    -> std.crypto.schnorr_verify
    # kernel
  privchain.build_transaction
    @ (inputs: list[bytes], outputs: list[bytes], kernels: list[kernel]) -> transaction
    + assembles a transaction from input commitments, output commitments, and kernels
    # transaction
  privchain.verify_transaction
    @ (tx: transaction) -> bool
    + returns true when commitments balance to excess and all kernels verify
    - returns false when sums differ or any kernel fails
    -> std.crypto.curve_add
    # verification
  privchain.aggregate
    @ (txs: list[transaction]) -> transaction
    + merges transactions and applies cut-through on matching input/output commitments
    ? ordering of inputs, outputs, and kernels within the aggregate is canonical
    # aggregation
