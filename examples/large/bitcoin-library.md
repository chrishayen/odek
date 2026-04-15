# Requirement: "a bitcoin library"

Key handling, address derivation, transaction construction, and signing. The cryptographic primitives live in std so the project layer stays focused on bitcoin semantics.

std
  std.hash
    std.hash.sha256
      fn (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest
      # hashing
    std.hash.ripemd160
      fn (data: bytes) -> bytes
      + returns the 20-byte RIPEMD-160 digest
      # hashing
  std.crypto
    std.crypto.secp256k1_derive_pubkey
      fn (private_key: bytes) -> result[bytes, string]
      + returns the 33-byte compressed public key for a secp256k1 private key
      - returns error when the private key is zero or out of range
      # elliptic_curve
    std.crypto.secp256k1_sign
      fn (private_key: bytes, digest: bytes) -> result[bytes, string]
      + returns a DER-encoded ECDSA signature with low-S normalization
      - returns error on invalid private key
      # signing
    std.crypto.secp256k1_verify
      fn (public_key: bytes, digest: bytes, signature: bytes) -> bool
      + returns true when the signature is valid for the key and digest
      # verification
  std.encoding
    std.encoding.base58_encode
      fn (data: bytes) -> string
      + encodes bytes using the base58 alphabet
      # encoding
    std.encoding.base58_decode
      fn (s: string) -> result[bytes, string]
      + decodes base58 into bytes
      - returns error on invalid characters
      # encoding
    std.encoding.bech32_encode
      fn (hrp: string, data: list[u8]) -> result[string, string]
      + encodes a bech32 string with the given human-readable prefix
      - returns error when data contains values above 31
      # encoding
    std.encoding.bech32_decode
      fn (s: string) -> result[tuple[string, list[u8]], string]
      + decodes a bech32 string into prefix and 5-bit groups
      - returns error when the checksum does not verify
      # encoding

btc
  btc.private_key_from_wif
    fn (wif: string) -> result[bytes, string]
    + decodes a WIF-encoded private key
    - returns error when the checksum does not match
    # key_import
    -> std.encoding.base58_decode
    -> std.hash.sha256
  btc.private_key_to_wif
    fn (private_key: bytes, compressed: bool) -> string
    + encodes a private key in WIF format with the optional compressed flag
    # key_export
    -> std.encoding.base58_encode
    -> std.hash.sha256
  btc.public_key_from_private
    fn (private_key: bytes) -> result[bytes, string]
    + returns the compressed public key for the given private key
    # key_derivation
    -> std.crypto.secp256k1_derive_pubkey
  btc.p2pkh_address
    fn (public_key: bytes, mainnet: bool) -> string
    + returns a legacy P2PKH address for the public key
    # addressing
    -> std.hash.sha256
    -> std.hash.ripemd160
    -> std.encoding.base58_encode
  btc.p2wpkh_address
    fn (public_key: bytes, mainnet: bool) -> result[string, string]
    + returns a native segwit bech32 address
    - returns error when bech32 encoding fails
    # addressing
    -> std.hash.sha256
    -> std.hash.ripemd160
    -> std.encoding.bech32_encode
  btc.parse_address
    fn (address: string) -> result[parsed_address, string]
    + parses a base58 or bech32 address and returns its script type and payload
    - returns error when the address format is unrecognized
    # addressing
    -> std.encoding.base58_decode
    -> std.encoding.bech32_decode
  btc.build_transaction
    fn (inputs: list[tx_input], outputs: list[tx_output], locktime: u32) -> transaction
    + constructs an unsigned transaction with the given inputs and outputs
    # transaction_building
  btc.sighash
    fn (tx: transaction, input_index: i32, script_code: bytes, sighash_type: u32) -> bytes
    + returns the 32-byte digest to be signed for the given input
    # signing
    -> std.hash.sha256
  btc.sign_input
    fn (tx: transaction, input_index: i32, private_key: bytes, script_code: bytes) -> result[transaction, string]
    + returns the transaction with a signature script populated for the input
    - returns error when the private key cannot sign
    # signing
    -> std.crypto.secp256k1_sign
  btc.serialize_transaction
    fn (tx: transaction) -> bytes
    + returns the canonical serialized transaction bytes
    # serialization
  btc.txid
    fn (tx: transaction) -> bytes
    + returns the double-SHA-256 transaction id
    # serialization
    -> std.hash.sha256
