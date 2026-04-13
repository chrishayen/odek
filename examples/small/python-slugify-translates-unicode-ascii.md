# Requirement: "a slugify library that translates unicode to ASCII"

Normalizes arbitrary text into a URL-friendly slug.

std: (all units exist)

slugify
  slugify.transliterate
    @ (input: string) -> string
    + returns input with diacritics stripped and common unicode replaced by ASCII equivalents
    ? unknown code points are dropped
    # transliteration
  slugify.slugify
    @ (input: string) -> string
    + returns a lowercase slug with non-alphanumeric runs collapsed to a single separator
    + strips leading and trailing separators
    - returns "" for input that contains no alphanumeric characters
    # slug
  slugify.slugify_with
    @ (input: string, separator: string, max_length: i32) -> string
    + returns a slug using the given separator and truncated to at most max_length
    + truncates on a separator boundary when possible
    # slug
