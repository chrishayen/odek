# Requirement: "a terminal string colorizer that highlights substring matches"

Wraps matched substrings with ANSI color escape codes. One function compiles a match spec, another applies it.

std: (all units exist)

marker
  marker.mark
    fn (input: string, needle: string, color: i32) -> string
    + wraps every occurrence of needle in input with the ANSI color code and a reset
    - returns input unchanged when needle is empty
    - returns input unchanged when needle does not occur
    ? color is an integer in the standard ANSI 0-255 palette
    # highlighting
  marker.mark_many
    fn (input: string, rules: list[mark_rule]) -> string
    + applies every rule in order, wrapping each matched substring with its color
    + later rules do not re-wrap text already inside an escape sequence
    # highlighting
