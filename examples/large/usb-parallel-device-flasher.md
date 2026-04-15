# Requirement: "a library for flashing multiple usb devices in parallel"

The project layer coordinates parallel writes to multiple block devices and reports progress per device. The UI layer is the caller's concern.

std
  std.fs
    std.fs.open_read
      fn (path: string) -> result[file_handle, string]
      + opens a file for reading
      - returns error when the file does not exist
      # filesystem
    std.fs.open_write
      fn (path: string) -> result[file_handle, string]
      + opens a path for writing in raw mode
      - returns error when the path cannot be opened
      # filesystem
    std.fs.read_chunk
      fn (handle: file_handle, size: i64) -> result[bytes, string]
      + reads up to size bytes from the current position
      + returns empty bytes at end of file
      # filesystem
    std.fs.write_all
      fn (handle: file_handle, data: bytes) -> result[i64, string]
      + writes all bytes and returns the number written
      - returns error on write failure
      # filesystem
    std.fs.file_size
      fn (path: string) -> result[i64, string]
      + returns the size of a file or block device in bytes
      # filesystem
    std.fs.close
      fn (handle: file_handle) -> void
      + releases the underlying file handle
      # filesystem
  std.concurrency
    std.concurrency.spawn
      fn (task: task_fn) -> task_handle
      + starts a task running concurrently
      # concurrency
    std.concurrency.join_all
      fn (handles: list[task_handle]) -> void
      + waits until all tasks have completed
      # concurrency
  std.hash
    std.hash.sha256
      fn (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest
      # cryptography

usb_flasher
  usb_flasher.list_targets
    fn () -> result[list[string], string]
    + enumerates removable block devices available for flashing
    - returns error when the device list cannot be queried
    ? enumeration is platform-specific; a stub may return an empty list in tests
    # enumeration
  usb_flasher.verify_image
    fn (image_path: string, expected_sha256: optional[bytes]) -> result[i64, string]
    + returns the image size when the file exists and optionally matches the expected digest
    - returns error when the digest does not match
    # verification
    -> std.fs.file_size
    -> std.fs.open_read
    -> std.fs.read_chunk
    -> std.hash.sha256
  usb_flasher.flash_one
    fn (image_path: string, target_path: string, progress: progress_sink) -> result[i64, string]
    + copies the image to the target device, emitting progress bytes to the sink
    - returns error when the target cannot be opened or a write fails
    # write
    -> std.fs.open_read
    -> std.fs.open_write
    -> std.fs.read_chunk
    -> std.fs.write_all
    -> std.fs.close
  usb_flasher.flash_parallel
    fn (image_path: string, targets: list[string], progress: map[string, progress_sink]) -> result[map[string, i64], string]
    + flashes the same image to multiple targets concurrently and returns per-target bytes written
    - returns error describing the first failed target when any write fails
    # orchestration
    -> std.concurrency.spawn
    -> std.concurrency.join_all
