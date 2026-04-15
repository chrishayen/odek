# Requirement: "a HyperLogLog cardinality estimator with sparse representation and bias correction"

Small sets use a sparse representation; once they exceed a threshold they are converted to a dense register array. Estimation uses LogLog-Beta bias correction.

std
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + returns the 64-bit FNV-1a hash of data
      # hashing
  std.math
    std.math.log
      fn (x: f64) -> f64
      + returns the natural log of x
      # math

hll
  hll.new
    fn (precision: i32) -> hll_state
    + creates a sketch with 2^precision registers; starts in sparse mode
    - errors when precision is outside [4, 18]
    # construction
  hll.add
    fn (state: hll_state, item: bytes) -> hll_state
    + hashes the item and updates the sketch, converting to dense if sparse size exceeds threshold
    # insertion
    -> std.hash.fnv64
  hll.estimate
    fn (state: hll_state) -> f64
    + returns the estimated cardinality using LogLog-Beta bias correction
    # estimation
    -> std.math.log
  hll.merge
    fn (a: hll_state, b: hll_state) -> result[hll_state, string]
    + returns a sketch representing the union of a and b
    - returns error when precisions differ
    # merge
