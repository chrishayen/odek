# Requirement: "a push notification server cluster"

A clustered push server accepts client subscriptions on topics and fans out published messages to all connected subscribers across nodes. State is in-memory per node; inter-node coordination flows through a thin pluggable transport.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.new_uuid
      @ () -> string
      + returns a fresh 128-bit identifier encoded as hex
      # identity

push_cluster
  push_cluster.new_node
    @ (node_id: string) -> node_state
    + creates an empty node with no subscribers and no peers
    # construction
  push_cluster.register_peer
    @ (state: node_state, peer_id: string, address: string) -> node_state
    + adds a peer so outgoing broadcasts will include it
    - returns unchanged state when peer_id already registered
    # cluster_membership
  push_cluster.subscribe
    @ (state: node_state, client_id: string, topic: string) -> tuple[string, node_state]
    + registers the client on the topic and returns a subscription token
    # subscription
    -> std.id.new_uuid
  push_cluster.unsubscribe
    @ (state: node_state, token: string) -> node_state
    + removes the subscription identified by the token
    - returns unchanged state when token is unknown
    # subscription
  push_cluster.publish_local
    @ (state: node_state, topic: string, payload: bytes) -> list[string]
    + returns the client_ids on this node that should receive the payload
    + records publish time for diagnostics
    # fanout
    -> std.time.now_millis
  push_cluster.encode_broadcast
    @ (topic: string, payload: bytes, origin_node: string) -> bytes
    + encodes a broadcast envelope for transmission to peer nodes
    # cluster_replication
  push_cluster.apply_broadcast
    @ (state: node_state, envelope: bytes) -> result[list[string], string]
    + decodes an incoming envelope and returns locally-matching client_ids
    - returns error when origin_node matches self (loop prevention)
    - returns error when envelope is malformed
    # cluster_replication
