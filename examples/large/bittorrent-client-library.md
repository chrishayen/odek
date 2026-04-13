# Requirement: "a BitTorrent client library"

Parses .torrent files, talks to trackers, negotiates peer connections, and exchanges piece blocks. The std layer provides bencode, SHA-1, and socket primitives; the project layer manages the swarm and piece store.

std
  std.bencode
    std.bencode.parse
      @ (raw: bytes) -> result[bencode_value, string]
      + parses a bencoded byte string
      - returns error on malformed input
      # serialization
    std.bencode.encode
      @ (value: bencode_value) -> bytes
      + encodes a bencode_value back to bytes
      # serialization
  std.crypto
    std.crypto.sha1
      @ (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest
      # hashing
  std.net
    std.net.tcp_connect
      @ (host: string, port: i32) -> result[tcp_handle, string]
      + opens a TCP connection
      - returns error on failure
      # networking
    std.net.tcp_read
      @ (handle: tcp_handle, n: i32) -> result[bytes, string]
      + reads up to n bytes
      # networking
    std.net.tcp_write
      @ (handle: tcp_handle, data: bytes) -> result[void, string]
      + writes all bytes
      # networking
    std.net.tcp_close
      @ (handle: tcp_handle) -> result[void, string]
      + closes the connection
      # networking
  std.http
    std.http.get
      @ (url: string) -> result[http_response, string]
      + performs a GET and returns status, headers, and body
      - returns error on network failure
      # networking
  std.fs
    std.fs.write_at
      @ (path: string, offset: i64, data: bytes) -> result[void, string]
      + writes data at offset, extending the file as needed
      - returns error on I/O failure
      # filesystem
    std.fs.read_at
      @ (path: string, offset: i64, n: i64) -> result[bytes, string]
      + reads n bytes starting at offset
      - returns error on I/O failure
      # filesystem

torrent
  torrent.parse_metainfo
    @ (raw: bytes) -> result[metainfo, string]
    + parses a .torrent file into metainfo (info_hash, piece_length, pieces, files, announce)
    - returns error when required keys are missing
    # metainfo
    -> std.bencode.parse
    -> std.crypto.sha1
  torrent.announce
    @ (info: metainfo, peer_id: bytes, port: i32, uploaded: i64, downloaded: i64, left: i64) -> result[tracker_response, string]
    + sends an HTTP announce and parses the response
    - returns error on network or decoding failure
    # tracker
    -> std.http.get
    -> std.bencode.parse
  torrent.handshake
    @ (handle: tcp_handle, info_hash: bytes, peer_id: bytes) -> result[peer_handshake, string]
    + exchanges the BitTorrent handshake
    - returns error on info-hash mismatch
    # peer_protocol
    -> std.net.tcp_write
    -> std.net.tcp_read
  torrent.read_peer_message
    @ (handle: tcp_handle) -> result[peer_message, string]
    + reads one length-prefixed peer message
    - returns error on short read
    # peer_protocol
    -> std.net.tcp_read
  torrent.write_peer_message
    @ (handle: tcp_handle, msg: peer_message) -> result[void, string]
    + writes a length-prefixed peer message
    # peer_protocol
    -> std.net.tcp_write
  torrent.new_piece_store
    @ (info: metainfo, base_path: string) -> piece_store
    + creates a piece store rooted at base_path
    # storage
  torrent.verify_piece
    @ (info: metainfo, piece_index: i32, data: bytes) -> bool
    + returns true when the SHA-1 of data matches the piece hash
    # storage
    -> std.crypto.sha1
  torrent.write_piece
    @ (store: piece_store, piece_index: i32, data: bytes) -> result[piece_store, string]
    + writes a verified piece at its file offset(s)
    - returns error when the piece fails verification
    # storage
    -> std.fs.write_at
  torrent.read_block
    @ (store: piece_store, piece_index: i32, offset: i32, length: i32) -> result[bytes, string]
    + reads a block from a stored piece
    # storage
    -> std.fs.read_at
  torrent.new_swarm
    @ (info: metainfo, peer_id: bytes) -> swarm_state
    + creates swarm state tracking peers, piece availability, and request queue
    # swarm
  torrent.add_peer
    @ (swarm: swarm_state, handle: tcp_handle, handshake: peer_handshake) -> swarm_state
    + registers a new peer connection
    # swarm
    -> std.net.tcp_connect
  torrent.update_bitfield
    @ (swarm: swarm_state, peer_index: i32, bitfield: bytes) -> swarm_state
    + records which pieces a peer has
    # swarm
  torrent.pick_next_request
    @ (swarm: swarm_state) -> optional[block_request]
    + chooses the next block to request using rarest-first
    - returns none when no rare block is available
    # strategy
  torrent.handle_message
    @ (swarm: swarm_state, store: piece_store, peer_index: i32, msg: peer_message) -> tuple[swarm_state, piece_store]
    + processes an incoming peer message and updates state
    # swarm
  torrent.close_peer
    @ (swarm: swarm_state, peer_index: i32) -> swarm_state
    + drops the peer and reclaims any outstanding requests
    # swarm
    -> std.net.tcp_close
  torrent.encode_announce_response
    @ (response: tracker_response) -> bytes
    + encodes a tracker response for persistence or testing
    # tracker
    -> std.bencode.encode
