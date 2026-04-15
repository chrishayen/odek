# Requirement: "a library of string manipulation and conversion functions"

Common transforms beyond what language primitives cover.

std: (all units exist)

strutil
  strutil.to_snake_case
    fn (s: string) -> string
    + converts CamelCase and kebab-case to snake_case
    + collapses runs of separators
    # case_conversion
  strutil.to_camel_case
    fn (s: string) -> string
    + converts snake_case and kebab-case to camelCase
    # case_conversion
  strutil.pad_left
    fn (s: string, width: i32, fill: string) -> string
    + pads s on the left to the target width using fill
    + returns s unchanged when already at or above width
    # padding
  strutil.truncate
    fn (s: string, max_len: i32, ellipsis: string) -> string
    + truncates to max_len codepoints, appending ellipsis if shortened
    # trimming
  strutil.strip_accents
    fn (s: string) -> string
    + removes combining diacritic marks
    # normalization
