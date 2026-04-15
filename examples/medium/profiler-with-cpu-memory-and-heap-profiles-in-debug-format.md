# Requirement: "a profiling library that captures CPU, memory, and heap profiles in a standard debug tool format"

Collects sample-based profiles of a running program and writes them in a widely supported profile format.

std
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time
  std.fs
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes a file
      # filesystem

profiler
  profiler.new_cpu_profiler
    fn (sample_hz: i32) -> profiler_state
    + creates a CPU profiler that samples the stack at sample_hz
    # construction
  profiler.new_heap_profiler
    fn () -> profiler_state
    + creates a heap profiler that records allocation samples
    # construction
  profiler.start
    fn (p: profiler_state) -> result[profiler_state, string]
    + begins collection
    - returns error when a collection is already running
    # lifecycle
    -> std.time.now_nanos
  profiler.stop
    fn (p: profiler_state) -> profile_data
    + stops collection and returns aggregated samples
    # lifecycle
    -> std.time.now_nanos
  profiler.record_sample
    fn (p: profiler_state, stack: list[string], weight: i64) -> profiler_state
    + adds a sample with the given stack frames and weight
    ? weight is a time delta for CPU and byte count for heap
    # sampling
  profiler.encode_profile
    fn (data: profile_data) -> bytes
    + serializes the profile to the standard profile wire format
    # encoding
  profiler.write_file
    fn (data: profile_data, path: string) -> result[void, string]
    + encodes and writes the profile to disk
    # io
    -> std.fs.write_all
