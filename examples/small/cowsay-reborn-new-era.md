# Requirement: "a library that wraps text in a speech balloon drawn by an ASCII character"

Word-wrap the message, frame it in a balloon, and attach a selectable ASCII figure below.

std
  std.text
    std.text.word_wrap
      @ (s: string, width: i32) -> list[string]
      + wraps text to lines no longer than width, breaking on spaces
      + a word longer than width is placed on its own line
      # text

cowsay
  cowsay.balloon
    @ (message: string, width: i32) -> string
    + draws a speech balloon around the wrapped message
    + uses single-line corners (/ \ < >) and dashes/pipes for edges
    # rendering
    -> std.text.word_wrap
  cowsay.figure
    @ (name: string) -> result[string, string]
    + returns the ASCII figure with the given name
    - returns error for unknown figure names
    # assets
  cowsay.say
    @ (message: string, figure_name: string, width: i32) -> result[string, string]
    + combines balloon and figure into the final ASCII output
    # composition
