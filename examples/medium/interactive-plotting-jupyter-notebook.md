# Requirement: "an interactive 2D plotting library for notebook environments"

Builds a plot model that can be serialized and handed to a rendering host. Interactivity is represented as a set of handlers the host can call.

std
  std.json
    std.json.encode
      @ (value: json_value) -> string
      + encodes a json value tree to a json string
      # serialization

plot
  plot.figure
    @ (width_px: i32, height_px: i32) -> figure_state
    + creates an empty figure with the given pixel dimensions
    # construction
  plot.add_line
    @ (fig: figure_state, xs: list[f64], ys: list[f64], label: string) -> figure_state
    + adds a line series; xs and ys must have equal length
    - returns the figure unchanged when lengths differ, flagging an error
    # series
  plot.add_scatter
    @ (fig: figure_state, xs: list[f64], ys: list[f64], label: string) -> figure_state
    + adds a scatter series
    # series
  plot.set_axes
    @ (fig: figure_state, x_label: string, y_label: string) -> figure_state
    + sets axis labels
    # axes
  plot.on_point_click
    @ (fig: figure_state, handler_id: string) -> figure_state
    + registers a named handler to be invoked on point-click events from the host
    # interaction
  plot.dispatch_event
    @ (fig: figure_state, event: plot_event) -> figure_state
    + routes an event from the host to the registered handlers
    - returns the figure unchanged when no handler matches
    # interaction
  plot.to_spec
    @ (fig: figure_state) -> string
    + serializes the figure as a JSON spec the host renders
    # serialization
    -> std.json.encode
