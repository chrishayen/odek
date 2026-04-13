# Requirement: "reusable components for a microcontroller board"

Thin abstractions over pin I/O and timers. The library does not perform real hardware access; it returns descriptors a host layer consumes.

std: (all units exist)

mcu
  mcu.pin_mode
    @ (pin: u8, mode: u8) -> pin_state
    + returns a pin_state tagged with the requested mode (0=input, 1=output)
    - returns error-tagged state when pin is outside 0..19
    # pin_configuration
  mcu.digital_write
    @ (state: pin_state, high: bool) -> pin_state
    + returns a new state with the output level set
    - returns unchanged state when the pin is in input mode
    # digital_output
  mcu.digital_read
    @ (state: pin_state) -> bool
    + returns the latched level for an input pin
    # digital_input
  mcu.delay_ticks
    @ (prescaler: u16, ticks: u32) -> u64
    + returns the number of clock cycles represented by the given prescaler and ticks
    # timing
