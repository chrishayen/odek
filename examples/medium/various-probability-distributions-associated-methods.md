# Requirement: "a library for probability distributions and their associated methods"

A handful of common continuous distributions, each with pdf, cdf, and sampling.

std
  std.math
    std.math.erf
      @ (x: f64) -> f64
      + returns the error function of x
      # math
    std.math.exp
      @ (x: f64) -> f64
      + returns e raised to x
      # math
    std.math.ln
      @ (x: f64) -> f64
      + returns the natural logarithm of x
      - returns -infinity for x <= 0
      # math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the square root of x
      # math
  std.rand
    std.rand.next_uniform
      @ () -> f64
      + returns a uniformly random f64 in [0, 1)
      # randomness

distributions
  distributions.normal_pdf
    @ (mean: f64, stddev: f64, x: f64) -> f64
    + returns the Gaussian density at x
    # distribution
    -> std.math.exp
    -> std.math.sqrt
  distributions.normal_cdf
    @ (mean: f64, stddev: f64, x: f64) -> f64
    + returns the Gaussian cumulative probability at x
    # distribution
    -> std.math.erf
    -> std.math.sqrt
  distributions.normal_sample
    @ (mean: f64, stddev: f64) -> f64
    + returns one sample drawn from the given normal distribution
    ? uses Box-Muller transform on two uniform draws
    # sampling
    -> std.rand.next_uniform
    -> std.math.ln
    -> std.math.sqrt
  distributions.exponential_pdf
    @ (rate: f64, x: f64) -> f64
    + returns the exponential density at x
    + returns 0 when x is negative
    # distribution
    -> std.math.exp
  distributions.exponential_cdf
    @ (rate: f64, x: f64) -> f64
    + returns 1 - exp(-rate * x) for x >= 0
    + returns 0 for negative x
    # distribution
    -> std.math.exp
  distributions.exponential_sample
    @ (rate: f64) -> f64
    + returns one sample drawn from the exponential distribution with the given rate
    # sampling
    -> std.rand.next_uniform
    -> std.math.ln
