# Requirement: "a multi-bar terminal progress display"

Tracks multiple named progress bars and renders them as a block of lines that can be refreshed in place.

std
  std.io
    std.io.write
      fn (data: string) -> void
      + writes text to standard output without a trailing newline
      # io
  std.term
    std.term.cursor_up
      fn (n: i32) -> string
      + returns the ANSI escape to move the cursor up n lines
      # terminal
    std.term.clear_line
      fn () -> string
      + returns the ANSI escape to clear the current line
      # terminal

progress
  progress.container_new
    fn (width: i32) -> progress_container
    + creates an empty container with the given render width
    # construction
  progress.add_bar
    fn (container: progress_container, name: string, total: i64) -> progress_container
    + adds a named bar with a total count
    # bars
  progress.update_bar
    fn (container: progress_container, name: string, current: i64) -> progress_container
    + sets the current count for a bar
    - no-op when the name is unknown
    # bars
  progress.complete_bar
    fn (container: progress_container, name: string) -> progress_container
    + marks a bar as finished at 100%
    # bars
  progress.render_frame
    fn (container: progress_container) -> string
    + returns the full multi-line rendered block
    ? each bar shows name, a filled portion, and percent
    # rendering
  progress.refresh
    fn (container: progress_container, previous_line_count: i32) -> i32
    + clears prior lines and writes the new frame, returning the new line count
    # rendering
    -> std.term.cursor_up
    -> std.term.clear_line
    -> std.io.write
