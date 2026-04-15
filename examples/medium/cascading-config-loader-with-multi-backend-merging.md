# Requirement: "a cascading configuration loader that merges values from multiple backends"

Each backend is consulted in order; later backends override earlier ones only for keys they actually define.

std: (all units exist)

config_cascade
  config_cascade.new
    fn () -> cascade_state
    + creates an empty cascade with no backends registered
    # construction
  config_cascade.add_backend
    fn (state: cascade_state, name: string, values: map[string, string]) -> cascade_state
    + appends a backend with the given name and its key/value pairs
    ? backend order determines override precedence; later wins
    # registration
  config_cascade.resolve
    fn (state: cascade_state) -> map[string, string]
    + returns the merged map where later backends override earlier ones for shared keys
    + preserves keys defined only in a single backend
    # merge
  config_cascade.get
    fn (state: cascade_state, key: string) -> optional[string]
    + returns the effective value for the key after cascading
    - returns none when no backend defines the key
    # query
  config_cascade.source_of
    fn (state: cascade_state, key: string) -> optional[string]
    + returns the name of the backend whose value ultimately wins for the key
    - returns none when no backend defines the key
    # provenance
