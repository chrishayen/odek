# Requirement: "a community profile card builder"

The input was a local user-group entry; interpreted as a simple profile card model.

std: (all units exist)

community
  community.new_profile
    fn (name: string, location: string) -> profile
    + creates a profile with the given name and location
    # construction
  community.format_profile
    fn (p: profile) -> string
    + returns "<name> — <location>"
    # formatting
