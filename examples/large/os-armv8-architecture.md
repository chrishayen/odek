# Requirement: "a minimal operating system kernel for a 64-bit ARM architecture"

Library-style kernel primitives: boot record, memory map, physical frame allocator, page table, scheduler, and a trap dispatcher. The caller wires these into an image.

std
  std.bits
    std.bits.set_bit
      @ (word: u64, index: i32) -> u64
      + returns word with the given bit set
      # bit_manipulation
    std.bits.clear_bit
      @ (word: u64, index: i32) -> u64
      + returns word with the given bit cleared
      # bit_manipulation
    std.bits.test_bit
      @ (word: u64, index: i32) -> bool
      + returns whether the given bit is set
      # bit_manipulation
  std.mem
    std.mem.zero_range
      @ (base: u64, length: u64) -> void
      + writes zero across the given physical range
      # memory
    std.mem.copy_range
      @ (dst: u64, src: u64, length: u64) -> void
      + copies bytes from one physical range to another
      # memory

kernel
  kernel.parse_boot_info
    @ (raw: bytes) -> result[boot_info, string]
    + parses the firmware-provided boot structure
    - returns error when the magic number does not match
    # boot
  kernel.build_memory_map
    @ (info: boot_info) -> memory_map
    + returns a memory map describing free and reserved physical regions
    # boot
  kernel.new_frame_allocator
    @ (map: memory_map) -> frame_allocator_state
    + constructs a bitmap-backed physical frame allocator from free regions
    # memory_management
    -> std.bits.set_bit
  kernel.alloc_frame
    @ (alloc: frame_allocator_state) -> result[tuple[u64, frame_allocator_state], string]
    + returns the physical address of a freshly allocated frame
    - returns error when no free frames remain
    # memory_management
    -> std.bits.test_bit
    -> std.bits.set_bit
  kernel.free_frame
    @ (alloc: frame_allocator_state, addr: u64) -> frame_allocator_state
    + marks the frame at addr as free
    # memory_management
    -> std.bits.clear_bit
  kernel.new_page_table
    @ () -> page_table_state
    + returns an empty top-level page table
    # memory_management
    -> std.mem.zero_range
  kernel.map_page
    @ (pt: page_table_state, virt: u64, phys: u64, flags: u64) -> result[page_table_state, string]
    + installs a translation from a virtual to a physical page
    - returns error when the translation already exists with different flags
    # memory_management
  kernel.unmap_page
    @ (pt: page_table_state, virt: u64) -> result[page_table_state, string]
    + removes a translation for a virtual page
    - returns error when the page is not mapped
    # memory_management
  kernel.new_task
    @ (entry: u64, stack_top: u64) -> task_state
    + returns a new task with the given entry point and kernel stack
    # scheduling
  kernel.schedule_next
    @ (runqueue: list[task_state]) -> tuple[task_state, list[task_state]]
    + returns the next task to run and the updated runqueue
    # scheduling
  kernel.handle_trap
    @ (task: task_state, cause: u64, value: u64) -> result[task_state, string]
    + updates a task based on the trap cause (syscall, fault, interrupt)
    - returns error on an unrecoverable fault, indicating the task must be killed
    # trap_dispatch
