# Requirement: "a framework for robotics, physical computing, and connected devices"

A robot is a bag of drivers attached to adaptors (the transport to hardware), plus a work loop that runs user-supplied behaviors. The project surface is a handful of small, composable pieces.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      @ (ms: i64) -> void
      + blocks the caller for the given duration
      # time
  std.io
    std.io.read_bytes
      @ (handle: i32, n: i32) -> result[bytes, string]
      + reads up to n bytes from an opaque handle
      - returns error when the handle is closed
      # io
    std.io.write_bytes
      @ (handle: i32, data: bytes) -> result[i32, string]
      + writes bytes to an opaque handle and returns the count written
      - returns error on transport failure
      # io

robotics
  robotics.new_adaptor
    @ (name: string, transport_handle: i32) -> adaptor_state
    + wraps a transport handle as a named adaptor
    # construction
  robotics.new_driver
    @ (name: string, adaptor: adaptor_state, address: i32) -> driver_state
    + binds a driver to an adaptor at a specific hardware address
    # construction
  robotics.driver_read
    @ (driver: driver_state, n_bytes: i32) -> result[bytes, string]
    + reads n_bytes from the underlying device
    # driver_io
    -> std.io.read_bytes
  robotics.driver_write
    @ (driver: driver_state, data: bytes) -> result[void, string]
    + writes data to the underlying device
    # driver_io
    -> std.io.write_bytes
  robotics.new_robot
    @ (name: string) -> robot_state
    + creates an empty robot with a display name
    # construction
  robotics.attach_driver
    @ (robot: robot_state, driver: driver_state) -> robot_state
    + registers a driver with the robot so the work loop can reach it
    # wiring
  robotics.register_event_handler
    @ (robot: robot_state, driver_name: string, event: string, handler: fn(bytes) -> void) -> robot_state
    + subscribes a handler to events emitted by a named driver
    # events
  robotics.emit_event
    @ (robot: robot_state, driver_name: string, event: string, payload: bytes) -> void
    + dispatches an event to every registered handler for (driver, event)
    # events
  robotics.schedule_every
    @ (robot: robot_state, interval_ms: i64, work: fn(robot_state) -> void) -> robot_state
    + registers a recurring job that runs on the work loop
    # scheduling
  robotics.run
    @ (robot: robot_state) -> void
    + drives the work loop until stopped: runs scheduled jobs, polls drivers, dispatches events
    # runtime
    -> std.time.now_millis
    -> std.time.sleep_millis
  robotics.stop
    @ (robot: robot_state) -> void
    + signals the work loop to exit after the current tick
    # runtime
