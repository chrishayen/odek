# Requirement: "a statistical data visualization library built on top of a plotting backend"

Computes the numbers a statistical plot needs, then hands a structured description to a pluggable plotting backend. The library itself draws nothing.

std: (all units exist)

statviz
  statviz.summary_stats
    fn (values: list[f64]) -> map[string, f64]
    + returns mean, median, stddev, min, max, q1, q3
    - returns a map of zeros when the input is empty
    # statistics
  statviz.histogram_bins
    fn (values: list[f64], bin_count: i32) -> list[tuple[f64, f64, i32]]
    + returns (lo, hi, count) bins spanning the data range
    - returns empty list when values is empty or bin_count <= 0
    # binning
  statviz.kernel_density
    fn (values: list[f64], sample_points: list[f64], bandwidth: f64) -> list[f64]
    + returns the density estimate at each sample point using a Gaussian kernel
    - returns zeros when values is empty
    # density
  statviz.linear_regression
    fn (xs: list[f64], ys: list[f64]) -> result[tuple[f64, f64], string]
    + returns (slope, intercept) minimizing squared error
    - returns error when xs and ys differ in length
    - returns error when there are fewer than two points
    # regression
  statviz.boxplot_spec
    fn (values: list[f64]) -> map[string, f64]
    + returns lower_whisker, q1, median, q3, upper_whisker
    ? whiskers are clipped to 1.5*IQR past the quartiles
    # plot_spec
  statviz.build_plot
    fn (kind: string, data: list[list[f64]], options: map[string, string]) -> result[plot_spec, string]
    + returns a plot_spec for "histogram", "scatter", "box", or "density"
    - returns error for unknown kinds
    # plot_spec
  statviz.render
    fn (spec: plot_spec, backend: plot_backend) -> result[void, string]
    + forwards the spec to the backend for drawing
    - returns error when the backend rejects the spec
    # render
