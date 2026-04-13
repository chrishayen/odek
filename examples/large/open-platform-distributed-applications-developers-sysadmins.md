# Requirement: "a container runtime library for building, running, and networking isolated processes"

Images, containers, and networks are managed independently. Storage and process isolation go through std primitives.

std
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest of data
      # cryptography
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      - returns error when the file is missing
      # filesystem
    std.fs.write_file
      @ (path: string, data: bytes) -> result[void, string]
      # filesystem
    std.fs.mount
      @ (source: string, target: string, options: map[string, string]) -> result[void, string]
      + mounts source at target with the given options
      - returns error when the kernel rejects the mount
      # filesystem
  std.process
    std.process.spawn_isolated
      @ (argv: list[string], namespaces: list[string], root: string) -> result[process_handle, string]
      + launches a process in the requested namespaces with the given root filesystem
      - returns error when namespace creation fails
      # process
    std.process.wait
      @ (handle: process_handle) -> result[i32, string]
      + blocks until the process exits and returns its exit code
      # process
  std.net
    std.net.create_bridge
      @ (name: string) -> result[void, string]
      - returns error when the bridge already exists
      # networking

runtime
  runtime.build_image
    @ (layers: list[bytes], config: map[string, string]) -> result[image_id, string]
    + returns a content-addressed identifier for the assembled image
    # image_build
    -> std.crypto.sha256
    -> std.fs.write_file
  runtime.pull_image
    @ (ref: string, fetch: layer_fetcher) -> result[image_id, string]
    + fetches missing layers and stores them locally
    - returns error when the reference cannot be resolved
    # image_pull
    -> std.fs.write_file
  runtime.list_images
    @ () -> list[image_summary]
    + returns every stored image with its tag and size
    # listing
  runtime.remove_image
    @ (id: image_id) -> result[void, string]
    - returns error when containers reference the image
    # deletion
  runtime.create_container
    @ (image: image_id, spec: container_spec) -> result[container_id, string]
    + prepares a rootfs and returns the container id
    - returns error when the image is unknown
    # creation
    -> std.fs.mount
  runtime.start_container
    @ (id: container_id) -> result[void, string]
    + launches the container's main process in an isolated namespace
    - returns error when the container is already running
    # lifecycle
    -> std.process.spawn_isolated
  runtime.stop_container
    @ (id: container_id, timeout_sec: i32) -> result[void, string]
    + sends a graceful stop signal and escalates after timeout_sec
    # lifecycle
    -> std.process.wait
  runtime.list_containers
    @ (only_running: bool) -> list[container_summary]
    # listing
  runtime.logs
    @ (id: container_id) -> result[bytes, string]
    + returns all captured stdout/stderr for the container
    - returns error when the container is unknown
    # observability
    -> std.fs.read_all
  runtime.create_network
    @ (name: string, cidr: string) -> result[network_id, string]
    - returns error when the CIDR is malformed
    # networking
    -> std.net.create_bridge
  runtime.attach_network
    @ (container: container_id, network: network_id, ip: optional[string]) -> result[string, string]
    + returns the IP assigned to the container on the network
    - returns error when no address is available
    # networking
  runtime.exec_in_container
    @ (id: container_id, argv: list[string]) -> result[process_handle, string]
    + runs a command inside a running container
    - returns error when the container is not running
    # execution
    -> std.process.spawn_isolated
