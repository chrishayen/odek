# Requirement: "a scalable, geo-replicated, transactional datastore"

A distributed key-value store with raft-based replication and MVCC transactions. std provides general replication and storage primitives.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + returns fnv-1a 64-bit hash of data
      # hashing
  std.net
    std.net.send_rpc
      fn (peer: string, method: string, payload: bytes) -> result[bytes, string]
      + sends an rpc request to a remote peer and returns the response bytes
      - returns error on connection failure or timeout
      # networking
  std.fs
    std.fs.append_record
      fn (path: string, record: bytes) -> result[void, string]
      + appends a length-prefixed record to a write-ahead log file
      - returns error on io failure
      # storage

datastore
  datastore.open
    fn (data_dir: string, node_id: string, peers: list[string]) -> result[datastore_state, string]
    + initializes a node with persistent storage and cluster membership
    - returns error when data directory cannot be opened
    # construction
    -> std.fs.append_record
  datastore.begin_txn
    fn (state: datastore_state) -> txn_handle
    + starts a new transaction with a read timestamp from the current clock
    # transactions
    -> std.time.now_millis
  datastore.txn_get
    fn (state: datastore_state, txn: txn_handle, key: bytes) -> result[optional[bytes], string]
    + reads the most recent value visible at the transaction's read timestamp
    - returns error when the key is locked by a conflicting writer
    # reads
  datastore.txn_put
    fn (state: datastore_state, txn: txn_handle, key: bytes, value: bytes) -> result[txn_handle, string]
    + stages a write in the transaction's intent buffer
    - returns error when another txn holds a write intent on the key
    # writes
  datastore.txn_commit
    fn (state: datastore_state, txn: txn_handle) -> result[datastore_state, string]
    + replicates the commit record through the raft log and acknowledges intents
    - returns error when consensus cannot be reached
    # commit
    -> std.time.now_millis
  datastore.txn_abort
    fn (state: datastore_state, txn: txn_handle) -> datastore_state
    + clears all intents belonging to the transaction
    # abort
  datastore.route_key
    fn (state: datastore_state, key: bytes) -> string
    + returns the node id responsible for the key's range
    ? range partitioning is based on a hash of the key
    # routing
    -> std.hash.fnv64
  datastore.raft_append
    fn (state: datastore_state, entries: list[bytes]) -> result[datastore_state, string]
    + appends entries to the raft log and replicates to followers
    - returns error when this node is not the leader
    # replication
    -> std.net.send_rpc
    -> std.fs.append_record
  datastore.raft_tick
    fn (state: datastore_state) -> datastore_state
    + advances election and heartbeat timers; triggers elections on timeout
    # consensus
    -> std.time.now_millis
  datastore.apply_committed
    fn (state: datastore_state) -> datastore_state
    + applies committed raft entries to the mvcc store
    # state_machine
  datastore.resolve_intents
    fn (state: datastore_state, txn_id: string, committed: bool) -> datastore_state
    + converts write intents into committed values or discards them on abort
    # intent_resolution
