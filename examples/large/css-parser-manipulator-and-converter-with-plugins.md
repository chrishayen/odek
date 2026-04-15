# Requirement: "a library to parse, manipulate, and convert modern CSS for older targets, with a plugin system"

The pipeline is tokenize -> parse AST -> run plugins over the tree -> serialize. Plugins are registered and invoked at well-defined visit points.

std
  std.strings
    std.strings.starts_with
      fn (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits s on every occurrence of sep
      # strings
  std.collections
    std.collections.list_append
      fn (xs: list[string], x: string) -> list[string]
      + returns xs with x appended
      # collections

csskit
  csskit.tokenize
    fn (source: string) -> result[list[css_token], string]
    + emits tokens for idents, strings, numbers, punctuation, at-keywords, and comments
    - returns error on unterminated strings or comments
    # tokenization
  csskit.parse
    fn (tokens: list[css_token]) -> result[css_stylesheet, string]
    + returns an AST with rules, at-rules, declarations, and nested selectors
    - returns error on mismatched braces
    # parsing
  csskit.serialize
    fn (sheet: css_stylesheet) -> string
    + round-trips a stylesheet to valid CSS text
    # serialization
  csskit.walk_rules
    fn (sheet: css_stylesheet, visit: visitor_fn) -> css_stylesheet
    + calls visit on every rule, allowing in-place replacement
    # traversal
  csskit.find_declarations
    fn (rule: css_rule, property: string) -> list[css_declaration]
    + returns every declaration whose property equals the given name
    # query
  csskit.set_declaration
    fn (rule: css_rule, property: string, value: string) -> css_rule
    + updates the declaration, or appends it when missing
    # mutation
  csskit.plugin_register
    fn (reg: plugin_registry, name: string, plugin: plugin_fn) -> plugin_registry
    + adds a plugin to the registry under the given name
    # plugins
  csskit.plugin_run_all
    fn (reg: plugin_registry, sheet: css_stylesheet) -> result[css_stylesheet, string]
    + runs every registered plugin in insertion order
    - returns the first plugin error, annotated with the plugin name
    # plugins
  csskit.convert_for_target
    fn (sheet: css_stylesheet, target: string) -> css_stylesheet
    + rewrites modern features to equivalents for the named target
    ? target is a short string such as "legacy" or "modern"; the mapping table is internal
    # compatibility
    -> std.strings.starts_with
  csskit.expand_nested_selectors
    fn (sheet: css_stylesheet) -> css_stylesheet
    + flattens nested rules into top-level rules by joining selectors
    # compatibility
    -> std.strings.split
    -> std.collections.list_append
  csskit.prefix_vendor_properties
    fn (sheet: css_stylesheet, prefixes: list[string]) -> css_stylesheet
    + adds vendor-prefixed copies for properties that need them
    # compatibility
