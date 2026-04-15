# Requirement: "a distributed file system that connects computing devices through a shared content-addressed file namespace"

Files are split into content-addressed blocks published on a peer network. Directories are merkle trees of block references.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns sha-256 digest
      # cryptography
  std.encoding
    std.encoding.multibase_encode
      fn (data: bytes) -> string
      + encodes bytes as a self-describing base-n string
      # encoding
    std.encoding.multibase_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes a multibase string
      - returns error on unsupported base
      # encoding
  std.net
    std.net.p2p_connect
      fn (peer_addr: string) -> result[peer_conn, string]
      + dials a peer
      - returns error on failure to connect
      # networking
    std.net.p2p_send
      fn (conn: peer_conn, msg: bytes) -> result[void, string]
      + sends a framed message
      # networking
    std.net.p2p_recv
      fn (conn: peer_conn) -> result[bytes, string]
      + receives the next framed message
      # networking
    std.net.p2p_listen
      fn (addr: string, handler: peer_handler) -> result[server_handle, string]
      + listens for incoming peer connections
      # networking
  std.store
    std.store.put
      fn (db: store_handle, key: bytes, value: bytes) -> result[void, string]
      + writes a key/value
      # storage
    std.store.get
      fn (db: store_handle, key: bytes) -> result[optional[bytes], string]
      + reads a key
      # storage

dht_fs
  dht_fs.new_node
    fn (listen_addr: string, store: store_handle) -> result[node_state, string]
    + returns a node that serves and stores blocks
    # construction
  dht_fs.start
    fn (state: node_state) -> result[node_state, string]
    + begins listening for peer connections
    # lifecycle
    -> std.net.p2p_listen
  dht_fs.connect_peer
    fn (state: node_state, peer_addr: string) -> result[node_state, string]
    + adds a peer to the connected set
    - returns error when the peer refuses the connection
    # peers
    -> std.net.p2p_connect
  dht_fs.block_cid
    fn (data: bytes) -> string
    + returns the content id for a block using sha-256 and multibase
    # addressing
    -> std.crypto.sha256
    -> std.encoding.multibase_encode
  dht_fs.put_block
    fn (state: node_state, data: bytes) -> result[string, string]
    + stores a block locally and announces availability; returns its cid
    # blocks
    -> std.store.put
  dht_fs.get_block
    fn (state: node_state, cid: string) -> result[bytes, string]
    + returns the block, fetching from peers when not present locally
    - returns error when no peer has the block
    # blocks
    -> std.store.get
    -> std.net.p2p_send
    -> std.net.p2p_recv
  dht_fs.chunk_file
    fn (data: bytes, chunk_size: i32) -> list[bytes]
    + splits a file into fixed-size chunks
    # chunking
  dht_fs.add_file
    fn (state: node_state, data: bytes) -> result[string, string]
    + chunks, stores, and builds a merkle root block; returns the root cid
    # files
  dht_fs.read_file
    fn (state: node_state, root_cid: string) -> result[bytes, string]
    + walks the merkle root, fetches chunks, and concatenates them
    - returns error when any chunk cannot be retrieved
    # files
  dht_fs.add_directory
    fn (state: node_state, entries: map[string, string]) -> result[string, string]
    + writes a directory block mapping names to child cids; returns the directory cid
    # directories
  dht_fs.resolve_path
    fn (state: node_state, root_cid: string, path: list[string]) -> result[string, string]
    + walks directory blocks from root_cid following path components; returns the final cid
    - returns error when any path component is missing
    # directories
  dht_fs.publish_name
    fn (state: node_state, name: string, target_cid: string, signing_key: bytes) -> result[void, string]
    + publishes a signed mapping from a mutable name to target_cid on the network
    # naming
