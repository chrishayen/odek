# Requirement: "a hierarchical configuration management library"

Stores configuration as a nested tree keyed by dot-separated paths. Multiple sources (defaults, file-loaded values, overrides) layer on top of each other, with later layers winning.

std
  std.text
    std.text.split
      fn (s: string, sep: string) -> list[string]
      + splits s on every occurrence of sep
      # text

config
  config.new_store
    fn () -> store_state
    + creates an empty config store with no layers
    # construction
  config.push_layer
    fn (store: store_state, name: string) -> store_state
    + adds a new empty layer with the given name on top of existing layers
    # layering
  config.set
    fn (store: store_state, layer_name: string, path: string, value: string) -> store_state
    + sets a value at the dotted path in the named layer
    - returns the store unchanged when layer_name does not exist
    # write
    -> std.text.split
  config.get
    fn (store: store_state, path: string) -> optional[string]
    + returns the value at path from the topmost layer that defines it
    - returns none when no layer defines path
    # read
    -> std.text.split
  config.get_subtree
    fn (store: store_state, prefix: string) -> map[string, string]
    + returns a flat map of every path under prefix, merged across layers with upper layers winning
    # read
  config.keys
    fn (store: store_state) -> list[string]
    + returns every defined dotted path in the merged view
    # query
  config.remove
    fn (store: store_state, layer_name: string, path: string) -> store_state
    + removes a value from the named layer
    # write
  config.watch
    fn (store: store_state, path: string, cookie: i64) -> bool
    + returns true when the value at path has changed since the given cookie
    ? cookies are store version numbers; callers obtain them via config.version
    # change_detection
  config.version
    fn (store: store_state) -> i64
    + returns a monotonic version number that increments on every write
    # change_detection
