# Requirement: "a library for protecting handcrafted SQL strings against injection"

Given a SQL string with named placeholders and a parameter map, validates that every placeholder is bound to a whitelisted type, then returns the safe query plus ordered values.

std: (all units exist)

sql_guard
  sql_guard.extract_placeholders
    fn (query: string) -> list[string]
    + returns all :name placeholders in the order they appear
    + returns an empty list when the query has none
    # parsing
  sql_guard.bind
    fn (query: string, params: map[string, sql_value]) -> result[prepared_query, string]
    + returns a prepared_query with ? markers in placeholder order and positional values
    - returns error when a placeholder has no matching parameter
    - returns error when a parameter is not one of the allowed primitive types
    # binding
  sql_guard.is_identifier_safe
    fn (ident: string) -> bool
    + returns true when ident matches ^[A-Za-z_][A-Za-z0-9_]*$
    - returns false when ident contains whitespace, quotes, or semicolons
    ? used for callers that must splice table or column names
    # validation
  sql_guard.quote_identifier
    fn (ident: string) -> result[string, string]
    + returns the identifier wrapped in double quotes with internal quotes doubled
    - returns error when the identifier is empty
    # escaping
