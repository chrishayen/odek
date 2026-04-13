# Requirement: "a distributed systems development framework"

A minimal toolkit for building services that talk to each other: service registration, peer discovery, and message dispatch over a pluggable transport.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns the current unix time in milliseconds
      # time
  std.random
    std.random.uuid_v4
      @ () -> string
      + returns a random version-4 UUID string
      # identifiers

framework
  framework.new_node
    @ (node_name: string) -> node_state
    + creates a node with the given human-readable name and a fresh node id
    # construction
    -> std.random.uuid_v4
  framework.register_service
    @ (state: node_state, service_name: string, handler_id: string) -> node_state
    + registers a handler id under a service name on this node
    # registration
  framework.announce
    @ (state: node_state, peer_address: string) -> result[node_state, string]
    + records a peer and its last-seen timestamp
    - returns error when peer_address is empty
    # discovery
    -> std.time.now_millis
  framework.lookup_service
    @ (state: node_state, service_name: string) -> list[string]
    + returns all known handler ids for a service across local and peer registries
    - returns an empty list when the service is unknown
    # discovery
  framework.dispatch
    @ (state: node_state, service_name: string, payload: bytes) -> result[bytes, string]
    + routes payload to a registered handler and returns its response
    - returns error when no handler is registered for the service
    # messaging
  framework.expire_stale_peers
    @ (state: node_state, max_age_millis: i64) -> node_state
    + drops peers whose last-seen timestamp is older than max_age_millis
    # discovery
    -> std.time.now_millis
