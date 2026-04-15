# Requirement: "a web browser engine that fetches, parses, lays out, and paints a page"

A high-level pipeline: network fetch, HTML parse, CSS parse, style resolution, layout, and paint into a display list.

std
  std.io
    std.io.http_get
      fn (url: string) -> result[http_response, string]
      + performs an HTTP GET and returns status, headers, and body
      - returns error on network failure
      # http
  std.encoding
    std.encoding.utf8_decode
      fn (raw: bytes) -> result[string, string]
      + decodes UTF-8 bytes to a string
      - returns error on invalid sequences
      # encoding

browser_engine
  browser_engine.tokenize_html
    fn (source: string) -> result[list[html_token], string]
    + splits HTML source into tags, text, comments, and doctype tokens
    - returns error on unterminated tags
    # html_tokenization
  browser_engine.parse_html
    fn (tokens: list[html_token]) -> result[dom_node, string]
    + builds a DOM tree with implicit element closing
    - returns error when the root document is empty
    # html_parsing
  browser_engine.tokenize_css
    fn (source: string) -> list[css_token]
    + splits CSS source into identifiers, punctuation, numbers, and strings
    # css_tokenization
  browser_engine.parse_css
    fn (tokens: list[css_token]) -> result[stylesheet, string]
    + builds a stylesheet of rules (selector + declarations)
    - returns error on mismatched braces
    # css_parsing
  browser_engine.match_selector
    fn (node: dom_node, selector: css_selector) -> bool
    + returns true when the node matches the selector (tag, id, class)
    # selector_matching
  browser_engine.compute_styles
    fn (root: dom_node, sheet: stylesheet) -> styled_tree
    + resolves declarations per node using selector specificity and inheritance
    # style_resolution
  browser_engine.build_layout_tree
    fn (styled: styled_tree) -> layout_box
    + creates a layout tree honoring display:block/inline/none
    # layout_tree
  browser_engine.compute_layout
    fn (root: layout_box, viewport_width: f64) -> layout_box
    + performs block-direction layout producing positions and sizes for each box
    # layout
  browser_engine.build_display_list
    fn (root: layout_box) -> list[paint_command]
    + walks the laid-out tree and emits rectangle and text paint commands
    # painting
  browser_engine.render_url
    fn (url: string) -> result[list[paint_command], string]
    + fetches, parses, styles, lays out, and returns a display list for a URL
    - returns error at any stage when the document cannot be processed
    # orchestration
    -> std.io.http_get
    -> std.encoding.utf8_decode
