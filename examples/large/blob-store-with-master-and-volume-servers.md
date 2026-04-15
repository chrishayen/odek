# Requirement: "a distributed blob store with a master and volume servers"

Clients upload blobs to a master, which assigns them to a volume server keyed by a content id. Each volume holds an append-only data file and an in-memory index mapping ids to offsets for O(1) seek on read.

std
  std.fs
    std.fs.open_append
      fn (path: string) -> result[file_handle, string]
      + opens or creates a file for appending
      # filesystem
    std.fs.append
      fn (f: file_handle, data: bytes) -> result[i64, string]
      + appends bytes and returns the starting offset
      # filesystem
    std.fs.read_at
      fn (path: string, offset: i64, length: i64) -> result[bytes, string]
      + reads length bytes starting at offset
      - returns error when offset+length exceeds file size
      # filesystem
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + returns a 64-bit FNV-1a hash
      # hashing
  std.net
    std.net.rpc_call
      fn (host: string, method: string, payload: bytes) -> result[bytes, string]
      + sends a request to a peer and returns the response body
      - returns error on connection failure or timeout
      # rpc

blob_store
  blob_store.new_volume
    fn (volume_id: i32, data_path: string) -> result[volume_state, string]
    + opens the data file and builds an in-memory index
    # volume_construction
    -> std.fs.open_append
  blob_store.volume_put
    fn (state: volume_state, blob_id: u64, data: bytes) -> result[volume_state, string]
    + appends data to the volume and records (blob_id -> offset, length)
    # volume_write
    -> std.fs.append
  blob_store.volume_get
    fn (state: volume_state, blob_id: u64) -> result[bytes, string]
    + reads a blob by id using a single disk seek
    - returns error when blob_id is unknown
    # volume_read
    -> std.fs.read_at
  blob_store.new_master
    fn () -> master_state
    + returns a master with no registered volumes
    # master_construction
  blob_store.register_volume
    fn (master: master_state, volume_id: i32, host: string) -> master_state
    + records the address of a volume server
    # master_registration
  blob_store.assign
    fn (master: master_state, blob_id: u64) -> result[assignment, string]
    + picks a volume for a blob id via consistent hashing
    - returns error when no volumes are registered
    # master_assignment
    -> std.hash.fnv64
  blob_store.lookup
    fn (master: master_state, blob_id: u64) -> result[assignment, string]
    + returns the volume host responsible for a blob id
    - returns error when no volumes are registered
    # master_lookup
    -> std.hash.fnv64
  blob_store.client_put
    fn (master: master_state, blob_id: u64, data: bytes) -> result[void, string]
    + resolves the volume via the master and uploads the blob
    - returns error when the chosen volume is unreachable
    # client_write
    -> std.net.rpc_call
  blob_store.client_get
    fn (master: master_state, blob_id: u64) -> result[bytes, string]
    + resolves the volume via the master and downloads the blob
    - returns error when no volume holds the blob
    # client_read
    -> std.net.rpc_call
