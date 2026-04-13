# Requirement: "a combined HTTP and peer-to-peer file sharing server"

Exposes files over HTTP and over a peer-to-peer torrent protocol; shared metadata and chunk store in the project, primitives in std.

std
  std.fs
    std.fs.read_range
      @ (path: string, offset: i64, length: i64) -> result[bytes, string]
      + reads length bytes from path starting at offset
      - returns error when offset+length exceeds file size
      # filesystem
    std.fs.file_size
      @ (path: string) -> result[i64, string]
      + returns the size of the file in bytes
      - returns error when path does not exist
      # filesystem
  std.crypto
    std.crypto.sha1
      @ (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest of data
      # cryptography
  std.http
    std.http.serve
      @ (addr: string, handler: fn(http_request) -> http_response) -> result[void, string]
      + binds addr and dispatches requests to handler
      - returns error when bind fails
      # networking
    std.http.response_bytes
      @ (status: i32, body: bytes, content_type: string) -> http_response
      + builds a binary response with the given status
      # networking
  std.net
    std.net.listen_tcp
      @ (addr: string) -> result[tcp_listener, string]
      + binds a TCP listener on addr
      - returns error when bind fails
      # networking
    std.net.accept_tcp
      @ (listener: tcp_listener) -> result[tcp_conn, string]
      + waits for and returns the next incoming connection
      - returns error when the listener is closed
      # networking

share_server
  share_server.add_file
    @ (path: string, piece_size: i32) -> result[file_manifest, string]
    + returns a manifest with per-piece SHA-1 hashes for the file
    - returns error when the file cannot be read
    # indexing
    -> std.fs.file_size
    -> std.fs.read_range
    -> std.crypto.sha1
  share_server.info_hash
    @ (manifest: file_manifest) -> bytes
    + returns a stable 20-byte identifier derived from the manifest
    # identity
    -> std.crypto.sha1
  share_server.new_registry
    @ () -> share_registry
    + returns an empty registry of shared files keyed by info_hash
    # construction
  share_server.register
    @ (reg: share_registry, manifest: file_manifest, path: string) -> void
    + records a file as available to both the HTTP and peer protocols
    # registration
  share_server.serve_http
    @ (reg: share_registry, addr: string) -> result[void, string]
    + serves directory listings and file downloads with HTTP range support
    - returns error when bind fails
    # http_interface
    -> std.http.serve
    -> std.http.response_bytes
  share_server.handle_peer_handshake
    @ (conn: tcp_conn, reg: share_registry) -> result[peer_session, string]
    + accepts a peer handshake and returns a session bound to the requested file
    - returns error on malformed handshake or unknown info_hash
    # peer_protocol
  share_server.handle_peer_request
    @ (session: peer_session, piece_index: i32) -> result[bytes, string]
    + returns the bytes of the requested piece
    - returns error on out-of-range index
    # peer_protocol
    -> std.fs.read_range
  share_server.serve_peers
    @ (reg: share_registry, addr: string) -> result[void, string]
    + accepts peer connections and serves pieces concurrently
    - returns error when the TCP listener cannot be bound
    # peer_interface
    -> std.net.listen_tcp
    -> std.net.accept_tcp
  share_server.start
    @ (reg: share_registry, http_addr: string, peer_addr: string) -> result[void, string]
    + starts both the HTTP server and the peer listener
    - returns error when either address cannot be bound
    # orchestration
