# Requirement: "a list-of-lists library"

Stores named lists of links and returns them. The whole point is the registry.

std: (all units exist)

awesome
  awesome.new
    fn () -> awesome_state
    + returns an empty registry
    # construction
  awesome.register
    fn (state: awesome_state, name: string, url: string) -> awesome_state
    + adds an entry to the registry
    # mutation
  awesome.all
    fn (state: awesome_state) -> list[tuple[string, string]]
    + returns every (name, url) pair in insertion order
    # read
