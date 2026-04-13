# Requirement: "a distributed sql database"

Nodes form a raft cluster. The leader applies sql statements in order via a mvcc key/value store that supports snapshot reads and serializable transactions.

std
  std.net
    std.net.rpc_serve
      @ (addr: string, handler: rpc_handler) -> result[server_handle, string]
      + starts an rpc server
      # networking
    std.net.rpc_call
      @ (addr: string, method: string, payload: bytes) -> result[bytes, string]
      + sends an rpc request
      # networking
  std.store
    std.store.put
      @ (db: store_handle, key: bytes, value: bytes) -> result[void, string]
      + writes a key/value pair
      # storage
    std.store.get
      @ (db: store_handle, key: bytes) -> result[optional[bytes], string]
      + reads a key
      # storage
    std.store.range
      @ (db: store_handle, start: bytes, end: bytes) -> list[tuple[bytes, bytes]]
      + scans a half-open key range
      # storage
    std.store.delete
      @ (db: store_handle, key: bytes) -> result[void, string]
      + removes a key
      # storage
  std.time
    std.time.now_millis
      @ () -> i64
      + current unix time in milliseconds
      # time

toydb
  toydb.raft_new_node
    @ (node_id: string, peers: list[string]) -> raft_state
    + returns a follower raft node with the given peers
    # raft
  toydb.raft_tick
    @ (state: raft_state, now_ms: i64) -> raft_state
    + drives election timeouts and heartbeats
    # raft
    -> std.time.now_millis
  toydb.raft_append_entry
    @ (state: raft_state, entry: bytes) -> result[raft_state, string]
    + the leader appends a command entry to its log and replicates it
    - returns error when the node is not the current leader
    # raft
  toydb.raft_handle_rpc
    @ (state: raft_state, msg: raft_msg) -> tuple[raft_state, list[raft_msg]]
    + applies an incoming raft rpc and returns any outgoing messages
    # raft
  toydb.raft_serve
    @ (state: raft_state, addr: string) -> result[server_handle, string]
    + serves raft rpcs on the given address
    # raft
    -> std.net.rpc_serve
  toydb.mvcc_begin
    @ (db: store_handle, mode: txn_mode) -> result[txn_state, string]
    + starts a transaction and assigns it a monotonically increasing timestamp
    # mvcc
  toydb.mvcc_get
    @ (txn: txn_state, key: bytes) -> result[optional[bytes], string]
    + returns the latest visible version for the given key
    # mvcc
    -> std.store.range
  toydb.mvcc_put
    @ (txn: txn_state, key: bytes, value: bytes) -> result[txn_state, string]
    + writes a new version tagged with the transaction timestamp
    - returns error on write/write conflict with a concurrent transaction
    # mvcc
    -> std.store.put
  toydb.mvcc_commit
    @ (txn: txn_state) -> result[void, string]
    + finalizes the transaction and makes its writes visible
    # mvcc
  toydb.mvcc_rollback
    @ (txn: txn_state) -> void
    + discards a transaction's buffered writes
    # mvcc
    -> std.store.delete
  toydb.parse_sql
    @ (sql: string) -> result[sql_stmt, string]
    + parses a single sql statement into an ast
    - returns error on syntax error
    # sql
  toydb.plan_statement
    @ (schema: catalog, stmt: sql_stmt) -> result[plan_node, string]
    + builds a logical plan tree for the statement
    - returns error when a referenced table or column does not exist
    # sql
  toydb.execute_plan
    @ (txn: txn_state, plan: plan_node) -> result[result_set, string]
    + runs a plan and returns rows or a row count
    # sql
  toydb.execute_sql
    @ (cluster: cluster_state, sql: string) -> result[result_set, string]
    + parses, plans, and executes sql against the cluster
    + reads are served from any node; writes are forwarded to the raft leader
    - returns error when the leader is unknown
    # orchestration
    -> std.net.rpc_call
