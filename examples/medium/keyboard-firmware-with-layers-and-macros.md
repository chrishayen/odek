# Requirement: "a keyboard firmware library with layers and macros"

Scans a key matrix, applies layered keymaps, and produces HID usage reports. Hardware access is abstracted behind std primitives.

std
  std.gpio
    std.gpio.set_output
      fn (pin: i32, high: bool) -> void
      + drives the pin high or low
      # hardware
    std.gpio.read_input
      fn (pin: i32) -> bool
      + returns the current logical level of the pin
      # hardware
  std.time
    std.time.now_micros
      fn () -> i64
      + returns current monotonic time in microseconds
      # time

keyboard
  keyboard.scan_matrix
    fn (rows: list[i32], cols: list[i32]) -> list[key_event]
    + strobes each row and samples columns, returning pressed/released transitions
    + debounces by requiring two consistent samples
    # input
    -> std.gpio.set_output
    -> std.gpio.read_input
    -> std.time.now_micros
  keyboard.new_keymap
    fn (layers: list[list[u16]]) -> keymap
    + builds a keymap with the given layer definitions
    # configuration
  keyboard.resolve_key
    fn (map: keymap, active_layer: i32, row: i32, col: i32) -> u16
    + returns the HID usage code for a position on a layer
    ? transparent cells fall through to the lower layer
    # resolution
  keyboard.expand_macro
    fn (code: u16) -> list[u16]
    + returns the sequence of usage codes emitted by a macro key
    + returns a single-element list for non-macro keys
    # macros
  keyboard.build_report
    fn (pressed: list[u16]) -> bytes
    + encodes up to six currently pressed usages into an HID boot report
    + the first byte is the modifier mask
    # hid
  keyboard.step
    fn (state: keyboard_state) -> tuple[bytes, keyboard_state]
    + scans, resolves, and returns the next HID report
    # runtime
