# Requirement: "a library for producing URL-friendly slugs from text with multi-language support"

Transliterates non-ASCII letters to ASCII using a per-language substitution table, then normalizes to lowercase hyphenated form.

std: (all units exist)

slugify
  slugify.transliterate
    fn (text: string, language: string) -> string
    + replaces accented and non-ASCII letters with ASCII equivalents using the per-language table
    + falls back to a generic table when the language is unknown
    ? language is an ISO 639-1 code like "en", "de", "pl"
    # transliteration
  slugify.normalize
    fn (text: string) -> string
    + lowercases the input and replaces runs of non-alphanumeric characters with a single hyphen
    + trims leading and trailing hyphens
    - empty input returns an empty string
    # normalization
  slugify.make
    fn (text: string, language: string) -> string
    + returns a slug by transliterating then normalizing
    + idempotent on already-slugified input
    # slugification
