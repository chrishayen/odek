# Requirement: "a memcache client library"

Connects to a memcache server, issues text-protocol commands, and returns parsed responses. Wire I/O and time live in std.

std
  std.net
    std.net.dial_tcp
      @ (host: string, port: i32) -> result[conn_state, string]
      + opens a TCP connection to host:port
      - returns error on connection failure
      # networking
    std.net.write_all
      @ (conn: conn_state, data: bytes) -> result[void, string]
      + writes the full buffer to the connection
      # networking
    std.net.read_line
      @ (conn: conn_state) -> result[string, string]
      + reads a line terminated by CRLF
      - returns error on EOF mid-line
      # networking
    std.net.read_n
      @ (conn: conn_state, n: i32) -> result[bytes, string]
      + reads exactly n bytes
      # networking

memcache
  memcache.connect
    @ (host: string, port: i32) -> result[memcache_client, string]
    + returns a client bound to the given memcache server
    # construction
    -> std.net.dial_tcp
  memcache.set
    @ (client: memcache_client, key: string, value: bytes, ttl_seconds: i32) -> result[void, string]
    + stores a value under key with an expiration
    - returns error when the server responds with a non-STORED status
    # write
    -> std.net.write_all
    -> std.net.read_line
  memcache.get
    @ (client: memcache_client, key: string) -> result[optional[bytes], string]
    + returns the value when present
    + returns none when the key is absent
    - returns error on protocol mismatch
    # read
    -> std.net.write_all
    -> std.net.read_line
    -> std.net.read_n
  memcache.delete
    @ (client: memcache_client, key: string) -> result[bool, string]
    + returns true when the key was deleted
    + returns false when the key was not present
    # delete
    -> std.net.write_all
    -> std.net.read_line
  memcache.increment
    @ (client: memcache_client, key: string, delta: i64) -> result[i64, string]
    + returns the new value after incrementing
    - returns error when the key is not numeric
    # atomic
    -> std.net.write_all
    -> std.net.read_line
  memcache.flush_all
    @ (client: memcache_client) -> result[void, string]
    + clears all entries on the server
    # admin
    -> std.net.write_all
    -> std.net.read_line
