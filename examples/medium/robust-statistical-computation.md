# Requirement: "a statistical computation library"

Descriptive statistics plus a small set of distribution primitives. Numerically careful summation and variance.

std: (all units exist)

stats
  stats.mean
    @ (xs: list[f64]) -> result[f64, string]
    + returns the arithmetic mean of the sample
    - returns error on empty input
    # descriptive
  stats.variance
    @ (xs: list[f64]) -> result[f64, string]
    + returns the sample variance using Welford's online algorithm
    - returns error on fewer than two values
    # descriptive
  stats.stddev
    @ (xs: list[f64]) -> result[f64, string]
    + returns the square root of the sample variance
    - returns error on fewer than two values
    # descriptive
  stats.quantile
    @ (xs: list[f64], q: f64) -> result[f64, string]
    + returns the linear-interpolation quantile for q in [0,1]
    - returns error when q is outside [0,1]
    - returns error on empty input
    # descriptive
  stats.normal_pdf
    @ (x: f64, mean: f64, stddev: f64) -> result[f64, string]
    + returns the normal probability density at x
    - returns error when stddev is not positive
    # distributions
  stats.normal_cdf
    @ (x: f64, mean: f64, stddev: f64) -> result[f64, string]
    + returns the normal cumulative distribution at x using an erf approximation
    - returns error when stddev is not positive
    # distributions
  stats.pearson_correlation
    @ (xs: list[f64], ys: list[f64]) -> result[f64, string]
    + returns Pearson's r for equal-length samples
    - returns error when lengths differ
    - returns error when either sample has zero variance
    # descriptive
