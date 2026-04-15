# Requirement: "a library to find overspecified function parameters that could be generalized with interface types"

Analyzes a parsed function signature and body to determine which methods are actually called on each parameter, then suggests a minimal interface that captures only those methods.

std
  std.set
    std.set.new
      fn () -> string_set
      + returns an empty string set
      # collection
    std.set.add
      fn (s: string_set, value: string) -> string_set
      + returns a set containing value
      # collection
    std.set.to_list
      fn (s: string_set) -> list[string]
      + returns the set members in stable order
      # collection

overspec
  overspec.collect_method_usage
    fn (body: function_body, param_name: string) -> list[string]
    + returns the distinct method names called on param_name within body
    + includes chained calls where the receiver is param_name
    ? field accesses are ignored; only method invocations count
    # usage_analysis
    -> std.set.new
    -> std.set.add
    -> std.set.to_list
  overspec.lookup_type_methods
    fn (type_name: string, type_index: type_index) -> list[string]
    + returns all methods defined on the named type
    - returns empty list when type is unknown
    # type_lookup
  overspec.is_overspecified
    fn (declared_methods: list[string], used_methods: list[string]) -> bool
    + returns true when used_methods is a strict subset of declared_methods
    - returns false when every declared method is used
    # classification
  overspec.synthesize_interface
    fn (used_methods: list[string]) -> interface_shape
    + returns an interface_shape listing exactly the used methods
    ? the resulting interface has no name; callers may label it
    # synthesis
  overspec.analyze_parameter
    fn (sig: function_signature, body: function_body, param_index: i32, types: type_index) -> optional[overspec_finding]
    + returns a finding describing the minimal interface and which methods to drop
    - returns none when the parameter is already minimal
    - returns none when no methods are called on the parameter
    # analysis
  overspec.analyze_function
    fn (sig: function_signature, body: function_body, types: type_index) -> list[overspec_finding]
    + returns findings for each overspecified parameter in order
    # analysis
  overspec.format_finding
    fn (finding: overspec_finding) -> string
    + renders a finding as a human-readable suggestion
    # rendering
