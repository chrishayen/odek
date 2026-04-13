# Requirement: "an embeddable in-memory key-value store with clustered replication and a RESP-style wire protocol"

Multiple data types, a replication peer layer, and a text-line network protocol.

std
  std.net
    std.net.tcp_listen
      @ (host: string, port: u16) -> result[listener_state, string]
      + binds and listens on the address
      # networking
    std.net.tcp_accept
      @ (listener: listener_state) -> result[conn_state, string]
      + blocks for the next connection
      # networking
    std.net.tcp_dial
      @ (host: string, port: u16) -> result[conn_state, string]
      + opens an outbound connection
      # networking
    std.net.read_bytes
      @ (conn: conn_state, n: i32) -> result[bytes, string]
      + reads exactly n bytes
      # networking
    std.net.write
      @ (conn: conn_state, data: bytes) -> result[void, string]
      + writes raw bytes to the connection
      # networking
  std.hash
    std.hash.crc32
      @ (data: bytes) -> u32
      + computes CRC32 of the input
      # hashing
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

kvstore
  kvstore.new
    @ () -> kv_state
    + creates an empty store supporting strings, lists, hashes, sets, and sorted sets
    # construction
  kvstore.set_string
    @ (state: kv_state, key: string, value: string, ttl_millis: optional[i64]) -> kv_state
    + stores a string value with an optional time-to-live
    # strings
    -> std.time.now_millis
  kvstore.get_string
    @ (state: kv_state, key: string) -> optional[string]
    + returns the value if the key exists and has not expired
    # strings
    -> std.time.now_millis
  kvstore.list_push
    @ (state: kv_state, key: string, value: string, at_head: bool) -> kv_state
    + prepends or appends to the list at key, creating it when missing
    # lists
  kvstore.list_pop
    @ (state: kv_state, key: string, from_head: bool) -> tuple[optional[string], kv_state]
    + removes and returns the head or tail element
    # lists
  kvstore.hash_set
    @ (state: kv_state, key: string, field: string, value: string) -> kv_state
    + sets a field inside the hash at key
    # hashes
  kvstore.hash_get
    @ (state: kv_state, key: string, field: string) -> optional[string]
    + returns the field value if present
    # hashes
  kvstore.set_add
    @ (state: kv_state, key: string, member: string) -> kv_state
    + adds a member to the set at key
    # sets
  kvstore.zset_add
    @ (state: kv_state, key: string, member: string, score: f64) -> kv_state
    + inserts or updates a scored member in the sorted set
    # sorted_sets
  kvstore.zset_range_by_score
    @ (state: kv_state, key: string, min_score: f64, max_score: f64) -> list[string]
    + returns members whose score lies in the inclusive range, sorted ascending
    # sorted_sets
  kvstore.delete
    @ (state: kv_state, key: string) -> kv_state
    + removes the key regardless of type
    # mutation
  kvstore.expire_sweep
    @ (state: kv_state) -> kv_state
    + removes keys whose ttl has elapsed
    # ttl
    -> std.time.now_millis
  kvstore.parse_command
    @ (raw: bytes) -> result[kv_command, string]
    + decodes a RESP-style array of bulk strings into a command
    - returns error on malformed frames
    # protocol
    -> std.net.read_bytes
  kvstore.encode_response
    @ (value: kv_response) -> bytes
    + serializes a response value in RESP form
    # protocol
  kvstore.handle_command
    @ (state: kv_state, cmd: kv_command) -> tuple[kv_response, kv_state]
    + dispatches the parsed command to the appropriate store operation
    # dispatch
  kvstore.start_server
    @ (state: kv_state, host: string, port: u16) -> result[server_state, string]
    + begins accepting client connections on the given address
    # lifecycle
    -> std.net.tcp_listen
  kvstore.join_cluster
    @ (state: kv_state, seed_host: string, seed_port: u16) -> result[cluster_state, string]
    + connects to a seed peer and exchanges membership information
    # clustering
    -> std.net.tcp_dial
    -> std.net.write
    -> std.hash.crc32
  kvstore.replicate_write
    @ (cluster: cluster_state, cmd: kv_command) -> result[void, string]
    + forwards a mutation command to every peer
    # replication
    -> std.net.write
