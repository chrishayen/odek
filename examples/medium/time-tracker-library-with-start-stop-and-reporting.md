# Requirement: "a time tracking library for starting, stopping, and reporting on task timers"

Generalized from a desktop app. The library manages task timers and produces reports; the UI is the caller's concern.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file's full contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes contents, replacing any existing file
      # filesystem
  std.json
    std.json.encode
      fn (value: json_value) -> string
      + serializes a dynamic value to a JSON string
      # serialization
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON string into a dynamic value
      - returns error on malformed JSON
      # serialization

time_tracker
  time_tracker.new
    fn () -> tracker_state
    + creates an empty tracker with no tasks and no active timer
    # construction
  time_tracker.add_task
    fn (state: tracker_state, name: string) -> tuple[tracker_state, i32]
    + creates a task with the given name and returns its id
    # tasks
  time_tracker.start
    fn (state: tracker_state, task_id: i32) -> result[tracker_state, string]
    + begins tracking time for the task, stopping any previously active one
    - returns error when the task id does not exist
    # timing
    -> std.time.now_seconds
  time_tracker.stop
    fn (state: tracker_state) -> result[tracker_state, string]
    + ends the active timer and appends the interval to the task
    - returns error when no timer is active
    # timing
    -> std.time.now_seconds
  time_tracker.total_for_task
    fn (state: tracker_state, task_id: i32) -> result[i64, string]
    + returns accumulated seconds for the task, including any active interval
    - returns error when the task id does not exist
    # reporting
    -> std.time.now_seconds
  time_tracker.report_between
    fn (state: tracker_state, from_unix: i64, to_unix: i64) -> list[task_report]
    + returns total seconds per task for intervals overlapping the window
    # reporting
  time_tracker.save
    fn (state: tracker_state, path: string) -> result[void, string]
    + writes the tracker state to disk as JSON
    # persistence
    -> std.json.encode
    -> std.fs.write_all
  time_tracker.load
    fn (path: string) -> result[tracker_state, string]
    + reads a tracker state from disk
    - returns error when the file is missing or malformed
    # persistence
    -> std.fs.read_all
    -> std.json.parse
