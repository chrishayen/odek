# Requirement: "a library that distributes hexagonal geo cells across virtual nodes using consistent hashing"

Cells are opaque 64-bit identifiers. The ring maps cells to owner nodes via virtual-node replicas for balance.

std
  std.hash
    std.hash.fnv64
      @ (data: bytes) -> u64
      + returns the FNV-1a 64-bit hash of the input
      # hashing

cell_ring
  cell_ring.new
    @ (replicas: i32) -> ring_state
    + creates an empty ring with the given virtual-node count per real node
    # construction
  cell_ring.add_node
    @ (state: ring_state, node_id: string) -> ring_state
    + inserts replicas virtual nodes onto the ring for node_id
    + replicas are placed at hash("node_id#i") for i in [0, replicas)
    # membership
    -> std.hash.fnv64
  cell_ring.remove_node
    @ (state: ring_state, node_id: string) -> ring_state
    + removes all virtual nodes belonging to node_id
    - leaves the ring unchanged when node_id is absent
    # membership
  cell_ring.owner_of
    @ (state: ring_state, cell: u64) -> result[string, string]
    + returns the node id that owns the cell via next-clockwise virtual node
    - returns error when the ring is empty
    # routing
    -> std.hash.fnv64
  cell_ring.owners_of
    @ (state: ring_state, cell: u64, count: i32) -> list[string]
    + returns the next count distinct owners clockwise from the cell position
    + result length may be less than count when fewer nodes exist
    # replication
    -> std.hash.fnv64
  cell_ring.rebalance_report
    @ (state: ring_state, cells: list[u64]) -> map[string, i32]
    + returns a count of cells assigned to each node for the given cell set
    # introspection
    -> std.hash.fnv64
