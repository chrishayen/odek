# Requirement: "a consistent hashing ring with configurable virtual replicas per node"

Each real node occupies multiple positions on the ring to smooth the key distribution.

std
  std.crypto
    std.crypto.fnv64
      fn (data: bytes) -> u64
      + returns a 64-bit fnv hash
      # hashing

replicated_ring
  replicated_ring.new
    fn (replicas_per_node: i32) -> replicated_ring_state
    + creates an empty ring configured with a replica count
    # construction
  replicated_ring.add_node
    fn (state: replicated_ring_state, node: string) -> replicated_ring_state
    + inserts replica_count virtual positions for node
    # membership
    -> std.crypto.fnv64
  replicated_ring.remove_node
    fn (state: replicated_ring_state, node: string) -> replicated_ring_state
    + removes all virtual positions for node
    # membership
  replicated_ring.get
    fn (state: replicated_ring_state, key: string) -> optional[string]
    + returns the real node responsible for key
    - returns none when the ring is empty
    # lookup
    -> std.crypto.fnv64
  replicated_ring.get_n
    fn (state: replicated_ring_state, key: string, n: i32) -> list[string]
    + returns up to n distinct real nodes walking clockwise
    # lookup
    -> std.crypto.fnv64
