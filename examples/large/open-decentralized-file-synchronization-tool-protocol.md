# Requirement: "a decentralized file synchronization protocol and library"

Devices advertise folder indexes, reconcile differences, and exchange blocks content-addressed by hash.

std
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest of data
      # cryptography
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      - returns error when the file is missing
      # filesystem
    std.fs.write_file
      @ (path: string, data: bytes) -> result[void, string]
      + writes data creating parent directories
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[file_entry], string]
      + recursively lists files with size and modification time
      # filesystem
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

sync
  sync.scan_folder
    @ (root: string, block_size: i32) -> result[folder_index, string]
    + returns an index of every file split into content-addressed blocks
    # indexing
    -> std.fs.walk
    -> std.fs.read_all
    -> std.crypto.sha256
  sync.diff_indexes
    @ (local: folder_index, remote: folder_index) -> index_diff
    + returns files to create, update, delete, and blocks to request
    ? newer modification time wins on conflict
    # reconciliation
    -> std.time.now_seconds
  sync.serialize_index
    @ (index: folder_index) -> bytes
    + encodes an index for transmission over the wire
    # protocol
  sync.parse_index
    @ (raw: bytes) -> result[folder_index, string]
    - returns error on truncated input
    - returns error on version mismatch
    # protocol
  sync.request_block
    @ (peer: peer_id, hash: string) -> result[bytes, string]
    + returns the block bytes when the peer has it
    - returns error when the peer lacks or refuses the block
    # block_exchange
  sync.verify_block
    @ (hash: string, data: bytes) -> bool
    + returns true when the SHA-256 of data matches hash
    # integrity
    -> std.crypto.sha256
  sync.apply_diff
    @ (root: string, diff: index_diff, fetch: block_fetcher) -> result[i32, string]
    + applies the diff to the local folder and returns files changed
    - returns error when a fetched block fails verification
    # application
    -> std.fs.write_file
  sync.announce
    @ (device: device_id, index: folder_index) -> bytes
    + returns an announcement message a peer can broadcast
    # discovery
  sync.handle_announce
    @ (peers: peer_table, msg: bytes) -> result[peer_table, string]
    + returns an updated peer table after ingesting an announcement
    - returns error on invalid message
    # discovery
  sync.authorize_peer
    @ (table: peer_table, device: device_id) -> bool
    + returns true when the device is in the authorized list
    # authorization
