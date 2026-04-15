# Requirement: "a library for walking a filesystem concurrently"

Walks a directory tree in parallel workers, applying a user-provided visitor to each entry. The caller tunes worker count and gets back completion stats.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[dir_entry], string]
      + returns the immediate entries in a directory
      - returns error when the path is not a directory
      # filesystem
    std.fs.is_dir
      fn (entry: dir_entry) -> bool
      + returns whether the entry refers to a directory
      # filesystem
  std.concurrency
    std.concurrency.spawn_workers
      fn (count: i32, task: function[string, void]) -> worker_pool
      + starts a pool of count workers sharing a single task function
      # concurrency
    std.concurrency.submit
      fn (pool: worker_pool, item: string) -> void
      + submits an item to be processed by the pool
      # concurrency
    std.concurrency.wait_all
      fn (pool: worker_pool) -> void
      + blocks until the pool has drained
      # concurrency

walker
  walker.new
    fn (workers: i32, visit: function[dir_entry, void]) -> walker_state
    + builds a walker configured with worker count and an entry visitor
    # configuration
  walker.walk
    fn (w: walker_state, root: string) -> result[walk_report, string]
    + walks root, visiting every entry exactly once across the pool, and returns a report with entry count and error count
    - returns error when root does not exist
    # traversal
    -> std.fs.list_dir
    -> std.fs.is_dir
    -> std.concurrency.spawn_workers
    -> std.concurrency.submit
    -> std.concurrency.wait_all
  walker.report_entry_count
    fn (report: walk_report) -> i64
    + returns the number of entries visited
    # reporting
