# Requirement: "a library for installing and managing scripts as background services on macOS"

Generates a launchd-style service descriptor from a specification and installs, starts, stops, or removes it through a pluggable service host.

std
  std.fs
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path atomically
      # filesystem
    std.fs.delete
      fn (path: string) -> result[void, string]
      + removes the file at the given path
      - returns error when the path does not exist
      # filesystem

macservice
  macservice.new_spec
    fn (label: string, script_path: string) -> service_spec
    + creates a minimal spec with a unique label and the script to run
    # construction
  macservice.with_log
    fn (spec: service_spec, stdout_path: string, stderr_path: string) -> service_spec
    + configures standard output and error log destinations
    # configuration
  macservice.render_descriptor
    fn (spec: service_spec) -> string
    + renders the spec as a launchd plist XML document
    # rendering
  macservice.install
    fn (spec: service_spec, descriptor_dir: string, host: service_host) -> result[void, string]
    + writes the rendered descriptor to the directory and asks the host to load it
    - returns error when the descriptor cannot be written
    # installation
    -> std.fs.write_all
  macservice.uninstall
    fn (label: string, descriptor_dir: string, host: service_host) -> result[void, string]
    + asks the host to unload the service and deletes the descriptor file
    - returns error when the descriptor is missing
    # removal
    -> std.fs.delete
