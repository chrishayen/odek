# Requirement: "a minimalist in-process message bus"

Topics map to subscriber ids. Publish returns the list of subscribers that should receive the payload; delivery is the caller's responsibility.

std: (all units exist)

bus
  bus.new
    fn () -> bus_state
    + returns an empty bus
    # construction
  bus.subscribe
    fn (state: bus_state, topic: string) -> tuple[i64, bus_state]
    + returns a new subscriber id and updated state
    ? ids are monotonically increasing from 1
    # subscription
  bus.unsubscribe
    fn (state: bus_state, sub_id: i64) -> bus_state
    + removes the subscription and returns new state
    - leaves state unchanged when sub_id is unknown
    # subscription
  bus.publish
    fn (state: bus_state, topic: string, payload: string) -> list[tuple[i64, string]]
    + returns (subscriber_id, payload) pairs for every listener on the topic
    - returns an empty list when no one is subscribed
    # dispatch
