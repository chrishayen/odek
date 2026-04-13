# Requirement: "a highly-available key-value store for shared configuration and service discovery"

A replicated key-value store with watches for service discovery. Consensus lives in a std primitive.

std
  std.raft
    std.raft.new_cluster
      @ (node_id: i64, peer_ids: list[i64]) -> raft_state
      + creates a raft cluster member in follower state
      # consensus
    std.raft.propose
      @ (state: raft_state, payload: bytes) -> result[raft_state, string]
      + appends a command to the leader log
      - returns error when the node is not leader
      # consensus
    std.raft.apply_committed
      @ (state: raft_state) -> tuple[list[bytes], raft_state]
      + returns newly committed payloads in order
      # consensus
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

kv_store
  kv_store.open
    @ (node_id: i64, peer_ids: list[i64]) -> kv_state
    + opens a replicated key-value store member
    # construction
    -> std.raft.new_cluster
  kv_store.put
    @ (kv: kv_state, key: string, value: bytes) -> result[kv_state, string]
    + replicates a put through consensus and applies it
    - returns error when the node is not leader
    # writes
    -> std.raft.propose
  kv_store.get
    @ (kv: kv_state, key: string) -> optional[bytes]
    + returns the value for a key from locally applied state
    - returns none when the key does not exist
    # reads
  kv_store.delete
    @ (kv: kv_state, key: string) -> result[kv_state, string]
    + replicates a delete and fires watchers
    # writes
    -> std.raft.propose
  kv_store.watch
    @ (kv: kv_state, prefix: string, callback_id: i64) -> kv_state
    + registers a watcher that fires on any key matching the prefix
    # watches
  kv_store.lease_grant
    @ (kv: kv_state, ttl_seconds: i64) -> tuple[i64, kv_state]
    + creates a lease id for service registration entries
    # leases
    -> std.time.now_seconds
  kv_store.tick
    @ (kv: kv_state) -> kv_state
    + expires leases, applies committed entries, and dispatches watch callbacks
    # maintenance
    -> std.raft.apply_committed
    -> std.time.now_seconds
