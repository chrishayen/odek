# Requirement: "a Bitcoin protocol library"

Key management, addresses, transaction construction and signing, script evaluation, and block header validation.

std
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte sha256 digest
      # cryptography
    std.crypto.ripemd160
      @ (data: bytes) -> bytes
      + returns the 20-byte ripemd160 digest
      # cryptography
    std.crypto.secp256k1_pubkey_from_priv
      @ (priv: bytes) -> result[bytes, string]
      + returns the compressed public key bytes
      - returns error when the private key is zero or out of range
      # cryptography
    std.crypto.secp256k1_sign
      @ (priv: bytes, msg_hash: bytes) -> result[bytes, string]
      + returns a der-encoded signature
      # cryptography
    std.crypto.secp256k1_verify
      @ (pub: bytes, msg_hash: bytes, sig: bytes) -> bool
      + returns true when the signature is valid for the public key and hash
      # cryptography
  std.encoding
    std.encoding.base58check_encode
      @ (version: u8, payload: bytes) -> string
      + returns the base58check encoded string
      # encoding
    std.encoding.base58check_decode
      @ (s: string) -> result[tuple[u8, bytes], string]
      + returns the version byte and payload
      - returns error when the checksum is invalid
      # encoding
    std.encoding.bech32_encode
      @ (hrp: string, data: bytes) -> string
      + returns the bech32 string for the given human-readable part
      # encoding
    std.encoding.bech32_decode
      @ (s: string) -> result[tuple[string, bytes], string]
      + returns the human-readable part and payload
      - returns error when the checksum fails
      # encoding

bitcoin
  bitcoin.double_sha256
    @ (data: bytes) -> bytes
    + returns sha256(sha256(data))
    # hashing
    -> std.crypto.sha256
  bitcoin.hash160
    @ (data: bytes) -> bytes
    + returns ripemd160(sha256(data))
    # hashing
    -> std.crypto.sha256
    -> std.crypto.ripemd160
  bitcoin.generate_key
    @ (seed: bytes) -> result[key_pair, string]
    + returns a deterministic key pair from the seed
    - returns error when the seed has fewer than 32 bytes
    # keys
    -> std.crypto.secp256k1_pubkey_from_priv
  bitcoin.p2pkh_address
    @ (pubkey: bytes, network: string) -> string
    + returns the pay-to-public-key-hash address
    # addresses
    -> std.encoding.base58check_encode
  bitcoin.p2wpkh_address
    @ (pubkey: bytes, network: string) -> string
    + returns the native segwit v0 address
    # addresses
    -> std.encoding.bech32_encode
  bitcoin.decode_address
    @ (address: string) -> result[decoded_address, string]
    + returns the script type and payload
    - returns error when the address is malformed
    # addresses
    -> std.encoding.base58check_decode
    -> std.encoding.bech32_decode
  bitcoin.build_transaction
    @ (inputs: list[tx_input], outputs: list[tx_output], locktime: u32) -> transaction
    + returns an unsigned transaction with the given inputs and outputs
    # transactions
  bitcoin.serialize_tx
    @ (tx: transaction, include_witness: bool) -> bytes
    + returns the transaction in its wire format
    # transactions
  bitcoin.parse_tx
    @ (raw: bytes) -> result[transaction, string]
    + returns the parsed transaction
    - returns error on truncated input
    # transactions
  bitcoin.sighash
    @ (tx: transaction, input_index: i32, script_pubkey: bytes, value: i64, sighash_type: u8) -> bytes
    + returns the message hash to be signed for the given input
    # transactions
  bitcoin.sign_input
    @ (tx: transaction, input_index: i32, priv: bytes, script_pubkey: bytes, value: i64, sighash_type: u8) -> result[transaction, string]
    + returns the transaction with the signed input script or witness
    - returns error when the script type is unsupported
    # signing
    -> std.crypto.secp256k1_sign
  bitcoin.verify_input
    @ (tx: transaction, input_index: i32, script_pubkey: bytes, value: i64) -> bool
    + returns true when the input unlocks its referenced output
    # verification
    -> std.crypto.secp256k1_verify
  bitcoin.evaluate_script
    @ (sig_script: bytes, pub_script: bytes, tx: transaction, input_index: i32, value: i64) -> result[bool, string]
    + returns the final boolean result of script execution
    - returns error on invalid opcodes
    # script
  bitcoin.block_header_hash
    @ (header: block_header) -> bytes
    + returns the double-sha256 of the serialized header
    # blocks
  bitcoin.validate_block_header
    @ (header: block_header, expected_target: bytes) -> result[void, string]
    + returns ok when the header hash meets the target
    - returns error when the hash exceeds the target
    # blocks
  bitcoin.txid
    @ (tx: transaction) -> bytes
    + returns the transaction id (legacy, non-witness)
    # transactions
