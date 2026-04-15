# Requirement: "a library that returns a view into a subdirectory of an embedded filesystem"

Treat an embedded filesystem as a map of paths to bytes and return a new one rooted at a subdirectory.

std: (all units exist)

subfs
  subfs.sub
    fn (files: map[string, bytes], prefix: string) -> result[map[string, bytes], string]
    + returns a new map containing only entries under prefix, with the prefix stripped from each key
    - returns error when no entries match the prefix
    ? trailing slash on prefix is optional and normalized internally
    # filesystem
