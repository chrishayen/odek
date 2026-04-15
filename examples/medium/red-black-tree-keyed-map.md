# Requirement: "a key-sorted map backed by a red-black tree"

An ordered map keyed on comparable values with logarithmic insertion, lookup, deletion, and in-order iteration.

std: (all units exist)

treemap
  treemap.new
    fn (compare: compare_fn) -> tree_state
    + creates an empty map using the provided key comparator
    # construction
  treemap.insert
    fn (t: tree_state, key: bytes, value: bytes) -> tree_state
    + inserts or overwrites the value for key, maintaining the red-black invariants via rotations and recoloring
    # mutation
  treemap.get
    fn (t: tree_state, key: bytes) -> optional[bytes]
    + returns the value for key when present
    - returns none when key is not in the map
    # query
  treemap.delete
    fn (t: tree_state, key: bytes) -> tree_state
    + removes the entry for key, rebalancing as needed
    ? deletion of a missing key is a no-op
    # mutation
  treemap.len
    fn (t: tree_state) -> i32
    + returns the number of entries
    # query
  treemap.entries
    fn (t: tree_state) -> list[tuple[bytes, bytes]]
    + returns all entries in ascending key order
    # iteration
  treemap.range
    fn (t: tree_state, low: bytes, high: bytes) -> list[tuple[bytes, bytes]]
    + returns entries whose keys fall in [low, high] in ascending order
    # iteration
