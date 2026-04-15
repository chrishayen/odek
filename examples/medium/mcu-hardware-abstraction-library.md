# Requirement: "a hardware abstraction library for embedded microcontrollers"

Peripheral drivers as data: GPIO, timers, UART. The library models register operations as pure functions over a memory-map state so logic can be tested without touching real hardware.

std: (all units exist)

mcu
  mcu.new_mmio
    fn () -> mmio_state
    + returns a zeroed memory-mapped IO state
    # construction
  mcu.read_reg
    fn (state: mmio_state, addr: u32) -> u32
    + returns the value at the register address
    # mmio
  mcu.write_reg
    fn (state: mmio_state, addr: u32, value: u32) -> mmio_state
    + stores value at the register address
    # mmio
  mcu.gpio_set_mode
    fn (state: mmio_state, pin: i32, mode: i32) -> mmio_state
    + configures a GPIO pin as input, output, or alternate function
    - returns unchanged state when pin is out of range
    # gpio
  mcu.gpio_write
    fn (state: mmio_state, pin: i32, high: bool) -> mmio_state
    + drives a pin high or low
    # gpio
  mcu.gpio_read
    fn (state: mmio_state, pin: i32) -> bool
    + returns the logical level of a pin
    # gpio
  mcu.timer_configure
    fn (state: mmio_state, timer_id: i32, period_ticks: u32) -> mmio_state
    + configures a hardware timer with the given period
    # timer
  mcu.timer_elapsed
    fn (state: mmio_state, timer_id: i32) -> bool
    + returns true when the timer has reached its configured period
    # timer
  mcu.uart_configure
    fn (state: mmio_state, uart_id: i32, baud: u32) -> mmio_state
    + configures a UART peripheral
    # uart
  mcu.uart_write_byte
    fn (state: mmio_state, uart_id: i32, byte: u8) -> mmio_state
    + writes a byte to the UART transmit register
    # uart
  mcu.uart_read_byte
    fn (state: mmio_state, uart_id: i32) -> optional[u8]
    + returns the next received byte when available
    - returns none when the receive register is empty
    # uart
