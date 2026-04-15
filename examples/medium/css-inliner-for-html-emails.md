# Requirement: "a css inliner for html emails"

Parses an HTML document and its embedded stylesheets, then folds matching CSS declarations onto each element's style attribute so the result renders in clients that strip <style> blocks.

std
  std.html
    std.html.parse
      fn (source: string) -> result[dom_node, string]
      + parses an HTML document into a DOM tree
      - returns error on malformed input
      # html_parse
    std.html.serialize
      fn (root: dom_node) -> string
      + serializes a DOM tree back to HTML
      # html_serialize
    std.html.get_attribute
      fn (node: dom_node, name: string) -> optional[string]
      + returns the attribute value if present
      # dom
    std.html.set_attribute
      fn (node: dom_node, name: string, value: string) -> dom_node
      + returns a node with the attribute set
      # dom
  std.css
    std.css.parse_stylesheet
      fn (source: string) -> result[list[css_rule], string]
      + parses CSS into a list of (selector, declarations) rules
      - returns error on unterminated block
      # css_parse
    std.css.match_selector
      fn (selector: string, node: dom_node) -> bool
      + returns true when the selector matches the element
      # css_match
    std.css.declarations_to_string
      fn (decls: list[css_declaration]) -> string
      + serializes declarations as a single "prop: value; ..." string
      # css_serialize

inliner
  inliner.extract_stylesheets
    fn (root: dom_node) -> tuple[dom_node, list[css_rule]]
    + removes every <style> element from the tree and returns the parsed rules
    + returns the root unchanged when no <style> elements exist
    # extraction
    -> std.html.parse
    -> std.css.parse_stylesheet
  inliner.apply_rules
    fn (root: dom_node, rules: list[css_rule]) -> dom_node
    + walks the tree and appends matching declarations to each element's existing style attribute
    + preserves existing inline declarations when a rule and the element both set the same property
    # application
    -> std.css.match_selector
    -> std.html.get_attribute
    -> std.html.set_attribute
    -> std.css.declarations_to_string
  inliner.inline
    fn (html: string) -> result[string, string]
    + parses, extracts stylesheets, applies them inline, and serializes the result
    - returns error on malformed HTML
    # pipeline
    -> std.html.serialize
