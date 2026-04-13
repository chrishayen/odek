# Requirement: "a library for controlling and monitoring keyboard and pointer input devices"

The project layer exposes synthesis and event-stream subscription; actual device access is behind an opaque device handle.

std: (all units exist)

input_devices
  input_devices.open
    @ () -> result[device_handle, string]
    + returns a handle bound to the host input subsystem
    - returns error when access is denied
    # construction
  input_devices.press_key
    @ (handle: device_handle, key: key_code) -> result[void, string]
    + synthesizes a key-down event
    - returns error when the key code is unknown
    # synthesis
  input_devices.release_key
    @ (handle: device_handle, key: key_code) -> result[void, string]
    + synthesizes a key-up event
    # synthesis
  input_devices.move_pointer
    @ (handle: device_handle, x: i32, y: i32) -> result[void, string]
    + moves the pointer to absolute screen coordinates
    # synthesis
  input_devices.click_button
    @ (handle: device_handle, button: pointer_button) -> result[void, string]
    + synthesizes a full press-release of a pointer button
    # synthesis
  input_devices.next_event
    @ (handle: device_handle) -> result[optional[input_event], string]
    + returns the next pending input event, or none when the queue is empty
    - returns error when the handle has been closed
    # monitoring
  input_devices.close
    @ (handle: device_handle) -> result[void, string]
    + releases the device handle
    # lifecycle
