# Requirement: "a compact reactive stream library focused on low overhead"

Minimal surface: emit, map, filter, subscribe. No merging, no time operators — keeps allocation and dispatch paths short.

std: (all units exist)

fast_stream
  fast_stream.create
    @ () -> stream_handle
    + returns a new empty stream handle
    # construction
  fast_stream.emit
    @ (handle: stream_handle, value: string) -> void
    + delivers value to every active subscriber in registration order
    # producer
  fast_stream.map
    @ (handle: stream_handle, f: fn(string) -> string) -> stream_handle
    + returns a derived handle whose values are f applied to values of the upstream
    # operator
  fast_stream.filter
    @ (handle: stream_handle, pred: fn(string) -> bool) -> stream_handle
    + returns a derived handle that forwards only values satisfying pred
    # operator
  fast_stream.subscribe
    @ (handle: stream_handle, observer: fn(string) -> void) -> i64
    + registers observer and returns a subscription id
    # subscription
  fast_stream.unsubscribe
    @ (handle: stream_handle, subscription_id: i64) -> void
    + removes the subscription with the given id
    - no-op when the id is unknown
    # subscription
