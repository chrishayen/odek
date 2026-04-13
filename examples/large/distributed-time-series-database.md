# Requirement: "a distributed time-series database"

Shards points by metric name, persists to segment files, replicates to peers, and answers range queries. Coordination, storage, and query planning are separated.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns current unix time in nanoseconds
      # time
  std.hash
    std.hash.fnv64
      @ (data: bytes) -> u64
      + computes 64-bit FNV-1a hash
      # hashing
  std.fs
    std.fs.append_file
      @ (path: string, data: bytes) -> result[void, string]
      + appends data to a file, creating it if missing
      - returns error on io failure
      # filesystem
    std.fs.read_file
      @ (path: string) -> result[bytes, string]
      + reads file contents
      - returns error when file is missing
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + lists directory entries
      # filesystem
  std.encoding
    std.encoding.varint_encode
      @ (v: i64) -> bytes
      + encodes a signed integer as a variable-length sequence
      # encoding
    std.encoding.varint_decode
      @ (data: bytes, offset: i32) -> result[tuple[i64, i32], string]
      + decodes a varint and returns (value, next_offset)
      - returns error on truncated input
      # encoding
  std.net
    std.net.rpc_call
      @ (peer: string, method: string, body: bytes) -> result[bytes, string]
      + sends an rpc request to a peer and returns the response body
      - returns error on connection or protocol failure
      # networking

tsdb
  tsdb.open
    @ (data_dir: string, shard_count: i32) -> result[tsdb_state, string]
    + opens or creates a database under data_dir with the given shard count
    - returns error when data_dir is not writable
    # construction
    -> std.fs.list_dir
  tsdb.shard_for
    @ (metric: string, shard_count: i32) -> i32
    + returns the shard index for a metric name via consistent hashing
    # routing
    -> std.hash.fnv64
  tsdb.write_point
    @ (state: tsdb_state, metric: string, tags: map[string,string], value: f64, ts_nanos: i64) -> result[void, string]
    + appends a point to the metric's active segment
    - returns error when the metric name is empty
    # ingestion
    -> std.encoding.varint_encode
    -> std.fs.append_file
  tsdb.flush_segment
    @ (state: tsdb_state, shard: i32) -> result[void, string]
    + seals the active segment on a shard and starts a new one
    # storage
    -> std.fs.append_file
  tsdb.query_range
    @ (state: tsdb_state, metric: string, start_nanos: i64, end_nanos: i64) -> result[list[tuple[i64,f64]], string]
    + returns points for metric within [start, end)
    - returns error when start >= end
    # query
    -> std.fs.read_file
    -> std.encoding.varint_decode
  tsdb.register_peer
    @ (state: tsdb_state, peer_addr: string) -> tsdb_state
    + adds a peer to the replication ring
    # cluster
  tsdb.replicate_write
    @ (state: tsdb_state, metric: string, payload: bytes) -> result[i32, string]
    + forwards a write to peers owning replicas of the metric's shard, returns ack count
    - returns error when no peers acknowledge
    # replication
    -> std.net.rpc_call
  tsdb.compact
    @ (state: tsdb_state, shard: i32) -> result[i32, string]
    + merges sealed segments on a shard, returns the number of segments merged
    # maintenance
    -> std.fs.read_file
    -> std.fs.append_file
  tsdb.close
    @ (state: tsdb_state) -> result[void, string]
    + flushes active segments and releases handles
    # lifecycle
    -> std.fs.append_file
