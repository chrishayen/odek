# Requirement: "a layer 4 load balancer with dynamic configuration loading"

A TCP-level proxy that watches a configuration source and rebuilds its backend pool without dropping existing connections.

std
  std.tcp
    std.tcp.listen
      @ (addr: string) -> result[tcp_listener, string]
      + binds to addr and returns a listener
      - returns error when the address is unavailable
      # network
    std.tcp.accept
      @ (listener: tcp_listener) -> result[tcp_conn, string]
      + blocks until a new client connection arrives
      # network
    std.tcp.dial
      @ (addr: string) -> result[tcp_conn, string]
      + opens a TCP connection to addr
      - returns error when the host is unreachable
      # network
    std.tcp.pipe
      @ (a: tcp_conn, b: tcp_conn) -> result[void, string]
      + copies bytes bidirectionally until either side closes
      # network
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads a file into a string
      - returns error when the file does not exist
      # filesystem
    std.fs.watch
      @ (path: string) -> file_watcher
      + creates a watcher that signals when path changes
      # filesystem
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

l4lb
  l4lb.parse_config
    @ (raw: string) -> result[lb_config, string]
    + parses a configuration blob into a frontend/backend spec
    - returns error on invalid structure
    # config
  l4lb.load_config
    @ (path: string) -> result[lb_config, string]
    + reads and parses a config file
    # config
    -> std.fs.read_all
  l4lb.new_pool
    @ (backends: list[string]) -> backend_pool
    + creates a pool of backend addresses
    # pool
  l4lb.pick_backend
    @ (pool: backend_pool) -> tuple[string, backend_pool]
    + returns the next backend using round-robin and advances state
    - returns an error-sentinel string when the pool is empty
    # selection
  l4lb.mark_unhealthy
    @ (pool: backend_pool, addr: string) -> backend_pool
    + removes a backend from the active rotation
    # health
  l4lb.health_check
    @ (addr: string) -> bool
    + returns true when a trivial TCP dial succeeds
    # health
    -> std.tcp.dial
  l4lb.serve
    @ (cfg: lb_config) -> result[void, string]
    + binds the frontend and forwards each connection to a selected backend
    - returns error when the frontend address cannot be bound
    # proxy
    -> std.tcp.listen
    -> std.tcp.accept
    -> std.tcp.dial
    -> std.tcp.pipe
  l4lb.reload_on_change
    @ (path: string, current: lb_config) -> result[lb_config, string]
    + reparses the config when the file has changed, leaving existing connections intact
    - returns error when the new config is invalid
    # reload
    -> std.fs.watch
    -> std.fs.read_all
