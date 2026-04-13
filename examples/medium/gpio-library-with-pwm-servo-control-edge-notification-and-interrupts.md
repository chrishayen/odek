# Requirement: "a GPIO library with PWM, servo control, edge notification, and interrupt handling"

A pin abstraction with digital, PWM, and servo modes, plus a callback-based edge-event stream. Hardware access is isolated behind std primitives.

std
  std.hardware
    std.hardware.mmio_write
      @ (address: u64, value: u32) -> void
      + writes a 32-bit value to a memory-mapped io register
      # hardware
    std.hardware.mmio_read
      @ (address: u64) -> u32
      + reads a 32-bit value from a memory-mapped io register
      # hardware
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic nanosecond timestamp
      # time

gpio
  gpio.open_pin
    @ (pin: i32, mode: i32) -> result[pin_handle, string]
    + returns a handle configured for input, output, or alt function
    - returns error when the pin number is outside the supported range
    # pin_setup
    -> std.hardware.mmio_write
  gpio.digital_write
    @ (handle: pin_handle, level: bool) -> void
    + drives the pin high or low
    # digital_output
    -> std.hardware.mmio_write
  gpio.digital_read
    @ (handle: pin_handle) -> bool
    + returns the current logic level of the pin
    # digital_input
    -> std.hardware.mmio_read
  gpio.set_pwm
    @ (handle: pin_handle, frequency_hz: i32, duty_cycle: f32) -> result[void, string]
    + configures hardware pwm with the given frequency and duty cycle in [0.0, 1.0]
    - returns error when duty cycle is out of range
    # pwm
    -> std.hardware.mmio_write
  gpio.set_servo_pulse
    @ (handle: pin_handle, pulse_us: i32) -> result[void, string]
    + emits a 50hz pwm with the given pulse width in microseconds
    - returns error when pulse_us is outside 500..2500
    # servo_control
    -> std.hardware.mmio_write
  gpio.watch_edges
    @ (handle: pin_handle, edge: i32) -> edge_watcher
    + returns a watcher that records rising, falling, or both edges with timestamps
    # edge_detection
    -> std.time.now_nanos
  gpio.poll_event
    @ (watcher: edge_watcher) -> optional[edge_event]
    + returns the next edge event with level and timestamp, or none if the queue is empty
    # event_dispatch
