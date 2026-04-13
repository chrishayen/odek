# Requirement: "a client driver for a distributed key-value store"

Connect, then get/put/delete by bucket and key over a pluggable transport.

std
  std.net
    std.net.dial_tcp
      @ (host: string, port: i32) -> result[conn_handle, string]
      + opens a TCP connection
      - returns error when the host is unreachable
      # networking
    std.net.write_all
      @ (conn: conn_handle, data: bytes) -> result[void, string]
      + writes all bytes
      # networking
    std.net.read_exact
      @ (conn: conn_handle, n: i32) -> result[bytes, string]
      + reads exactly n bytes
      - returns error on short read
      # networking

kv_driver
  kv_driver.connect
    @ (host: string, port: i32) -> result[session_state, string]
    + opens a session to the server
    - returns error when the connection fails
    # session
    -> std.net.dial_tcp
  kv_driver.put
    @ (session: session_state, bucket: string, key: string, value: bytes) -> result[void, string]
    + stores value under (bucket, key)
    - returns error on server-side failures
    # write
    -> std.net.write_all
    -> std.net.read_exact
  kv_driver.get
    @ (session: session_state, bucket: string, key: string) -> result[optional[bytes], string]
    + returns the stored value, or none when the key is absent
    - returns error on server-side failures
    # read
    -> std.net.write_all
    -> std.net.read_exact
  kv_driver.delete
    @ (session: session_state, bucket: string, key: string) -> result[void, string]
    + removes the entry when present
    - returns error on server-side failures
    # write
    -> std.net.write_all
    -> std.net.read_exact
  kv_driver.close
    @ (session: session_state) -> result[void, string]
    + closes the underlying connection
    # session
