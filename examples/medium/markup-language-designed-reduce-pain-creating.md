# Requirement: "a markup language compiler that produces responsive HTML email"

Parses a structured markup source into a tree and emits table-based responsive HTML that renders consistently in email clients.

std: (all units exist)

email_markup
  email_markup.parse
    @ (source: string) -> result[markup_node, string]
    + parses the markup source into a tree of components with their attributes
    - returns error on unknown tags or malformed attributes
    # parsing
  email_markup.validate
    @ (tree: markup_node) -> result[void, list[string]]
    + checks that parent-child relationships match the component spec
    - returns the list of validation errors when the tree is invalid
    # validation
  email_markup.resolve_styles
    @ (tree: markup_node) -> markup_node
    + resolves component attributes into concrete CSS properties on each node
    # styling
  email_markup.render_component
    @ (node: markup_node) -> string
    + renders a single node as an HTML table fragment with inlined styles
    # rendering
  email_markup.render
    @ (tree: markup_node) -> string
    + renders the full document to a responsive HTML email with media queries
    # rendering
  email_markup.compile
    @ (source: string) -> result[string, string]
    + convenience pipeline: parse, validate, resolve styles, render
    - returns error when any stage fails
    # pipeline
