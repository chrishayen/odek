# Requirement: "a power management library exposing a message-bus interface"

Reads and writes power policy knobs (profile, CPU governor, backlight), exposes them over a message bus, and emits change events. Filesystem and bus primitives are std.

std
  std.fs
    std.fs.read_file
      @ (path: string) -> result[bytes, string]
      + reads file contents
      - returns error when file is missing or unreadable
      # filesystem
    std.fs.write_file
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to file, replacing existing contents
      - returns error when path is not writable
      # filesystem
  std.bus
    std.bus.connect
      @ (address: string) -> result[bus_conn, string]
      + opens a message-bus connection
      - returns error when the socket cannot be reached
      # messaging
    std.bus.export_object
      @ (conn: bus_conn, path: string, methods: list[string]) -> result[void, string]
      + exposes an object path with the given method names
      # messaging
    std.bus.emit_signal
      @ (conn: bus_conn, path: string, name: string, payload: bytes) -> result[void, string]
      + emits a signal on an object path
      # messaging
  std.encoding
    std.encoding.parse_int
      @ (s: string) -> result[i64, string]
      + parses a base-10 integer
      - returns error on non-numeric input
      # encoding

power
  power.start
    @ (bus_addr: string) -> result[power_state, string]
    + connects to the bus and exposes the power object at /org/power/Management
    - returns error when the bus is unreachable
    # lifecycle
    -> std.bus.connect
    -> std.bus.export_object
  power.get_profile
    @ (state: power_state) -> result[string, string]
    + returns the active power profile name ("balanced", "performance", "battery")
    # query
    -> std.fs.read_file
  power.set_profile
    @ (state: power_state, profile: string) -> result[void, string]
    + applies a power profile and emits a ProfileChanged signal
    - returns error when profile is not one of the recognized names
    # mutation
    -> std.fs.write_file
    -> std.bus.emit_signal
  power.get_cpu_governor
    @ (state: power_state, cpu_index: i32) -> result[string, string]
    + returns the governor of a specific cpu
    # query
    -> std.fs.read_file
  power.set_cpu_governor
    @ (state: power_state, governor: string) -> result[void, string]
    + sets the governor for every cpu
    - returns error when governor is unknown
    # mutation
    -> std.fs.write_file
  power.get_backlight
    @ (state: power_state) -> result[i32, string]
    + returns the current backlight percentage
    # query
    -> std.fs.read_file
    -> std.encoding.parse_int
  power.set_backlight
    @ (state: power_state, percent: i32) -> result[void, string]
    + sets backlight percentage and emits BacklightChanged
    - returns error when percent is outside [0, 100]
    # mutation
    -> std.fs.write_file
    -> std.bus.emit_signal
  power.on_battery
    @ (state: power_state) -> result[bool, string]
    + returns true when the system is running on battery
    # query
    -> std.fs.read_file
  power.subscribe_events
    @ (state: power_state, sink: string) -> result[void, string]
    + registers a subscriber path for power-change signals
    # subscription
  power.stop
    @ (state: power_state) -> result[void, string]
    + releases the bus name and closes the connection
    # lifecycle
