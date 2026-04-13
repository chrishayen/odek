# Requirement: "a framework for building fast and reliable proxy-style network services"

A programmable proxy with routing, filters, upstream pools, and health checking.

std
  std.net
    std.net.listen_tcp
      @ (addr: string) -> result[listener, string]
      + binds a TCP listener on the given address
      - returns error on bind failure
      # network
    std.net.accept
      @ (l: listener) -> result[tcp_conn, string]
      + blocks until a new connection arrives
      # network
    std.net.dial_tcp
      @ (addr: string) -> result[tcp_conn, string]
      + opens a TCP connection to the given address
      # network
    std.net.read
      @ (c: tcp_conn, max: i32) -> result[bytes, string]
      + reads up to max bytes
      - returns error on closed connection
      # network
    std.net.write
      @ (c: tcp_conn, data: bytes) -> result[i32, string]
      + writes bytes and returns the count written
      # network
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

proxyfw
  proxyfw.new
    @ () -> proxy_state
    + creates a proxy with no routes or upstreams
    # construction
  proxyfw.add_upstream
    @ (state: proxy_state, name: string, addrs: list[string]) -> proxy_state
    + registers a named pool of upstream addresses
    # upstream
  proxyfw.add_route
    @ (state: proxy_state, match: route_match, upstream: string) -> proxy_state
    + routes matching requests to a named upstream
    # routing
  proxyfw.add_filter
    @ (state: proxy_state, phase: string, filter: callback) -> proxy_state
    + installs a filter at a named phase (request, response)
    # filters
  proxyfw.pick_upstream
    @ (state: proxy_state, name: string) -> result[string, string]
    + returns the next healthy address from the named pool
    - returns error when the pool is empty or all hosts are unhealthy
    # load_balancing
  proxyfw.health_check
    @ (state: proxy_state, name: string) -> proxy_state
    + probes each address in a pool and updates health flags
    # health
    -> std.net.dial_tcp
    -> std.time.now_millis
  proxyfw.handle_connection
    @ (state: proxy_state, client: tcp_conn) -> result[void, string]
    + reads a request from a client, runs filters, and proxies to an upstream
    - returns error when upstream dial fails
    # request_handling
    -> std.net.read
    -> std.net.write
    -> std.net.dial_tcp
  proxyfw.serve
    @ (state: proxy_state, addr: string) -> result[void, string]
    + listens on addr and dispatches connections in a loop
    # server
    -> std.net.listen_tcp
    -> std.net.accept
  proxyfw.graceful_shutdown
    @ (state: proxy_state, deadline_ms: i64) -> result[void, string]
    + stops accepting new connections and waits for in-flight work up to a deadline
    # lifecycle
    -> std.time.now_millis
