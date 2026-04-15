# Requirement: "a network simulation and testing library"

A virtual network with nodes, links, packets, and a clock-driven scheduler for deterministic tests.

std
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns a monotonic nanosecond timestamp
      # time
  std.random
    std.random.new
      fn (seed: u64) -> rng_state
      + creates a deterministic pseudo-random source from a seed
      # randomness
    std.random.next_f64
      fn (state: rng_state) -> tuple[f64, rng_state]
      + returns a uniform sample in [0, 1) and the advanced state
      # randomness

netsim
  netsim.new_network
    fn () -> network_state
    + returns an empty network with no nodes or links
    # construction
  netsim.add_node
    fn (net: network_state, name: string) -> result[network_state, string]
    + registers a node with a unique name
    - returns error when a node with the same name already exists
    # topology
  netsim.add_link
    fn (net: network_state, from: string, to: string, latency_ms: i32, loss: f64) -> result[network_state, string]
    + registers a directed link with latency and loss probability
    - returns error when either endpoint is unknown
    - returns error when loss is outside [0, 1]
    # topology
  netsim.send
    fn (net: network_state, from: string, to: string, payload: bytes, now_ms: i64) -> result[network_state, string]
    + schedules a packet along the shortest path with accumulated latency
    - returns error when no path exists between endpoints
    # transmission
  netsim.tick
    fn (net: network_state, now_ms: i64) -> tuple[list[delivered_packet], network_state]
    + returns packets whose scheduled delivery time is at or before now_ms
    + drops packets sampled as lost per link loss probability
    # simulation
    -> std.random.next_f64
  netsim.partition
    fn (net: network_state, group_a: list[string], group_b: list[string]) -> network_state
    + disables all links that cross the two groups
    # fault_injection
  netsim.heal
    fn (net: network_state) -> network_state
    + re-enables all disabled links
    # fault_injection
  netsim.inbox
    fn (net: network_state, node: string) -> list[delivered_packet]
    + returns packets already delivered to a node
    - returns empty list when node is unknown
    # inspection
  netsim.pending_count
    fn (net: network_state) -> i32
    + returns the number of packets still in flight
    # inspection
  netsim.seed
    fn (net: network_state, seed: u64) -> network_state
    + replaces the loss sampling rng with a fresh seeded one
    # determinism
    -> std.random.new
