# Requirement: "a JSON pretty-printer with configurable colors and indentation"

Formats a JSON value with indentation and optional ANSI color codes by token type.

std: (all units exist)

jpretty
  jpretty.new_style
    @ () -> style_state
    + returns a default style with 2-space indentation and no colors
    # construction
  jpretty.set_indent
    @ (style: style_state, spaces: i32) -> style_state
    + sets the indentation width in spaces
    # configuration
  jpretty.set_color
    @ (style: style_state, token: string, ansi: string) -> style_state
    + sets the ANSI escape sequence for a token kind ("key", "string", "number", "bool", "null")
    # configuration
  jpretty.format
    @ (raw: string, style: style_state) -> result[string, string]
    + returns a pretty-printed version of the JSON input using the style
    - returns error on invalid JSON
    ? colors are applied only when the style has non-empty sequences for the token
    # formatting
