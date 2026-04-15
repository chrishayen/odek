# Requirement: "a library to access record fields using a dotted path expression"

Navigates a nested map structure using a dotted path like `user.address.city`.

std: (all units exist)

pathget
  pathget.split_path
    fn (path: string) -> list[string]
    + splits a dotted path into its segments
    + returns an empty list for an empty string
    # parsing
  pathget.get
    fn (root: map[string, value], path: string) -> result[value, string]
    + walks root segment by segment and returns the final value
    - returns error when any intermediate segment is missing
    - returns error when an intermediate segment is not a map
    # lookup
    -> pathget.split_path
  pathget.set
    fn (root: map[string, value], path: string, new_value: value) -> result[map[string, value], string]
    + returns a new root with the value at path replaced, creating intermediate maps as needed
    - returns error when an intermediate segment exists but is not a map
    # mutation
    -> pathget.split_path
