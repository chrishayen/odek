# Requirement: "a modern unix-like operating system with a monolithic kernel"

A kernel library exposing the core subsystems a monolithic unix-like OS needs. Hardware drivers and the boot loader are out of scope; this is the in-memory kernel data model and syscall layer.

std
  std.memory
    std.memory.alloc_pages
      fn (count: u64) -> result[u64, string]
      + returns the physical base address of a contiguous page run
      - returns error when no free run of the requested size exists
      # memory
    std.memory.free_pages
      fn (base: u64, count: u64) -> void
      + marks the page run as free
      # memory
  std.sync
    std.sync.spinlock_acquire
      fn (lock: u64) -> void
      + busy-waits until the lock word transitions from 0 to 1
      # synchronization
    std.sync.spinlock_release
      fn (lock: u64) -> void
      + stores 0 into the lock word
      # synchronization
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns a monotonic tick count in nanoseconds
      # time

kernel
  kernel.boot_init
    fn (mem_size: u64) -> kernel_state
    + initializes page allocator, process table, and vfs root
    # bootstrap
    -> std.memory.alloc_pages
  kernel.process_create
    fn (state: kernel_state, parent_pid: i32, image: bytes) -> result[i32, string]
    + allocates a new pid and a fresh address space
    - returns error when the process table is full
    # process_management
    -> std.memory.alloc_pages
  kernel.process_exit
    fn (state: kernel_state, pid: i32, code: i32) -> void
    + marks the process as zombie and wakes its parent
    # process_management
  kernel.scheduler_tick
    fn (state: kernel_state) -> i32
    + picks the next runnable pid using round-robin among ready processes
    + returns the previously running pid for context-switch bookkeeping
    # scheduling
    -> std.time.now_nanos
  kernel.vm_map
    fn (state: kernel_state, pid: i32, vaddr: u64, len: u64, flags: u32) -> result[void, string]
    + installs page table entries covering the requested virtual range
    - returns error when the range overlaps an existing mapping
    # virtual_memory
    -> std.memory.alloc_pages
  kernel.vm_unmap
    fn (state: kernel_state, pid: i32, vaddr: u64, len: u64) -> result[void, string]
    + removes page table entries and releases backing pages
    - returns error when the range is not fully mapped
    # virtual_memory
    -> std.memory.free_pages
  kernel.vfs_mount
    fn (state: kernel_state, path: string, fs_kind: string) -> result[void, string]
    + attaches a filesystem at the given path
    - returns error when the path already has a mount
    # filesystem
  kernel.vfs_open
    fn (state: kernel_state, pid: i32, path: string, flags: i32) -> result[i32, string]
    + returns a new per-process file descriptor
    - returns error when the path does not resolve
    # filesystem
  kernel.vfs_read
    fn (state: kernel_state, pid: i32, fd: i32, max: u64) -> result[bytes, string]
    + returns up to max bytes from the file position
    - returns error when the fd is not open for reading
    # filesystem
  kernel.vfs_write
    fn (state: kernel_state, pid: i32, fd: i32, data: bytes) -> result[u64, string]
    + returns the number of bytes written
    - returns error when the fd is not open for writing
    # filesystem
  kernel.syscall_dispatch
    fn (state: kernel_state, pid: i32, nr: i32, args: list[u64]) -> result[i64, i32]
    + routes numbered syscalls to the appropriate kernel handler
    - returns errno when the syscall number is unknown
    # syscalls
    -> std.sync.spinlock_acquire
    -> std.sync.spinlock_release
  kernel.irq_handle
    fn (state: kernel_state, irq: i32) -> void
    + dispatches a hardware interrupt to registered handlers and yields the scheduler
    # interrupts
