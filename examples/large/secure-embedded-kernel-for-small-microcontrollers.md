# Requirement: "a secure embedded operating system kernel for small microcontrollers"

A cooperative kernel for memory-constrained devices: preemptive scheduling, isolated process grants, and a syscall dispatcher.

std
  std.mem
    std.mem.zero
      fn (buf: bytes) -> bytes
      + returns a buffer of the same length with all zero bytes
      # memory
    std.mem.copy
      fn (dst: bytes, src: bytes) -> bytes
      + copies src into dst up to min length
      # memory
  std.time
    std.time.monotonic_us
      fn () -> i64
      + returns a monotonic microsecond counter
      # time

kernel
  kernel.new
    fn (ram_size: i32, flash_size: i32) -> kernel_state
    + initializes a kernel with empty process table and sized memory pools
    # construction
    -> std.mem.zero
  kernel.register_process
    fn (state: kernel_state, name: string, entry: u32, stack_size: i32, grant_size: i32) -> result[process_id, string]
    + allocates a stack and grant region, returns a process id
    - returns error when there is no remaining RAM for the allocation
    # process_table
  kernel.schedule_next
    fn (state: kernel_state) -> optional[process_id]
    + picks the next runnable process using round-robin over the process table
    - returns none when no process is runnable
    # scheduling
  kernel.tick
    fn (state: kernel_state, elapsed_us: i64) -> kernel_state
    + advances all software timers and wakes processes whose timer expired
    # scheduling
    -> std.time.monotonic_us
  kernel.handle_syscall
    fn (state: kernel_state, caller: process_id, number: u32, args: list[u32]) -> result[syscall_result, string]
    + dispatches numbered syscalls to the appropriate kernel handler
    - returns error when the syscall number is unknown
    - returns error when the caller is not currently scheduled
    # syscall
  kernel.grant_alloc
    fn (state: kernel_state, pid: process_id, size: i32) -> result[u32, string]
    + allocates size bytes inside the process grant region and returns the base address
    - returns error when the grant region is exhausted
    # memory
  kernel.capability_issue
    fn (state: kernel_state, pid: process_id, resource: string) -> capability
    + issues an unforgeable token that authorizes access to a named resource
    ? capabilities are opaque to processes and only kernel code can dereference them
    # security
  kernel.capability_check
    fn (state: kernel_state, cap: capability, resource: string) -> bool
    + returns true when the capability grants access to the requested resource
    - returns false when the capability was revoked
    # security
  kernel.ipc_send
    fn (state: kernel_state, from: process_id, to: process_id, msg: bytes) -> result[kernel_state, string]
    + enqueues a message on the target process inbox
    - returns error when the target has no inbox capacity left
    # ipc
    -> std.mem.copy
  kernel.ipc_receive
    fn (state: kernel_state, pid: process_id) -> optional[bytes]
    + pops the next pending message for the given process
    - returns none when the inbox is empty
    # ipc
  kernel.fault
    fn (state: kernel_state, pid: process_id, reason: string) -> kernel_state
    + marks the process as faulted and frees its resources
    # fault_handling
