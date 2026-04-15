# Requirement: "a virtual machine monitor that runs guest operating systems in a fast, secure virtualized environment"

A lightweight VMM facade: configure a VM, allocate memory, attach virtio devices, and drive vCPUs. The host's KVM-like interface is abstracted behind a std primitive.

std
  std.hv
    std.hv.open
      fn () -> result[hv_handle, string]
      + opens the host hypervisor interface
      - returns error when the host does not expose a hypervisor
      # hypervisor
    std.hv.create_vm
      fn (hv: hv_handle) -> result[vm_handle, string]
      + creates a new virtual machine on the hypervisor
      # hypervisor
    std.hv.create_vcpu
      fn (vm: vm_handle, index: u32) -> result[vcpu_handle, string]
      + creates a vCPU with the given index
      # hypervisor
    std.hv.run_vcpu
      fn (vcpu: vcpu_handle) -> result[vm_exit, string]
      + runs the vCPU until it exits to the host and returns the exit reason
      - returns error on vcpu fault
      # hypervisor
  std.mm
    std.mm.map_guest_memory
      fn (vm: vm_handle, guest_addr: u64, size: u64) -> result[bytes, string]
      + allocates host-backed memory and maps it into the guest at the given address
      - returns error when the region overlaps an existing mapping
      # memory
    std.mm.unmap_guest_memory
      fn (vm: vm_handle, guest_addr: u64) -> result[void, string]
      + removes a previously-created guest mapping
      # memory
  std.io
    std.io.eventfd_new
      fn () -> result[event_fd, string]
      + creates an edge-triggered notification file descriptor
      # io
    std.io.eventfd_signal
      fn (fd: event_fd) -> void
      + wakes any waiter on the fd
      # io

vmm
  vmm.new
    fn (memory_bytes: u64, vcpu_count: u32) -> result[vmm_state, string]
    + creates a VMM configured with the requested memory size and vCPU count
    - returns error when memory_bytes is zero or vcpu_count is zero
    # construction
    -> std.hv.open
    -> std.hv.create_vm
  vmm.load_kernel
    fn (state: vmm_state, kernel: bytes, load_addr: u64) -> result[vmm_state, string]
    + copies the kernel image into guest memory at the given address
    - returns error when the image would not fit
    # kernel_loading
    -> std.mm.map_guest_memory
  vmm.attach_virtio_block
    fn (state: vmm_state, backing: bytes) -> result[vmm_state, string]
    + attaches a virtio-block device backed by the given buffer
    # device_attach
    -> std.io.eventfd_new
  vmm.attach_virtio_net
    fn (state: vmm_state, mac: bytes) -> result[vmm_state, string]
    + attaches a virtio-net device with the given MAC address
    - returns error when mac is not 6 bytes
    # device_attach
    -> std.io.eventfd_new
  vmm.attach_virtio_console
    fn (state: vmm_state) -> vmm_state
    + attaches a virtio-console device that buffers guest output
    # device_attach
  vmm.start_vcpus
    fn (state: vmm_state) -> result[vmm_state, string]
    + creates vCPU handles for all configured cores and marks them runnable
    # scheduling
    -> std.hv.create_vcpu
  vmm.run
    fn (state: vmm_state) -> result[vm_exit_reason, string]
    + advances all runnable vCPUs until any one exits to the host and returns the reason
    - returns error on vCPU fault
    # scheduling
    -> std.hv.run_vcpu
  vmm.handle_mmio
    fn (state: vmm_state, addr: u64, data: bytes, is_write: bool) -> result[bytes, string]
    + dispatches an MMIO access to the device that owns the address range
    - returns error when no device claims the address
    # mmio_dispatch
  vmm.inject_irq
    fn (state: vmm_state, line: u32) -> result[void, string]
    + raises an interrupt line to the guest
    - returns error when the line is out of range
    # interrupts
    -> std.io.eventfd_signal
  vmm.shutdown
    fn (state: vmm_state) -> void
    + unmaps guest memory and releases all device and vCPU handles
    # teardown
    -> std.mm.unmap_guest_memory
