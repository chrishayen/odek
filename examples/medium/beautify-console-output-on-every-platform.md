# Requirement: "a cross-platform library of composable styled console components"

Styling and layout primitives live in the project; only terminal width detection is a std utility.

std
  std.term
    std.term.width
      @ () -> i32
      + returns the current terminal width in columns
      ? returns 80 when width cannot be detected
      # terminal

console_style
  console_style.colorize
    @ (text: string, fg: rgb_color, bg: rgb_color) -> string
    + wraps text in ANSI truecolor foreground and background escapes followed by a reset
    # styling
  console_style.bold
    @ (text: string) -> string
    + wraps text in the ANSI bold on/off sequences
    # styling
  console_style.render_box
    @ (content: string, title: string) -> string
    + returns content framed in a rounded Unicode box with an optional title
    ? content lines wider than the terminal are truncated
    # layout
    -> std.term.width
  console_style.render_table
    @ (headers: list[string], rows: list[list[string]]) -> string
    + returns an aligned text table with a divider under the header row
    - returns empty string when rows is empty and headers is empty
    # layout
  console_style.render_bullet_list
    @ (items: list[string]) -> string
    + returns items joined by newlines, each prefixed with a bullet glyph
    # layout
  console_style.render_progress_bar
    @ (label: string, fraction: f64) -> string
    + returns "label [####------] 40%" sized to the terminal width
    ? fraction outside [0, 1] is clamped
    # layout
    -> std.term.width
