# Requirement: "a library for visualizing a repository's commit history as a heatmap of activity by day"

Given a list of commits with timestamps, bucket them by day, compute a heatmap grid, and render it as a textual calendar.

std
  std.time
    std.time.day_of_week
      @ (unix_seconds: i64) -> i32
      + returns 0 for Sunday through 6 for Saturday
      # time
    std.time.start_of_day
      @ (unix_seconds: i64) -> i64
      + returns the unix timestamp at local midnight for that day
      # time
    std.time.days_between
      @ (from_unix: i64, to_unix: i64) -> i32
      + returns the number of whole days between two timestamps
      # time

commit_heatmap
  commit_heatmap.bucket_by_day
    @ (commit_times: list[i64]) -> map[i64, i32]
    + returns a map from day-start timestamps to commit counts
    # aggregation
    -> std.time.start_of_day
  commit_heatmap.build_grid
    @ (buckets: map[i64, i32], start: i64, end: i64) -> heatmap_grid
    + produces a grid with one column per week and one row per weekday
    + cells outside the range are marked empty
    # grid
    -> std.time.days_between
    -> std.time.day_of_week
  commit_heatmap.intensity
    @ (count: i32, max_count: i32) -> i32
    + returns 0 through 4 for a five-level intensity scale
    # rendering
  commit_heatmap.render_text
    @ (grid: heatmap_grid) -> string
    + renders the grid using space, period, colon, plus, and hash for intensities 0..4
    + rows are separated by newlines
    # rendering
  commit_heatmap.top_days
    @ (buckets: map[i64, i32], limit: i32) -> list[i64]
    + returns the day timestamps with the highest counts, up to limit entries
    # analysis
