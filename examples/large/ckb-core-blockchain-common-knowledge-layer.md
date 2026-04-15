# Requirement: "a public permissionless blockchain acting as a common knowledge layer for a broader network"

Minimal but complete chain core: cell-based UTXO model, proof-of-work consensus, transaction validation with scripts, chain state, and a gossip-ready mempool.

std
  std.crypto
    std.crypto.blake2b_256
      fn (data: bytes) -> bytes
      + returns the 32-byte BLAKE2b-256 digest of data
      # hashing
    std.crypto.ecdsa_verify
      fn (public_key: bytes, message: bytes, signature: bytes) -> bool
      + returns true when the signature is valid for (public_key, message)
      # cryptography
  std.encoding
    std.encoding.molecule_encode
      fn (value: schema_value) -> bytes
      + serializes a typed value using a length-prefixed schema layout
      # encoding
    std.encoding.molecule_decode
      fn (schema: schema_id, data: bytes) -> result[schema_value, string]
      + deserializes data against the schema
      - returns error on length mismatches or unknown fields
      # encoding
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

ckb_core
  ckb_core.cell_new
    fn (capacity: u64, lock_hash: bytes, type_hash: optional[bytes], data: bytes) -> cell
    + constructs an unspent cell with the given lock and optional type script
    # model
  ckb_core.cell_id
    fn (cell: cell) -> bytes
    + returns the BLAKE2b-256 hash of the canonical cell encoding
    # model
    -> std.crypto.blake2b_256
    -> std.encoding.molecule_encode
  ckb_core.transaction_new
    fn (inputs: list[out_point], outputs: list[cell], witnesses: list[bytes]) -> transaction
    + assembles a transaction from input references, output cells, and witnesses
    # model
  ckb_core.transaction_hash
    fn (tx: transaction) -> bytes
    + returns the hash of the transaction excluding witnesses
    # model
    -> std.crypto.blake2b_256
    -> std.encoding.molecule_encode
  ckb_core.verify_scripts
    fn (tx: transaction, resolved_inputs: list[cell]) -> result[void, string]
    + runs lock scripts and type scripts against their cells and witnesses
    - returns error when any script rejects the transaction
    # validation
    -> std.crypto.ecdsa_verify
  ckb_core.verify_capacity
    fn (tx: transaction, resolved_inputs: list[cell]) -> result[void, string]
    + checks that output capacity does not exceed input capacity
    - returns error on overspend
    - returns error when any output is below the minimum capacity for its data length
    # validation
  ckb_core.block_new
    fn (parent_hash: bytes, number: u64, txs: list[transaction], nonce: u64) -> block
    + constructs a block with the given header fields and transactions
    # model
    -> std.time.now_millis
  ckb_core.block_hash
    fn (b: block) -> bytes
    + returns the hash of the block header
    # model
    -> std.crypto.blake2b_256
    -> std.encoding.molecule_encode
  ckb_core.pow_target_from_difficulty
    fn (difficulty: u64) -> bytes
    + returns the 32-byte target threshold that corresponds to the given difficulty
    # consensus
  ckb_core.pow_verify
    fn (header_hash: bytes, target: bytes) -> bool
    + returns true when header_hash is lexicographically below target
    # consensus
  ckb_core.chain_new
    fn (genesis: block) -> chain_state
    + initializes a chain state seeded with the genesis block and its outputs as the UTXO set
    # state
  ckb_core.chain_apply_block
    fn (state: chain_state, b: block) -> result[chain_state, string]
    + validates and appends a block, updating the UTXO set
    - returns error when parent_hash does not match the tip
    - returns error when any transaction is invalid
    - returns error when the proof-of-work is below target
    # state
  ckb_core.mempool_new
    fn () -> mempool_state
    + creates an empty mempool
    # mempool
  ckb_core.mempool_submit
    fn (mempool: mempool_state, state: chain_state, tx: transaction) -> result[mempool_state, string]
    + validates tx against the current UTXO set and inserts it
    - returns error when inputs are missing or validation fails
    # mempool
  ckb_core.mempool_take
    fn (mempool: mempool_state, max: i32) -> tuple[list[transaction], mempool_state]
    + removes and returns up to max transactions in descending fee order
    # mempool
