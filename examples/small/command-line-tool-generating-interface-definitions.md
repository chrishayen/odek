# Requirement: "a library that derives interface definitions from a set of method signatures"

Given a list of method signatures, produce a minimal interface declaration covering all of them.

std: (all units exist)

interfaces
  interfaces.parse_signature
    @ (text: string) -> result[method_sig, string]
    + parses "name(arg_type, ...) -> return_type" into a structured signature
    - returns error when parentheses are unbalanced
    # parsing
  interfaces.collect
    @ (signatures: list[string]) -> result[list[method_sig], string]
    + parses every signature and returns them in declaration order
    - returns error on the first malformed line, identifying its index
    # collection
    -> interfaces.parse_signature
  interfaces.render_interface
    @ (name: string, methods: list[method_sig]) -> string
    + renders a single interface block containing every method in order
    ? duplicate method names are kept as-is; dedup is the caller's job
    # rendering
