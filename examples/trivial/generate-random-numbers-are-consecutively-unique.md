# Requirement: "a library that generates random integers within a range such that no value is equal to the one immediately before it"

Returns a stateful generator whose successive outputs are never equal.

std
  std.random
    std.random.int_in_range
      @ (lo: i32, hi: i32) -> i32
      + returns a uniformly distributed integer in [lo, hi]
      # randomness

unique_random
  unique_random.next
    @ (state: unique_random_state, lo: i32, hi: i32) -> tuple[i32, unique_random_state]
    + returns a value in [lo, hi] that differs from the previous value
    + on first call may return any value in the range
    ? when lo equals hi the function returns that single value and does not loop
    # rng
    -> std.random.int_in_range
