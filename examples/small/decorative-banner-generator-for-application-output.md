# Requirement: "a library for adding decorative banners to application output"

Renders a multi-line banner with a border around a title. Printing is the caller's job.

std: (all units exist)

banner
  banner.render
    fn (title: string, width: i32) -> string
    + returns a multi-line string with a box border surrounding the title
    + centers the title when width exceeds title length
    - returns error-like marker string when width is smaller than title length + padding
    ? uses ASCII box characters; caller can print the result
    # rendering
  banner.render_styled
    fn (title: string, width: i32, border_char: string) -> string
    + returns a banner using the given single-character border glyph
    ? border_char must be exactly one visible character
    # rendering
