# Requirement: "a rolling log file writer"

Writes bytes to a log file, rotating to a timestamped backup when a size threshold is exceeded. Retains a bounded number of backups.

std
  std.fs
    std.fs.append_bytes
      fn (path: string, data: bytes) -> result[i64, string]
      + appends bytes to the file, creating it when absent, returns new file size
      - returns error when the path is not writable
      # filesystem
    std.fs.rename
      fn (from: string, to: string) -> result[void, string]
      + atomically renames a file
      - returns error when the destination is invalid
      # filesystem
    std.fs.remove
      fn (path: string) -> result[void, string]
      + removes the file; no error when absent
      # filesystem
    std.fs.list_dir
      fn (dir: string) -> result[list[string], string]
      + returns names of entries in the directory
      # filesystem
  std.time
    std.time.now_millis
      fn () -> i64
      + returns the current unix time in milliseconds
      # time

rolling_log
  rolling_log.open
    fn (path: string, max_bytes: i64, max_backups: i32) -> rolling_log_state
    + creates a rolling log handle bound to the target path and limits
    # construction
  rolling_log.write
    fn (state: rolling_log_state, data: bytes) -> result[rolling_log_state, string]
    + appends data to the active file
    + rotates and prunes old backups when the file exceeds max_bytes
    - returns error when the underlying filesystem operation fails
    # rotation
    -> std.fs.append_bytes
    -> std.fs.rename
    -> std.fs.list_dir
    -> std.fs.remove
    -> std.time.now_millis
  rolling_log.close
    fn (state: rolling_log_state) -> void
    + finalizes the handle; subsequent writes are not allowed
    # lifecycle
