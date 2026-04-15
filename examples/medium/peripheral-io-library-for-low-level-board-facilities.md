# Requirement: "a peripherals I/O library to interface with low-level board facilities"

Exposes GPIO, I2C, and SPI operations on a single-board computer through a thin typed layer.

std
  std.fs
    std.fs.read_text
      fn (path: string) -> result[string, string]
      + returns the contents of a sysfs-style file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_text
      fn (path: string, content: string) -> result[void, string]
      + writes the contents of a sysfs-style file
      - returns error when the file cannot be written
      # filesystem
  std.time
    std.time.sleep_micros
      fn (us: i64) -> void
      + suspends the calling context for the given number of microseconds
      # time

periph
  periph.list_gpio_pins
    fn () -> result[list[gpio_pin], string]
    + returns the gpio pins exposed by the current board
    - returns error when the board cannot be detected
    # discovery
  periph.set_pin_direction
    fn (pin: gpio_pin, direction: string) -> result[void, string]
    + configures the pin as "in" or "out"
    - returns error when direction is neither "in" nor "out"
    # gpio
    -> std.fs.write_text
  periph.read_pin
    fn (pin: gpio_pin) -> result[bool, string]
    + returns the current logic level of the pin
    - returns error when the pin is not configured as input
    # gpio
    -> std.fs.read_text
  periph.write_pin
    fn (pin: gpio_pin, level: bool) -> result[void, string]
    + drives the pin high or low
    - returns error when the pin is not configured as output
    # gpio
    -> std.fs.write_text
  periph.i2c_open
    fn (bus: i32, address: u8) -> result[i2c_device, string]
    + opens an I2C device at the given 7-bit address on the numbered bus
    - returns error when the bus does not exist
    # i2c
  periph.i2c_read
    fn (dev: i2c_device, n: i32) -> result[bytes, string]
    + reads n bytes from the device
    - returns error on bus error or NACK
    # i2c
  periph.i2c_write
    fn (dev: i2c_device, data: bytes) -> result[void, string]
    + writes data to the device
    - returns error on bus error or NACK
    # i2c
  periph.spi_open
    fn (bus: i32, chip_select: i32, speed_hz: i32) -> result[spi_device, string]
    + opens an SPI device at the given bus and chip select
    - returns error when the bus does not exist
    # spi
  periph.spi_transfer
    fn (dev: spi_device, tx: bytes) -> result[bytes, string]
    + performs a full-duplex transfer of the same length as tx
    - returns error on transfer failure
    # spi
    -> std.time.sleep_micros
