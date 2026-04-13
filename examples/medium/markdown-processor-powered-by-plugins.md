# Requirement: "a plugin-driven markdown processor that parses, transforms, and serializes"

A pipeline that parses markdown into a syntax tree, runs visitor plugins over it, and serializes back to markdown or HTML.

std: (all units exist)

mdproc
  mdproc.new
    @ () -> mdproc_state
    + creates an empty processor with no plugins attached
    # construction
  mdproc.use
    @ (state: mdproc_state, plugin: fn(md_node) -> md_node) -> mdproc_state
    + registers a transform that will run over the parsed tree
    # extension
  mdproc.parse
    @ (state: mdproc_state, source: string) -> result[md_node, string]
    + parses markdown source into a tree of md_node values
    - returns error on unrecoverable syntax problems
    # parsing
  mdproc.run
    @ (state: mdproc_state, tree: md_node) -> md_node
    + applies every registered transform to the tree in registration order
    # transformation
  mdproc.stringify
    @ (tree: md_node) -> string
    + serializes the tree back to markdown
    # serialization
  mdproc.to_html
    @ (tree: md_node) -> string
    + serializes the tree to HTML
    # serialization
  mdproc.process
    @ (state: mdproc_state, source: string) -> result[md_node, string]
    + convenience: parse then run all transforms
    - returns error when parsing fails
    # pipeline
