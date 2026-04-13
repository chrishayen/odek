# Requirement: "an event emitter with wildcard topics, predicate filters, and cancellation"

Subscribers match by topic patterns and optional predicates; each subscription can be cancelled.

std: (all units exist)

emitter
  emitter.new
    @ () -> emitter_state
    + creates an empty emitter with no subscribers
    # construction
  emitter.on
    @ (state: emitter_state, pattern: string) -> tuple[emitter_state, subscription_id]
    + registers a subscription for topics matching the pattern and returns its id
    ? pattern supports "*" for one segment and "**" for any tail
    # subscription
  emitter.on_where
    @ (state: emitter_state, pattern: string, predicate: fn(string, bytes) -> bool) -> tuple[emitter_state, subscription_id]
    + registers a subscription that only fires when the predicate returns true
    # subscription
  emitter.off
    @ (state: emitter_state, id: subscription_id) -> emitter_state
    + cancels a subscription; subsequent emits skip it
    + no-op when the id is not present
    # cancellation
  emitter.emit
    @ (state: emitter_state, topic: string, payload: bytes) -> list[delivery]
    + returns the deliveries for every subscription whose pattern matches the topic and whose predicate accepts the payload
    - returns an empty list when no subscription matches
    # dispatch
  emitter.match_pattern
    @ (pattern: string, topic: string) -> bool
    + returns true when pattern matches topic using "*" and "**" wildcards
    - returns false when segment counts differ and pattern has no "**"
    # matching
