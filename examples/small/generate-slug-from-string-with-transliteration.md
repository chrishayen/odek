# Requirement: "a library that turns a freeform string into a URL slug, transliterating non-ASCII characters to ASCII"

Normalizes, transliterates, then joins words with a separator.

std: (all units exist)

speakingurl
  speakingurl.transliterate
    @ (text: string) -> string
    + replaces accented letters with their base ASCII equivalents
    + replaces common non-latin letters with ASCII equivalents when a mapping exists
    + leaves already-ASCII characters unchanged
    # transliteration
  speakingurl.tokenize_words
    @ (text: string) -> list[string]
    + splits on any run of non-alphanumeric characters
    + drops empty tokens
    # tokenizing
  speakingurl.slugify
    @ (text: string, separator: string) -> string
    + lowercases, transliterates, tokenizes, then joins with the separator
    + returns "" when the input has no alphanumeric characters
    # slugify
