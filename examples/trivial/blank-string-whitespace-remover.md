# Requirement: "a library to verify or remove blank and whitespace characters from strings"

std: (all units exist)

blank
  blank.is_blank
    @ (s: string) -> bool
    + returns true when s is empty or contains only whitespace characters
    - returns false when s contains any non-whitespace character
    # check
  blank.remove_whitespace
    @ (s: string) -> string
    + returns s with every whitespace character removed
    + returns "" when the input contains only whitespace
    # transform
