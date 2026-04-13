# Requirement: "a distributed graph database with node and edge storage, graph traversals, and simple query support"

Nodes and edges live in a pluggable storage layer and are sharded across peers by node id. The project layer handles identity, local CRUD, shard routing, and breadth-first traversal.

std
  std.hash
    std.hash.fnv64
      @ (data: bytes) -> u64
      + computes FNV-64 hash
      # hashing
  std.kv
    std.kv.get
      @ (store: kv_store, key: bytes) -> optional[bytes]
      + returns the stored value when present
      # storage
    std.kv.put
      @ (store: kv_store, key: bytes, value: bytes) -> result[void, string]
      + writes a key/value pair
      # storage
    std.kv.delete
      @ (store: kv_store, key: bytes) -> result[void, string]
      + removes a key
      - returns error when the key does not exist
      # storage
    std.kv.scan_prefix
      @ (store: kv_store, prefix: bytes) -> list[tuple[bytes,bytes]]
      + returns all key/value pairs whose key begins with prefix
      # storage
  std.net
    std.net.rpc_call
      @ (peer: string, method: string, payload: bytes) -> result[bytes, string]
      + invokes a method on a remote peer and returns its response
      - returns error on network failure
      # network
  std.serde
    std.serde.encode_record
      @ (fields: map[string,string]) -> bytes
      + encodes a string field map for storage
      # serialization
    std.serde.decode_record
      @ (data: bytes) -> result[map[string,string], string]
      + decodes a stored field map
      - returns error on corrupt input
      # serialization

graphdb
  graphdb.node_id_new
    @ () -> node_id
    + returns a fresh 16-byte node id
    # identity
  graphdb.edge_key
    @ (from: node_id, label: string, to: node_id) -> bytes
    + returns the canonical storage key for an edge
    # identity
  graphdb.new
    @ (peers: list[string], self_peer: string, local_store: kv_store) -> graphdb_state
    + constructs a graph database node bound to a peer ring
    # construction
  graphdb.route_for
    @ (state: graphdb_state, id: node_id) -> string
    + returns the peer responsible for the given node id
    # routing
    -> std.hash.fnv64
  graphdb.put_node_local
    @ (state: graphdb_state, id: node_id, props: map[string,string]) -> result[void, string]
    + writes a node into the local store
    # local_write
    -> std.serde.encode_record
    -> std.kv.put
  graphdb.get_node_local
    @ (state: graphdb_state, id: node_id) -> optional[map[string,string]]
    + reads a node from the local store
    # local_read
    -> std.kv.get
    -> std.serde.decode_record
  graphdb.put_edge_local
    @ (state: graphdb_state, from: node_id, label: string, to: node_id, props: map[string,string]) -> result[void, string]
    + writes an edge keyed by (from, label, to) locally
    # local_write
    -> std.kv.put
    -> std.serde.encode_record
  graphdb.neighbors_local
    @ (state: graphdb_state, id: node_id, label: string) -> list[node_id]
    + returns all outgoing neighbors on the given label from the local store
    # local_read
    -> std.kv.scan_prefix
  graphdb.put_node
    @ (state: graphdb_state, id: node_id, props: map[string,string]) -> result[void, string]
    + routes to the owning peer and writes there
    - returns error when the owning peer is unreachable
    # routed_write
    -> std.net.rpc_call
  graphdb.get_node
    @ (state: graphdb_state, id: node_id) -> result[optional[map[string,string]], string]
    + routes to the owning peer and reads
    - returns error when the owning peer is unreachable
    # routed_read
    -> std.net.rpc_call
  graphdb.neighbors
    @ (state: graphdb_state, id: node_id, label: string) -> result[list[node_id], string]
    + routes to the owning peer and returns its neighbors list
    - returns error when the owning peer is unreachable
    # routed_read
    -> std.net.rpc_call
  graphdb.bfs
    @ (state: graphdb_state, start: node_id, label: string, max_depth: i32) -> result[list[node_id], string]
    + traverses the graph breadth-first on one edge label and returns visited nodes in order
    - returns error when a required peer is unreachable mid-traversal
    # traversal
  graphdb.delete_node
    @ (state: graphdb_state, id: node_id) -> result[void, string]
    + deletes a node and its outgoing edges on the owning peer
    - returns error when the node does not exist
    # routed_write
    -> std.net.rpc_call
    -> std.kv.delete
