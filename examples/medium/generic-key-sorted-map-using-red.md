# Requirement: "a key-sorted map backed by a red-black tree"

An ordered map keyed on comparable values with logarithmic insertion, lookup, deletion, and in-order iteration.

std: (all units exist)

treemap
  treemap.new
    @ (compare: compare_fn) -> tree_state
    + creates an empty map using the provided key comparator
    # construction
  treemap.insert
    @ (t: tree_state, key: bytes, value: bytes) -> tree_state
    + inserts or overwrites the value for key, maintaining the red-black invariants via rotations and recoloring
    # mutation
  treemap.get
    @ (t: tree_state, key: bytes) -> optional[bytes]
    + returns the value for key when present
    - returns none when key is not in the map
    # query
  treemap.delete
    @ (t: tree_state, key: bytes) -> tree_state
    + removes the entry for key, rebalancing as needed
    ? deletion of a missing key is a no-op
    # mutation
  treemap.len
    @ (t: tree_state) -> i32
    + returns the number of entries
    # query
  treemap.entries
    @ (t: tree_state) -> list[tuple[bytes, bytes]]
    + returns all entries in ascending key order
    # iteration
  treemap.range
    @ (t: tree_state, low: bytes, high: bytes) -> list[tuple[bytes, bytes]]
    + returns entries whose keys fall in [low, high] in ascending order
    # iteration
