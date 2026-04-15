# Requirement: "a dns proxy that blocks ad domains"

Accepts dns queries, checks each name against a set of blocklists, and either returns a sinkhole answer or forwards to an upstream resolver. Answers are cached by ttl.

std
  std.net
    std.net.udp_listen
      fn (addr: string) -> result[udp_socket, string]
      + binds a udp socket on addr for receiving datagrams
      - returns error when the address is already in use
      # network
    std.net.udp_recv
      fn (sock: udp_socket) -> result[tuple[bytes, string], string]
      + returns the next datagram and its source address
      # network
    std.net.udp_send
      fn (sock: udp_socket, to: string, data: bytes) -> result[void, string]
      + sends a datagram to to
      # network
  std.http
    std.http.get
      fn (url: string, headers: map[string, string]) -> result[http_response, string]
      + fetches a url
      # network
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns unix time in seconds
      # time

dns_proxy
  dns_proxy.decode_query
    fn (packet: bytes) -> result[dns_query, string]
    + parses a dns query packet into transaction id, name, and type
    - returns error on truncated or malformed packets
    # protocol
  dns_proxy.encode_response
    fn (query: dns_query, answer: dns_answer) -> bytes
    + encodes a dns response packet for a query
    # protocol
  dns_proxy.load_blocklist
    fn (source: string) -> result[blocklist, string]
    + loads a hosts-format blocklist from a url or local path
    - returns error when the source is unreachable
    # blocklist
    -> std.http.get
  dns_proxy.merge_blocklists
    fn (lists: list[blocklist]) -> blocklist
    + combines multiple blocklists into a single set
    # blocklist
  dns_proxy.is_blocked
    fn (list: blocklist, name: string) -> bool
    + returns true when name or any parent domain is in the list
    # matching
  dns_proxy.new_cache
    fn (max_entries: i32) -> dns_cache
    + creates an empty lru cache
    # caching
  dns_proxy.cache_lookup
    fn (cache: dns_cache, key: string) -> optional[dns_answer]
    + returns the cached answer if present and not expired
    # caching
    -> std.time.now_seconds
  dns_proxy.cache_store
    fn (cache: dns_cache, key: string, answer: dns_answer, ttl_seconds: i32) -> dns_cache
    + inserts an answer with an expiration
    # caching
    -> std.time.now_seconds
  dns_proxy.resolve_upstream
    fn (upstream_addr: string, query: dns_query) -> result[dns_answer, string]
    + forwards the query to upstream_addr and returns the answer
    - returns error on timeout or malformed response
    # forwarding
    -> std.net.udp_send
    -> std.net.udp_recv
  dns_proxy.sinkhole_answer
    fn (query: dns_query) -> dns_answer
    + returns an A/AAAA answer pointing at 0.0.0.0 / ::
    # blocking
  dns_proxy.handle_packet
    fn (state: dns_proxy_state, packet: bytes, source: string) -> tuple[dns_proxy_state, bytes]
    + decodes the query, decides block vs forward vs cache hit, and returns the response packet
    # request_loop
  dns_proxy.serve
    fn (state: dns_proxy_state, listen_addr: string, upstream: string) -> result[void, string]
    + runs the proxy loop on listen_addr, never returning unless the socket dies
    # server
    -> std.net.udp_listen
    -> std.net.udp_recv
    -> std.net.udp_send
