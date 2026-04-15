# Requirement: "a library for mounting a filesystem image over a network share"

Exposes a local filesystem image as a network share that a client host can mount. The project layer orchestrates; primitives live in std.

std
  std.fs
    std.fs.exists
      fn (path: string) -> bool
      + true when the path exists
      # filesystem
  std.net
    std.net.free_port
      fn () -> result[u16, string]
      + returns an unused local TCP port
      - returns error when no port can be bound
      # network
    std.net.wait_listening
      fn (host: string, port: u16, timeout_ms: i32) -> result[void, string]
      + returns void once a TCP connection succeeds
      - returns error on timeout
      # network
  std.process
    std.process.spawn
      fn (program: string, args: list[string]) -> result[process_handle, string]
      + starts the process and returns a handle
      - returns error when the program is not found
      # process
    std.process.wait
      fn (handle: process_handle) -> result[i32, string]
      + returns the exit code
      # process

netmount
  netmount.prepare_export
    fn (image_path: string, mount_point: string) -> result[export_spec, string]
    + returns an export spec with image path, internal mount, and a free host port
    - returns error when the image does not exist
    # preparation
    -> std.fs.exists
    -> std.net.free_port
  netmount.start_server
    fn (spec: export_spec) -> result[server_handle, string]
    + spawns the share server and waits for it to listen
    - returns error when the server exits before listening
    # serving
    -> std.process.spawn
    -> std.net.wait_listening
  netmount.mount_client
    fn (host: string, port: u16, local_path: string) -> result[void, string]
    + mounts the remote share at the local path
    - returns error when the mount call fails
    # mounting
    -> std.process.spawn
    -> std.process.wait
  netmount.unmount_client
    fn (local_path: string) -> result[void, string]
    + unmounts the local path
    - returns error when the path is not mounted
    # mounting
    -> std.process.spawn
    -> std.process.wait
  netmount.stop_server
    fn (handle: server_handle) -> result[void, string]
    + terminates the server and waits for exit
    # teardown
    -> std.process.wait
