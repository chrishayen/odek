# Requirement: "a library for multiple simultaneous individually controllable terminal spinners"

Tracks a group of named spinners and renders a multi-line frame on demand. Output goes through a thin std writer so callers control where it lands.

std
  std.io
    std.io.write_string
      fn (text: string) -> result[void, string]
      + writes text to the standard output stream
      - returns error when the stream is closed
      # io

spinners
  spinners.new
    fn () -> spinners_state
    + creates an empty spinner group
    # construction
  spinners.add
    fn (state: spinners_state, name: string, label: string) -> spinners_state
    + adds a spinner in the "running" status with frame index 0
    # registration
  spinners.set_status
    fn (state: spinners_state, name: string, status: string) -> result[spinners_state, string]
    + updates the named spinner's status (running, success, fail)
    - returns error when the name is not registered
    # control
  spinners.render_frame
    fn (state: spinners_state) -> tuple[string, spinners_state]
    + returns the multi-line frame for all spinners and advances each running spinner's frame index
    ? finished spinners render with a static marker instead of a spinner glyph
    # rendering
    -> std.io.write_string
