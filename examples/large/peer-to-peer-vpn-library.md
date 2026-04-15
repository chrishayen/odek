# Requirement: "a peer-to-peer VPN library"

Peer discovery, authenticated key exchange, symmetric encryption of tunneled frames, and a routing table mapping virtual addresses to peers. Packet I/O is delegated to the caller; this library operates on byte-level frames.

std
  std.crypto
    std.crypto.x25519_keypair
      fn () -> tuple[bytes, bytes]
      + returns a new (private, public) X25519 keypair
      # cryptography
    std.crypto.x25519_shared
      fn (private: bytes, peer_public: bytes) -> bytes
      + computes the shared secret
      # cryptography
    std.crypto.aead_encrypt
      fn (key: bytes, nonce: bytes, plaintext: bytes, ad: bytes) -> bytes
      + encrypts and authenticates plaintext with associated data
      # cryptography
    std.crypto.aead_decrypt
      fn (key: bytes, nonce: bytes, ciphertext: bytes, ad: bytes) -> result[bytes, string]
      + decrypts and verifies ciphertext
      - returns error when authentication tag is invalid
      # cryptography
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.net
    std.net.parse_addr
      fn (raw: string) -> result[net_addr, string]
      + parses a host:port address
      - returns error on malformed input
      # networking

vpn
  vpn.new_node
    fn (virtual_addr: string) -> node_state
    + creates a VPN node with an assigned virtual address
    # construction
    -> std.crypto.x25519_keypair
  vpn.local_public_key
    fn (state: node_state) -> bytes
    + returns this node's public key for distribution
    # identity
  vpn.add_peer
    fn (state: node_state, virtual_addr: string, public_key: bytes, transport_addr: string) -> result[node_state, string]
    + registers a peer and derives the shared session key
    - returns error when the transport address is malformed
    # peering
    -> std.crypto.x25519_shared
    -> std.net.parse_addr
  vpn.remove_peer
    fn (state: node_state, virtual_addr: string) -> node_state
    + removes a peer and its session state
    # peering
  vpn.encrypt_frame
    fn (state: node_state, dest_virtual_addr: string, payload: bytes) -> result[bytes, string]
    + returns the encrypted outbound frame for the destination
    - returns error when destination is unknown
    # tunneling
    -> std.crypto.aead_encrypt
    -> std.time.now_millis
  vpn.decrypt_frame
    fn (state: node_state, source_virtual_addr: string, frame: bytes) -> result[bytes, string]
    + returns the decrypted payload from an authenticated frame
    - returns error when source is unknown or tag is invalid
    # tunneling
    -> std.crypto.aead_decrypt
  vpn.route_lookup
    fn (state: node_state, dest_virtual_addr: string) -> optional[string]
    + returns the transport address for the destination virtual address
    - returns none when the destination is not routable
    # routing
  vpn.peer_list
    fn (state: node_state) -> list[peer_summary]
    + returns all configured peers with virtual and transport addresses
    # inspection
