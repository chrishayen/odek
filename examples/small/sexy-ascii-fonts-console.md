# Requirement: "render text as large ASCII-art banners"

Looks up each character in a font table and composes the per-row output side by side.

std: (all units exist)

asciiart
  asciiart.builtin_fonts
    @ () -> list[string]
    + returns the names of the built-in fonts
    # fonts
  asciiart.load_font
    @ (name: string) -> result[font_state, string]
    + loads a built-in font by name
    - returns error when the font is unknown
    # fonts
  asciiart.render
    @ (text: string, font: font_state) -> string
    + returns a multi-line string where each character in text is rendered in the font and concatenated horizontally
    + returns a newline between font rows
    - unsupported characters are replaced with blank space of the same width
    # rendering
  asciiart.render_colored
    @ (text: string, font: font_state, color: string) -> string
    + like render but wraps each line with an ANSI color escape
    - returns the uncolored rendering when color is empty
    # rendering
