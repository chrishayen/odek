# Requirement: "a structured log handler that writes color-tinted lines"

Given a log record, produce a single formatted line with ANSI color per level. Color output is controlled by a flag so tests can compare plain strings.

std: (all units exist)

tint
  tint.format_line
    @ (level: string, time_rfc3339: string, message: string, attrs: map[string, string], use_color: bool) -> string
    + returns a line of the form "<time> <LEVEL> <message> key=value key=value"
    + wraps the level token in an ANSI color code matching the level when use_color is true
    + emits attributes in key-sorted order for stable output
    + escapes attribute values containing spaces by wrapping them in double quotes
    # formatting
  tint.color_for_level
    @ (level: string) -> string
    + returns a distinct ANSI foreground code for "debug", "info", "warn", and "error"
    + returns the reset code for any other level
    # formatting
