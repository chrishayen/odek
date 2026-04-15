# Requirement: "a library for parsing and validating x.509 proxy certificates"

std
  std.encoding
    std.encoding.pem_decode
      fn (text: string) -> result[list[pem_block], string]
      + returns each pem block with its type and der bytes
      - returns error on malformed pem
      # encoding
    std.encoding.der_parse_certificate
      fn (data: bytes) -> result[certificate, string]
      + parses a der-encoded x.509 certificate
      - returns error on malformed der
      # encoding
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.crypto
    std.crypto.verify_signature
      fn (issuer_public_key: bytes, tbs_bytes: bytes, signature: bytes, algorithm: string) -> bool
      + returns true when the signature verifies under the issuer's public key
      # cryptography

proxy_cert
  proxy_cert.load_chain
    fn (pem_text: string) -> result[list[certificate], string]
    + decodes a pem bundle into an ordered list of certificates, leaf first
    - returns error when no certificates are found
    # loading
    -> std.encoding.pem_decode
    -> std.encoding.der_parse_certificate
  proxy_cert.is_proxy
    fn (cert: certificate) -> bool
    + returns true when the certificate carries a proxy certificate info extension
    # classification
  proxy_cert.validate_chain
    fn (chain: list[certificate], trust_roots: list[certificate]) -> result[void, string]
    + verifies each signature, checks validity periods, and ensures proxies only delegate from an end-entity
    - returns error when any signature fails to verify
    - returns error when a certificate is expired or not yet valid
    - returns error when a proxy certificate is issued by another proxy with disallowed delegation
    # validation
    -> std.crypto.verify_signature
    -> std.time.now_seconds
  proxy_cert.effective_identity
    fn (chain: list[certificate]) -> result[string, string]
    + returns the subject distinguished name of the end-entity that issued the first proxy
    - returns error when no end-entity is present
    # identity
