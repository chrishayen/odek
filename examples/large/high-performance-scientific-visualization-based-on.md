# Requirement: "a scientific visualization library backed by a GPU shader pipeline"

Builds scene graphs of data primitives and submits draw calls to a GPU backend. The shader pipeline is treated as an opaque backend interface.

std
  std.math
    std.math.min_max_f64
      @ (xs: list[f64]) -> tuple[f64, f64]
      + returns the minimum and maximum of a non-empty list
      - returns (0, 0) when xs is empty
      # math
    std.math.lerp
      @ (a: f64, b: f64, t: f64) -> f64
      + returns a + (b - a) * t
      # math
  std.color
    std.color.rgba_from_hex
      @ (hex: string) -> result[rgba, string]
      + parses "#rrggbb" or "#rrggbbaa"
      - returns error on invalid length or non-hex characters
      # color

viz
  viz.new_scene
    @ () -> scene
    + creates an empty scene with identity transform
    # construction
  viz.add_line_plot
    @ (s: scene, xs: list[f64], ys: list[f64], color: rgba) -> scene
    + adds a line-strip primitive over the given points
    - returns the scene unchanged when xs and ys differ in length
    # primitives
  viz.add_scatter
    @ (s: scene, xs: list[f64], ys: list[f64], color: rgba, size: f32) -> scene
    + adds a point-cloud primitive with per-vertex size
    # primitives
  viz.add_image
    @ (s: scene, pixels: bytes, width: u32, height: u32) -> scene
    + adds a textured quad from raw RGBA pixels
    # primitives
  viz.fit_view
    @ (s: scene) -> scene
    + sets the view transform so all primitives fit within the viewport with a small margin
    # camera
    -> std.math.min_max_f64
  viz.set_colormap
    @ (s: scene, stops: list[rgba]) -> scene
    + stores an interpolation-ready colormap used by value-to-color lookups
    # color
  viz.value_to_color
    @ (s: scene, value: f64, vmin: f64, vmax: f64) -> rgba
    + maps value linearly onto the active colormap
    # color
    -> std.math.lerp
  viz.render
    @ (s: scene, backend: render_backend, width: u32, height: u32) -> result[bytes, string]
    + walks the scene and submits draw calls, returning the framebuffer pixels
    - returns error when backend initialization fails
    # rendering
