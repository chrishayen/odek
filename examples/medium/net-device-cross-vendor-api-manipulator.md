# Requirement: "a cross-vendor api to manipulate network devices"

A uniform driver interface for configuring network devices from multiple vendors. Drivers are pluggable; the core exposes a single facade.

std
  std.net
    std.net.ssh_connect
      fn (host: string, port: u16, user: string, password: string) -> result[ssh_session, string]
      + opens an ssh session to the device
      - returns error on auth failure or unreachable host
      # transport
    std.net.ssh_exec
      fn (session: ssh_session, command: string) -> result[string, string]
      + sends a command and returns combined stdout/stderr
      - returns error when the session is closed
      # transport
    std.net.ssh_close
      fn (session: ssh_session) -> void
      + releases the underlying socket
      # transport

net_device
  net_device.register_driver
    fn (vendor: string, driver: device_driver) -> void
    + associates a vendor identifier with a driver implementation
    + later calls to open with the same vendor resolve to this driver
    # driver_registry
  net_device.open
    fn (vendor: string, host: string, user: string, password: string) -> result[device_handle, string]
    + returns a connected handle bound to the driver registered for the vendor
    - returns error when no driver is registered for the vendor
    - returns error on connection failure
    # connection
    -> std.net.ssh_connect
  net_device.get_config
    fn (handle: device_handle) -> result[string, string]
    + retrieves the device's running configuration as text via the bound driver
    - returns error when the driver rejects the read
    # configuration_read
    -> std.net.ssh_exec
  net_device.load_merge_config
    fn (handle: device_handle, config: string) -> result[void, string]
    + stages configuration changes for later commit without replacing the running config
    - returns error on syntax rejection by the device
    # configuration_write
    -> std.net.ssh_exec
  net_device.commit
    fn (handle: device_handle) -> result[void, string]
    + atomically applies staged configuration changes
    - returns error when no changes are staged
    - returns error when the device rejects the commit
    # commit
    -> std.net.ssh_exec
  net_device.rollback
    fn (handle: device_handle) -> result[void, string]
    + reverts the most recent committed change
    - returns error when no rollback point is available
    # rollback
    -> std.net.ssh_exec
  net_device.close
    fn (handle: device_handle) -> void
    + releases the device connection
    # connection
    -> std.net.ssh_close
