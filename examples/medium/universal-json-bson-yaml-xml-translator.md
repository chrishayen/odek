# Requirement: "a format translator that parses JSON, YAML, or XML into a common tree and renders it through a user-supplied template"

One in-memory tree, format-specific parsers and renderers, and a small template engine that walks the tree.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[generic_tree, string]
      - returns error on invalid JSON
      # parsing
    std.json.render
      @ (tree: generic_tree) -> string
      # serialization
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[generic_tree, string]
      - returns error on invalid YAML
      # parsing
    std.yaml.render
      @ (tree: generic_tree) -> string
      # serialization
  std.xml
    std.xml.parse
      @ (raw: string) -> result[generic_tree, string]
      - returns error on invalid XML
      # parsing
    std.xml.render
      @ (tree: generic_tree) -> string
      # serialization

translator
  translator.parse
    @ (format: string, raw: string) -> result[generic_tree, string]
    + dispatches on format name ("json", "yaml", "xml")
    - returns error on unknown format
    # parsing
    -> std.json.parse
    -> std.yaml.parse
    -> std.xml.parse
  translator.render
    @ (format: string, tree: generic_tree) -> result[string, string]
    + dispatches on format name
    - returns error on unknown format
    # rendering
    -> std.json.render
    -> std.yaml.render
    -> std.xml.render
  translator.compile_template
    @ (template: string) -> result[compiled_template, string]
    + parses a template with "{{path.to.value}}" placeholders and "{{#each list}}" loops
    - returns error on malformed template
    # templates
  translator.apply_template
    @ (compiled: compiled_template, tree: generic_tree) -> string
    + walks the tree and renders the template
    ? missing paths render as the empty string
    # templates
  translator.translate
    @ (from_format: string, raw: string, template: string) -> result[string, string]
    + one-shot parse + template rendering
    - returns error on parse or template failure
    # pipeline
    -> translator.parse
    -> translator.compile_template
    -> translator.apply_template
