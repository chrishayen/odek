# Requirement: "a confidential smart contract execution layer inside a trusted enclave"

Runs deterministic contract code inside an attested enclave, seals state so only the enclave can read it, and produces signed attestations over each state transition for an outer consensus layer.

std
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte digest
      # cryptography
    std.crypto.sign_ed25519
      @ (private_key: bytes, message: bytes) -> bytes
      + returns a 64-byte signature
      # cryptography
    std.crypto.verify_ed25519
      @ (public_key: bytes, message: bytes, sig: bytes) -> bool
      + returns true only for valid signatures
      # cryptography
    std.crypto.encrypt_aead
      @ (key: bytes, plaintext: bytes, nonce: bytes) -> bytes
      + authenticated encryption
      # cryptography
    std.crypto.decrypt_aead
      @ (key: bytes, ciphertext: bytes, nonce: bytes) -> result[bytes, string]
      + returns plaintext on valid tag
      - returns error on tampering
      # cryptography
  std.enclave
    std.enclave.measure_self
      @ () -> bytes
      + returns the hash of the currently running enclave code
      # attestation
    std.enclave.sealing_key
      @ () -> bytes
      + returns a stable key derived from the enclave identity
      # attestation
    std.enclave.quote
      @ (report_data: bytes) -> bytes
      + returns an attestation quote embedding report_data
      # attestation

confidential_contract
  confidential_contract.new_runtime
    @ (genesis_state: bytes) -> runtime_state
    + initializes the contract runtime with a sealed genesis state
    # construction
    -> std.enclave.sealing_key
    -> std.crypto.encrypt_aead
  confidential_contract.deploy
    @ (state: runtime_state, code: bytes, initializer: bytes) -> result[tuple[contract_id, runtime_state], string]
    + stores new contract code and runs its initializer inside the enclave
    - returns error when code fails determinism checks
    # deployment
    -> std.crypto.sha256
  confidential_contract.invoke
    @ (state: runtime_state, cid: contract_id, caller: bytes, input: bytes) -> result[tuple[bytes, runtime_state], string]
    + executes the contract method and returns its output
    - returns error on runtime trap
    # execution
  confidential_contract.seal_state
    @ (state: runtime_state) -> bytes
    + returns the encrypted blob representing current storage
    # persistence
    -> std.enclave.sealing_key
    -> std.crypto.encrypt_aead
  confidential_contract.unseal_state
    @ (blob: bytes) -> result[runtime_state, string]
    + restores runtime state from a sealed blob
    - returns error when blob was sealed by a different enclave identity
    # persistence
    -> std.enclave.sealing_key
    -> std.crypto.decrypt_aead
  confidential_contract.attest_transition
    @ (state: runtime_state, prev_root: bytes, new_root: bytes) -> bytes
    + returns a signed attestation binding prev_root to new_root under the enclave identity
    # attestation
    -> std.enclave.measure_self
    -> std.enclave.quote
    -> std.crypto.sign_ed25519
  confidential_contract.verify_attestation
    @ (prev_root: bytes, new_root: bytes, quote: bytes, expected_measurement: bytes) -> bool
    + returns true when the quote binds the roots under the expected enclave code
    # verification
    -> std.crypto.verify_ed25519
  confidential_contract.state_root
    @ (state: runtime_state) -> bytes
    + returns the merkle root of the current storage
    # persistence
    -> std.crypto.sha256
  confidential_contract.simulate_query
    @ (state: runtime_state, cid: contract_id, input: bytes) -> result[bytes, string]
    + runs a read-only method and returns its output without mutating state
    # query
  confidential_contract.encrypt_input
    @ (state: runtime_state, caller_pub: bytes, plaintext: bytes) -> bytes
    + encrypts an input for an enclave under a shared key derived from caller_pub
    # confidentiality
    -> std.crypto.encrypt_aead
  confidential_contract.decrypt_input
    @ (state: runtime_state, caller_pub: bytes, ciphertext: bytes) -> result[bytes, string]
    + recovers a caller-encrypted input inside the enclave
    - returns error on bad tag
    # confidentiality
    -> std.crypto.decrypt_aead
