# Requirement: "a pipeline library for filter, map, and reduce over line-delimited input"

Compose stateless transformations over an input stream of lines and return the transformed output.

std: (all units exist)

linepipe
  linepipe.filter
    @ (lines: list[string], predicate_source: string) -> result[list[string], string]
    + returns only lines for which the predicate evaluates truthy
    - returns error when the predicate source has a syntax error
    # filter
  linepipe.map
    @ (lines: list[string], expr_source: string) -> result[list[string], string]
    + returns each line transformed by the expression
    - returns error when the expression source has a syntax error
    # map
  linepipe.reduce
    @ (lines: list[string], init: string, expr_source: string) -> result[string, string]
    + folds the lines into a single string using the expression
    - returns error when the expression source has a syntax error
    # reduce
