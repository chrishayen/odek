# Requirement: "a private overlay mesh network built on an encrypted-tunnel datagram protocol"

Models the mesh as a membership registry plus a peer routing table. The library manages invites, membership, peer endpoints, and produces tunnel configuration artifacts for each node.

std
  std.crypto
    std.crypto.generate_keypair
      @ () -> tuple[bytes, bytes]
      + returns a (public, private) asymmetric key pair
      # cryptography
    std.crypto.random_bytes
      @ (length: i32) -> bytes
      + returns cryptographically random bytes of the requested length
      # cryptography
  std.net
    std.net.parse_cidr
      @ (text: string) -> result[cidr_range, string]
      + parses a CIDR range such as "10.0.0.0/16"
      - returns error on malformed text
      # networking
    std.net.next_address_in
      @ (range: cidr_range, used: list[string]) -> result[string, string]
      + returns the next unused address within the range
      - returns error when the range is exhausted
      # networking
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

mesh
  mesh.create_network
    @ (name: string, range: string) -> result[network_state, string]
    + creates a new mesh network with a given CIDR range
    - returns error on invalid CIDR
    # administration
    -> std.net.parse_cidr
  mesh.add_peer
    @ (net: network_state, peer_name: string, endpoint: string) -> result[tuple[network_state, peer_record], string]
    + assigns an overlay address to a peer and records a public endpoint and fresh keypair
    - returns error when the peer name already exists
    - returns error when the address range is exhausted
    # membership
    -> std.crypto.generate_keypair
    -> std.net.next_address_in
  mesh.remove_peer
    @ (net: network_state, peer_name: string) -> result[network_state, string]
    + removes a peer from the network
    - returns error when the peer is not a member
    # membership
  mesh.generate_invite
    @ (net: network_state, peer_name: string, ttl_seconds: i64) -> result[invite_token, string]
    + mints a single-use invite for a peer with an expiry timestamp
    - returns error when the peer is already a member
    # onboarding
    -> std.crypto.random_bytes
    -> std.time.now_seconds
  mesh.redeem_invite
    @ (net: network_state, token: invite_token, public_key: bytes, endpoint: string) -> result[tuple[network_state, peer_record], string]
    + consumes an invite and registers the redeeming peer's public key and endpoint
    - returns error when the token is expired or unknown
    # onboarding
    -> std.time.now_seconds
  mesh.render_peer_config
    @ (net: network_state, peer_name: string) -> result[string, string]
    + returns the rendered tunnel configuration for a peer, listing its siblings as routes
    - returns error when the peer is not a member
    # configuration
  mesh.list_peers
    @ (net: network_state) -> list[peer_record]
    + returns all peer records in the network
    # inspection
  mesh.peer_for_address
    @ (net: network_state, address: string) -> optional[peer_record]
    + returns the peer that owns an overlay address if any
    # routing
  mesh.rotate_peer_key
    @ (net: network_state, peer_name: string) -> result[tuple[network_state, peer_record], string]
    + rotates the peer's keypair, leaving the address unchanged
    - returns error when the peer is not a member
    # key_management
    -> std.crypto.generate_keypair
