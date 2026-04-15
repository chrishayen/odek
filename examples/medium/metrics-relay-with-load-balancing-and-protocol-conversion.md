# Requirement: "a relay that load-balances monitoring events across backends and converts them to a metrics protocol"

Accepts structured events, distributes them across a pool of downstream sinks, and optionally translates each event into a line-oriented metric record.

std
  std.hash
    std.hash.fnv1a_64
      fn (data: bytes) -> u64
      + computes 64-bit FNV-1a over the input
      # hashing

relay
  relay.new
    fn (backends: list[string]) -> relay_state
    + creates a relay with an ordered list of backend identifiers
    # construction
  relay.pick_backend
    fn (state: relay_state, event_key: string) -> string
    + returns the backend assigned to the given event key using consistent hashing
    # load_balancing
    -> std.hash.fnv1a_64
  relay.forward
    fn (state: relay_state, event: metric_event, sink: event_sink) -> result[void, string]
    + writes the event to the selected backend via the sink callback
    - returns error when the sink rejects the write
    # forwarding
  relay.to_metric_line
    fn (event: metric_event) -> string
    + renders an event as "path value timestamp" in the metric line format
    + replaces spaces in the path with dots
    # conversion
  relay.mark_backend_down
    fn (state: relay_state, backend: string) -> relay_state
    + removes the backend from rotation until it is restored
    # failover
