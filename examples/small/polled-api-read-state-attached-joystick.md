# Requirement: "a polled API for reading joystick state"

Open a joystick device by index and sample its axes and buttons on demand.

std: (all units exist)

joystick
  joystick.open
    @ (index: i32) -> result[device_state, string]
    + opens the joystick at the given index
    - returns error when no device exists at index
    # construction
  joystick.read_state
    @ (state: device_state) -> result[reading, string]
    + returns a snapshot of all axes and buttons at the current moment
    - returns error when the device has been disconnected
    # polling
  joystick.axis_count
    @ (state: device_state) -> i32
    + returns the number of axes the device exposes
    # introspection
  joystick.button_count
    @ (state: device_state) -> i32
    + returns the number of buttons the device exposes
    # introspection
  joystick.close
    @ (state: device_state) -> void
    + releases the device handle
    # teardown
