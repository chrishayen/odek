# Requirement: "a fan-out sum microbenchmark"

The classic skynet benchmark, as a library: a recursive split that fans out, each leaf contributing its index, with a final reduction. Time measurement is a std primitive.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time

skynet
  skynet.sum
    @ (num: i64, size: i64, div: i64) -> i64
    + returns the total sum of leaf indices in [num, num + size)
    + returns num when size is 1
    ? recursive split: when size > 1, subdivide into div chunks and combine their sums
    # computation
  skynet.run
    @ (total_leaves: i64, div: i64) -> benchmark_result
    + measures wall time and returns it alongside the computed total
    + total equals total_leaves * (total_leaves - 1) / 2
    # benchmark
    -> std.time.now_nanos
    -> skynet.sum
