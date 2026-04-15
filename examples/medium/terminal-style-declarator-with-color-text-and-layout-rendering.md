# Requirement: "declaratively define styles for color, text format, and layout, then render styled strings for a terminal"

A style is a value that accumulates foreground color, bold/italic flags, and padding. Rendering produces an ANSI-escaped string.

std: (all units exist)

termstyle
  termstyle.new_style
    fn () -> style_state
    + creates a style with no foreground, no decorations, and zero padding
    # construction
  termstyle.foreground
    fn (style: style_state, rgb: u32) -> style_state
    + returns a style with the given 24-bit foreground color
    # color
  termstyle.bold
    fn (style: style_state, on: bool) -> style_state
    + returns a style with bold enabled or disabled
    # decoration
  termstyle.italic
    fn (style: style_state, on: bool) -> style_state
    + returns a style with italic enabled or disabled
    # decoration
  termstyle.padding
    fn (style: style_state, left: i32, right: i32) -> style_state
    + returns a style with horizontal padding in spaces
    ? negative values are clamped to zero
    # layout
  termstyle.render
    fn (style: style_state, text: string) -> string
    + wraps text with ANSI SGR sequences for the style, then padding, then a reset
    + returns the text unchanged when no attributes are set
    # rendering
