# Requirement: "a small charting and visualization library with a partial grammar-of-graphics spec"

Given a dataset and a spec describing marks and scales, compute the rendered geometry. The library does not emit pixels; it hands back positions so the caller can draw them.

std: (all units exist)

gust
  gust.new_dataset
    fn (rows: list[map[string, f64]]) -> result[dataset, string]
    + returns a dataset with the given numeric rows
    - returns error when rows have inconsistent keys
    # data
  gust.linear_scale
    fn (domain_min: f64, domain_max: f64, range_min: f64, range_max: f64) -> result[scale, string]
    + returns a linear scale mapping domain to range
    - returns error when domain_min equals domain_max
    # scales
  gust.apply_scale
    fn (s: scale, value: f64) -> f64
    + returns the range value corresponding to value
    + clamps values outside the domain to the range endpoints
    # scales
  gust.bar_marks
    fn (data: dataset, x_field: string, y_field: string, x_scale: scale, y_scale: scale) -> result[list[rect], string]
    + returns one rectangle per row with positions taken from the scaled fields
    - returns error when a field is missing from the dataset
    # marks
    -> gust.apply_scale
  gust.line_marks
    fn (data: dataset, x_field: string, y_field: string, x_scale: scale, y_scale: scale) -> result[list[point], string]
    + returns scaled points in row order suitable for a connected line
    - returns error when a field is missing from the dataset
    # marks
    -> gust.apply_scale
  gust.axis_ticks
    fn (s: scale, count: i32) -> list[tick]
    + returns count evenly spaced ticks between the domain endpoints with their scaled positions
    + returns two ticks for count less than two
    # axes
