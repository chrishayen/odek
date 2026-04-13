# Requirement: "a case conversion library that preserves common initialisms"

Convert between snake_case, camelCase, and PascalCase while keeping recognized initialisms (HTTP, URL, ID, etc.) intact rather than treating them as separate words.

std: (all units exist)

case_convert
  case_convert.split_words
    @ (s: string, initialisms: list[string]) -> list[string]
    + splits s into lowercase word tokens, recognizing underscores, hyphens, and case transitions as boundaries
    + keeps each initialism in initialisms as a single token regardless of surrounding casing
    # tokenization
  case_convert.to_snake
    @ (s: string, initialisms: list[string]) -> string
    + joins the tokens of s with underscores, all lowercase
    + preserves initialisms as uppercase runs when round-tripping
    # snake_case
  case_convert.to_camel
    @ (s: string, initialisms: list[string]) -> string
    + lowercases the first token and capitalizes subsequent tokens
    + emits each initialism token in uppercase when it is not the first token
    # camelCase
  case_convert.to_pascal
    @ (s: string, initialisms: list[string]) -> string
    + capitalizes every token
    + emits each initialism token in uppercase
    # PascalCase
