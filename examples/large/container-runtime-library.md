# Requirement: "a container runtime library"

A library that manages container images, containers, and networks: pulls images to a local store, creates and starts containers from an image, and attaches them to a network. Actual syscalls and network I/O are the host's responsibility; this owns state and orchestration.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file
      # filesystem
    std.fs.make_dir
      fn (path: string) -> result[void, string]
      + creates a directory, including parents
      # filesystem
  std.hash
    std.hash.sha256_hex
      fn (data: bytes) -> string
      + returns the lowercase hex SHA-256 digest of the input
      # hashing
  std.encoding
    std.encoding.json_encode_any
      fn (value: json_value) -> string
      + encodes an arbitrary JSON value
      # serialization
    std.encoding.json_decode_any
      fn (raw: string) -> result[json_value, string]
      + decodes arbitrary JSON
      - returns error on invalid JSON
      # serialization
  std.random
    std.random.uuid_v4
      fn () -> string
      + returns a random UUID v4 string
      # random

container
  container.new_runtime
    fn (root: string) -> runtime_state
    + creates a runtime rooted at a storage directory
    # construction
    -> std.fs.make_dir
  container.add_image
    fn (state: runtime_state, reference: string, layers: list[bytes], config: string) -> result[runtime_state, string]
    + stores an image under a reference, computing its content-addressed digest
    - returns error when the config JSON is invalid
    # image_store
    -> std.hash.sha256_hex
    -> std.encoding.json_decode_any
    -> std.fs.write_all
  container.lookup_image
    fn (state: runtime_state, reference: string) -> result[image_record, string]
    + returns an image record by reference
    - returns error when the reference is unknown
    # image_store
  container.create_container
    fn (state: runtime_state, image_ref: string, name: string, command: list[string]) -> result[tuple[string, runtime_state], string]
    + creates a stopped container from an image and returns its id
    - returns error when the image reference is unknown
    - returns error when the container name is already in use
    # lifecycle
    -> std.random.uuid_v4
    -> std.fs.make_dir
  container.start
    fn (state: runtime_state, id: string) -> result[runtime_state, string]
    + transitions a container from created or stopped to running
    - returns error when the container is unknown
    - returns error when the container is already running
    # lifecycle
  container.stop
    fn (state: runtime_state, id: string) -> result[runtime_state, string]
    + transitions a running container to stopped
    - returns error when the container is unknown
    # lifecycle
  container.remove
    fn (state: runtime_state, id: string) -> result[runtime_state, string]
    + deletes a stopped container from the runtime
    - returns error when the container is running
    # lifecycle
  container.list
    fn (state: runtime_state) -> list[container_record]
    + returns a snapshot of all known containers
    # inventory
  container.create_network
    fn (state: runtime_state, name: string, cidr: string) -> result[runtime_state, string]
    + registers a named network with an address range
    - returns error when the CIDR is malformed
    # networking
  container.attach_network
    fn (state: runtime_state, container_id: string, network_name: string) -> result[tuple[string, runtime_state], string]
    + assigns an address from the network and attaches the container
    - returns error when the network or container is unknown
    - returns error when the network has no free addresses
    # networking
  container.persist_state
    fn (state: runtime_state) -> result[void, string]
    + writes the runtime's metadata to disk
    # persistence
    -> std.encoding.json_encode_any
    -> std.fs.write_all
  container.load_state
    fn (root: string) -> result[runtime_state, string]
    + reconstructs a runtime from disk
    - returns error when the metadata file is missing or corrupt
    # persistence
    -> std.fs.read_all
    -> std.encoding.json_decode_any
