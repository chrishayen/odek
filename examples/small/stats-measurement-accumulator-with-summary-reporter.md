# Requirement: "a statistics accumulator that records measurements and reports summary values"

Maintains a running set of numeric samples and exposes count, mean, and standard deviation.

std: (all units exist)

stats
  stats.new
    fn () -> stats_state
    + creates an empty accumulator
    # construction
  stats.record
    fn (state: stats_state, value: f64) -> stats_state
    + returns a new state with the value added
    ? uses Welford's online algorithm so variance stays numerically stable
    # recording
  stats.count
    fn (state: stats_state) -> i64
    + returns the number of recorded samples
    # reporting
  stats.mean
    fn (state: stats_state) -> f64
    + returns the arithmetic mean of recorded samples
    - returns 0.0 when no samples have been recorded
    # reporting
  stats.stddev
    fn (state: stats_state) -> f64
    + returns the sample standard deviation
    - returns 0.0 when fewer than two samples have been recorded
    # reporting
