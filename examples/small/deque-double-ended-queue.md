# Requirement: "a double-ended queue"

A ring-buffer-backed deque with O(1) amortized operations at both ends.

std: (all units exist)

deque
  deque.new
    @ () -> deque_state
    + returns an empty deque
    # construction
  deque.push_front
    @ (dq: deque_state, value: i64) -> deque_state
    + prepends a value; grows the underlying buffer when full
    # mutation
  deque.push_back
    @ (dq: deque_state, value: i64) -> deque_state
    + appends a value; grows the underlying buffer when full
    # mutation
  deque.pop_front
    @ (dq: deque_state) -> tuple[optional[i64], deque_state]
    + returns the front value and the new deque
    - returns none when the deque is empty
    # mutation
  deque.pop_back
    @ (dq: deque_state) -> tuple[optional[i64], deque_state]
    + returns the back value and the new deque
    - returns none when the deque is empty
    # mutation
  deque.len
    @ (dq: deque_state) -> i32
    + returns the number of elements currently in the deque
    # introspection
