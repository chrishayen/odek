# Requirement: "a shielded-transaction cryptocurrency protocol with zero-knowledge proofs"

Implements note commitments, nullifiers, and zero-knowledge proofs for a privacy-preserving value transfer protocol.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + computes SHA-256 of data and returns 32 bytes
      # hashing
    std.crypto.blake2b
      fn (data: bytes, personalization: bytes) -> bytes
      + computes BLAKE2b with a personalization tag
      # hashing
    std.crypto.ed25519_sign
      fn (key: bytes, message: bytes) -> bytes
      + signs a message with an Ed25519 private key
      # signatures
    std.crypto.ed25519_verify
      fn (public_key: bytes, message: bytes, signature: bytes) -> bool
      + verifies an Ed25519 signature
      # signatures
  std.random
    std.random.bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
  std.math
    std.math.mod_exp
      fn (base: bytes, exponent: bytes, modulus: bytes) -> bytes
      + computes modular exponentiation on big-integer byte encodings
      # arithmetic

shielded
  shielded.generate_spending_key
    fn () -> spending_key
    + creates a new random spending key
    # keys
    -> std.random.bytes
  shielded.derive_viewing_key
    fn (spending: spending_key) -> viewing_key
    + derives a read-only viewing key from a spending key
    # keys
    -> std.crypto.blake2b
  shielded.derive_address
    fn (viewing: viewing_key) -> shielded_address
    + derives a payment address from a viewing key
    # keys
    -> std.crypto.blake2b
  shielded.create_note
    fn (address: shielded_address, value: u64) -> note
    + creates a note assigning value to the address with a fresh randomness seed
    # notes
    -> std.random.bytes
  shielded.note_commitment
    fn (n: note) -> bytes
    + returns the 32-byte commitment binding the note's address, value, and randomness
    # notes
    -> std.crypto.blake2b
  shielded.note_nullifier
    fn (n: note, spending: spending_key) -> bytes
    + returns the nullifier that will be revealed when the note is spent
    # notes
    -> std.crypto.blake2b
  shielded.build_shielded_tx
    fn (inputs: list[note], outputs: list[note], spending: spending_key) -> result[shielded_tx, string]
    + builds a transaction consuming input notes and producing output notes
    - returns error when input and output values do not balance
    # transactions
    -> shielded.note_commitment
    -> shielded.note_nullifier
  shielded.prove_transfer
    fn (tx: shielded_tx, spending: spending_key) -> result[zk_proof, string]
    + generates a zero-knowledge proof that the transaction is valid without revealing notes
    - returns error when witnesses are inconsistent
    # proofs
    -> std.math.mod_exp
    -> std.crypto.blake2b
  shielded.verify_transfer
    fn (tx: shielded_tx, proof: zk_proof) -> result[bool, string]
    + verifies a zero-knowledge proof against public commitments and nullifiers
    - returns false when the proof does not satisfy the circuit
    # verification
    -> std.math.mod_exp
    -> std.crypto.sha256
  shielded.sign_tx
    fn (tx: shielded_tx, spending: spending_key) -> bytes
    + produces a binding signature committing to the transaction
    # signatures
    -> std.crypto.ed25519_sign
  shielded.verify_tx_signature
    fn (tx: shielded_tx, signature: bytes, public_key: bytes) -> bool
    + verifies the binding signature on a transaction
    # signatures
    -> std.crypto.ed25519_verify
  shielded.apply_to_state
    fn (state: chain_state, tx: shielded_tx, proof: zk_proof) -> result[chain_state, string]
    + inserts commitments and nullifiers into chain state after verifying the proof
    - returns error when a nullifier has already been seen (double spend)
    # state
    -> shielded.verify_transfer
