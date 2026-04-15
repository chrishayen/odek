# Requirement: "a text wrapping library that breaks paragraphs at line-width boundaries"

Breaks a paragraph into lines no longer than a given width, splitting on word boundaries.

std: (all units exist)

textwrap
  textwrap.wrap
    fn (text: string, width: i32) -> list[string]
    + splits text into lines no longer than width, breaking at whitespace
    + returns a single-line list when text fits within width
    - returns an empty list when text is empty
    ? words longer than width are placed on their own line rather than split
    # wrapping
  textwrap.fill
    fn (text: string, width: i32) -> string
    + joins wrapped lines with newline separators
    # formatting
    -> textwrap.wrap
