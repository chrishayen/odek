# Requirement: "an HTML and CSS rendering engine with PDF export"

Parses HTML and CSS into a styled box tree, runs a simple block-flow layout, and serializes the result to PDF pages.

std
  std.strings
    std.strings.to_lower
      fn (s: string) -> string
      + returns the lowercased form of s
      # text
    std.strings.trim
      fn (s: string) -> string
      + strips leading and trailing whitespace
      # text
  std.collections
    std.collections.hashmap_new
      fn () -> map[string, string]
      + returns an empty string-to-string map
      # collections

render_engine
  render_engine.parse_html
    fn (source: string) -> result[dom_node, string]
    + parses HTML into a DOM tree
    - returns error on malformed start tag
    # html
    -> std.strings.to_lower
  render_engine.parse_css
    fn (source: string) -> result[list[css_rule], string]
    + parses a stylesheet into rules of selectors and declarations
    - returns error on unterminated block
    # css
    -> std.strings.trim
  render_engine.match_selector
    fn (rule: css_rule, node: dom_node) -> bool
    + returns true when the selector matches the node
    - supports tag, id, and class selectors
    # styling
  render_engine.compute_styles
    fn (dom: dom_node, rules: list[css_rule]) -> styled_node
    + applies all matching rules to each node, producing a styled tree with cascading
    # styling
    -> std.collections.hashmap_new
  render_engine.build_box_tree
    fn (styled: styled_node) -> box_node
    + converts styled nodes into block, inline, or anonymous boxes per display
    # layout
  render_engine.layout
    fn (root: box_node, viewport_width_px: f64) -> layout_box
    + assigns positions and sizes using a block-flow algorithm
    + wraps inline boxes to fit viewport width
    # layout
  render_engine.paint_to_display_list
    fn (layout: layout_box) -> list[draw_command]
    + emits a flat list of draw commands (rectangles, glyph runs)
    # painting
  render_engine.emit_pdf
    fn (commands: list[draw_command], page_width_pt: f64, page_height_pt: f64) -> bytes
    + writes a minimal PDF containing one or more pages with the commands
    # pdf
  render_engine.render_to_pdf
    fn (html: string, css: string, page_width_pt: f64, page_height_pt: f64) -> result[bytes, string]
    + top-level pipeline: parse, style, lay out, paint, and emit PDF
    - returns error at the first failing stage
    # pipeline
