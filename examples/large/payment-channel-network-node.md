# Requirement: "a payment channel network node"

A node that opens, updates, and cooperatively closes bidirectional payment channels, routes hashlocked payments across a graph of peers, and broadcasts settlement transactions to an underlying ledger.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest of the input
      # cryptography
    std.crypto.sign_ecdsa
      fn (private_key: bytes, message: bytes) -> bytes
      + returns a compact ECDSA signature over the message
      - returns empty bytes if the private key is not 32 bytes
      # cryptography
    std.crypto.verify_ecdsa
      fn (public_key: bytes, message: bytes, sig: bytes) -> bool
      + returns true only when signature is valid for the public key
      # cryptography
  std.net
    std.net.dial_tcp
      fn (host: string, port: u16) -> result[conn_state, string]
      + returns a connected socket handle
      - returns error when host is unreachable
      # networking
    std.net.send
      fn (conn: conn_state, payload: bytes) -> result[void, string]
      + writes the full payload to the socket
      # networking
    std.net.recv
      fn (conn: conn_state, max: u32) -> result[bytes, string]
      + reads up to max bytes
      - returns error when the peer closed
      # networking
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

channel_node
  channel_node.new
    fn (identity_key: bytes) -> channel_node_state
    + creates a node bound to a long-lived identity keypair
    # construction
  channel_node.connect_peer
    fn (state: channel_node_state, peer_pubkey: bytes, host: string, port: u16) -> result[channel_node_state, string]
    + establishes an authenticated transport to a remote peer
    - returns error when the peer identity does not match
    # peering
    -> std.net.dial_tcp
    -> std.crypto.verify_ecdsa
  channel_node.open_channel
    fn (state: channel_node_state, peer_pubkey: bytes, capacity: u64) -> result[tuple[channel_id, channel_node_state], string]
    + constructs the funding transaction and exchanges commitment signatures
    - returns error when the peer is not connected
    # channel_lifecycle
    -> std.crypto.sign_ecdsa
  channel_node.send_htlc
    fn (state: channel_node_state, cid: channel_id, amount: u64, payment_hash: bytes, expiry: i64) -> result[channel_node_state, string]
    + updates the commitment to add a hashlocked output
    - returns error when local balance is insufficient
    # htlc
    -> std.crypto.sha256
  channel_node.settle_htlc
    fn (state: channel_node_state, cid: channel_id, preimage: bytes) -> result[channel_node_state, string]
    + releases funds once the preimage for an outstanding htlc is revealed
    - returns error when the preimage does not hash to any pending htlc
    # htlc
    -> std.crypto.sha256
  channel_node.route_payment
    fn (state: channel_node_state, destination: bytes, amount: u64) -> result[list[channel_id], string]
    + returns the sequence of channels chosen to forward the payment
    - returns error when no path with enough capacity exists
    # routing
  channel_node.forward_htlc
    fn (state: channel_node_state, incoming: channel_id, outgoing: channel_id, amount: u64, payment_hash: bytes) -> result[channel_node_state, string]
    + relays an htlc along a multi-hop route
    - returns error when fee policy on the outgoing channel is not met
    # routing
  channel_node.close_channel
    fn (state: channel_node_state, cid: channel_id) -> result[bytes, string]
    + returns the signed cooperative closing transaction
    - returns error when the channel has pending htlcs
    # channel_lifecycle
    -> std.crypto.sign_ecdsa
  channel_node.force_close
    fn (state: channel_node_state, cid: channel_id) -> result[bytes, string]
    + returns the latest commitment transaction for unilateral broadcast
    # channel_lifecycle
  channel_node.handle_peer_message
    fn (state: channel_node_state, from: bytes, msg: bytes) -> result[channel_node_state, string]
    + dispatches incoming protocol messages to the matching handler
    - returns error on unknown message types
    # protocol
    -> std.net.recv
  channel_node.tick
    fn (state: channel_node_state) -> channel_node_state
    + expires htlcs whose cltv deadline has passed and retries pending updates
    # timekeeping
    -> std.time.now_seconds
