# Requirement: "a library that rewrites function bodies by inserting zero-value return statements matching each function's declared return types"

Operates on a simple language-agnostic function AST.

std: (all units exist)

zero_returns
  zero_returns.parse_signature
    @ (source: string) -> result[func_sig, string]
    + returns a parsed signature with name and ordered return types
    - returns error on malformed signature text
    # parsing
  zero_returns.zero_value_for
    @ (type_name: string) -> string
    + returns the canonical zero-value literal for known types
    ? known types: integer, float, string, bool, list, map, optional
    - returns empty string for unknown types
    # type_defaults
  zero_returns.build_return_stmt
    @ (types: list[string]) -> string
    + returns a single return statement with one zero value per type, comma-separated
    + returns empty string when types is empty
    # codegen
    -> zero_returns.zero_value_for
  zero_returns.needs_return
    @ (body: string) -> bool
    + returns true when the function body lacks a terminating return
    # analysis
  zero_returns.rewrite_function
    @ (source: string) -> result[string, string]
    + returns source with a zero-value return statement appended to bodies that need one
    + leaves already-terminated functions untouched
    - returns error when the signature cannot be parsed
    # rewriting
    -> zero_returns.parse_signature
    -> zero_returns.needs_return
    -> zero_returns.build_return_stmt
