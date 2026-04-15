# Requirement: "an atomic file write stream"

Writes buffered data to a temporary file and renames it into place on close so readers never see a partial file.

std
  std.fs
    std.fs.open_temp
      fn (dir: string, prefix: string) -> result[file_handle, string]
      + creates a new unique temp file in dir and returns an open handle
      - returns error when dir is not writable
      # filesystem
    std.fs.write
      fn (h: file_handle, data: bytes) -> result[void, string]
      + appends data to the open file
      - returns error on write failure
      # filesystem
    std.fs.fsync
      fn (h: file_handle) -> result[void, string]
      + flushes buffered writes to disk
      # filesystem
    std.fs.close
      fn (h: file_handle) -> result[void, string]
      + closes the handle
      # filesystem
    std.fs.rename
      fn (src: string, dst: string) -> result[void, string]
      + atomically renames src to dst, replacing dst if present
      - returns error when src does not exist or the rename crosses filesystems
      # filesystem
    std.fs.remove
      fn (path: string) -> result[void, string]
      + deletes a file
      # filesystem

atomic_write
  atomic_write.open
    fn (target_path: string) -> result[stream_state, string]
    + opens a temp file alongside target_path and returns the stream state
    - returns error when the parent directory is not writable
    # lifecycle
    -> std.fs.open_temp
  atomic_write.write
    fn (state: stream_state, chunk: bytes) -> result[stream_state, string]
    + appends a chunk to the underlying temp file
    - returns error on write failure
    # writing
    -> std.fs.write
  atomic_write.commit
    fn (state: stream_state) -> result[void, string]
    + fsyncs, closes, and renames the temp file onto the target path
    - returns error when any step fails and removes the temp file
    # lifecycle
    -> std.fs.fsync
    -> std.fs.close
    -> std.fs.rename
  atomic_write.abort
    fn (state: stream_state) -> result[void, string]
    + closes and removes the temp file without touching the target
    - returns error on removal failure
    # lifecycle
    -> std.fs.close
    -> std.fs.remove
