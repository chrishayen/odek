# Requirement: "a kernel subsystem library providing processes, scheduling, memory, and virtual file system"

Library abstractions of core OS subsystems. No hardware specifics — the pieces a kernel project would compose.

std: (all units exist)

kernel
  kernel.create_process
    fn (name: string, entry: u64) -> process_id
    + returns a new process with a fresh id and the given entry address
    # process
  kernel.kill_process
    fn (id: process_id) -> result[void, string]
    + removes the process from all queues
    - returns error when the id is not known
    # process
  kernel.list_processes
    fn () -> list[process_id]
    + returns all live process ids
    # process
  kernel.schedule_next
    fn () -> optional[process_id]
    + returns the next process to run using round-robin
    + returns none when no processes are runnable
    # scheduling
  kernel.mark_blocked
    fn (id: process_id, reason: string) -> result[void, string]
    + moves the process off the runnable queue
    - returns error when the id is not known
    # scheduling
  kernel.mark_runnable
    fn (id: process_id) -> result[void, string]
    + moves the process back onto the runnable queue
    - returns error when the id is not known
    # scheduling
  kernel.allocate_page
    fn (pid: process_id) -> result[u64, string]
    + returns the physical address of a freshly mapped page
    - returns error when no free pages remain
    # memory
  kernel.free_page
    fn (pid: process_id, address: u64) -> result[void, string]
    + releases a page previously allocated to the process
    - returns error when the page is not owned by the process
    # memory
  kernel.map_region
    fn (pid: process_id, virt: u64, phys: u64, length: u64) -> result[void, string]
    + installs a virtual-to-physical mapping for the process
    - returns error when length is not page-aligned
    # memory
  kernel.vfs_mount
    fn (path: string, fs_driver: fs_driver_handle) -> result[void, string]
    + mounts the driver at the given path
    - returns error when the path is already a mount point
    # vfs
  kernel.vfs_open
    fn (pid: process_id, path: string, flags: i32) -> result[file_descriptor, string]
    + opens the path and returns a descriptor owned by the process
    - returns error when no mount matches the path
    # vfs
  kernel.vfs_read
    fn (fd: file_descriptor, length: u64) -> result[bytes, string]
    + reads up to length bytes
    - returns error when fd is closed
    # vfs
  kernel.vfs_write
    fn (fd: file_descriptor, data: bytes) -> result[u64, string]
    + writes the bytes and returns the number written
    - returns error when fd is read-only
    # vfs
  kernel.vfs_close
    fn (fd: file_descriptor) -> result[void, string]
    + releases the descriptor
    - returns error when fd was already closed
    # vfs
