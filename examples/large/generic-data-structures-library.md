# Requirement: "a generic data structures library"

Generic container types parameterized by element type. Each container is a small cluster of construction plus core operations.

std: (all units exist)

vector
  vector.new
    fn () -> vector_state
    + returns an empty vector
    # construction
  vector.push
    fn (state: vector_state, value: bytes) -> vector_state
    + appends value to the end
    # mutation
  vector.at
    fn (state: vector_state, index: i32) -> result[bytes, string]
    + returns the element at index
    - returns error when index is out of range
    # access
  vector.len
    fn (state: vector_state) -> i32
    + returns the element count
    # inspection

deque
  deque.new
    fn () -> deque_state
    + returns an empty double-ended queue
    # construction
  deque.push_front
    fn (state: deque_state, value: bytes) -> deque_state
    + prepends value
    # mutation
  deque.push_back
    fn (state: deque_state, value: bytes) -> deque_state
    + appends value
    # mutation
  deque.pop_front
    fn (state: deque_state) -> result[tuple[bytes, deque_state], string]
    + returns the front value and the new state
    - returns error when empty
    # mutation

hashset
  hashset.new
    fn () -> hashset_state
    + returns an empty set
    # construction
  hashset.add
    fn (state: hashset_state, value: bytes) -> hashset_state
    + inserts value when absent
    # mutation
  hashset.contains
    fn (state: hashset_state, value: bytes) -> bool
    + returns true when present
    # query

ordered_map
  ordered_map.new
    fn () -> ordered_map_state
    + returns an empty map with insertion-order iteration
    # construction
  ordered_map.put
    fn (state: ordered_map_state, key: bytes, value: bytes) -> ordered_map_state
    + inserts or updates, preserving insertion order for new keys
    # mutation
  ordered_map.get
    fn (state: ordered_map_state, key: bytes) -> optional[bytes]
    + returns the value when present
    # query
  ordered_map.keys_in_order
    fn (state: ordered_map_state) -> list[bytes]
    + returns keys in insertion order
    # iteration
