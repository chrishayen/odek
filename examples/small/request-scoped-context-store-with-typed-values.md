# Requirement: "a library for storing and retrieving typed values in a request-scoped context"

A shallow immutable key-value store keyed by string, with typed accessors for common value kinds.

std: (all units exist)

context_store
  context_store.new
    @ () -> context_state
    + creates an empty context
    # construction
  context_store.with_value
    @ (ctx: context_state, key: string, value: string) -> context_state
    + returns a new context with the key bound
    + overwrites any existing binding for that key
    # mutation
  context_store.get
    @ (ctx: context_state, key: string) -> optional[string]
    + returns the bound value for a key
    - returns none when the key is absent
    # lookup
  context_store.get_int
    @ (ctx: context_state, key: string) -> optional[i64]
    + parses the bound value as a 64-bit integer and returns it
    - returns none when the key is absent or the value is not a valid integer
    # lookup
