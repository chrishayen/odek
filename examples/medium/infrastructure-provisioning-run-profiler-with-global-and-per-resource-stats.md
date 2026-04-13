# Requirement: "a profiler for infrastructure provisioning runs that generates global and per-resource stats"

Parses a provisioning run log into timed events and computes aggregate statistics.

std
  std.text
    std.text.split_lines
      @ (raw: string) -> list[string]
      + splits on LF, trimming trailing CR
      # text
  std.time
    std.time.parse_iso8601
      @ (s: string) -> result[i64, string]
      + returns unix milliseconds for the parsed timestamp
      - returns error on malformed input
      # time

run_profiler
  run_profiler.parse_log
    @ (raw: string) -> result[list[event], string]
    + returns timestamped events with resource id and lifecycle phase
    - returns error when a line is missing a timestamp
    # parsing
    -> std.text.split_lines
    -> std.time.parse_iso8601
  run_profiler.resource_timings
    @ (events: list[event]) -> map[string, i64]
    + returns per-resource elapsed milliseconds from creation_start to creation_end
    # per_resource_stats
  run_profiler.global_stats
    @ (events: list[event]) -> global_stats
    + returns total duration, number of resources, and wall-clock overlap
    # global_stats
  run_profiler.slowest
    @ (events: list[event], k: i32) -> list[tuple[string, i64]]
    + returns the k resources with largest elapsed time
    # ranking
  run_profiler.render_histogram
    @ (events: list[event], buckets: i32) -> list[tuple[i64, i64, i64]]
    + returns (bucket_start_ms, bucket_end_ms, count) tuples
    # visualization
  run_profiler.render_gantt_rows
    @ (events: list[event]) -> list[tuple[string, i64, i64]]
    + returns (resource_id, start_ms, end_ms) rows for gantt rendering
    # visualization
