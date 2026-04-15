# Requirement: "a design-by-contract library that synchronizes code contracts with their documentation"

Parses contract annotations from doc comments and produces runtime checks plus a sync report.

std
  std.strings
    std.strings.split_lines
      fn (s: string) -> list[string]
      + splits on \n preserving order
      # strings
    std.strings.trim
      fn (s: string) -> string
      + removes leading and trailing whitespace
      # strings

contract
  contract.parse_doc
    fn (doc: string) -> result[contract_spec, string]
    + extracts pre, post, and invariant clauses from a doc comment
    - returns error on malformed clauses
    # parsing
    -> std.strings.split_lines
    -> std.strings.trim
  contract.check_preconditions
    fn (spec: contract_spec, bindings: map[string, string]) -> result[void, string]
    + returns ok when all preconditions evaluate true
    - returns a descriptive error naming the first failing precondition
    # runtime_checking
  contract.check_postconditions
    fn (spec: contract_spec, bindings: map[string, string]) -> result[void, string]
    + returns ok when all postconditions evaluate true
    - returns a descriptive error naming the first failing postcondition
    # runtime_checking
  contract.diff_against_source
    fn (doc_spec: contract_spec, source_spec: contract_spec) -> list[string]
    + returns a list of human-readable mismatches between doc and source
    + returns an empty list when the two are in sync
    # synchronization
