# Requirement: "a library that dynamically generates self-signed certificates and certificate authorities for tests"

The project layer is two entry points; the heavy lifting lives in std cryptographic primitives.

std
  std.crypto
    std.crypto.rsa_generate_key
      fn (bits: i32) -> result[rsa_keypair, string]
      + returns a fresh RSA keypair of the given size
      - returns error when bits is not a supported size
      # cryptography
    std.crypto.rsa_sign
      fn (key: rsa_keypair, data: bytes) -> bytes
      + returns an RSA signature over the data
      # cryptography
  std.pki
    std.pki.encode_certificate
      fn (subject: string, issuer: string, public_key: bytes, signature: bytes, not_before: i64, not_after: i64) -> bytes
      + returns a DER-encoded X.509 certificate
      # pki
    std.pki.encode_pem
      fn (label: string, der: bytes) -> string
      + returns a PEM block with the given label around the DER bytes
      # pki
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

testcerts
  testcerts.generate_ca
    fn (common_name: string, validity_days: i32) -> result[tuple[string, string], string]
    + returns (certificate_pem, private_key_pem) for a new self-signed CA
    - returns error on key generation failure
    # certificate_authority
    -> std.crypto.rsa_generate_key
    -> std.crypto.rsa_sign
    -> std.pki.encode_certificate
    -> std.pki.encode_pem
    -> std.time.now_seconds
  testcerts.issue_leaf
    fn (ca_cert_pem: string, ca_key_pem: string, common_name: string, sans: list[string], validity_days: i32) -> result[tuple[string, string], string]
    + returns (certificate_pem, private_key_pem) for a leaf certificate signed by the CA
    - returns error when the CA certificate or key fails to parse
    # leaf_issuance
    -> std.crypto.rsa_generate_key
    -> std.crypto.rsa_sign
    -> std.pki.encode_certificate
    -> std.pki.encode_pem
    -> std.time.now_seconds
