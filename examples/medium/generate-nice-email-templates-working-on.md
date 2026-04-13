# Requirement: "a library that renders responsive email templates compatible with a wide range of mail clients"

Parses a component-based email markup, lowers it to table-based HTML that renders consistently in legacy clients, and inlines styles.

std
  std.xml
    std.xml.parse
      @ (source: string) -> result[xml_node, string]
      + returns the root element with children and attributes
      - returns error on unbalanced tags
      # parsing
    std.xml.serialize
      @ (node: xml_node) -> string
      + emits the tree as text with attributes in original order
      # serialization
  std.css
    std.css.parse_rules
      @ (source: string) -> list[css_rule]
      + returns selector + declaration pairs
      # css

mrml
  mrml.parse_mjml
    @ (source: string) -> result[mjml_tree, string]
    + parses a component markup document into a tree of known component nodes
    - returns error on unknown components
    # parsing
    -> std.xml.parse
  mrml.lower_to_html
    @ (tree: mjml_tree) -> xml_node
    + lowers each component into legacy-safe table HTML
    + preserves attribute-derived sizing
    # lowering
  mrml.inline_styles
    @ (doc: xml_node, rules: list[css_rule]) -> xml_node
    + applies matching rules as inline style attributes
    + leaves elements without matches unchanged
    # style_inlining
  mrml.render
    @ (source: string) -> result[string, string]
    + parses, lowers, inlines styles, and serializes the result
    - returns error at the first failing stage
    # orchestration
    -> std.css.parse_rules
    -> std.xml.serialize
