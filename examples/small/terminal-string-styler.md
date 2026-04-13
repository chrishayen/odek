# Requirement: "a terminal string styling library"

Wraps strings in ANSI escape codes for color and text attributes. Styles compose by chaining.

std: (all units exist)

styler
  styler.new
    @ () -> style
    + creates a style with no attributes set
    # construction
  styler.with_color
    @ (s: style, color: i8) -> style
    + adds a foreground color (0-7 standard, 8-15 bright)
    # composition
  styler.with_bg
    @ (s: style, color: i8) -> style
    + adds a background color
    # composition
  styler.with_attr
    @ (s: style, attr: i8) -> style
    + adds a text attribute (0=bold, 1=dim, 2=italic, 3=underline)
    # composition
  styler.apply
    @ (s: style, text: string) -> string
    + wraps text in the ANSI escape sequences for the style and resets at the end
    + returns text unchanged when no attributes are set
    # rendering
