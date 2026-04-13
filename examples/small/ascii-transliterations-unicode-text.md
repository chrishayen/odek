# Requirement: "a library that transliterates Unicode text to a plain ASCII approximation"

Each code point is mapped to an ASCII replacement via a static table; unknown code points collapse to an empty string.

std
  std.unicode
    std.unicode.code_points
      @ (text: string) -> list[i32]
      + decodes the string into a sequence of Unicode scalar values
      # unicode

unidecode
  unidecode.transliterate
    @ (text: string) -> string
    + returns an ASCII approximation built by concatenating per-code-point replacements
    + returns an empty string for empty input
    ? unknown code points contribute nothing
    # transliteration
    -> std.unicode.code_points
  unidecode.replacement_for
    @ (code_point: i32) -> string
    + returns the ASCII replacement for a single code point, or "" if none
    + ASCII code points below 128 return themselves
    # lookup
