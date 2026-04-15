# Requirement: "a desktop automation library for mouse, keyboard, and screen reading"

Exposes primitives to move the cursor, synthesize clicks and keystrokes, and sample screen pixels. Assumes an underlying OS bridge.

std: (all units exist)

desktop
  desktop.move_mouse
    fn (x: i32, y: i32) -> result[void, string]
    + moves the cursor to absolute screen coordinates
    - returns error when coordinates are outside any monitor
    # mouse
  desktop.click
    fn (button: string) -> result[void, string]
    + presses and releases the given button at the cursor's current position
    - returns error for unknown button name
    # mouse
  desktop.type_text
    fn (text: string) -> result[void, string]
    + synthesizes keystrokes for each character in the string
    # keyboard
  desktop.press_key
    fn (key: string, modifiers: list[string]) -> result[void, string]
    + presses a named key with optional modifier keys held down
    - returns error for unknown key name
    # keyboard
  desktop.read_pixel
    fn (x: i32, y: i32) -> result[pixel_rgb, string]
    + returns the RGB value of the pixel at the given screen coordinate
    - returns error when the coordinate is off-screen
    # screen
  desktop.capture_region
    fn (x: i32, y: i32, w: i32, h: i32) -> result[bitmap_image, string]
    + captures a rectangular screen region as a bitmap
    - returns error when the rectangle extends beyond the screen
    # screen
  desktop.screen_size
    fn () -> pair_i32
    + returns the width and height of the primary display
    # screen
