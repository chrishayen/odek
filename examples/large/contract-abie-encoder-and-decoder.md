# Requirement: "encode and decode smart contract invocations"

A library that turns function names and typed arguments into a packed calldata payload and decodes return data according to a declared signature.

std
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + returns the lowercase hex representation without prefix
      # encoding
    std.encoding.hex_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes hex, tolerating an optional "0x" prefix
      - returns error on non-hex characters or odd length
      # encoding
    std.encoding.encode_u256_be
      fn (value: bytes) -> bytes
      + returns a 32-byte big-endian representation, left-padding with zeros
      - returns error when value is longer than 32 bytes
      # encoding
  std.crypto
    std.crypto.keccak256
      fn (data: bytes) -> bytes
      + returns the 32-byte Keccak-256 hash of data
      # cryptography

contract_abi
  contract_abi.parse_signature
    fn (signature: string) -> result[parsed_signature, string]
    + parses strings like "transfer(address,uint256)" into name and parameter types
    - returns error on unbalanced parentheses
    - returns error on unknown parameter type keywords
    # parsing
  contract_abi.function_selector
    fn (signature: string) -> result[bytes, string]
    + returns the first 4 bytes of the Keccak-256 hash of the canonical signature
    - returns error when signature cannot be parsed
    # selectors
    -> std.crypto.keccak256
  contract_abi.encode_value
    fn (type_name: string, value: bytes) -> result[bytes, string]
    + returns the padded, type-specific encoding of a single argument
    - returns error when the value does not fit the declared type
    # encoding
    -> std.encoding.encode_u256_be
  contract_abi.encode_call
    fn (signature: string, args: list[tuple[string, bytes]]) -> result[bytes, string]
    + returns selector followed by the packed argument encoding
    - returns error when args count does not match the signature
    - returns error when any arg fails to encode
    # call_encoding
  contract_abi.encode_call_hex
    fn (signature: string, args: list[tuple[string, bytes]]) -> result[string, string]
    + returns the 0x-prefixed hex string for the packed call
    # call_encoding
    -> std.encoding.hex_encode
  contract_abi.decode_value
    fn (type_name: string, data: bytes, offset: i32) -> result[tuple[bytes, i32], string]
    + returns the decoded value and the next offset
    - returns error when data runs out before the slot is complete
    # decoding
  contract_abi.decode_return
    fn (types: list[string], data: bytes) -> result[list[bytes], string]
    + returns the decoded return values in declared order
    - returns error when data length is not a multiple of 32 for non-dynamic types
    # decoding
  contract_abi.decode_event_log
    fn (topic0_signature: string, topics: list[bytes], data: bytes) -> result[map[string, bytes], string]
    + returns indexed and non-indexed event parameters keyed by name
    - returns error when topic0 does not match the signature hash
    # events
    -> std.crypto.keccak256
