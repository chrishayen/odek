# Requirement: "a distributed relational database that replicates a single-node SQL engine across nodes using a consensus log"

The project layer owns the cluster state, replication log application, and a thin SQL facade. std provides the raft-style consensus primitives, a local SQL store, and monotonic time.

std
  std.time
    std.time.now_monotonic_ms
      @ () -> i64
      + returns a monotonically non-decreasing millisecond counter
      # time
  std.consensus
    std.consensus.new_log
      @ (node_id: string) -> log_state
      + creates an empty replicated log for the given node
      # consensus
    std.consensus.append
      @ (log: log_state, entry: bytes) -> result[u64, string]
      + appends an entry, returning its index, when the node is leader
      - returns error when the node is not leader
      # consensus
    std.consensus.commit_index
      @ (log: log_state) -> u64
      + returns the highest index known committed by a majority
      # consensus
    std.consensus.read_entry
      @ (log: log_state, index: u64) -> result[bytes, string]
      + returns the committed entry at the given index
      - returns error when index is not yet committed
      # consensus
  std.sql
    std.sql.open_local
      @ (path: string) -> result[sql_db, string]
      + opens or creates a local single-node SQL database at path
      - returns error when the path is not writable
      # storage
    std.sql.execute_local
      @ (db: sql_db, stmt: string, params: list[string]) -> result[sql_result, string]
      + executes a parameterized statement on the local database
      - returns error when the statement is rejected
      # storage
    std.sql.classify
      @ (stmt: string) -> string
      + returns "read" for SELECT-class statements and "write" otherwise
      # parsing

cluster
  cluster.new
    @ (node_id: string, data_dir: string) -> cluster_state
    + opens a local database and an empty consensus log
    - returns error when data_dir is not writable
    # construction
    -> std.sql.open_local
    -> std.consensus.new_log
  cluster.join
    @ (c: cluster_state, peer_id: string) -> result[void, string]
    + registers a peer for replication
    - returns error when peer_id is already known
    # membership
  cluster.leave
    @ (c: cluster_state, peer_id: string) -> bool
    + removes a peer and returns true if it existed
    # membership
  cluster.is_leader
    @ (c: cluster_state) -> bool
    + returns true when this node currently holds leadership
    # consensus
  cluster.execute
    @ (c: cluster_state, stmt: string, params: list[string]) -> result[sql_result, string]
    + reads bypass the log; writes are proposed to the log and applied on commit
    - returns error when a write is issued on a non-leader node
    - returns error when the statement is rejected by the local store
    # execution
    -> std.sql.classify
    -> std.consensus.append
    -> std.sql.execute_local
  cluster.apply_committed
    @ (c: cluster_state) -> result[i32, string]
    + applies all newly committed entries to the local store, returning how many were applied
    - returns error when a committed entry fails to apply
    # replication
    -> std.consensus.commit_index
    -> std.consensus.read_entry
    -> std.sql.execute_local
  cluster.snapshot
    @ (c: cluster_state) -> result[bytes, string]
    + returns a serialized snapshot of the local store at the current commit index
    # replication
  cluster.restore
    @ (c: cluster_state, snapshot: bytes, last_index: u64) -> result[void, string]
    + replaces local state with the snapshot and advances apply index to last_index
    - returns error when the snapshot is malformed
    # replication
  cluster.status
    @ (c: cluster_state) -> cluster_status
    + returns node id, role, commit index, applied index, peer count
    # inspection
    -> std.time.now_monotonic_ms
