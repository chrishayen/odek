# Requirement: "top-K sketches with a sliding window variant based on the HeavyKeeper algorithm"

A probabilistic heavy-hitter sketch with a fixed-K leaderboard, plus a sliding-window view that decays old counts.

std
  std.hash
    std.hash.fnv64
      @ (data: bytes) -> u64
      + returns the 64-bit FNV-1a hash of data
      # hashing
  std.random
    std.random.next_unit_f64
      @ () -> f64
      + returns a pseudo-random number in [0.0, 1.0)
      # randomness

topk
  topk.new
    @ (k: i32, width: i32, depth: i32, decay: f64) -> result[topk_state, string]
    + creates a sketch with width buckets per row, depth rows, and the given decay base
    - returns error when k, width, or depth is non-positive
    - returns error when decay is not in (0, 1]
    # construction
  topk.observe
    @ (state: topk_state, item: bytes, count: i64) -> topk_state
    + increments the minimum-counter estimate for item
    + probabilistically replaces the colliding bucket owner per HeavyKeeper when count dominates
    + promotes item into the top-K heap when its estimate exceeds the smallest kept
    # update
    -> std.hash.fnv64
    -> std.random.next_unit_f64
  topk.estimate
    @ (state: topk_state, item: bytes) -> i64
    + returns the current estimated count for item
    + returns 0 for an unseen item
    # query
    -> std.hash.fnv64
  topk.list
    @ (state: topk_state) -> list[item_count]
    + returns the current top-K items ordered by estimated count descending
    # query
  topk.new_windowed
    @ (k: i32, width: i32, depth: i32, window_size: i32) -> result[topk_state, string]
    + creates a sliding-window variant that keeps counts from the last window_size buckets
    - returns error when window_size is non-positive
    # construction
  topk.tick
    @ (state: topk_state) -> topk_state
    + advances the sliding window by one bucket, discarding the oldest counts
    # update
