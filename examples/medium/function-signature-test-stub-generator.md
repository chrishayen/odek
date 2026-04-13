# Requirement: "a test stub generator from function signatures"

Parses a source snippet and produces test stubs naming each function and argument, so callers can fill in assertions.

std
  std.strings
    std.strings.trim
      @ (s: string) -> string
      + returns s without leading or trailing whitespace
      # text
    std.strings.join
      @ (items: list[string], sep: string) -> string
      + concatenates items separated by sep
      + returns "" for an empty list
      # text

stub_gen
  stub_gen.parse_signatures
    @ (source: string) -> result[list[function_signature], string]
    + extracts function name, argument names with types, and return types
    + ignores comments and whitespace
    - returns error on unbalanced parentheses in a signature
    # parsing
    -> std.strings.trim
  stub_gen.signature_to_stub
    @ (sig: function_signature) -> string
    + produces a test stub that declares inputs, calls the function, and marks an assertion todo
    + includes one section per argument
    # emission
    -> std.strings.join
  stub_gen.render_stubs
    @ (signatures: list[function_signature]) -> string
    + concatenates stubs for all signatures into a single test file body
    + returns "" when signatures is empty
    # emission
    -> std.strings.join
