# Requirement: "a lightweight time tracking library tied to version control commits"

Records work intervals against commits and summarizes time per file or author.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.vcs
    std.vcs.head_commit
      @ (repo_dir: string) -> result[string, string]
      + returns the current head commit hash
      - returns error when the path is not a repository
      # version_control
    std.vcs.commit_author
      @ (repo_dir: string, commit: string) -> result[string, string]
      + returns the author name for the commit
      # version_control
    std.vcs.commit_files
      @ (repo_dir: string, commit: string) -> result[list[string], string]
      + returns file paths changed in the commit
      # version_control

time_tracker
  time_tracker.new_log
    @ () -> tracker_state
    + creates an empty tracker
    # construction
  time_tracker.start_interval
    @ (state: tracker_state, repo_dir: string) -> result[tracker_state, string]
    + records an interval start tied to the current head commit
    - returns error when the repository has no commits
    # intervals
    -> std.time.now_seconds
    -> std.vcs.head_commit
  time_tracker.stop_interval
    @ (state: tracker_state) -> result[tracker_state, string]
    + closes the open interval using the current time
    - returns error when there is no open interval
    # intervals
    -> std.time.now_seconds
  time_tracker.total_by_commit
    @ (state: tracker_state) -> map[string, i64]
    + returns total seconds recorded per commit hash
    # reporting
  time_tracker.total_by_file
    @ (state: tracker_state, repo_dir: string) -> result[map[string, i64], string]
    + distributes each interval's seconds evenly across the commit's changed files
    - returns error when any recorded commit cannot be inspected
    # reporting
    -> std.vcs.commit_files
  time_tracker.total_by_author
    @ (state: tracker_state, repo_dir: string) -> result[map[string, i64], string]
    + returns total seconds per commit author
    - returns error when any recorded commit cannot be inspected
    # reporting
    -> std.vcs.commit_author
