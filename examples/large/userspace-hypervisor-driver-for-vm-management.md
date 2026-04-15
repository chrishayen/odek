# Requirement: "a userspace library that drives a kernel hypervisor to configure, run, and manage virtual machines"

Opens a handle to the kernel hypervisor, configures vCPUs and memory, attaches emulated devices, and steps the VM through its lifecycle.

std
  std.syscall
    std.syscall.open_device
      fn (path: string) -> result[fd, string]
      + opens a character device and returns its handle
      - returns error when the device does not exist
      # syscall
    std.syscall.ioctl
      fn (fd: fd, cmd: u32, arg: bytes) -> result[bytes, string]
      + issues an ioctl and returns the reply buffer
      - returns error with errno text on failure
      # syscall
    std.syscall.mmap
      fn (fd: fd, length: u64) -> result[mapped_region, string]
      + maps a region from the file into the process address space
      # syscall
    std.syscall.close_fd
      fn (fd: fd) -> void
      + closes the file descriptor
      # syscall

vmm
  vmm.create
    fn (name: string) -> result[vm_handle, string]
    + asks the hypervisor to create a VM instance with the given name
    - returns error when a VM with the same name exists
    # lifecycle
    -> std.syscall.open_device
    -> std.syscall.ioctl
  vmm.destroy
    fn (vm: vm_handle) -> result[void, string]
    + releases the VM and any mapped memory
    # lifecycle
    -> std.syscall.ioctl
    -> std.syscall.close_fd
  vmm.set_memory
    fn (vm: vm_handle, bytes_count: u64) -> result[vm_handle, string]
    + reserves guest physical memory of the requested size and maps it
    - returns error when the size is not page-aligned
    # memory
    -> std.syscall.ioctl
    -> std.syscall.mmap
  vmm.load_image
    fn (vm: vm_handle, guest_phys: u64, image: bytes) -> result[void, string]
    + copies the image bytes into guest memory at the given physical address
    - returns error when the range exceeds allocated memory
    # memory
  vmm.add_vcpu
    fn (vm: vm_handle, vcpu_id: i32) -> result[vm_handle, string]
    + creates a vCPU with the given id
    - returns error when the id is already in use
    # vcpu
    -> std.syscall.ioctl
  vmm.set_register
    fn (vm: vm_handle, vcpu_id: i32, reg: string, value: u64) -> result[void, string]
    + sets a named register on the vCPU
    - returns error when the register name is not recognized
    # vcpu
    -> std.syscall.ioctl
  vmm.attach_device
    fn (vm: vm_handle, spec: device_spec) -> result[vm_handle, string]
    + wires an emulated device (serial, block, net) into the VM
    - returns error when the device kind is unknown
    # devices
  vmm.run_vcpu
    fn (vm: vm_handle, vcpu_id: i32) -> result[vm_exit, string]
    + enters the guest until the vCPU yields; returns the exit reason
    - returns error when the vCPU has not been initialized
    # execution
    -> std.syscall.ioctl
  vmm.pause
    fn (vm: vm_handle) -> result[void, string]
    + suspends all vCPUs
    # lifecycle
    -> std.syscall.ioctl
  vmm.resume
    fn (vm: vm_handle) -> result[void, string]
    + resumes all vCPUs
    # lifecycle
    -> std.syscall.ioctl
