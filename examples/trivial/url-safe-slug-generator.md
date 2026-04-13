# Requirement: "a slugify function that turns a string into a URL-safe slug"

std: (all units exist)

slugify
  slugify.to_slug
    @ (input: string) -> string
    + lowercases, strips diacritics, and replaces runs of non-alphanumeric characters with a single hyphen
    + trims leading and trailing hyphens
    - returns "" for an input with no alphanumeric characters
    # slugification
