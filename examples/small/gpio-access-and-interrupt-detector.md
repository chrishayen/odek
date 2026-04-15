# Requirement: "a GPIO access and interrupt detection library"

Read and write digital pins and register callbacks fired on edge transitions. Hardware I/O flows through thin std primitives so tests can substitute a fake device.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file as bytes
      - returns error when the path does not exist
      # io
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the file, truncating any existing content
      - returns error when the path is not writable
      # io

gpio
  gpio.open
    fn (pin_number: i32, direction: string) -> result[pin_handle, string]
    + opens the pin and configures it as "in" or "out"
    - returns error when direction is not "in" or "out"
    - returns error when the pin is already exported
    # construction
    -> std.fs.write_all
  gpio.read
    fn (handle: pin_handle) -> result[bool, string]
    + returns true for high, false for low
    - returns error when the pin was opened as an output
    # input
    -> std.fs.read_all
  gpio.write
    fn (handle: pin_handle, value: bool) -> result[void, string]
    + drives the pin high or low
    - returns error when the pin was opened as an input
    # output
    -> std.fs.write_all
  gpio.watch
    fn (handle: pin_handle, edge: string, callback: closure[bool]) -> result[void, string]
    + invokes the callback with the new pin value on each matching edge
    ? edge is "rising", "falling", or "both"
    - returns error on unknown edge string
    # interrupts
  gpio.close
    fn (handle: pin_handle) -> void
    + unexports the pin and releases any registered watchers
    # lifecycle
