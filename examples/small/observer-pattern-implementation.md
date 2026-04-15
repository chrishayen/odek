# Requirement: "an implementation of the observer pattern"

A typed event bus: subscribers register by event name and receive notifications when publish is called.

std: (all units exist)

observer
  observer.new
    fn () -> bus_state
    + returns an empty bus with no subscribers
    # construction
  observer.subscribe
    fn (b: bus_state, event: string, handler_id: string) -> bus_state
    + returns a new bus with handler_id added to the event's subscriber list
    ? duplicate handler_id on the same event is a no-op
    # subscription
  observer.unsubscribe
    fn (b: bus_state, event: string, handler_id: string) -> bus_state
    + returns a new bus with the handler removed from the event
    ? unknown event or handler is a no-op
    # subscription
  observer.publish
    fn (b: bus_state, event: string) -> list[string]
    + returns the list of subscribed handler_ids for the event in subscription order
    - returns [] when the event has no subscribers
    ? the caller dispatches; the bus only tracks routing
    # notification
