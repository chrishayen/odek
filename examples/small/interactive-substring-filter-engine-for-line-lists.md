# Requirement: "an interactive substring filtering engine for a list of lines"

Given a list of candidate lines and a query string, returns the matching lines in stable order. The interaction loop (reading keystrokes, redrawing) is the caller's responsibility.

std: (all units exist)

filter
  filter.new
    @ (candidates: list[string]) -> filter_state
    + creates a filter state holding the candidate lines and an empty query
    # construction
  filter.set_query
    @ (state: filter_state, query: string) -> filter_state
    + returns a state with the query replaced
    # query
  filter.matches
    @ (state: filter_state) -> list[string]
    + returns candidates containing the query as a case-insensitive substring, preserving order
    + returns every candidate when the query is empty
    # filtering
  filter.select
    @ (state: filter_state, index: i32) -> result[string, string]
    + returns the nth matching line (0-indexed)
    - returns error when index is out of range
    # selection
    -> filter.matches
