# Requirement: "a container runtime lifecycle manager"

Manages the lifecycle of local container runtime instances: create, start, stop, and list. The host interaction is abstracted behind a hypervisor handle type.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file's contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, replacing any existing content
      # filesystem
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[map[string, string], string]
      + parses a YAML document as a flat string-to-string map
      - returns error on invalid YAML
      # serialization
    std.yaml.encode
      @ (data: map[string, string]) -> string
      + serializes a string-to-string map as YAML
      # serialization

runtime_manager
  runtime_manager.create_instance
    @ (name: string, cpus: i32, memory_mb: i32, disk_gb: i32) -> result[instance_spec, string]
    + returns a spec with a generated disk image path and the given resource limits
    - returns error when cpus, memory_mb, or disk_gb is not positive
    # construction
  runtime_manager.save_spec
    @ (spec: instance_spec, config_dir: string) -> result[void, string]
    + writes the spec to "{config_dir}/{name}.yaml"
    # persistence
    -> std.yaml.encode
    -> std.fs.write_all
  runtime_manager.load_spec
    @ (name: string, config_dir: string) -> result[instance_spec, string]
    + reads and parses the spec for a named instance
    - returns error when the file is missing
    # persistence
    -> std.fs.read_all
    -> std.yaml.parse
  runtime_manager.start
    @ (spec: instance_spec, hypervisor: hypervisor_handle) -> result[instance_handle, string]
    + boots the instance via the hypervisor handle and returns a running handle
    - returns error when the instance is already running
    # lifecycle
  runtime_manager.stop
    @ (handle: instance_handle, hypervisor: hypervisor_handle) -> result[void, string]
    + gracefully shuts down the instance
    - returns error when the instance is not running
    # lifecycle
  runtime_manager.list_instances
    @ (config_dir: string) -> result[list[instance_spec], string]
    + returns all specs found in the config directory
    # discovery
  runtime_manager.delete_instance
    @ (name: string, config_dir: string) -> result[void, string]
    + removes the spec file and any associated disk image
    - returns error when the instance is still running
    # lifecycle
