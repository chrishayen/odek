# Requirement: "a decentralized smart-contract platform designed for low-end devices"

The library surface is the node-internals engineer's view: state storage, transaction validation, block production, and a bytecode VM. Networking is the caller's responsibility.

std
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest
      # hashing
    std.crypto.ed25519_verify
      @ (public_key: bytes, message: bytes, signature: bytes) -> bool
      + returns true when the signature is valid
      - returns false on length mismatches
      # signatures
    std.crypto.ed25519_sign
      @ (secret_key: bytes, message: bytes) -> bytes
      + returns a 64-byte signature
      # signatures
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + returns the lowercase hex encoding
      # encoding
  std.serialize
    std.serialize.encode
      @ (value: serializable) -> bytes
      + returns a deterministic binary encoding
      # serialization
    std.serialize.decode
      @ (raw: bytes) -> result[serializable, string]
      + returns the decoded value
      - returns error on truncated input
      # serialization

chain
  chain.new_state
    @ () -> chain_state
    + returns an empty world state with zero accounts
    # construction
  chain.apply_transaction
    @ (state: chain_state, tx: transaction) -> result[chain_state, string]
    + applies a signed transaction and returns the new state
    - returns error when the signature is invalid
    - returns error when the sender balance is insufficient
    # execution
    -> std.crypto.ed25519_verify
  chain.account_balance
    @ (state: chain_state, address: bytes) -> u64
    + returns the balance, or 0 for unknown addresses
    # state
  chain.produce_block
    @ (state: chain_state, mempool: list[transaction], parent_hash: bytes, height: u64) -> result[block, string]
    + applies valid transactions and returns a block with merkle root and header
    - returns error when mempool contains only invalid transactions
    # block_production
    -> chain.apply_transaction
    -> std.crypto.sha256
  chain.validate_block
    @ (state: chain_state, blk: block) -> result[void, string]
    + re-executes transactions and verifies the merkle root matches the header
    - returns error on any header or execution mismatch
    # validation
    -> std.crypto.sha256
  chain.block_hash
    @ (blk: block) -> bytes
    + returns the SHA-256 hash of the canonical block encoding
    # hashing
    -> std.serialize.encode
    -> std.crypto.sha256
  chain.deploy_contract
    @ (state: chain_state, sender: bytes, code: bytes) -> result[tuple[chain_state, bytes], string]
    + stores bytecode under a deterministic contract address
    - returns error when sender balance cannot cover the deploy fee
    # contracts
    -> std.crypto.sha256
  chain.call_contract
    @ (state: chain_state, caller: bytes, contract: bytes, input: bytes, gas: u64) -> result[tuple[chain_state, bytes], string]
    + runs contract bytecode and returns the updated state and output
    - returns error when gas is exhausted
    - returns error when the contract does not exist
    # contracts
    -> chain.vm_execute
  chain.vm_execute
    @ (code: bytes, input: bytes, gas: u64) -> result[tuple[bytes, u64], string]
    + executes bytecode in a minimal stack VM and returns output and gas used
    - returns error on unknown opcodes
    - returns error when the stack underflows
    # vm
  chain.address_from_pubkey
    @ (public_key: bytes) -> bytes
    + returns the address as the first 20 bytes of sha256(public_key)
    # accounts
    -> std.crypto.sha256
  chain.sign_transaction
    @ (tx: transaction, secret_key: bytes) -> transaction
    + returns tx with its signature field populated
    # signing
    -> std.crypto.ed25519_sign
    -> std.serialize.encode
