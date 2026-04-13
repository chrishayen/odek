# Requirement: "a task time tracking library"

Starts and stops timed entries for named tasks, persists them, and summarizes time spent per task and per day.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
    std.time.start_of_day
      @ (unix_seconds: i64) -> i64
      + returns unix seconds at local midnight preceding the given instant
      # time
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error on missing file
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating or truncating
      # filesystem

time_tracker
  time_tracker.load
    @ (store_path: string) -> result[tracker_state, string]
    + loads recorded entries from the store, or starts empty when the file is absent
    # persistence
    -> std.fs.read_all
  time_tracker.save
    @ (state: tracker_state, store_path: string) -> result[void, string]
    + writes the current entries to the store
    # persistence
    -> std.fs.write_all
  time_tracker.start
    @ (state: tracker_state, task: string) -> result[tracker_state, string]
    + begins a new entry for the given task with now as start time
    - returns error when another entry is still running
    # start
    -> std.time.now_seconds
  time_tracker.stop
    @ (state: tracker_state) -> result[tracker_state, string]
    + closes the currently running entry with now as end time
    - returns error when no entry is running
    # stop
    -> std.time.now_seconds
  time_tracker.total_for_task
    @ (state: tracker_state, task: string) -> i64
    + returns the total seconds recorded for the given task across all entries
    # aggregation
  time_tracker.totals_for_day
    @ (state: tracker_state, day_unix_seconds: i64) -> map[string, i64]
    + returns per-task seconds for entries that overlap the given day
    # aggregation
    -> std.time.start_of_day
