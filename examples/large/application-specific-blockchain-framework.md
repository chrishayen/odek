# Requirement: "a framework for building application-specific blockchains"

A framework for building purpose-built blockchains: modules register message handlers, a state store tracks account balances and arbitrary key-value data, and a consensus-agnostic block processor sequences transactions deterministically.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest
      # cryptography
    std.crypto.verify_signature
      fn (pubkey: bytes, message: bytes, signature: bytes) -> bool
      + returns true when the signature is valid for the message under the pubkey
      - returns false on any malformed input or mismatch
      # cryptography
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding
  std.collections
    std.collections.map_keys_sorted
      fn (m: map[string,bytes]) -> list[string]
      + returns the keys of m sorted lexicographically
      ? used for deterministic iteration over state
      # collections

chain
  chain.new
    fn (chain_id: string) -> chain_state
    + creates an empty chain with height 0 and no registered modules
    # construction
  chain.register_module
    fn (state: chain_state, module_name: string, message_types: list[string]) -> chain_state
    + associates a module with the set of message type tags it owns
    - returns unchanged state if module_name is already registered
    # module_registry
  chain.account_set_balance
    fn (state: chain_state, address: string, denom: string, amount: i64) -> chain_state
    + writes balance for (address, denom) into state
    # state_store
  chain.account_get_balance
    fn (state: chain_state, address: string, denom: string) -> i64
    + returns the balance or 0 when absent
    # state_store
  chain.account_transfer
    fn (state: chain_state, from: string, to: string, denom: string, amount: i64) -> result[chain_state, string]
    + debits from, credits to, and returns the updated state
    - returns error "insufficient funds" when from's balance is less than amount
    - returns error "invalid amount" when amount is non-positive
    # accounts
  chain.store_set
    fn (state: chain_state, module_name: string, key: string, value: bytes) -> chain_state
    + writes a value under a module-scoped key
    # state_store
  chain.store_get
    fn (state: chain_state, module_name: string, key: string) -> optional[bytes]
    + returns the stored value or none
    # state_store
  chain.verify_tx
    fn (state: chain_state, tx_bytes: bytes, signer_pubkey: bytes, signature: bytes) -> result[void, string]
    + returns ok when the signature is valid over tx_bytes
    - returns error "bad signature" when verification fails
    # tx_validation
    -> std.crypto.verify_signature
  chain.deliver_tx
    fn (state: chain_state, module_name: string, message_type: string, payload: bytes) -> result[chain_state, string]
    + routes the message to the registered module and applies its effect
    - returns error "unknown module" when module_name has no registration
    - returns error "unknown message type" when message_type is not owned by the module
    # tx_execution
  chain.begin_block
    fn (state: chain_state, proposer: string, timestamp_secs: i64) -> chain_state
    + advances the internal block context to a new height
    # block_lifecycle
  chain.end_block
    fn (state: chain_state) -> chain_state
    + finalizes per-block invariants and increments height
    # block_lifecycle
  chain.app_hash
    fn (state: chain_state) -> bytes
    + computes a deterministic digest over all state at the current height
    ? keys are walked in sorted order so the hash is reproducible across nodes
    # state_commitment
    -> std.collections.map_keys_sorted
    -> std.crypto.sha256
  chain.address_to_string
    fn (addr: bytes) -> string
    + renders a binary address as a printable hex string
    # formatting
    -> std.encoding.hex_encode
