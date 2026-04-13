# Requirement: "a recursive DNS resolver with DNSSEC validation"

Wire-format encode/decode for DNS messages, iterative resolution from the root, cached answers with per-record TTL, and DNSSEC signature chain validation. Transport (UDP/TCP) is the caller's responsibility.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.crypto
    std.crypto.rsa_verify
      @ (public_key: bytes, message: bytes, signature: bytes) -> bool
      + returns true when RSA signature is valid
      # cryptography
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns SHA-256 digest
      # cryptography
  std.net
    std.net.parse_ipv4
      @ (raw: string) -> result[bytes, string]
      + parses dotted-quad into 4 bytes
      - returns error on malformed input
      # networking

dns
  dns.encode_query
    @ (name: string, qtype: i32, qid: u16) -> bytes
    + encodes a DNS query in wire format
    # wire_protocol
  dns.decode_message
    @ (wire: bytes) -> result[dns_message, string]
    + decodes a DNS wire-format message
    - returns error on truncated or malformed input
    # wire_protocol
  dns.encode_response
    @ (msg: dns_message) -> bytes
    + encodes a DNS message as wire-format bytes
    # wire_protocol
  dns.new_resolver
    @ (root_hints: list[string]) -> resolver_state
    + creates a resolver seeded with root server addresses
    # construction
    -> std.net.parse_ipv4
  dns.resolve
    @ (state: resolver_state, name: string, qtype: i32) -> result[tuple[list[dns_record], resolver_state], string]
    + iteratively resolves a name from the nearest authoritative server
    - returns error on NXDOMAIN or network failure
    # recursion
    -> std.time.now_seconds
  dns.cache_get
    @ (state: resolver_state, name: string, qtype: i32) -> optional[list[dns_record]]
    + returns cached records when present and within TTL
    - returns none when absent or expired
    # caching
    -> std.time.now_seconds
  dns.cache_put
    @ (state: resolver_state, name: string, qtype: i32, records: list[dns_record]) -> resolver_state
    + inserts records with TTL into the cache
    # caching
    -> std.time.now_seconds
  dns.validate_chain
    @ (records: list[dns_record], rrsigs: list[dns_record], trust_anchor: bytes) -> result[void, string]
    + validates the DNSSEC signature chain from records up to the trust anchor
    - returns error when any signature fails or the chain is broken
    # dnssec
    -> std.crypto.rsa_verify
    -> std.crypto.sha256
  dns.minimize_qname
    @ (name: string) -> list[string]
    + returns the sequence of labels queried during qname-minimization
    # privacy
