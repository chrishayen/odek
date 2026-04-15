# Requirement: "a vanity address generator for blockchain wallets"

Generates keypairs until the derived address matches a user pattern. The cryptographic primitives live in std; the project layer is the search loop and pattern check.

std
  std.rand
    std.rand.bytes
      fn (n: u32) -> bytes
      + returns n cryptographically-secure random bytes
      # randomness
  std.crypto
    std.crypto.secp256k1_derive_public
      fn (private_key: bytes) -> result[bytes, string]
      + derives the uncompressed public key from a 32-byte private key
      - returns error when private_key is not 32 bytes or is zero
      # cryptography
    std.crypto.keccak256
      fn (data: bytes) -> bytes
      + computes the Keccak-256 digest (32 bytes)
      # cryptography
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes as lowercase hex without a prefix
      # encoding

vanity
  vanity.derive_address
    fn (private_key: bytes) -> result[string, string]
    + returns the 40-character hex address for a private key
    - returns error when the private key is invalid
    # derivation
    -> std.crypto.secp256k1_derive_public
    -> std.crypto.keccak256
    -> std.encoding.hex_encode
  vanity.matches_pattern
    fn (address: string, prefix: string, suffix: string) -> bool
    + returns true when the address starts with prefix and ends with suffix
    + empty prefix or suffix matches anything in that position
    # matching
  vanity.find
    fn (prefix: string, suffix: string, max_attempts: u64) -> result[tuple[bytes, string], string]
    + generates random keypairs until one matches, returning (private_key, address)
    - returns error when max_attempts is exhausted without a match
    - returns error when prefix or suffix contains non-hex characters
    # search
    -> std.rand.bytes
