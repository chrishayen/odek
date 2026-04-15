# Requirement: "a library for reading from and writing to serial ports"

Opens named serial devices with a baud configuration, reads and writes bytes, and closes them. Flow control details are passed through as opaque options.

std: (all units exist)

serial
  serial.list_ports
    fn () -> result[list[string], string]
    + returns the names of serial devices currently available on the host
    - returns error when the host enumeration fails
    # discovery
  serial.open
    fn (device: string, baud: i32, data_bits: i32, stop_bits: i32, parity: string) -> result[port, string]
    + opens the named device with the given line settings
    - returns error when the device does not exist
    - returns error when the baud rate is unsupported
    ? parity is one of "none", "even", "odd"
    # lifecycle
  serial.read
    fn (p: port, max_bytes: i32, timeout_ms: i32) -> result[bytes, string]
    + returns up to max_bytes read from the port
    + returns an empty byte string when the timeout elapses with no data
    - returns error when the port is closed
    # io
  serial.write
    fn (p: port, data: bytes) -> result[i32, string]
    + returns the number of bytes actually written
    - returns error when the port is closed
    # io
  serial.flush
    fn (p: port) -> result[void, string]
    + waits for all buffered output bytes to be transmitted
    - returns error when the port is closed
    # io
  serial.close
    fn (p: port) -> result[void, string]
    + closes the port and releases the underlying device handle
    + safe to call twice; second call is a no-op
    # lifecycle
