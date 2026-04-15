# Requirement: "a string case conversion and slugification library"

Converts arbitrary strings between camel case, snake case, kebab case, and URL-safe slugs. All functions share a common word-split primitive.

std: (all units exist)

stringcase
  stringcase.split_words
    fn (input: string) -> list[string]
    + splits on whitespace, punctuation, and camel-case boundaries into lowercase words
    + treats consecutive uppercase letters as a single word
    - returns an empty list for empty input
    # tokenization
  stringcase.to_camel_case
    fn (input: string) -> string
    + first word lowercase, subsequent words capitalized, no separators
    # formatting
    -> stringcase.split_words
  stringcase.to_snake_case
    fn (input: string) -> string
    + all lowercase, words joined by underscores
    # formatting
    -> stringcase.split_words
  stringcase.to_kebab_case
    fn (input: string) -> string
    + all lowercase, words joined by hyphens
    # formatting
    -> stringcase.split_words
  stringcase.to_slug
    fn (input: string) -> string
    + lowercase, hyphen-joined, ASCII-only, non-alphanumerics dropped
    + collapses multiple hyphens into one
    # formatting
    -> stringcase.split_words
