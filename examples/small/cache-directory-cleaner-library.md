# Requirement: "a library for cleaning software cache directories"

Enumerates known cache directories, computes their total size, and removes entries older than a threshold.

std
  std.fs
    std.fs.list_dir
      @ (path: string) -> result[list[dir_entry], string]
      + returns entries with name, size, modified_unix_seconds, and is_dir
      - returns error when path does not exist or is not a directory
      # filesystem
    std.fs.remove_recursive
      @ (path: string) -> result[void, string]
      + deletes a file or directory tree rooted at path
      - returns error when removal fails on any entry
      # filesystem
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

cache_cleanup
  cache_cleanup.scan
    @ (roots: list[string]) -> result[list[cache_entry], string]
    + walks each root and returns every file with its path, size, and modified time
    - returns error on the first unreadable root
    # scanning
    -> std.fs.list_dir
  cache_cleanup.total_bytes
    @ (entries: list[cache_entry]) -> i64
    + sums the size of every entry
    # reporting
  cache_cleanup.filter_older_than
    @ (entries: list[cache_entry], max_age_seconds: i64) -> list[cache_entry]
    + keeps entries whose modified time is older than now minus the threshold
    # selection
    -> std.time.now_seconds
  cache_cleanup.purge
    @ (entries: list[cache_entry]) -> result[i64, string]
    + removes every entry and returns the total bytes reclaimed
    - returns error on the first removal failure
    # cleanup
    -> std.fs.remove_recursive
