# Requirement: "a library for driving microcontroller peripherals from a small embedded runtime"

The source item was a tutorial, not a library; interpreted as a tiny peripheral-driver library for GPIO, PWM, and ADC on a microcontroller.

std: (all units exist)

mcu
  mcu.gpio_set_mode
    @ (pin: i32, mode: pin_mode) -> result[void, string]
    + configures a pin as input, output, or input-pull-up
    - returns error when pin is out of range
    # gpio
  mcu.gpio_write
    @ (pin: i32, high: bool) -> result[void, string]
    + drives an output pin high or low
    # gpio
  mcu.gpio_read
    @ (pin: i32) -> result[bool, string]
    + reads the logical level of an input pin
    # gpio
  mcu.pwm_write
    @ (pin: i32, duty_0_to_255: i32) -> result[void, string]
    + sets PWM duty cycle on a PWM-capable pin
    - returns error when the pin does not support PWM
    # pwm
  mcu.adc_read
    @ (channel: i32) -> result[i32, string]
    + returns a 12-bit sample from the ADC channel
    # adc
