# Requirement: "a library for advanced ANSI style and color support in terminal applications"

Produces ANSI escape sequences for styled text and detects terminal color support. Does not print.

std
  std.os
    std.os.get_env
      fn (name: string) -> optional[string]
      + returns the environment variable value when set
      # system
  std.io
    std.io.is_tty
      fn (fd: i32) -> bool
      + returns true when the file descriptor refers to a terminal
      # io

ansi_style
  ansi_style.detect_profile
    fn () -> color_profile
    + returns the highest color profile the environment supports (none, 16, 256, truecolor)
    + returns none when output is not a terminal
    # detection
    -> std.os.get_env
    -> std.io.is_tty
  ansi_style.rgb_to_ansi256
    fn (r: u8, g: u8, b: u8) -> u8
    + returns the closest 256-color palette index for the given RGB
    # color_conversion
  ansi_style.rgb_to_ansi16
    fn (r: u8, g: u8, b: u8) -> u8
    + returns the closest 16-color ANSI code
    # color_conversion
  ansi_style.foreground
    fn (profile: color_profile, r: u8, g: u8, b: u8) -> string
    + returns the ANSI sequence to set foreground to the given RGB, adapted to the profile
    + returns empty string when profile is none
    # sequences
    -> ansi_style.rgb_to_ansi256
    -> ansi_style.rgb_to_ansi16
  ansi_style.background
    fn (profile: color_profile, r: u8, g: u8, b: u8) -> string
    + returns the ANSI sequence to set background
    # sequences
    -> ansi_style.rgb_to_ansi256
    -> ansi_style.rgb_to_ansi16
  ansi_style.style
    fn (bold: bool, italic: bool, underline: bool) -> string
    + returns the ANSI sequence enabling the requested attributes
    + returns empty string when all attributes are false
    # sequences
  ansi_style.reset
    fn () -> string
    + returns the ANSI reset sequence
    # sequences
  ansi_style.wrap
    fn (text: string, prefix: string) -> string
    + returns prefix + text + reset
    # sequences
    -> ansi_style.reset
