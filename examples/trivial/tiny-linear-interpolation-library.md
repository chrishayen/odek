# Requirement: "a tiny linear interpolation library"

Interpolates y-values along a piecewise linear curve defined by sorted (x, y) sample points.

std: (all units exist)

interpolate
  interpolate.at
    @ (xs: list[f64], ys: list[f64], x: f64) -> f64
    + returns the linearly interpolated y for x between two sample points
    + returns the endpoint y when x falls outside the sample range (clamped)
    - returns 0.0 when xs and ys have different lengths or fewer than two points
    ? xs must be sorted ascending; caller's responsibility
    # interpolation
