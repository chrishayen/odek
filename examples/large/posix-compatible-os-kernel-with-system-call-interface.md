# Requirement: "a general-purpose OS kernel with a posix-compatible system call interface"

Core kernel abstractions: processes, virtual memory, a VFS, and a syscall table. Scope is deliberately a single module; the ABI names are represented by opaque identifiers.

std
  std.mem
    std.mem.alloc_pages
      @ (count: i32) -> result[u64, string]
      + reserves count contiguous 4KiB pages and returns the base physical address
      - returns error when not enough free pages remain
      # memory
    std.mem.free_pages
      @ (base: u64, count: i32) -> void
      + returns the pages to the free list
      # memory
  std.hash
    std.hash.fnv1a_64
      @ (data: bytes) -> u64
      + computes FNV-1a 64-bit hash
      # hashing

kernel
  kernel.boot
    @ (total_pages: i32) -> kernel_state
    + initializes page allocator, scheduler, and vfs root
    # construction
    -> std.mem.alloc_pages
  kernel.spawn_process
    @ (state: kernel_state, image: bytes, argv: list[string]) -> result[process_id, string]
    + loads the image into a fresh address space and enqueues it
    - returns error on image parse failure
    # process_table
  kernel.exit_process
    @ (state: kernel_state, pid: process_id, code: i32) -> kernel_state
    + marks the process as exited and reclaims its pages
    # process_table
    -> std.mem.free_pages
  kernel.context_switch
    @ (state: kernel_state) -> kernel_state
    + chooses the next runnable process by priority then round-robin
    # scheduling
  kernel.map_region
    @ (state: kernel_state, pid: process_id, vaddr: u64, size: i64, prot: u8) -> result[kernel_state, string]
    + creates a virtual memory mapping in the process address space
    - returns error when vaddr overlaps an existing mapping
    # virtual_memory
    -> std.mem.alloc_pages
  kernel.unmap_region
    @ (state: kernel_state, pid: process_id, vaddr: u64, size: i64) -> kernel_state
    + removes the mapping and frees backing pages
    # virtual_memory
    -> std.mem.free_pages
  kernel.vfs_mount
    @ (state: kernel_state, path: string, fs: filesystem) -> result[kernel_state, string]
    + attaches a filesystem at the given mount point
    - returns error when path is not a directory
    # vfs
  kernel.vfs_open
    @ (state: kernel_state, pid: process_id, path: string, flags: u32) -> result[tuple[kernel_state, file_descriptor], string]
    + resolves path across mounts and returns a new file descriptor
    - returns error when the path does not exist and create flag is unset
    # vfs
    -> std.hash.fnv1a_64
  kernel.vfs_read
    @ (state: kernel_state, pid: process_id, fd: file_descriptor, len: i64) -> result[bytes, string]
    + reads up to len bytes from the current file offset and advances it
    - returns error when the fd is closed
    # vfs
  kernel.vfs_write
    @ (state: kernel_state, pid: process_id, fd: file_descriptor, data: bytes) -> result[i64, string]
    + writes data at the current offset and returns bytes written
    # vfs
  kernel.syscall
    @ (state: kernel_state, pid: process_id, number: u32, args: list[u64]) -> result[tuple[kernel_state, i64], string]
    + dispatches a numbered syscall to the appropriate kernel handler
    - returns error when the syscall number is unknown
    # syscall
  kernel.signal
    @ (state: kernel_state, pid: process_id, sig: i32) -> kernel_state
    + posts a signal to the target process
    # signals
