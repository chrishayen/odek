# Requirement: "convert a terminal recording into an animated SVG"

Consumes a terminal recording (header + timed write events), replays it onto a virtual screen, and emits an animated SVG.

std
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

terminal_svg
  terminal_svg.parse_recording
    fn (raw: string) -> result[recording, string]
    + parses a terminal recording (header line with width/height, followed by one timed event per line)
    - returns error on invalid header
    - returns error on malformed event lines
    # parsing
    -> std.json.parse_object
  terminal_svg.new_screen
    fn (cols: i32, rows: i32) -> screen_state
    + creates a blank virtual screen with the given dimensions
    # virtual_terminal
  terminal_svg.apply
    fn (screen: screen_state, data: string) -> screen_state
    + feeds raw terminal output into the screen, handling printable chars and ANSI SGR/cursor escapes
    # virtual_terminal
  terminal_svg.replay
    fn (rec: recording) -> list[screen_frame]
    + replays a recording into a list of (timestamp, screen) frames
    # replay
  terminal_svg.render_svg
    fn (frames: list[screen_frame], cols: i32, rows: i32) -> string
    + returns an SVG document whose groups animate between frames using begin/dur attributes
    # rendering
  terminal_svg.render_static
    fn (screen: screen_state) -> string
    + returns a static SVG snapshot of a single screen
    # rendering
