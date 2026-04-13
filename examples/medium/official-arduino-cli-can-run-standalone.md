# Requirement: "a library for managing microcontroller sketches: compile, upload, and board discovery"

Embeddable library for working with microcontroller toolchains. Uses thin std primitives for process invocation and serial port enumeration.

std
  std.process
    std.process.run
      @ (cmd: string, args: list[string]) -> result[process_output, string]
      + runs a command and returns stdout, stderr, and exit code
      - returns error when the binary cannot be located
      # process
  std.serial
    std.serial.list_ports
      @ () -> list[serial_port]
      + returns all serial ports currently attached to the host
      # hardware
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full contents of the file
      - returns error when the file is missing
      # filesystem

sketch
  sketch.load
    @ (path: string) -> result[sketch_state, string]
    + loads a sketch directory into an internal representation
    - returns error when the main file is missing
    # loading
    -> std.fs.read_all
  sketch.list_boards
    @ () -> list[board_info]
    + returns every supported board with its fqbn and name
    # board_catalog
  sketch.discover_devices
    @ () -> list[attached_device]
    + returns attached boards with their port and detected fqbn
    # discovery
    -> std.serial.list_ports
  sketch.compile
    @ (state: sketch_state, fqbn: string, build_dir: string) -> result[compiled_artifact, string]
    + returns artifact paths on success
    - returns error with diagnostics when compilation fails
    # compilation
    -> std.process.run
  sketch.upload
    @ (artifact: compiled_artifact, port: string, fqbn: string) -> result[void, string]
    + writes the artifact to the board at the given port
    - returns error when the port is unavailable
    # upload
    -> std.process.run
  sketch.install_library
    @ (name: string, version: string) -> result[void, string]
    + installs a sketch library by name and version
    - returns error when the library or version is unknown
    # library_management
    -> std.process.run
