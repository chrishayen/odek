# Requirement: "a script tree interpreter and wallet for a blockchain virtual machine"

A stack-machine script interpreter paired with a minimal wallet that can build, sign, and verify transactions. Cryptographic primitives live in std so the project layer stays focused on script evaluation and wallet bookkeeping.

std
  std.crypto
    std.crypto.blake2b_256
      fn (data: bytes) -> bytes
      + returns a 32-byte blake2b digest
      # cryptography
    std.crypto.schnorr_sign
      fn (private_key: bytes, message: bytes) -> bytes
      + produces a schnorr signature over message under the given key
      # cryptography
    std.crypto.schnorr_verify
      fn (public_key: bytes, message: bytes, signature: bytes) -> bool
      + returns true when the signature is valid
      - returns false for tampered message or signature
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
  std.encoding
    std.encoding.base58_encode
      fn (data: bytes) -> string
      + encodes bytes to base58 with no padding
      # encoding
    std.encoding.base58_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes base58
      - returns error on non-alphabet characters
      # encoding

scripttree
  scripttree.opcode
    fn (tag: string, immediate: optional[bytes]) -> opcode
    + constructs an opcode with optional immediate payload
    # instruction
  scripttree.compile
    fn (source: string) -> result[list[opcode], string]
    + parses textual script source into an opcode list
    - returns error on unknown mnemonic or malformed literal
    # compilation
  scripttree.evaluate
    fn (program: list[opcode], context: eval_context) -> result[bool, string]
    + runs the program and returns the top-of-stack boolean
    - returns error on stack underflow or type mismatch
    # evaluation
  scripttree.new_context
    fn (height: i64, inputs: list[bytes], outputs: list[bytes]) -> eval_context
    + builds an evaluation context exposing transaction data to scripts
    # context
  wallet.new
    fn (seed: bytes) -> wallet_state
    + derives a wallet from a seed
    ? key derivation is deterministic; callers reuse a seed to restore state
    # construction
    -> std.crypto.blake2b_256
  wallet.address
    fn (state: wallet_state) -> string
    + returns the encoded address for the wallet's default key
    -> std.encoding.base58_encode
    # addressing
  wallet.build_transaction
    fn (state: wallet_state, recipient: string, amount: i64, fee: i64) -> result[transaction, string]
    + builds an unsigned transaction spending owned outputs
    - returns error when funds are insufficient
    # transaction_building
  wallet.sign_transaction
    fn (state: wallet_state, tx: transaction) -> result[transaction, string]
    + returns the transaction with signatures attached to each input
    - returns error when a required key is not held
    # signing
    -> std.crypto.schnorr_sign
    -> std.crypto.blake2b_256
  wallet.verify_transaction
    fn (tx: transaction) -> bool
    + returns true when every input signature is valid and scripts evaluate to true
    - returns false when any signature fails or script rejects
    # verification
    -> std.crypto.schnorr_verify
  wallet.add_output
    fn (state: wallet_state, output_id: bytes, value: i64, script: list[opcode]) -> wallet_state
    + records a new unspent output owned by the wallet
    # bookkeeping
