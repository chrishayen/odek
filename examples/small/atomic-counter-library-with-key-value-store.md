# Requirement: "a library that uses a key-value store as an atomic counter"

Treats a named row in a pluggable store as a counter. Updates go through the store's conditional-update primitive so increments never race.

std
  std.kv
    std.kv.atomic_add
      @ (table: string, key: string, delta: i64) -> result[i64, string]
      + atomically adds delta to the numeric value at key and returns the new value
      + creates the row with initial value 0 + delta when missing
      - returns error when the stored value is not numeric
      # storage
    std.kv.get_number
      @ (table: string, key: string) -> result[i64, string]
      + returns the current numeric value at key
      - returns error when the row does not exist
      # storage

atomic_counter
  atomic_counter.increment
    @ (table: string, key: string) -> result[i64, string]
    + adds 1 and returns the post-update value
    # counting
    -> std.kv.atomic_add
  atomic_counter.add
    @ (table: string, key: string, delta: i64) -> result[i64, string]
    + adds an arbitrary positive or negative delta
    - returns error on store failure
    # counting
    -> std.kv.atomic_add
  atomic_counter.read
    @ (table: string, key: string) -> result[i64, string]
    + returns the current value without modifying it
    - returns error when the counter does not exist
    # counting
    -> std.kv.get_number
