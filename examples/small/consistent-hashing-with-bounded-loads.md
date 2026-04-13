# Requirement: "a consistent hashing ring with bounded loads"

Distributes keys to nodes on a hash ring but caps each node's share so no node exceeds an average-plus-epsilon load.

std
  std.crypto
    std.crypto.fnv64
      @ (data: bytes) -> u64
      + returns a 64-bit fnv hash
      # hashing

bounded_ring
  bounded_ring.new
    @ (epsilon: f64) -> bounded_ring_state
    + creates an empty ring with load factor 1 + epsilon
    # construction
  bounded_ring.add_node
    @ (state: bounded_ring_state, node: string) -> bounded_ring_state
    + places a node on the ring at its hashed position
    # membership
    -> std.crypto.fnv64
  bounded_ring.remove_node
    @ (state: bounded_ring_state, node: string) -> bounded_ring_state
    + removes a node and reassigns its keys on the next lookup
    # membership
  bounded_ring.assign
    @ (state: bounded_ring_state, key: string) -> tuple[string, bounded_ring_state]
    + returns the chosen node, walking clockwise past any node at or above the bound
    - returns empty string when the ring is empty
    # assignment
    -> std.crypto.fnv64
  bounded_ring.release
    @ (state: bounded_ring_state, key: string) -> bounded_ring_state
    + decrements the load for the node that owns key
    # accounting
