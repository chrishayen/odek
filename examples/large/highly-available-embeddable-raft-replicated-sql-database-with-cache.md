# Requirement: "a highly-available, embeddable, raft-based replicated SQL database with cache"

Raft replication of a SQL log plus a local query cache. Real work is in the raft primitives and log storage; the project surface is small.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.sync
    std.sync.new_rwlock
      fn () -> rwlock_handle
      + returns an unlocked read-write lock
      # concurrency
  std.fs
    std.fs.append_bytes
      fn (path: string, data: bytes) -> result[void, string]
      + appends bytes to a file, creating it if missing
      - returns error on I/O failure
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file into memory
      - returns error when the file does not exist
      # filesystem
  std.encoding
    std.encoding.encode_log_entry
      fn (term: i64, index: i64, payload: bytes) -> bytes
      + serializes a log entry with a length prefix and crc32 checksum
      # serialization
    std.encoding.decode_log_entries
      fn (buf: bytes) -> result[list[tuple[i64, i64, bytes]], string]
      + decodes a sequence of length-prefixed entries
      - returns error on truncated input or bad checksum
      # serialization
  std.raft
    std.raft.new_node
      fn (node_id: i64, peer_ids: list[i64]) -> raft_state
      + initializes a raft node in follower state
      # consensus
    std.raft.step
      fn (state: raft_state, msg: raft_message) -> tuple[raft_state, list[raft_message]]
      + advances the raft state machine with an incoming message and returns outgoing messages
      # consensus
    std.raft.tick
      fn (state: raft_state) -> tuple[raft_state, list[raft_message]]
      + advances election and heartbeat timers by one tick
      # consensus
    std.raft.propose
      fn (state: raft_state, payload: bytes) -> result[raft_state, string]
      + appends a command to the leader's log
      - returns error when the node is not leader
      # consensus

replicated_sql
  replicated_sql.open
    fn (data_dir: string, node_id: i64, peer_ids: list[i64]) -> result[db_state, string]
    + opens or creates a replicated database under data_dir
    - returns error when the data directory cannot be read
    # construction
    -> std.fs.read_all
    -> std.encoding.decode_log_entries
    -> std.raft.new_node
  replicated_sql.execute
    fn (db: db_state, sql: string) -> result[db_state, string]
    + proposes a write statement through raft and applies it on commit
    - returns error when the node is not leader
    # writes
    -> std.raft.propose
    -> std.encoding.encode_log_entry
    -> std.fs.append_bytes
  replicated_sql.query
    fn (db: db_state, sql: string) -> result[list[map[string, string]], string]
    + runs a read query against the local applied state
    - returns error on SQL parse failure
    # reads
  replicated_sql.cache_get
    fn (db: db_state, key: string) -> optional[bytes]
    + returns a cached query result if present and not stale
    - returns none when the cache entry was invalidated by a recent write
    # caching
    -> std.time.now_millis
  replicated_sql.cache_put
    fn (db: db_state, key: string, value: bytes, ttl_ms: i64) -> db_state
    + stores a query result with an expiration timestamp
    # caching
    -> std.time.now_millis
  replicated_sql.apply_committed
    fn (db: db_state) -> db_state
    + applies all raft-committed entries to the local SQL state and invalidates affected cache keys
    # application
  replicated_sql.handle_message
    fn (db: db_state, msg: raft_message) -> tuple[db_state, list[raft_message]]
    + routes an incoming raft message and returns messages to send
    # replication
    -> std.raft.step
  replicated_sql.tick
    fn (db: db_state) -> tuple[db_state, list[raft_message]]
    + advances timers; called periodically by the caller
    # replication
    -> std.raft.tick
  replicated_sql.leader
    fn (db: db_state) -> optional[i64]
    + returns the current leader's node id if known
    # introspection
  replicated_sql.close
    fn (db: db_state) -> result[void, string]
    + flushes pending state and releases resources
    # lifecycle
