# Requirement: "an in-process event bus with topic subscriptions and async delivery"

Subscribers register against topics; publishers enqueue events; the bus drains in delivery order.

std: (all units exist)

bus
  bus.new
    @ () -> bus_state
    + creates an empty bus with no subscriptions and no pending events
    # construction
  bus.subscribe
    @ (b: bus_state, topic: string, subscriber_id: string) -> void
    + registers the subscriber against the topic; duplicate registrations are idempotent
    # subscription
  bus.unsubscribe
    @ (b: bus_state, topic: string, subscriber_id: string) -> bool
    + removes the subscription and returns true if it existed
    # subscription
  bus.publish
    @ (b: bus_state, topic: string, payload: bytes) -> i32
    + enqueues the event for every current subscriber of topic, returning the count
    + returns 0 when no subscribers are registered
    # publication
  bus.drain
    @ (b: bus_state) -> list[delivery]
    + returns all pending deliveries in FIFO order and clears the queue
    ? each delivery pairs (subscriber_id, topic, payload)
    # dispatch
