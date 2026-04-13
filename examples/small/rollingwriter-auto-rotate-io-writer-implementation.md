# Requirement: "an auto-rotating file writer with multiple rotation policies"

Stateful writer: appends bytes, and rotates when a policy (size, time, or line count) is exceeded. Actual file I/O is delegated to the host; this library returns rotation decisions.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

rollwriter
  rollwriter.new
    @ (base_path: string, policy: rotation_policy) -> writer_state
    + creates a writer anchored at base_path with the given policy and zero counters
    # construction
    -> std.time.now_seconds
  rollwriter.size_policy
    @ (max_bytes: i64) -> rotation_policy
    + returns a policy that rotates when accumulated bytes exceed max_bytes
    # policy
  rollwriter.time_policy
    @ (interval_seconds: i64) -> rotation_policy
    + returns a policy that rotates when now - last_rotation exceeds interval
    # policy
  rollwriter.line_policy
    @ (max_lines: i64) -> rotation_policy
    + returns a policy that rotates when line count exceeds max_lines
    # policy
  rollwriter.append
    @ (state: writer_state, data: bytes) -> tuple[writer_state, optional[string]]
    + returns updated counters and, if rotation is due, the new rotated path
    # appending
    -> std.time.now_seconds
