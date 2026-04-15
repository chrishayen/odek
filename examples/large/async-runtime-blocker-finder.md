# Requirement: "a library for finding what is blocking an async runtime"

Attributes long-running synchronous calls on async worker threads to the async task that owns them. Sampling primitives and symbolization live in std; the project layer aggregates samples into reports.

std
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time
  std.sampling
    std.sampling.capture_stack
      fn (thread_id: i64) -> result[list[i64], string]
      + returns the current instruction pointers of the target thread, outermost first
      - returns error when the thread cannot be inspected
      # profiling
    std.sampling.resolve_symbol
      fn (address: i64) -> optional[symbol_info]
      + returns function name, file, and line for an instruction pointer
      - returns none when the address has no debug info
      # symbolization
  std.threads
    std.threads.list_worker_threads
      fn () -> list[i64]
      + returns the thread ids registered as async runtime workers
      # threading
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file into memory
      - returns error when the file does not exist
      # filesystem

block_profiler
  block_profiler.new
    fn (interval_ms: i32) -> profiler_state
    + creates a profiler sampling every interval_ms milliseconds
    # construction
  block_profiler.register_task
    fn (state: profiler_state, task_id: i64, thread_id: i64) -> profiler_state
    + records that task_id is currently executing on thread_id
    # task_tracking
  block_profiler.unregister_task
    fn (state: profiler_state, task_id: i64) -> profiler_state
    + removes the task-to-thread mapping
    # task_tracking
  block_profiler.sample_once
    fn (state: profiler_state) -> profiler_state
    + captures one stack per worker thread and stores a stack_sample
    ? each sample is tagged with the task currently assigned to that worker
    # sampling
    -> std.threads.list_worker_threads
    -> std.sampling.capture_stack
    -> std.time.now_nanos
  block_profiler.mark_blocking
    fn (sample: stack_sample) -> bool
    + returns true when any frame in the sample is a known blocking syscall
    + recognizes file, network, and mutex wait frames
    # classification
    -> std.sampling.resolve_symbol
  block_profiler.aggregate
    fn (samples: list[stack_sample]) -> list[task_blocking_report]
    + groups samples by task id and counts blocking occurrences
    + reports total sampled time and blocked time per task
    # aggregation
  block_profiler.rank_by_blocked_time
    fn (reports: list[task_blocking_report]) -> list[task_blocking_report]
    + returns reports sorted from most to least blocked time
    # ranking
  block_profiler.top_blocking_stacks
    fn (samples: list[stack_sample], n: i32) -> list[stack_frequency]
    + returns the n most frequently observed blocking stacks
    + each entry includes a resolved symbol chain and sample count
    # hotspots
    -> std.sampling.resolve_symbol
  block_profiler.format_report
    fn (reports: list[task_blocking_report], hotspots: list[stack_frequency]) -> string
    + renders a human-readable profiling summary
    # rendering
  block_profiler.load_symbol_map
    fn (path: string) -> result[symbol_map, string]
    + loads a cached symbol map from disk
    - returns error when the file is missing or malformed
    # symbolization
    -> std.fs.read_all
