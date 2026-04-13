# Requirement: "a cloud infrastructure orchestration library"

Manage compute, networking, and storage resources against a pluggable driver. Driver calls are abstracted so the library never talks to a specific cloud.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

cloud
  cloud.new_project
    @ (name: string) -> project_state
    + returns an empty project with the given name
    # project
  cloud.register_driver
    @ (p: project_state, name: string, driver: driver_handle) -> project_state
    + binds a named driver for compute/network/storage calls
    - rejects a name already registered
    # drivers
  cloud.create_instance
    @ (p: project_state, image: string, size: string) -> result[instance_id, string]
    + asks the driver to launch an instance and records it in the project
    - returns error when the image name is empty
    # compute
    -> std.json.encode_object
  cloud.terminate_instance
    @ (p: project_state, id: instance_id) -> result[project_state, string]
    + removes the instance from the project and asks the driver to destroy it
    - returns error when the id is not tracked
    # compute
  cloud.list_instances
    @ (p: project_state) -> list[instance_id]
    + returns all tracked instances
    # compute
  cloud.create_network
    @ (p: project_state, cidr: string) -> result[network_id, string]
    + creates a virtual network with the given CIDR block
    - returns error when the CIDR is malformed
    # network
  cloud.attach_network
    @ (p: project_state, instance: instance_id, network: network_id) -> result[project_state, string]
    + attaches the instance to the network
    - returns error when either id is not tracked
    # network
  cloud.create_volume
    @ (p: project_state, size_gib: i64) -> result[volume_id, string]
    + creates a storage volume
    - returns error when size is not positive
    # storage
  cloud.attach_volume
    @ (p: project_state, volume: volume_id, instance: instance_id) -> result[project_state, string]
    + attaches the volume to the instance
    - returns error when the volume is already attached
    # storage
  cloud.describe_instance
    @ (p: project_state, id: instance_id) -> result[map[string, string], string]
    + returns instance attributes as a flat map
    - returns error when the id is not tracked
    # introspection
    -> std.json.parse_object
  cloud.snapshot_state
    @ (p: project_state) -> string
    + returns a JSON encoding of the entire project
    -> std.json.encode_object
    -> std.time.now_seconds
    # persistence
