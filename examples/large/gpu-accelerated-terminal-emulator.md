# Requirement: "a GPU-accelerated terminal emulator"

A large system: PTY host, ANSI parser, grid model, font shaping, and a GPU renderer abstraction.

std
  std.process
    std.process.spawn_pty
      @ (program: string, args: list[string]) -> result[pty_handle, string]
      + spawns a child process attached to a pseudo-terminal
      - returns error when the program cannot be executed
      # process
    std.process.pty_read
      @ (handle: pty_handle, buf_size: i32) -> result[bytes, string]
      + reads up to buf_size bytes from the pty output
      # process
    std.process.pty_write
      @ (handle: pty_handle, data: bytes) -> result[i32, string]
      + writes data to the pty input
      # process
    std.process.pty_resize
      @ (handle: pty_handle, cols: i32, rows: i32) -> result[void, string]
      + resizes the pty window
      # process
  std.gpu
    std.gpu.create_device
      @ (surface_id: i64) -> result[gpu_device, string]
      + acquires a GPU device bound to a native surface
      # gpu
    std.gpu.create_texture
      @ (device: gpu_device, width: i32, height: i32) -> result[gpu_texture, string]
      + allocates a 2D RGBA texture on the device
      # gpu
    std.gpu.draw_quads
      @ (device: gpu_device, texture: gpu_texture, quads: bytes) -> result[void, string]
      + submits a batch of textured quads for rendering
      # gpu
    std.gpu.present
      @ (device: gpu_device) -> result[void, string]
      + presents the current frame to the surface
      # gpu
  std.text
    std.text.shape_glyphs
      @ (font: bytes, codepoints: list[i32], size_px: f32) -> list[glyph_rect]
      + shapes a run of codepoints into positioned glyph rectangles
      # text_shaping

terminal
  terminal.new_grid
    @ (cols: i32, rows: i32) -> grid_state
    + creates an empty cell grid with the given dimensions
    # grid
  terminal.parse_ansi
    @ (state: grid_state, input: bytes) -> grid_state
    + applies incoming bytes to the grid, interpreting ANSI escape sequences
    + handles cursor movement, colors, erase, and scroll region
    - unknown escape sequences are ignored
    # ansi_parser
  terminal.resize
    @ (state: grid_state, cols: i32, rows: i32) -> grid_state
    + reflows the grid to new dimensions preserving content where possible
    # grid
  terminal.scroll
    @ (state: grid_state, lines: i32) -> grid_state
    + scrolls the visible region by lines (positive = down)
    # grid
  terminal.build_glyph_atlas
    @ (font: bytes, size_px: f32) -> glyph_atlas
    + rasterizes ASCII glyphs into a texture atlas
    # rendering
    -> std.text.shape_glyphs
  terminal.render_frame
    @ (device: gpu_device, atlas: glyph_atlas, state: grid_state) -> result[void, string]
    + rasterizes the grid into textured quads and presents a frame
    # rendering
    -> std.gpu.draw_quads
    -> std.gpu.present
  terminal.attach_pty
    @ (cols: i32, rows: i32, program: string, args: list[string]) -> result[terminal_session, string]
    + spawns a child process and returns a session binding a grid to the pty
    # session
    -> std.process.spawn_pty
  terminal.pump_io
    @ (session: terminal_session) -> result[terminal_session, string]
    + reads pending pty output and feeds it through the ANSI parser
    # session
    -> std.process.pty_read
  terminal.send_keys
    @ (session: terminal_session, data: bytes) -> result[void, string]
    + writes keystrokes to the pty input
    # session
    -> std.process.pty_write
