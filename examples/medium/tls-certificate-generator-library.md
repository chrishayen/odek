# Requirement: "a tls certificate generation library"

Issues a CA, signs leaf certificates under it, and serializes both to PEM. Cryptography primitives live in std.

std
  std.crypto
    std.crypto.generate_rsa_key
      fn (bits: i32) -> result[rsa_key, string]
      + returns a fresh RSA keypair
      - returns error when bits is below 2048
      # cryptography
    std.crypto.sign_rsa_sha256
      fn (key: rsa_key, data: bytes) -> bytes
      + returns an RSA-SHA256 signature
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n bytes from a cryptographic RNG
      # randomness
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.encoding
    std.encoding.pem_encode
      fn (label: string, der: bytes) -> string
      + wraps DER bytes in a PEM block with the given label
      # encoding
    std.encoding.der_encode_certificate
      fn (fields: cert_fields, signature: bytes) -> bytes
      + returns the ASN.1 DER encoding of an x509 certificate
      # encoding

tls_certs
  tls_certs.new_ca
    fn (common_name: string, valid_days: i32) -> result[ca_state, string]
    + returns a self-signed CA with a fresh key and the given validity
    - returns error when valid_days is not positive
    # ca
    -> std.crypto.generate_rsa_key
    -> std.crypto.random_bytes
    -> std.time.now_seconds
    -> std.encoding.der_encode_certificate
    -> std.crypto.sign_rsa_sha256
  tls_certs.issue_leaf
    fn (ca: ca_state, common_name: string, dns_names: list[string], valid_days: i32) -> result[leaf_cert, string]
    + returns a leaf certificate signed by the CA
    - returns error when any dns name is malformed
    # issuance
    -> std.crypto.generate_rsa_key
    -> std.crypto.random_bytes
    -> std.time.now_seconds
    -> std.encoding.der_encode_certificate
    -> std.crypto.sign_rsa_sha256
  tls_certs.ca_to_pem
    fn (ca: ca_state) -> string
    + returns the CA certificate encoded as PEM
    # encoding
    -> std.encoding.pem_encode
  tls_certs.leaf_to_pem
    fn (leaf: leaf_cert) -> string
    + returns the leaf certificate encoded as PEM
    # encoding
    -> std.encoding.pem_encode
  tls_certs.key_to_pem
    fn (key: rsa_key) -> string
    + returns the private key encoded as a PEM block
    # encoding
    -> std.encoding.pem_encode
  tls_certs.verify_chain
    fn (leaf: leaf_cert, ca: ca_state) -> result[void, string]
    + returns ok when the leaf's signature verifies under the CA and dates are valid
    - returns error when the leaf has expired
    - returns error when the signature does not match
    # verification
    -> std.time.now_seconds
