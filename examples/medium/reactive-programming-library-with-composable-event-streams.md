# Requirement: "a reactive programming library with composable event streams"

Push-based streams with map, filter, and merge operators plus subscription. No scheduler magic — callers drive emission.

std: (all units exist)

reactive
  reactive.source
    fn () -> stream
    + returns an empty source stream that can be pushed to
    # construction
  reactive.push
    fn (s: stream, value: string) -> stream
    + delivers value to every current subscriber of s
    # producer
  reactive.map
    fn (s: stream, f: fn(string) -> string) -> stream
    + returns a stream whose values are f applied to values of s
    # operator
  reactive.filter
    fn (s: stream, pred: fn(string) -> bool) -> stream
    + returns a stream that forwards only values where pred returns true
    # operator
  reactive.merge
    fn (a: stream, b: stream) -> stream
    + returns a stream that emits values from both a and b
    # operator
  reactive.subscribe
    fn (s: stream, observer: fn(string) -> void) -> subscription
    + registers observer to receive future values and returns a handle for unsubscribe
    # subscription
  reactive.unsubscribe
    fn (sub: subscription) -> void
    + stops delivery to the observer associated with sub
    # subscription
