# Requirement: "an X.509 certificate validation library for web PKI"

Parses DER-encoded certificates, verifies signature chains up to a trust anchor, and checks name constraints and validity periods.

std
  std.asn1
    std.asn1.parse_der
      @ (data: bytes) -> result[asn1_node, string]
      + parses a DER-encoded structure into a node tree
      - returns error on malformed length prefixes
      - returns error on truncated input
      # parsing
    std.asn1.read_oid
      @ (node: asn1_node) -> result[string, string]
      + returns the dotted string form of an object identifier
      - returns error when the node is not an OID
      # parsing
    std.asn1.read_integer
      @ (node: asn1_node) -> result[bytes, string]
      + returns the big-endian bytes of an integer node
      # parsing
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
    std.time.parse_utctime
      @ (text: string) -> result[i64, string]
      + parses a UTCTime string to seconds since epoch
      - returns error on invalid formats
      # time
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the sha256 digest
      # hashing
    std.crypto.verify_rsa_pkcs1
      @ (pubkey: bytes, signed: bytes, signature: bytes) -> bool
      + returns true when signature is valid under the RSA public key
      - returns false when the signature does not match
      # signature_verification
    std.crypto.verify_ecdsa
      @ (pubkey: bytes, signed: bytes, signature: bytes) -> bool
      + returns true when the ECDSA signature is valid
      - returns false on mismatch
      # signature_verification

x509
  x509.parse_certificate
    @ (der: bytes) -> result[certificate, string]
    + returns a certificate with subject, issuer, validity, key, and extensions
    - returns error when the ASN.1 structure is not a certificate
    # parsing
    -> std.asn1.parse_der
    -> std.asn1.read_oid
    -> std.asn1.read_integer
  x509.is_time_valid
    @ (cert: certificate, now_seconds: i64) -> bool
    + returns true when now falls within notBefore and notAfter
    # validation
    -> std.time.parse_utctime
  x509.verify_signature
    @ (child: certificate, parent_pubkey: bytes) -> bool
    + returns true when the child's signature verifies under the parent's key
    - returns false when the algorithm is unsupported
    # validation
    -> std.crypto.verify_rsa_pkcs1
    -> std.crypto.verify_ecdsa
  x509.match_name_constraints
    @ (cert: certificate, hostname: string) -> bool
    + returns true when the hostname matches a SAN or the CN
    - returns false when no names match
    # validation
  x509.build_chain
    @ (leaf: certificate, intermediates: list[certificate], roots: list[certificate]) -> result[list[certificate], string]
    + returns an ordered chain from leaf to a trust anchor
    - returns error when no chain can be constructed
    # path_building
  x509.verify
    @ (leaf: certificate, intermediates: list[certificate], roots: list[certificate], hostname: string) -> result[void, string]
    + returns ok when the chain verifies, names match, and all periods are valid
    - returns error when the hostname does not match the leaf
    - returns error when any certificate in the chain is expired
    - returns error when no trusted chain exists
    # validation
    -> std.time.now_seconds
  x509.fingerprint_sha256
    @ (cert: certificate) -> bytes
    + returns the sha256 digest of the DER encoding
    # inspection
    -> std.crypto.sha256
