# Requirement: "a cross-chain crypto asset custody ledger"

A ledger that tracks deposit receipts from external chains and allows withdrawals, with a pluggable chain-adapter for proofs.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest of input
      # cryptography
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + returns lowercase hexadecimal representation of bytes
      # encoding

ledger
  ledger.new
    fn () -> ledger_state
    + creates an empty ledger with no registered chains
    # construction
  ledger.register_chain
    fn (state: ledger_state, chain_id: string, adapter: chain_adapter) -> ledger_state
    + registers an external chain and its proof-verification adapter
    - replaces any existing adapter for the same chain_id
    # registration
  ledger.record_deposit
    fn (state: ledger_state, chain_id: string, owner: string, amount: u64, proof: bytes) -> result[ledger_state, string]
    + credits the owner's balance when the proof validates against the chain adapter
    - returns error when chain_id is not registered
    - returns error when the adapter rejects the proof
    # deposits
    -> std.crypto.sha256
  ledger.request_withdrawal
    fn (state: ledger_state, chain_id: string, owner: string, amount: u64) -> result[tuple[ledger_state, string], string]
    + debits the balance and returns (new_state, receipt_id) on success
    - returns error when the owner's balance is insufficient
    - returns error when chain_id is not registered
    # withdrawals
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  ledger.balance_of
    fn (state: ledger_state, chain_id: string, owner: string) -> u64
    + returns the owner's current balance on the given chain
    + returns 0 when the owner has no recorded balance
    # query
  ledger.pending_withdrawals
    fn (state: ledger_state, chain_id: string) -> list[withdrawal_record]
    + returns all withdrawal receipts that have not yet been settled
    # query
