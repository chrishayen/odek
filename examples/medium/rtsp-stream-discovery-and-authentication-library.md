# Requirement: "a library to discover and authenticate against rtsp video streams"

Scans a host for RTSP endpoints, probes common stream paths, and attempts authentication against a credential list to identify accessible streams.

std
  std.net
    std.net.tcp_connect
      @ (host: string, port: u16, timeout_ms: i32) -> result[connection, string]
      + returns an open connection to the host:port
      - returns error on refused, unreachable, or timeout
      # networking
    std.net.send
      @ (conn: connection, data: bytes) -> result[void, string]
      + writes data to the connection
      - returns error when the peer has closed the connection
      # networking
    std.net.recv
      @ (conn: connection, max_bytes: i32) -> result[bytes, string]
      + reads up to max_bytes from the connection
      - returns error on read failure
      # networking
  std.crypto
    std.crypto.md5
      @ (data: bytes) -> bytes
      + returns the MD5 digest of data
      # cryptography

rtsp_scan
  rtsp_scan.probe_port
    @ (host: string, port: u16) -> result[bool, string]
    + returns true when the port responds to an RTSP OPTIONS request with a valid status line
    - returns error when the network is unreachable
    # discovery
    -> std.net.tcp_connect
    -> std.net.send
    -> std.net.recv
  rtsp_scan.describe_stream
    @ (host: string, port: u16, path: string, credentials: optional[credentials]) -> result[stream_status, string]
    + returns a status indicating ok, auth_required with realm/nonce, or not_found
    - returns error on a malformed response
    # probing
    -> std.net.send
    -> std.net.recv
  rtsp_scan.build_digest_response
    @ (user: string, password: string, realm: string, nonce: string, method: string, uri: string) -> string
    + returns the MD5-based HTTP Digest response string for the given challenge
    # auth
    -> std.crypto.md5
  rtsp_scan.try_credentials
    @ (host: string, port: u16, path: string, candidates: list[credentials]) -> result[optional[credentials], string]
    + returns the first credential pair that successfully authenticates, or none if every candidate fails
    - returns error when the target becomes unreachable mid-scan
    # auth
  rtsp_scan.scan_host
    @ (host: string, ports: list[u16], paths: list[string], candidates: list[credentials]) -> list[stream_finding]
    + returns one finding per reachable (port, path) combination, tagged with the winning credential if any
    # orchestration
