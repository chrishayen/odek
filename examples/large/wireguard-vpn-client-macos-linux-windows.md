# Requirement: "a VPN client library using the WireGuard protocol"

Parses config, manages peers and the tunnel interface, and routes encrypted traffic.

std
  std.crypto
    std.crypto.curve25519_keypair
      @ (private_key: bytes) -> bytes
      + derives the public key from a 32-byte private key
      # cryptography
    std.crypto.curve25519_shared
      @ (private_key: bytes, peer_public: bytes) -> bytes
      + returns a 32-byte shared secret via X25519
      # cryptography
    std.crypto.chacha20poly1305_seal
      @ (key: bytes, nonce: bytes, plaintext: bytes, aad: bytes) -> bytes
      + encrypts and authenticates plaintext
      # cryptography
    std.crypto.chacha20poly1305_open
      @ (key: bytes, nonce: bytes, ciphertext: bytes, aad: bytes) -> result[bytes, string]
      + decrypts and verifies ciphertext
      - returns error when the tag does not verify
      # cryptography
    std.crypto.blake2s
      @ (data: bytes) -> bytes
      + returns a 32-byte BLAKE2s digest
      # cryptography
  std.encoding
    std.encoding.base64_decode
      @ (s: string) -> result[bytes, string]
      + decodes standard base64
      - returns error on invalid characters
      # encoding
  std.net
    std.net.udp_open
      @ (port: i32) -> result[udp_socket, string]
      + binds a UDP socket on the given port
      - returns error when the port is in use
      # network
    std.net.udp_send
      @ (sock: udp_socket, addr: string, port: i32, data: bytes) -> result[void, string]
      + sends a datagram
      - returns error on send failure
      # network
    std.net.udp_recv
      @ (sock: udp_socket) -> result[udp_packet, string]
      + receives the next datagram
      - returns error on socket failure
      # network
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

vpn
  vpn.parse_config
    @ (text: string) -> result[vpn_config, string]
    + parses an interface/peer config into a typed struct
    - returns error on missing required fields
    # configuration
    -> std.encoding.base64_decode
  vpn.derive_public_key
    @ (private_key: bytes) -> bytes
    + returns the curve25519 public key for a private key
    # keys
    -> std.crypto.curve25519_keypair
  vpn.new_tunnel
    @ (cfg: vpn_config) -> result[tunnel_state, string]
    + creates an uninitialized tunnel from parsed config
    - returns error when no peers are defined
    # construction
  vpn.handshake_init
    @ (state: tunnel_state, peer_index: i32) -> result[tuple[tunnel_state, bytes], string]
    + performs the initiator half of the Noise IK handshake
    - returns error when the peer is unknown
    # handshake
    -> std.crypto.curve25519_shared
    -> std.crypto.blake2s
    -> std.time.now_millis
  vpn.handshake_response
    @ (state: tunnel_state, peer_index: i32, msg: bytes) -> result[tunnel_state, string]
    + consumes the responder message and installs session keys
    - returns error on mac failure
    # handshake
    -> std.crypto.blake2s
  vpn.encrypt_packet
    @ (state: tunnel_state, peer_index: i32, plaintext: bytes) -> result[tuple[tunnel_state, bytes], string]
    + encapsulates an IP packet as a Transport message
    - returns error when no session exists for the peer
    # dataplane
    -> std.crypto.chacha20poly1305_seal
  vpn.decrypt_packet
    @ (state: tunnel_state, ciphertext: bytes) -> result[tuple[tunnel_state, bytes], string]
    + decapsulates a Transport message back to an IP packet
    - returns error on replay or tag failure
    # dataplane
    -> std.crypto.chacha20poly1305_open
  vpn.route_outbound
    @ (state: tunnel_state, sock: udp_socket, plaintext: bytes) -> result[tunnel_state, string]
    + encrypts and sends a packet to the matching peer's endpoint
    - returns error when no peer matches the destination
    # routing
    -> std.net.udp_send
  vpn.route_inbound
    @ (state: tunnel_state, sock: udp_socket) -> result[tuple[tunnel_state, bytes], string]
    + receives and decrypts the next datagram
    - returns error on decryption failure
    # routing
    -> std.net.udp_recv
