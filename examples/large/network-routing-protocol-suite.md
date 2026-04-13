# Requirement: "a suite of network routing protocols"

A routing information base plus minimal OSPF-style and BGP-style neighbor state machines, producing route updates.

std
  std.net
    std.net.send_packet
      @ (iface: string, dest: string, payload: bytes) -> result[void, string]
      + sends a raw packet out a named interface
      - returns error when the interface is down
      # networking
    std.net.receive_packet
      @ (iface: string) -> result[tuple[string, bytes], string]
      + blocks until a packet arrives on the interface
      # networking
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns the current unix time in seconds
      # time
  std.encoding
    std.encoding.encode_tlv
      @ (type_id: i32, value: bytes) -> bytes
      + encodes a type-length-value record
      # serialization
    std.encoding.decode_tlv
      @ (buf: bytes) -> result[tuple[i32, bytes, i32], string]
      + decodes one TLV and returns type, value, and bytes consumed
      - returns error on truncated input
      # serialization

routing
  routing.new_rib
    @ () -> rib_state
    + creates an empty routing information base
    # construction
  routing.install_route
    @ (rib: rib_state, prefix: string, next_hop: string, metric: i32, protocol: string) -> rib_state
    + installs a route with a source protocol tag
    + replaces an existing route only when the new one has a better metric
    # rib
  routing.remove_route
    @ (rib: rib_state, prefix: string, protocol: string) -> rib_state
    + removes all routes for a prefix sourced by the given protocol
    # rib
  routing.lookup
    @ (rib: rib_state, address: string) -> optional[string]
    + returns the next hop for an address by longest prefix match
    - returns none when no route matches
    # rib
  routing.new_link_state_node
    @ (router_id: string) -> ls_state
    + creates a link-state protocol node
    # link_state
  routing.link_state_advertise
    @ (ls: ls_state, neighbor_id: string, cost: i32) -> ls_state
    + records a direct link to a neighbor with a cost
    # link_state
  routing.link_state_handle_lsa
    @ (ls: ls_state, origin: string, neighbors: list[tuple[string, i32]]) -> ls_state
    + ingests a received LSA and updates the topology database
    # link_state
  routing.link_state_compute_spf
    @ (ls: ls_state) -> list[tuple[string, string, i32]]
    + runs shortest path first and returns (prefix, next_hop, metric) tuples
    # link_state
  routing.new_path_vector_peer
    @ (local_as: i32, peer_as: i32) -> pv_state
    + creates a path-vector neighbor in idle state
    # path_vector
  routing.path_vector_advance
    @ (pv: pv_state, event: string) -> pv_state
    + advances the neighbor state machine on an event like open, keepalive, update
    # path_vector
  routing.path_vector_handle_update
    @ (pv: pv_state, prefix: string, as_path: list[i32], next_hop: string) -> tuple[bool, pv_state]
    + processes an update and returns true when the route should be installed
    - returns (false, unchanged_state) when the as_path contains a loop
    # path_vector
  routing.tick
    @ (rib: rib_state, ls: ls_state, pv: pv_state) -> tuple[rib_state, ls_state, pv_state]
    + runs periodic timers: refresh LSAs, send keepalives, expire dead neighbors
    # maintenance
    -> std.time.now_seconds
  routing.encode_lsa
    @ (origin: string, neighbors: list[tuple[string, i32]]) -> bytes
    + serializes an LSA to bytes for transmission
    # wire
    -> std.encoding.encode_tlv
  routing.decode_lsa
    @ (buf: bytes) -> result[tuple[string, list[tuple[string, i32]]], string]
    + deserializes an LSA
    - returns error on truncated input
    # wire
    -> std.encoding.decode_tlv
