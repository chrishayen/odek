# Requirement: "a data structures library with containers, sets, lists, stacks, maps, bidimaps, trees, and hash sets"

A collection of generic container types. Each structure is its own small module with construction plus the handful of operations that define it.

std: (all units exist)

list
  list.new
    fn () -> list_state
    + returns an empty list
    # construction
  list.append
    fn (state: list_state, value: string) -> list_state
    + appends value to the end
    # mutation
  list.get
    fn (state: list_state, index: i32) -> result[string, string]
    + returns the value at index
    - returns error when index is out of range
    # access
  list.size
    fn (state: list_state) -> i32
    + returns the current element count
    # inspection

stack
  stack.new
    fn () -> stack_state
    + returns an empty stack
    # construction
  stack.push
    fn (state: stack_state, value: string) -> stack_state
    + pushes value onto the top
    # mutation
  stack.pop
    fn (state: stack_state) -> result[tuple[string, stack_state], string]
    + returns the top value and the new state
    - returns error when the stack is empty
    # mutation

set
  set.new
    fn () -> set_state
    + returns an empty set
    # construction
  set.add
    fn (state: set_state, value: string) -> set_state
    + inserts value if not already present
    # mutation
  set.contains
    fn (state: set_state, value: string) -> bool
    + returns true when value is present
    - returns false when value is absent
    # query
  set.remove
    fn (state: set_state, value: string) -> set_state
    + removes value if present
    # mutation

hashmap
  hashmap.new
    fn () -> hashmap_state
    + returns an empty map
    # construction
  hashmap.put
    fn (state: hashmap_state, key: string, value: string) -> hashmap_state
    + inserts or overwrites the entry
    # mutation
  hashmap.get
    fn (state: hashmap_state, key: string) -> optional[string]
    + returns the value when the key is present
    - returns none when the key is absent
    # query
  hashmap.delete
    fn (state: hashmap_state, key: string) -> hashmap_state
    + removes the entry when present
    # mutation

bidimap
  bidimap.new
    fn () -> bidimap_state
    + returns an empty bidirectional map
    # construction
  bidimap.put
    fn (state: bidimap_state, key: string, value: string) -> bidimap_state
    + inserts the pair, replacing any existing entry for either side
    ? bidirectional uniqueness: both key and value are unique within their side
    # mutation
  bidimap.get_by_key
    fn (state: bidimap_state, key: string) -> optional[string]
    + returns the value mapped from the key
    # query
  bidimap.get_by_value
    fn (state: bidimap_state, value: string) -> optional[string]
    + returns the key mapped from the value
    # query

bst
  bst.new
    fn () -> bst_state
    + returns an empty binary search tree
    # construction
  bst.insert
    fn (state: bst_state, value: i64) -> bst_state
    + inserts value in sorted position
    # mutation
  bst.contains
    fn (state: bst_state, value: i64) -> bool
    + returns true when value is in the tree
    # query
  bst.in_order
    fn (state: bst_state) -> list[i64]
    + returns all values in sorted order
    # traversal
