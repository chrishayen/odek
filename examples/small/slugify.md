# Requirement: "slugify a string (lowercase, hyphens for spaces, strip punctuation)"

One pure function. All the cleaning fits inline — no helpers.

std: (all units exist)

slug
  slug.from_string
    @ (input: string) -> string
    + returns "hello-world" for "Hello World"
    + returns "foo-bar-baz" for "  foo   bar   baz  "
    + collapses consecutive whitespace into a single hyphen
    + strips characters that are not alphanumeric or hyphen
    + trims leading and trailing hyphens
    + output is lowercase ascii
    ? non-ascii characters are stripped, not transliterated
    # normalization
