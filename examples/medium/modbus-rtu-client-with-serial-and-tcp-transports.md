# Requirement: "a Modbus-RTU client supporting both serial and TCP transports"

Encode and decode Modbus function codes 1, 3, 5, 6, 15, 16 with a CRC16 frame layer for serial and an MBAP header for TCP.

std
  std.net
    std.net.dial_tcp
      fn (host: string, port: i32) -> result[conn_handle, string]
      + returns an outbound tcp connection
      - returns error when the host cannot be reached
      # network
    std.net.write
      fn (c: conn_handle, data: bytes) -> result[i32, string]
      + returns the number of bytes written
      # network
    std.net.read
      fn (c: conn_handle, max: i32) -> result[bytes, string]
      + returns up to max bytes from the connection
      # network
  std.serial
    std.serial.open
      fn (device: string, baud: i32) -> result[conn_handle, string]
      + returns an open serial port handle
      - returns error when the device cannot be opened
      # serial
    std.serial.write
      fn (c: conn_handle, data: bytes) -> result[i32, string]
      + returns the number of bytes written
      # serial
    std.serial.read
      fn (c: conn_handle, max: i32, timeout_ms: i32) -> result[bytes, string]
      + returns bytes received within the timeout
      - returns error on timeout
      # serial

modbus
  modbus.crc16
    fn (data: bytes) -> u16
    + returns the Modbus-RTU crc16 of the input
    # framing
  modbus.encode_read_holding
    fn (unit_id: u8, address: u16, quantity: u16) -> bytes
    + returns the pdu for function code 3
    # pdu
  modbus.encode_write_single
    fn (unit_id: u8, address: u16, value: u16) -> bytes
    + returns the pdu for function code 6
    # pdu
  modbus.encode_write_multiple
    fn (unit_id: u8, address: u16, values: list[u16]) -> result[bytes, string]
    + returns the pdu for function code 16
    - returns error when values is empty or exceeds 123 registers
    # pdu
  modbus.decode_response
    fn (pdu: bytes) -> result[modbus_response, string]
    + returns the parsed response body
    - returns error when the function code byte has the error flag set
    - returns error on truncated pdu
    # pdu
  modbus.frame_rtu
    fn (unit_id: u8, pdu: bytes) -> bytes
    + returns the rtu frame with crc16 appended
    # framing
  modbus.unframe_rtu
    fn (frame: bytes) -> result[bytes, string]
    + returns the pdu after checking crc16
    - returns error when the crc does not match
    # framing
  modbus.frame_tcp
    fn (transaction_id: u16, unit_id: u8, pdu: bytes) -> bytes
    + returns the tcp frame with the mbap header
    # framing
  modbus.unframe_tcp
    fn (frame: bytes) -> result[bytes, string]
    + returns the pdu after validating the mbap header length
    - returns error when the header is truncated
    # framing
  modbus.request_serial
    fn (c: conn_handle, unit_id: u8, pdu: bytes, timeout_ms: i32) -> result[bytes, string]
    + writes an rtu frame and returns the response pdu
    - returns error on crc mismatch
    # transport
    -> std.serial.write
    -> std.serial.read
  modbus.request_tcp
    fn (c: conn_handle, transaction_id: u16, unit_id: u8, pdu: bytes) -> result[bytes, string]
    + writes a tcp frame and returns the response pdu
    - returns error when the response transaction id does not match
    # transport
    -> std.net.write
    -> std.net.read
