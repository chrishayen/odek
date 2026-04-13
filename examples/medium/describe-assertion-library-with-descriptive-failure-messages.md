# Requirement: "an assertion library that produces descriptive failure messages by inspecting the expression"

Evaluates a boolean expression and, on failure, walks the expression tree to report each subexpression's value.

std: (all units exist)

describe_assert
  describe_assert.parse_expression
    @ (source: string) -> result[expr_node, string]
    + parses a boolean expression into an AST of literals, identifiers, calls, and operators
    - returns error on malformed source
    # parsing
  describe_assert.evaluate
    @ (node: expr_node, bindings: map[string, value]) -> result[value, string]
    + evaluates the AST against the provided variable bindings
    - returns error when an identifier is unbound
    # evaluation
  describe_assert.trace
    @ (node: expr_node, bindings: map[string, value]) -> list[sub_result]
    + returns each subexpression's source text alongside its evaluated value
    # introspection
  describe_assert.format_failure
    @ (source: string, trace: list[sub_result]) -> string
    + renders a multi-line failure message underlining each subexpression with its value
    # formatting
  describe_assert.check
    @ (source: string, bindings: map[string, value]) -> result[void, string]
    + parses, evaluates, and returns ok when the expression is truthy
    - returns a formatted descriptive failure when the expression evaluates to false
    # assertion
