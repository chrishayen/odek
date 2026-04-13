# Requirement: "a client library for reading and writing records on a distributed ledger"

Generalized from a blockchain SDK. The library signs transactions, submits them, and queries ledger state. Cryptography primitives live in std.

std
  std.crypto
    std.crypto.sign_detached
      @ (private_key: bytes, data: bytes) -> bytes
      + produces a detached signature over the given bytes
      # cryptography
    std.crypto.keccak256
      @ (data: bytes) -> bytes
      + returns a 32-byte hash of the input
      # cryptography
  std.json
    std.json.encode
      @ (value: json_value) -> string
      + serializes a dynamic value to a compact JSON string
      # serialization
  std.net
    std.net.http_post
      @ (url: string, body: string) -> result[string, string]
      + performs an HTTP POST and returns the response body
      - returns error on any network or status failure
      # network

ledger_client
  ledger_client.new
    @ (endpoint_url: string, private_key: bytes) -> ledger_client_state
    + creates a client bound to an endpoint and a signing key
    # construction
  ledger_client.derive_address
    @ (private_key: bytes) -> string
    + returns the canonical textual address for the given key
    # addressing
    -> std.crypto.keccak256
  ledger_client.build_transaction
    @ (state: ledger_client_state, to_address: string, payload: bytes, nonce: i64) -> bytes
    + returns the canonical serialized unsigned transaction bytes
    # transaction
    -> std.json.encode
  ledger_client.sign_transaction
    @ (state: ledger_client_state, unsigned: bytes) -> bytes
    + returns the serialized signed transaction
    # transaction
    -> std.crypto.sign_detached
  ledger_client.submit_transaction
    @ (state: ledger_client_state, signed: bytes) -> result[string, string]
    + submits the transaction and returns the accepted transaction id
    - returns error when the endpoint rejects the transaction
    # submission
    -> std.net.http_post
  ledger_client.get_record
    @ (state: ledger_client_state, key: string) -> result[optional[bytes], string]
    + returns the current value stored under the key or none
    - returns error when the endpoint is unreachable
    # query
    -> std.net.http_post
  ledger_client.put_record
    @ (state: ledger_client_state, key: string, value: bytes) -> result[string, string]
    + builds, signs, and submits a write transaction and returns its id
    # write
