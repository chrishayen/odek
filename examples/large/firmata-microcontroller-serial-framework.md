# Requirement: "a framework for controlling microcontrollers over a Firmata-style serial protocol"

Wraps a serial transport in a typed API for reading/writing digital and analog pins, PWM, and querying firmware capabilities. The framing and sysex primitives live in std; the project layer exposes pin-oriented operations and a board state cache.

std
  std.serial
    std.serial.open
      fn (port: string, baud: i32) -> result[serial_handle, string]
      + opens a serial port at the given baud rate
      - returns error when port cannot be opened
      # io
    std.serial.read_byte
      fn (handle: serial_handle) -> result[u8, string]
      + reads one byte, blocking until available
      - returns error on disconnect
      # io
    std.serial.write_bytes
      fn (handle: serial_handle, data: bytes) -> result[void, string]
      + writes all bytes to the serial port
      - returns error on disconnect
      # io
  std.bytes
    std.bytes.seven_bit_encode
      fn (data: bytes) -> bytes
      + encodes arbitrary bytes as pairs of 7-bit LSB/MSB values
      # encoding
    std.bytes.seven_bit_decode
      fn (data: bytes) -> bytes
      + reverses seven_bit_encode
      # encoding

firmata
  firmata.connect
    fn (port: string) -> result[board_state, string]
    + opens the port, reads the firmware report, and returns a ready board_state
    - returns error when handshake fails
    # connection
    -> std.serial.open
  firmata.read_message
    fn (state: board_state) -> result[firmata_message, string]
    + reads the next complete message (command, sysex, or protocol version)
    - returns error on unexpected byte
    # framing
    -> std.serial.read_byte
  firmata.encode_command
    fn (cmd: firmata_command) -> bytes
    + encodes a command into wire bytes
    # framing
  firmata.encode_sysex
    fn (sysex_id: u8, payload: bytes) -> bytes
    + wraps payload in the sysex start/end envelope
    # framing
    -> std.bytes.seven_bit_encode
  firmata.decode_sysex
    fn (data: bytes) -> result[tuple[u8, bytes], string]
    + returns sysex id and decoded payload
    - returns error when start or end markers are missing
    # framing
    -> std.bytes.seven_bit_decode
  firmata.set_pin_mode
    fn (state: board_state, pin: i32, mode: pin_mode) -> result[board_state, string]
    + sends a pin-mode command and records the new mode
    - returns error when pin does not support mode
    # pin_control
    -> std.serial.write_bytes
  firmata.digital_write
    fn (state: board_state, pin: i32, high: bool) -> result[board_state, string]
    + sends a digital-write command for the pin
    - returns error when pin is not in output mode
    # pin_control
  firmata.digital_read
    fn (state: board_state, pin: i32) -> result[bool, string]
    + returns the last reported digital value for the pin
    - returns error when no report has arrived
    # pin_control
  firmata.analog_write
    fn (state: board_state, pin: i32, value: i32) -> result[board_state, string]
    + sends a PWM write with value 0..255
    - returns error when value is out of range
    # pin_control
  firmata.analog_read
    fn (state: board_state, pin: i32) -> result[i32, string]
    + returns the last reported analog value (0..1023)
    # pin_control
  firmata.query_capabilities
    fn (state: board_state) -> result[capability_map, string]
    + sends the capability query sysex and parses the response
    - returns error on malformed response
    # discovery
  firmata.apply_incoming
    fn (state: board_state, msg: firmata_message) -> board_state
    + updates pin caches from a digital-port or analog-report message
    # state_tracking
  firmata.disconnect
    fn (state: board_state) -> result[void, string]
    + closes the underlying serial port
    # connection
