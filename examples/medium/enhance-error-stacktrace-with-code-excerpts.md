# Requirement: "enhance an error stack trace with source code excerpts"

Parses a raw stack trace, locates each frame's source file, and attaches a window of surrounding lines.

std
  std.fs
    std.fs.read_all_text
      @ (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when the file does not exist
      # filesystem

stack_enhancer
  stack_enhancer.parse_frames
    @ (raw: string) -> list[frame]
    + returns one frame per "at function (file:line:col)" line found in the trace
    + returns empty list when no frame lines match
    # parsing
  stack_enhancer.load_source_window
    @ (path: string, line: i32, radius: i32) -> result[source_window, string]
    + returns the lines from (line - radius) to (line + radius) with line numbers
    - returns error when the file cannot be read
    - returns error when line exceeds the file length
    # source_lookup
    -> std.fs.read_all_text
  stack_enhancer.attach_excerpts
    @ (frames: list[frame], radius: i32) -> list[enhanced_frame]
    + returns each frame with its source window when available
    + frames whose file cannot be read still appear but without a window
    # enrichment
  stack_enhancer.enhance
    @ (raw: string, radius: i32) -> list[enhanced_frame]
    + parses the trace and attaches source windows for every frame
    # enrichment
  stack_enhancer.render
    @ (frames: list[enhanced_frame]) -> string
    + returns a human-readable multi-line rendering of the enhanced trace
    # rendering
