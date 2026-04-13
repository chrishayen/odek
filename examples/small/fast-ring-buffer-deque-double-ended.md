# Requirement: "a ring-buffer deque (double-ended queue)"

A growable ring buffer storing generic byte payloads; indices wrap modulo the backing capacity.

std: (all units exist)

deque
  deque.new
    @ (initial_capacity: i32) -> deque_state
    + creates an empty deque with the given initial capacity
    ? capacity is rounded up to a power of two
    # construction
  deque.push_front
    @ (state: deque_state, item: bytes) -> deque_state
    + prepends an item, growing the ring if full
    # push
  deque.push_back
    @ (state: deque_state, item: bytes) -> deque_state
    + appends an item, growing the ring if full
    # push
  deque.pop_front
    @ (state: deque_state) -> result[tuple[bytes, deque_state], string]
    + returns the first item and a state without it
    - returns error when the deque is empty
    # pop
  deque.pop_back
    @ (state: deque_state) -> result[tuple[bytes, deque_state], string]
    + returns the last item and a state without it
    - returns error when the deque is empty
    # pop
  deque.len
    @ (state: deque_state) -> i32
    + returns the number of stored items
    # inspection
