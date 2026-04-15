# Requirement: "a library for talking to IoT sensor and actuator devices"

Generic register-based interface over an opaque bus (I2C-like) with typed read and write helpers.

std: (all units exist)

iodev
  iodev.open_bus
    fn (bus_path: string) -> result[bus_state, string]
    + opens an underlying device bus for reads and writes
    - returns error when the bus path cannot be opened
    # bus
  iodev.close_bus
    fn (bus: bus_state) -> void
    + releases bus resources
    # bus
  iodev.write_register
    fn (bus: bus_state, address: u8, register: u8, value: u8) -> result[void, string]
    + writes a byte to a register on the device at address
    - returns error on bus IO failure
    # io
  iodev.read_register
    fn (bus: bus_state, address: u8, register: u8) -> result[u8, string]
    + reads a byte from a register on the device at address
    - returns error on bus IO failure
    # io
  iodev.read_block
    fn (bus: bus_state, address: u8, register: u8, length: i32) -> result[bytes, string]
    + reads a consecutive block of bytes starting at register
    - returns error when length is not positive
    # io
  iodev.scan_addresses
    fn (bus: bus_state) -> result[list[u8], string]
    + probes each 7-bit address on the bus and returns those that acknowledge
    # discovery
