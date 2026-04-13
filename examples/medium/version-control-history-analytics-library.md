# Requirement: "a version-control history analytics library"

Given a repository adapter that yields commits, produce contributor, churn, and activity metrics.

std
  std.time
    std.time.unix_to_date
      @ (unix_seconds: i64) -> string
      + returns an "YYYY-MM-DD" date string in UTC
      # time

repo_insights
  repo_insights.load_history
    @ (repo: repo_adapter) -> result[history, string]
    + drains the adapter and returns an in-memory commit history
    - returns error when the adapter reports a read failure
    # ingestion
  repo_insights.commit_count
    @ (h: history) -> i32
    + returns the total number of commits
    # basic_stats
  repo_insights.contributors
    @ (h: history) -> list[contributor_record]
    + returns contributors sorted by commit count descending
    + each record includes author email and total line additions and deletions
    # contributors
  repo_insights.commits_per_day
    @ (h: history) -> map[string, i32]
    + returns a date-to-count map of commit activity
    # activity
    -> std.time.unix_to_date
  repo_insights.file_churn
    @ (h: history) -> list[file_churn_record]
    + returns files sorted by total change volume descending
    + each record has path, additions, deletions, and commit count
    # churn
  repo_insights.bus_factor
    @ (h: history, path: string) -> i32
    + returns the minimum number of authors responsible for more than half the lines of a file
    - returns 0 when no commits touch the path
    # risk
  repo_insights.first_commit_time
    @ (h: history) -> optional[i64]
    + returns the unix timestamp of the earliest commit
    - returns none when the history is empty
    # basic_stats
  repo_insights.longest_gap
    @ (h: history) -> i64
    + returns the longest gap in seconds between consecutive commits
    + returns 0 when there are fewer than two commits
    # activity
