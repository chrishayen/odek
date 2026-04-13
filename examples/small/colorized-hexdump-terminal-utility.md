# Requirement: "a colorized hexdump formatter"

Renders bytes as an offset/hex/ASCII hexdump with color codes applied per byte class.

std: (all units exist)

hexdump
  hexdump.format
    @ (data: bytes, bytes_per_line: i32) -> string
    + returns a multi-line hexdump with offset, hex columns, and ASCII gutter
    + defaults to 16 bytes per line when bytes_per_line is zero
    # formatting
  hexdump.format_colored
    @ (data: bytes, bytes_per_line: i32) -> string
    + wraps printable bytes, control bytes, and nulls in distinct ANSI color codes
    # formatting
  hexdump.classify_byte
    @ (b: u8) -> string
    + returns "null", "control", "printable", or "high" for a byte
    # classification
