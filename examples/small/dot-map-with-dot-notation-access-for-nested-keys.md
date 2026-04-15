# Requirement: "a map type with dot-notation access for nested keys"

Wraps a nested map so callers can read and write "a.b.c" paths without manually drilling.

std: (all units exist)

dot_map
  dot_map.new
    fn () -> dot_map_state
    + creates an empty dot map
    # construction
  dot_map.from_map
    fn (source: map[string, string]) -> dot_map_state
    + wraps an existing flat map
    # construction
  dot_map.get
    fn (state: dot_map_state, path: string) -> optional[string]
    + returns the value at a dotted path or none when any segment is missing
    ? segments are split on "." and walked left-to-right
    # access
  dot_map.set
    fn (state: dot_map_state, path: string, value: string) -> dot_map_state
    + writes value at a dotted path, creating intermediate nested maps as needed
    # mutation
  dot_map.delete
    fn (state: dot_map_state, path: string) -> dot_map_state
    + removes the value at a dotted path if present
    # mutation
  dot_map.has
    fn (state: dot_map_state, path: string) -> bool
    + returns true when the full path exists
    - returns false when any intermediate segment is missing
    # access
