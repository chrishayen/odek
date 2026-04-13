# Requirement: "a library for installing and managing scripts as unix system services with syslog logging"

Generates a systemd-style unit file from a specification and installs, starts, stops, or removes it through a pluggable service host.

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

unixservice
  unixservice.new_spec
    @ (name: string, exec_path: string) -> service_spec
    + creates a minimal spec with a unit name and the script or executable to run
    # construction
  unixservice.with_syslog
    @ (spec: service_spec, identifier: string) -> service_spec
    + routes standard output and error to syslog under the given identifier
    # configuration
  unixservice.render_unit
    @ (spec: service_spec) -> string
    + renders the spec as a system service unit file
    # rendering
  unixservice.install
    @ (spec: service_spec, unit_dir: string, host: service_host) -> result[void, string]
    + writes the rendered unit to the unit directory and asks the host to enable and start it
    - returns error when the unit cannot be written
    # installation
    -> std.fs.write_all
  unixservice.uninstall
    @ (name: string, unit_dir: string, host: service_host) -> result[void, string]
    + asks the host to stop and disable the service, then deletes the unit file
    - returns error when the unit is missing
    # removal
    -> std.fs.delete
