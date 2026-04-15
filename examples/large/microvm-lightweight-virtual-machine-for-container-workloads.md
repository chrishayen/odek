# Requirement: "a lightweight virtual machine for container workloads"

Configures and launches minimal virtual machines with vcpus, memory, a kernel image, a root filesystem, and network interfaces. Exposes a small control API.

std
  std.fs
    std.fs.exists
      fn (path: string) -> bool
      + returns true when path exists on the host
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the full contents of the file
      - returns error when the file cannot be read
      # filesystem
  std.vmm
    std.vmm.create_vm
      fn (vcpus: i32, memory_mib: i32) -> result[vm_handle, string]
      + allocates a new virtual machine with the requested cpu and memory
      - returns error when the host cannot satisfy the request
      # virtualization
    std.vmm.load_kernel
      fn (vm: vm_handle, image: bytes, cmdline: string) -> result[void, string]
      + loads the kernel image and sets the boot cmdline
      - returns error on invalid image format
      # virtualization
    std.vmm.attach_block
      fn (vm: vm_handle, host_path: string, is_root: bool) -> result[void, string]
      + attaches a block device backed by a host file
      - returns error when host_path does not exist
      # virtualization
    std.vmm.attach_net
      fn (vm: vm_handle, tap_name: string) -> result[void, string]
      + attaches a tap-backed network interface
      # virtualization
    std.vmm.start
      fn (vm: vm_handle) -> result[void, string]
      + boots the guest
      - returns error when required devices are missing
      # virtualization
    std.vmm.stop
      fn (vm: vm_handle) -> result[void, string]
      + halts the guest
      # virtualization

microvm
  microvm.new_config
    fn () -> vm_config
    + creates an empty config with zero vcpus and zero memory
    # construction
  microvm.set_machine
    fn (cfg: vm_config, vcpus: i32, memory_mib: i32) -> result[vm_config, string]
    + records vcpu and memory requirements
    - returns error when vcpus or memory_mib is not positive
    # configuration
  microvm.set_kernel
    fn (cfg: vm_config, path: string, cmdline: string) -> result[vm_config, string]
    + records the kernel image path and boot cmdline
    - returns error when path does not exist
    # configuration
    -> std.fs.exists
  microvm.add_drive
    fn (cfg: vm_config, host_path: string, is_root: bool) -> result[vm_config, string]
    + adds a block device backed by the host file
    - returns error when host_path does not exist
    # configuration
    -> std.fs.exists
  microvm.add_network
    fn (cfg: vm_config, tap_name: string) -> vm_config
    + adds a tap-backed network interface
    # configuration
  microvm.launch
    fn (cfg: vm_config) -> result[vm_handle, string]
    + creates the vm, loads the kernel, attaches devices, and starts it
    - returns error when any stage fails
    # lifecycle
    -> std.fs.read_all
    -> std.vmm.create_vm
    -> std.vmm.load_kernel
    -> std.vmm.attach_block
    -> std.vmm.attach_net
    -> std.vmm.start
  microvm.shutdown
    fn (vm: vm_handle) -> result[void, string]
    + stops the guest cleanly
    # lifecycle
    -> std.vmm.stop
  microvm.handle_control
    fn (cfg: vm_config, verb: string, path: string, body: string) -> result[vm_config, string]
    + applies PUT /machine-config, /boot-source, /drives/{id}, /network-interfaces to the config
    - returns error on unknown route or malformed body
    # control_api
