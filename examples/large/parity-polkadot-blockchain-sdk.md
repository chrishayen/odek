# Requirement: "a blockchain SDK for building sovereign chains with pluggable consensus"

A core blockchain runtime: blocks and transactions, a state trie with cryptographic root, a pluggable consensus interface, and a peer-to-peer gossip layer. std supplies hashing, signing, and a clock.

std
  std.hash
    std.hash.blake2_256
      @ (data: bytes) -> bytes
      + returns the 32-byte BLAKE2b-256 digest
      # cryptography
    std.hash.keccak_256
      @ (data: bytes) -> bytes
      + returns the 32-byte Keccak-256 digest
      # cryptography
  std.crypto
    std.crypto.ed25519_sign
      @ (private_key: bytes, message: bytes) -> bytes
      + returns the 64-byte Ed25519 signature
      # cryptography
    std.crypto.ed25519_verify
      @ (public_key: bytes, message: bytes, signature: bytes) -> bool
      + returns true when the signature matches the message under the public key
      - returns false on any mismatch
      # cryptography
  std.encoding
    std.encoding.scale_encode_u64
      @ (value: u64) -> bytes
      + returns the compact SCALE encoding of the given unsigned integer
      # encoding
    std.encoding.scale_decode_u64
      @ (data: bytes, offset: i64) -> result[tuple[u64, i64], string]
      + returns (value, next_offset) after decoding a compact u64
      - returns error when the encoding is truncated or malformed
      # encoding
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

chain
  chain.new_account
    @ (public_key: bytes) -> account_id
    + returns an account identifier derived from the public key
    # identity
    -> std.hash.blake2_256
  chain.new_transaction
    @ (from: account_id, to: account_id, amount: u64, nonce: u64) -> transaction
    + returns an unsigned transaction with the given fields
    # transactions
  chain.sign_transaction
    @ (tx: transaction, private_key: bytes) -> signed_transaction
    + returns the transaction with an attached Ed25519 signature
    # transactions
    -> std.crypto.ed25519_sign
    -> std.encoding.scale_encode_u64
  chain.verify_transaction
    @ (signed: signed_transaction, public_key: bytes) -> bool
    + returns true when the attached signature is valid for the transaction body
    - returns false on any mismatch
    # transactions
    -> std.crypto.ed25519_verify
  chain.new_state
    @ () -> chain_state
    + returns an empty state with zero accounts
    # state
  chain.apply_transaction
    @ (state: chain_state, signed: signed_transaction) -> result[chain_state, string]
    + returns an updated state after transferring the amount
    - returns error when the sender balance is insufficient
    - returns error when the transaction nonce does not match the sender
    # state
  chain.state_root
    @ (state: chain_state) -> bytes
    + returns the 32-byte cryptographic root over all accounts
    # state
    -> std.hash.blake2_256
  chain.new_block
    @ (parent_hash: bytes, height: u64, transactions: list[signed_transaction]) -> block
    + returns an unsealed block with the given parent, height, and payload
    # blocks
    -> std.time.now_millis
  chain.block_hash
    @ (block: block) -> bytes
    + returns the hash of the block header
    # blocks
    -> std.hash.blake2_256
  chain.apply_block
    @ (state: chain_state, block: block) -> result[chain_state, string]
    + returns a state with every transaction in the block applied in order
    - returns error when any transaction fails to apply
    # execution
  chain.new_consensus
    @ (propose_fn: fn(chain_state, i64) -> block, validate_fn: fn(block) -> bool) -> consensus_engine
    + returns a consensus engine with pluggable propose and validate hooks
    # consensus
  chain.propose_block
    @ (engine: consensus_engine, state: chain_state) -> block
    + returns the block produced by the engine's propose hook
    # consensus
    -> std.time.now_millis
  chain.validate_block
    @ (engine: consensus_engine, block: block) -> bool
    + returns the engine's validate hook result
    # consensus
  chain.new_network
    @ () -> network_state
    + returns a network with no peers
    # networking
  chain.add_peer
    @ (net: network_state, peer_id: string, address: string) -> network_state
    + returns a network with the peer registered
    # networking
  chain.gossip_block
    @ (net: network_state, block: block) -> list[string]
    + returns the peer ids the block was forwarded to
    # networking
  chain.gossip_transaction
    @ (net: network_state, tx: signed_transaction) -> list[string]
    + returns the peer ids the transaction was forwarded to
    # networking
