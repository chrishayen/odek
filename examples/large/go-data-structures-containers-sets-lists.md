# Requirement: "a data structures library with containers, sets, lists, stacks, maps, bidimaps, trees, and hash sets"

A collection of generic container types. Each structure is its own small module with construction plus the handful of operations that define it.

std: (all units exist)

list
  list.new
    @ () -> list_state
    + returns an empty list
    # construction
  list.append
    @ (state: list_state, value: string) -> list_state
    + appends value to the end
    # mutation
  list.get
    @ (state: list_state, index: i32) -> result[string, string]
    + returns the value at index
    - returns error when index is out of range
    # access
  list.size
    @ (state: list_state) -> i32
    + returns the current element count
    # inspection

stack
  stack.new
    @ () -> stack_state
    + returns an empty stack
    # construction
  stack.push
    @ (state: stack_state, value: string) -> stack_state
    + pushes value onto the top
    # mutation
  stack.pop
    @ (state: stack_state) -> result[tuple[string, stack_state], string]
    + returns the top value and the new state
    - returns error when the stack is empty
    # mutation

set
  set.new
    @ () -> set_state
    + returns an empty set
    # construction
  set.add
    @ (state: set_state, value: string) -> set_state
    + inserts value if not already present
    # mutation
  set.contains
    @ (state: set_state, value: string) -> bool
    + returns true when value is present
    - returns false when value is absent
    # query
  set.remove
    @ (state: set_state, value: string) -> set_state
    + removes value if present
    # mutation

hashmap
  hashmap.new
    @ () -> hashmap_state
    + returns an empty map
    # construction
  hashmap.put
    @ (state: hashmap_state, key: string, value: string) -> hashmap_state
    + inserts or overwrites the entry
    # mutation
  hashmap.get
    @ (state: hashmap_state, key: string) -> optional[string]
    + returns the value when the key is present
    - returns none when the key is absent
    # query
  hashmap.delete
    @ (state: hashmap_state, key: string) -> hashmap_state
    + removes the entry when present
    # mutation

bidimap
  bidimap.new
    @ () -> bidimap_state
    + returns an empty bidirectional map
    # construction
  bidimap.put
    @ (state: bidimap_state, key: string, value: string) -> bidimap_state
    + inserts the pair, replacing any existing entry for either side
    ? bidirectional uniqueness: both key and value are unique within their side
    # mutation
  bidimap.get_by_key
    @ (state: bidimap_state, key: string) -> optional[string]
    + returns the value mapped from the key
    # query
  bidimap.get_by_value
    @ (state: bidimap_state, value: string) -> optional[string]
    + returns the key mapped from the value
    # query

bst
  bst.new
    @ () -> bst_state
    + returns an empty binary search tree
    # construction
  bst.insert
    @ (state: bst_state, value: i64) -> bst_state
    + inserts value in sorted position
    # mutation
  bst.contains
    @ (state: bst_state, value: i64) -> bool
    + returns true when value is in the tree
    # query
  bst.in_order
    @ (state: bst_state) -> list[i64]
    + returns all values in sorted order
    # traversal
