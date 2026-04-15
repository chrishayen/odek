# Requirement: "a property-based testing library with generators, shrinking, and reproducible seeds"

Properties are functions from generated inputs to bool. The runner draws many examples, and when one fails it shrinks the input toward a minimal counterexample. A seeded PRNG makes runs reproducible.

std
  std.random
    std.random.new
      fn (seed: u64) -> rng_state
      + creates a seeded deterministic PRNG
      # random
    std.random.next_u64
      fn (state: rng_state) -> tuple[u64, rng_state]
      + returns the next 64-bit value and the advanced state
      # random
    std.random.next_in_range
      fn (state: rng_state, lo: i64, hi: i64) -> tuple[i64, rng_state]
      + returns a value in [lo, hi] and the advanced state
      # random

proptest
  proptest.int_gen
    fn (lo: i64, hi: i64) -> generator[i64]
    + creates a generator drawing integers in [lo, hi]
    # generator
    -> std.random.next_in_range
  proptest.list_gen
    fn (inner: generator[i64], max_len: i32) -> generator[list[i64]]
    + creates a generator producing lists up to max_len
    # generator
    -> std.random.next_in_range
  proptest.string_gen
    fn (max_len: i32) -> generator[string]
    + creates a generator for ASCII strings up to max_len
    # generator
  proptest.shrink_int
    fn (value: i64) -> list[i64]
    + returns candidate simpler integers (zero, negation, halves)
    # shrinking
  proptest.shrink_list
    fn (value: list[i64]) -> list[list[i64]]
    + returns candidates by removing each element and halving the list
    # shrinking
  proptest.run
    fn (gen: generator[i64], prop: fn(i64) -> bool, count: i32, seed: u64) -> run_result
    + draws count examples; returns ok if all pass
    - returns a failing case when the property returns false
    # runner
    -> std.random.new
  proptest.shrink
    fn (prop: fn(i64) -> bool, failing: i64) -> i64
    + returns the minimal failing value reachable by repeated shrinking
    # minimization
  proptest.check
    fn (gen: generator[i64], prop: fn(i64) -> bool, count: i32, seed: u64) -> check_result
    + runs the property and, on failure, returns the shrunken counterexample and original seed
    # orchestration
  proptest.check_list
    fn (gen: generator[list[i64]], prop: fn(list[i64]) -> bool, count: i32, seed: u64) -> check_result
    + list-valued variant of check with list shrinking
    # orchestration
  proptest.reproduce
    fn (seed: u64, gen: generator[i64], index: i32) -> i64
    + returns the example produced at the given index under the given seed
    ? enables reproducing a failure reported from a prior run
    # reproducibility
    -> std.random.new
