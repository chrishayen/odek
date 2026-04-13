# Requirement: "a library for installing and managing scripts as native windows services with event log integration"

Generates a windows service descriptor from a specification and installs, starts, stops, or removes it through a pluggable service control host.

std
  std.fs
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path atomically
      # filesystem
    std.fs.delete
      @ (path: string) -> result[void, string]
      + removes the file at the given path
      - returns error when the path does not exist
      # filesystem

winservice
  winservice.new_spec
    @ (name: string, display_name: string, exec_path: string) -> service_spec
    + creates a spec with an internal name, a display name, and the executable to run
    # construction
  winservice.with_event_log
    @ (spec: service_spec, source: string) -> service_spec
    + configures the event log source name used for service output
    # configuration
  winservice.render_descriptor
    @ (spec: service_spec) -> string
    + renders the spec as a service wrapper descriptor document
    # rendering
  winservice.install
    @ (spec: service_spec, descriptor_dir: string, host: service_host) -> result[void, string]
    + writes the descriptor to the directory and asks the host to register the service
    - returns error when the descriptor cannot be written
    # installation
    -> std.fs.write_all
  winservice.uninstall
    @ (name: string, descriptor_dir: string, host: service_host) -> result[void, string]
    + asks the host to unregister the service and deletes the descriptor file
    - returns error when the descriptor is missing
    # removal
    -> std.fs.delete
