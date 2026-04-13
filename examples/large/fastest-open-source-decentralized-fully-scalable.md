# Requirement: "a block DAG consensus library for a decentralized ledger"

Nodes append blocks that reference multiple parents, forming a DAG. A topological ordering procedure selects a main chain via ghost-like weight and finalizes transactions once a block has enough cumulative descendants.

std
  std.hash
    std.hash.sha256
      @ (data: bytes) -> bytes
      + returns a 32-byte digest
      # hashing
  std.crypto
    std.crypto.ed25519_verify
      @ (pubkey: bytes, msg: bytes, sig: bytes) -> bool
      + returns true when the signature is valid
      # cryptography
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.codec
    std.codec.encode_block
      @ (b: block) -> bytes
      + returns a canonical binary encoding of a block
      # serialization
    std.codec.decode_block
      @ (data: bytes) -> result[block, string]
      + decodes a canonical block
      - returns error on malformed data
      # serialization

block_dag
  block_dag.block_id
    @ (b: block) -> block_hash
    + returns sha256 of the canonical encoding of b
    # identity
    -> std.hash.sha256
    -> std.codec.encode_block
  block_dag.new_dag
    @ (genesis: block) -> dag_state
    + returns a DAG containing only the genesis block
    # construction
  block_dag.validate_block
    @ (state: dag_state, b: block) -> result[void, string]
    + checks signature, timestamp bounds, and parent existence
    - returns error when any parent is unknown
    - returns error when the signature is invalid
    # validation
    -> std.crypto.ed25519_verify
    -> std.time.now_seconds
  block_dag.add_block
    @ (state: dag_state, b: block) -> result[dag_state, string]
    + inserts a validated block and updates descendant counts for its parents
    - returns error when the block is already known
    # insertion
  block_dag.select_main_chain
    @ (state: dag_state) -> list[block_hash]
    + returns the chain of blocks with greatest cumulative subtree weight
    # ordering
  block_dag.finalized_blocks
    @ (state: dag_state, confirmations: i32) -> list[block_hash]
    + returns main-chain blocks with at least confirmations descendants
    # finality
  block_dag.ordered_transactions
    @ (state: dag_state) -> list[transaction]
    + returns transactions in the order implied by the main chain
    + skips transactions appearing in non-main branches
    # ordering
  block_dag.tips
    @ (state: dag_state) -> list[block_hash]
    + returns blocks that currently have no descendants
    # tips
  block_dag.serialize
    @ (state: dag_state) -> bytes
    + returns a canonical encoding of the current DAG state
    # persistence
    -> std.codec.encode_block
  block_dag.load
    @ (data: bytes) -> result[dag_state, string]
    + restores a previously serialized DAG
    - returns error on corrupt data
    # persistence
    -> std.codec.decode_block
