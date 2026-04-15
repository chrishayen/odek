# Requirement: "a minimal serial monitor with an in-terminal plotter"

The library exposes a small monitor state: it ingests text lines from a serial source, parses numeric samples, and renders an ASCII plot. Actual serial IO and terminal writes are the caller's concern.

std
  std.text
    std.text.split_lines
      fn (input: string) -> list[string]
      + splits on \n, dropping a trailing empty line
      # text
    std.text.parse_f64
      fn (s: string) -> optional[f64]
      + parses decimal numbers, optionally signed
      - returns empty on non-numeric input
      # parsing

serial_plot
  serial_plot.new
    fn (capacity: i32) -> plotter_state
    + returns a plotter that retains the most recent capacity samples
    ? older samples are dropped as a ring buffer
    # construction
  serial_plot.ingest_line
    fn (state: plotter_state, line: string) -> plotter_state
    + appends a sample when the line is numeric
    - leaves state unchanged when the line is not numeric
    # ingestion
    -> std.text.parse_f64
  serial_plot.ingest_chunk
    fn (state: plotter_state, chunk: string) -> plotter_state
    + splits the chunk into lines and ingests each
    # ingestion
    -> std.text.split_lines
  serial_plot.samples
    fn (state: plotter_state) -> list[f64]
    + returns samples in insertion order, oldest first
    # accessor
  serial_plot.render
    fn (state: plotter_state, width: i32, height: i32) -> string
    + returns a multi-line ASCII plot scaled to the buffer's min/max
    + returns a blank canvas of the requested size when no samples exist
    ? each row is exactly width characters wide, rows separated by \n
    # rendering
