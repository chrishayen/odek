# Requirement: "a decentralized bitcoin bridge to a relay chain"

A two-way peg library: bitcoin transactions are verified via SPV proofs, wrapped assets are issued on the relay chain, and collateral-backed vaults secure redemption.

std
  std.crypto
    std.crypto.sha256d
      @ (data: bytes) -> bytes
      + returns double-SHA256 of input
      + result is 32 bytes
      # cryptography
    std.crypto.verify_ecdsa
      @ (pubkey: bytes, msg: bytes, sig: bytes) -> bool
      + returns true when signature is valid for pubkey over msg
      - returns false on malformed signature
      # cryptography
  std.encoding
    std.encoding.hex_decode
      @ (s: string) -> result[bytes, string]
      + decodes lowercase or uppercase hex
      - returns error on odd length or non-hex character
      # encoding
    std.encoding.varint_decode
      @ (data: bytes, offset: i64) -> result[tuple[i64, i64], string]
      + returns (value, new_offset) for bitcoin-style compact size
      - returns error on truncated input
      # encoding
  std.merkle
    std.merkle.verify_branch
      @ (leaf: bytes, branch: list[bytes], index: i64, root: bytes) -> bool
      + returns true when the branch proves leaf is under root at index
      - returns false on mismatched root
      # merkle_proof

bridge
  bridge.parse_block_header
    @ (raw: bytes) -> result[block_header, string]
    + decodes an 80-byte bitcoin block header
    - returns error when raw is not exactly 80 bytes
    # header_parsing
    -> std.crypto.sha256d
  bridge.verify_header_chain
    @ (headers: list[block_header], start_hash: bytes) -> result[bytes, string]
    + returns the tip hash when each header links to the previous
    - returns error on broken parent linkage
    - returns error on insufficient proof-of-work
    # chain_verification
    -> std.crypto.sha256d
  bridge.verify_inclusion
    @ (tx_hash: bytes, branch: list[bytes], index: i64, block_root: bytes) -> bool
    + returns true when the transaction is in the block
    # spv_proof
    -> std.merkle.verify_branch
  bridge.parse_transaction
    @ (raw: bytes) -> result[btc_transaction, string]
    + decodes a bitcoin transaction with inputs and outputs
    - returns error on truncated data
    # transaction_parsing
    -> std.encoding.varint_decode
  bridge.open_vault
    @ (operator: string, collateral: u64) -> result[vault_id, string]
    + registers a new vault with locked collateral
    - returns error when collateral is below minimum
    # vault_management
  bridge.request_issue
    @ (vault: vault_id, amount: u64, requester: string) -> result[issue_request, string]
    + reserves vault capacity for an incoming bitcoin deposit
    - returns error when vault has insufficient free capacity
    # issue_request
  bridge.execute_issue
    @ (req: issue_request, proof: inclusion_proof, tx: btc_transaction) -> result[u64, string]
    + mints wrapped tokens to the requester when the deposit is proven
    - returns error when the proof does not match the requested amount
    - returns error when the transaction does not pay the vault address
    # issue_execution
    -> std.crypto.sha256d
  bridge.request_redeem
    @ (holder: string, amount: u64, btc_address: string) -> result[redeem_request, string]
    + burns wrapped tokens and schedules a vault payout
    - returns error when holder balance is below amount
    # redeem_request
  bridge.execute_redeem
    @ (req: redeem_request, proof: inclusion_proof, tx: btc_transaction) -> result[void, string]
    + releases collateral back to the vault when the payout is proven
    - returns error when the payout amount or address mismatches
    # redeem_execution
  bridge.liquidate_vault
    @ (vault: vault_id, reason: string) -> result[u64, string]
    + slashes vault collateral and returns the amount redistributed
    - returns error when vault is already liquidated
    # liquidation
  bridge.collateral_ratio
    @ (vault: vault_id, btc_price: u64) -> f64
    + returns the collateral-to-issued-value ratio for the vault
    ? caller supplies the price oracle value
    # collateralization
