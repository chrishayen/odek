# Requirement: "a data analysis library with statistics, visualization, and columnar file support"

A dataframe-centered analysis library: columnar storage, descriptive statistics, grouping, plotting into an image buffer, and columnar file IO.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file fully into memory
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to a file, replacing any existing contents
      # filesystem
  std.math
    std.math.sqrt
      fn (x: f64) -> f64
      + returns the square root of x
      - returns NaN for negative input
      # math
    std.math.sort_f64
      fn (values: list[f64]) -> list[f64]
      + returns the values sorted ascending
      # math
  std.image
    std.image.new_canvas
      fn (width: u32, height: u32) -> canvas_state
      + creates a blank RGBA canvas
      # image
    std.image.draw_line
      fn (canvas: canvas_state, x0: u32, y0: u32, x1: u32, y1: u32, color: u32) -> canvas_state
      + draws a straight line on the canvas
      # image
    std.image.encode_png
      fn (canvas: canvas_state) -> bytes
      + encodes the canvas as PNG
      # image

insyra
  insyra.column_from_f64
    fn (name: string, values: list[f64]) -> column
    + creates a numeric column
    # construction
  insyra.column_from_string
    fn (name: string, values: list[string]) -> column
    + creates a string column
    # construction
  insyra.new_dataframe
    fn (columns: list[column]) -> result[dataframe, string]
    + creates a dataframe from equal-length columns
    - returns error when columns have mismatched lengths
    # construction
  insyra.select
    fn (df: dataframe, names: list[string]) -> result[dataframe, string]
    + returns a dataframe with only the requested columns
    - returns error when a name does not exist
    # projection
  insyra.filter
    fn (df: dataframe, column: string, predicate: fn(cell) -> bool) -> result[dataframe, string]
    + returns only rows where the predicate holds on the named column
    - returns error when the column does not exist
    # filtering
  insyra.describe
    fn (df: dataframe, column: string) -> result[stats_summary, string]
    + returns count, mean, stddev, min, median, and max for a numeric column
    - returns error when the column is non-numeric
    # statistics
    -> std.math.sqrt
    -> std.math.sort_f64
  insyra.group_by_mean
    fn (df: dataframe, group_column: string, value_column: string) -> result[dataframe, string]
    + returns a two-column dataframe with group keys and per-group means
    - returns error when either column does not exist
    # aggregation
  insyra.plot_line
    fn (df: dataframe, x_column: string, y_column: string, width: u32, height: u32) -> result[bytes, string]
    + renders a line chart as PNG bytes
    - returns error when columns are missing or non-numeric
    # visualization
    -> std.image.new_canvas
    -> std.image.draw_line
    -> std.image.encode_png
  insyra.read_columnar_file
    fn (path: string) -> result[dataframe, string]
    + loads a dataframe from a columnar file
    - returns error when the file is not a valid columnar format
    # io
    -> std.fs.read_all
  insyra.write_columnar_file
    fn (df: dataframe, path: string) -> result[void, string]
    + writes a dataframe to a columnar file
    # io
    -> std.fs.write_all
