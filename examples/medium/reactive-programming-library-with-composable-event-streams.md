# Requirement: "a reactive programming library with composable event streams"

Push-based streams with map, filter, and merge operators plus subscription. No scheduler magic — callers drive emission.

std: (all units exist)

reactive
  reactive.source
    @ () -> stream
    + returns an empty source stream that can be pushed to
    # construction
  reactive.push
    @ (s: stream, value: string) -> stream
    + delivers value to every current subscriber of s
    # producer
  reactive.map
    @ (s: stream, f: fn(string) -> string) -> stream
    + returns a stream whose values are f applied to values of s
    # operator
  reactive.filter
    @ (s: stream, pred: fn(string) -> bool) -> stream
    + returns a stream that forwards only values where pred returns true
    # operator
  reactive.merge
    @ (a: stream, b: stream) -> stream
    + returns a stream that emits values from both a and b
    # operator
  reactive.subscribe
    @ (s: stream, observer: fn(string) -> void) -> subscription
    + registers observer to receive future values and returns a handle for unsubscribe
    # subscription
  reactive.unsubscribe
    @ (sub: subscription) -> void
    + stops delivery to the observer associated with sub
    # subscription
