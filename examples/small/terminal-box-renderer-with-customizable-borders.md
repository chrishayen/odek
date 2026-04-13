# Requirement: "a library that renders customizable boxes around text for terminal output"

Given content and a box style, produce the framed multi-line string. Style covers border glyphs, padding, and optional title.

std: (all units exist)

box
  box.default_style
    @ () -> box_style
    + returns a single-line ASCII border style with one-cell padding
    # defaults
  box.with_title
    @ (style: box_style, title: string) -> box_style
    + returns a new style whose top border embeds the given title
    # customization
  box.render
    @ (content: string, style: box_style) -> string
    + returns the content surrounded by a border matching the style
    + splits content on newlines and pads each line to the widest line's width
    + accounts for wide characters when computing visible width
    # rendering
  box.render_lines
    @ (lines: list[string], style: box_style) -> list[string]
    + returns the rendered box as a list of lines without a trailing newline
    # rendering
