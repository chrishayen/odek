# Requirement: "a library that notifies subscribers about pluggable activity events on a machine"

Users register activity sources, subscribe handlers by event type, and drive the system by pushing events through the library.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

activity_bus
  activity_bus.new
    @ () -> bus_state
    + returns an empty activity bus
    # construction
  activity_bus.register_source
    @ (bus: bus_state, name: string) -> result[bus_state, string]
    + registers a named activity source
    - returns error when the name is already registered
    # registration
  activity_bus.subscribe
    @ (bus: bus_state, event_type: string, handler: function[activity_event, void]) -> bus_state
    + adds a handler for events of the given type
    # subscription
  activity_bus.publish
    @ (bus: bus_state, source: string, event_type: string, payload: map[string, string]) -> result[bus_state, string]
    + dispatches an event with a timestamp to every subscriber for its type
    - returns error when source is not registered
    # dispatch
    -> std.time.now_millis
  activity_bus.recent
    @ (bus: bus_state, limit: i32) -> list[activity_event]
    + returns up to limit most recent events, newest first
    # inspection
