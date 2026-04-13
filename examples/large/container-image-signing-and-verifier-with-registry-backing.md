# Requirement: "a container image signing and verification library backed by an image registry"

Signatures are computed over the image digest, stored as attached artifacts in the registry, and verified by looking them up.

std
  std.crypto
    std.crypto.sign_ed25519
      @ (private_key: bytes, message: bytes) -> bytes
      + computes a 64-byte Ed25519 signature
      # cryptography
    std.crypto.verify_ed25519
      @ (public_key: bytes, message: bytes, signature: bytes) -> bool
      + returns true when the signature matches
      - returns false on any tampering
      # cryptography
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest
      # cryptography
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + encodes bytes to standard base64
      # encoding
    std.encoding.base64_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes standard base64
      - returns error on invalid input
      # encoding
  std.json
    std.json.encode_object
      @ (fields: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid input
      # serialization

image_sign
  image_sign.compute_image_digest
    @ (manifest: bytes) -> string
    + returns "sha256:<hex>" over the manifest bytes
    # digest
    -> std.crypto.sha256
  image_sign.build_payload
    @ (digest: string, identity: string) -> string
    + returns a JSON payload binding the digest to a signer identity
    # payload
    -> std.json.encode_object
  image_sign.sign_payload
    @ (payload: string, private_key: bytes) -> string
    + returns the base64-encoded Ed25519 signature of the payload
    # signing
    -> std.crypto.sign_ed25519
    -> std.encoding.base64_encode
  image_sign.store_signature
    @ (registry: registry_handle, digest: string, payload: string, signature: string) -> result[void, string]
    + uploads the payload and signature as an attached artifact tagged by digest
    - returns error when the registry refuses the upload
    # storage
  image_sign.fetch_signatures
    @ (registry: registry_handle, digest: string) -> result[list[signature_record], string]
    + returns all signature records attached to the given digest
    - returns error when no attachment is found
    # retrieval
    -> std.json.parse_object
  image_sign.verify_signature
    @ (record: signature_record, public_key: bytes, expected_digest: string) -> result[void, string]
    + returns ok when the signature is valid and the payload binds to expected_digest
    - returns error when the payload digest does not match
    - returns error when the Ed25519 signature does not verify
    # verification
    -> std.encoding.base64_decode
    -> std.crypto.verify_ed25519
  image_sign.verify_image
    @ (registry: registry_handle, manifest: bytes, public_keys: list[bytes]) -> result[signature_record, string]
    + returns the first signature whose payload binds to the manifest digest and verifies under any key
    - returns error when no signature verifies
    # verification
