# Requirement: "a library for drawing layout diagrams declared as code"

A diagram is a set of labeled boxes with edges between them. The renderer lays boxes out on a grid and emits an ascii-art drawing.

std: (all units exist)

diagram
  diagram.new
    fn () -> diagram_state
    + creates an empty diagram
    # construction
  diagram.add_box
    fn (state: diagram_state, id: string, label: string) -> diagram_state
    + adds a box with the given id and display label
    + replaces an existing box with the same id
    # definition
  diagram.add_edge
    fn (state: diagram_state, from_id: string, to_id: string) -> result[diagram_state, string]
    + records a directed edge between two existing boxes
    - returns error when either endpoint id is unknown
    # definition
  diagram.layout
    fn (state: diagram_state) -> layout_state
    + assigns each box a (column, row) position using a layered algorithm
    + edges run top-to-bottom; nodes without predecessors land on row 0
    # layout
  diagram.render_ascii
    fn (layout: layout_state) -> string
    + returns an ascii-art drawing of the laid-out diagram with boxed labels and connecting lines
    # rendering
  diagram.render_svg
    fn (layout: layout_state) -> string
    + returns an SVG document for the laid-out diagram
    # rendering
