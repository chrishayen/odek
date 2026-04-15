# Requirement: "a proxy that simulates network and system conditions for automated tests"

A tcp pass-through with configurable fault injectors (latency, bandwidth caps, slow close, data corruption).

std
  std.net
    std.net.listen_tcp
      fn (host: string, port: i32) -> result[listener_handle, string]
      + returns a listener bound to the given address
      - returns error when the address is already in use
      # network
    std.net.accept
      fn (l: listener_handle) -> result[conn_handle, string]
      + returns the next inbound connection
      # network
    std.net.dial_tcp
      fn (host: string, port: i32) -> result[conn_handle, string]
      + returns an outbound connection to the address
      - returns error when the upstream cannot be reached
      # network
    std.net.read
      fn (c: conn_handle, max: i32) -> result[bytes, string]
      + returns up to max bytes read from the connection
      # network
    std.net.write
      fn (c: conn_handle, data: bytes) -> result[i32, string]
      + returns the number of bytes written
      # network
  std.time
    std.time.sleep_millis
      fn (ms: i64) -> void
      + blocks the current worker for the given duration
      # time

toxiproxy
  toxiproxy.new
    fn (listen_host: string, listen_port: i32, upstream_host: string, upstream_port: i32) -> proxy_state
    + returns a proxy configured with the given endpoints and no toxics
    # construction
  toxiproxy.add_toxic_latency
    fn (p: proxy_state, name: string, direction: string, latency_ms: i64, jitter_ms: i64) -> result[proxy_state, string]
    + returns the proxy with a latency toxic attached
    - returns error when direction is not "upstream" or "downstream"
    # toxic
  toxiproxy.add_toxic_bandwidth
    fn (p: proxy_state, name: string, direction: string, rate_kib_per_sec: i64) -> result[proxy_state, string]
    + returns the proxy with a bandwidth cap attached
    - returns error when the rate is zero or negative
    # toxic
  toxiproxy.add_toxic_slow_close
    fn (p: proxy_state, name: string, delay_ms: i64) -> proxy_state
    + returns the proxy that delays connection close
    # toxic
  toxiproxy.remove_toxic
    fn (p: proxy_state, name: string) -> result[proxy_state, string]
    + returns the proxy with the named toxic removed
    - returns error when no toxic with that name exists
    # toxic
  toxiproxy.run
    fn (p: proxy_state) -> result[void, string]
    + accepts connections, dials upstream, and copies bytes through the configured toxics
    - returns error when the listener cannot bind
    # lifecycle
    -> std.net.listen_tcp
    -> std.net.accept
    -> std.net.dial_tcp
    -> std.net.read
    -> std.net.write
    -> std.time.sleep_millis
