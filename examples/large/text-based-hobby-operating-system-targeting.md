# Requirement: "a text-mode hobby operating system kernel for 64-bit PCs booted from firmware"

A kernel is a large system. The project rune set covers the core subsystems; std holds generic primitives for byte buffers, UTF-8, and a clock primitive.

std
  std.bytes
    std.bytes.copy
      @ (src: bytes, dst_offset: i64, dst: bytes) -> bytes
      + copies src into dst at the given offset and returns the updated buffer
      # bytes
  std.text
    std.text.utf8_decode
      @ (raw: bytes) -> result[string, string]
      + decodes a UTF-8 byte sequence into a string
      - returns error on invalid UTF-8
      # text
  std.time
    std.time.now_ticks
      @ () -> i64
      + returns a monotonic tick count since boot
      # time

kernel
  kernel.boot
    @ (memory_map: list[memory_region]) -> kernel_state
    + initialises the kernel from a firmware-provided memory map
    # boot
  kernel.allocate_frame
    @ (state: kernel_state) -> result[u64, string]
    + reserves a physical frame and returns its base address
    - returns error when no free frames remain
    # memory
  kernel.free_frame
    @ (state: kernel_state, frame: u64) -> result[void, string]
    + returns a previously allocated frame to the free list
    - returns error when the frame was not allocated
    # memory
  kernel.register_interrupt
    @ (state: kernel_state, vector: u8, handler_id: i32) -> kernel_state
    + associates a handler id with an interrupt vector
    # interrupts
  kernel.dispatch_interrupt
    @ (state: kernel_state, vector: u8) -> kernel_state
    + invokes the handler registered for the vector and returns the updated state
    # interrupts
  kernel.spawn_task
    @ (state: kernel_state, entry: u64) -> tuple[kernel_state, i32]
    + creates a new task starting at the given entry address and returns its task id
    # scheduling
  kernel.schedule_next
    @ (state: kernel_state) -> tuple[kernel_state, optional[i32]]
    + returns the id of the next runnable task, or none if idle
    # scheduling
    -> std.time.now_ticks
  kernel.read_keyboard
    @ (state: kernel_state) -> tuple[kernel_state, optional[u8]]
    + returns the next scancode from the keyboard buffer, if any
    # input
  kernel.vga_write
    @ (state: kernel_state, row: i32, col: i32, text: string) -> kernel_state
    + writes text to a text-mode display buffer at the given cell
    - out-of-range coordinates are clipped
    # output
    -> std.bytes.copy
  kernel.fs_open
    @ (state: kernel_state, path: string) -> result[i32, string]
    + opens a path in the in-memory filesystem and returns a file handle
    - returns error when the path does not exist
    # filesystem
  kernel.fs_read
    @ (state: kernel_state, handle: i32, count: i64) -> result[bytes, string]
    + reads up to count bytes from the handle's current offset
    - returns error on an invalid handle
    # filesystem
  kernel.shell_execute
    @ (state: kernel_state, command_line: string) -> tuple[kernel_state, string]
    + parses a command line and returns the shell's textual response
    # shell
    -> std.text.utf8_decode
