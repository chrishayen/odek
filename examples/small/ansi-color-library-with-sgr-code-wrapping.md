# Requirement: "an ANSI terminal color library that wraps strings with SGR codes"

Returns color-wrapped strings. Callers pass the wrapped strings to whatever printing facility they use.

std: (all units exist)

ansi_color
  ansi_color.colorize
    fn (text: string, fg: i32) -> string
    + returns text wrapped in an SGR sequence setting the foreground color and a reset at the end
    ? fg follows the standard 0-7 palette plus 8-15 bright range
    # color
  ansi_color.colorize_bg
    fn (text: string, bg: i32) -> string
    + returns text wrapped in an SGR sequence setting the background color and a reset at the end
    # color
  ansi_color.bold
    fn (text: string) -> string
    + returns text wrapped in the bold SGR sequence and a reset
    # style
  ansi_color.underline
    fn (text: string) -> string
    + returns text wrapped in the underline SGR sequence and a reset
    # style
  ansi_color.strip
    fn (text: string) -> string
    + returns text with all SGR sequences removed
    + returns the input unchanged when it contains no escapes
    # stripping
