# Requirement: "a realtime 2D, 3D, and image plotting library"

Maintains plot state that can be updated at high frequency and queried for the current frame. Rendering is the host's job.

std: (all units exist)

plot
  plot.new_line_plot
    fn (capacity: i32) -> plot_state
    + creates a line plot with a ring buffer of the given capacity
    ? capacity bounds memory so realtime updates do not grow unbounded
    # construction
  plot.new_image_plot
    fn (width: i32, height: i32) -> plot_state
    + creates an image plot backed by a width*height pixel buffer
    # construction
  plot.new_volume_plot
    fn (nx: i32, ny: i32, nz: i32) -> plot_state
    + creates a 3D scalar-field plot
    # construction
  plot.push_point
    fn (state: plot_state, x: f64, y: f64) -> plot_state
    + appends a point to a line plot; oldest point drops when at capacity
    - returns state unchanged when called on a non-line plot
    # update
  plot.set_pixel
    fn (state: plot_state, x: i32, y: i32, value: f32) -> plot_state
    + writes one pixel in an image plot
    - returns state unchanged when coordinates are out of range
    # update
  plot.set_voxel
    fn (state: plot_state, x: i32, y: i32, z: i32, value: f32) -> plot_state
    + writes one voxel in a volume plot
    # update
  plot.snapshot
    fn (state: plot_state) -> plot_frame
    + returns an immutable frame the renderer can consume
    # query
  plot.auto_range
    fn (state: plot_state) -> tuple[f64, f64]
    + returns (min, max) of the current data for axis scaling
    # analysis
