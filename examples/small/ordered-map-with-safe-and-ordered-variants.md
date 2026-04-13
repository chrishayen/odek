# Requirement: "a generic map interface with safe and ordered variants"

A single map abstraction with insertion-order iteration and concurrency-safe access.

std: (all units exist)

ordered_map
  ordered_map.new
    @ () -> ordered_map_state
    + creates an empty ordered map
    # construction
  ordered_map.set
    @ (state: ordered_map_state, key: string, value: string) -> ordered_map_state
    + inserts a new key at the end of the iteration order
    + overwrites an existing key without changing its position
    # mutation
  ordered_map.get
    @ (state: ordered_map_state, key: string) -> optional[string]
    + returns the value for an existing key
    - returns none for a missing key
    # lookup
  ordered_map.delete
    @ (state: ordered_map_state, key: string) -> ordered_map_state
    + removes a key and closes the gap in iteration order
    # mutation
  ordered_map.keys_in_order
    @ (state: ordered_map_state) -> list[string]
    + returns keys in the order they were first inserted
    # iteration
  ordered_map.with_lock
    @ (state: ordered_map_state) -> ordered_map_state
    + wraps the map so concurrent set/get/delete operations are serialized
    ? safety is a property of the wrapper, not a duplicate type
    # concurrency
