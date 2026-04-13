# Requirement: "a library of validators and sanitizers for common input types"

A collection of focused, independent checks and transforms. Each rune does one thing.

std: (all units exist)

validate
  validate.is_email
    @ (s: string) -> bool
    + returns true for strings matching a single-at-sign local@domain shape with a non-empty TLD
    - returns false when there is no '@' or no '.' in the domain part
    # validation
  validate.is_url
    @ (s: string) -> bool
    + returns true when the string starts with http:// or https:// and has a non-empty host
    - returns false on missing scheme
    # validation
  validate.is_int_in_range
    @ (s: string, min_val: i64, max_val: i64) -> bool
    + returns true when s parses as an integer inside [min_val, max_val]
    - returns false on non-numeric input or out-of-range values
    # validation
  validate.trim_whitespace
    @ (s: string) -> string
    + returns the string with leading and trailing whitespace removed
    # sanitization
  validate.strip_html_tags
    @ (s: string) -> string
    + removes anything between '<' and '>' characters
    + preserves text content between tags
    # sanitization
  validate.normalize_email
    @ (s: string) -> string
    + lowercases the domain part and leaves the local part untouched
    - returns the input unchanged when there is no '@'
    # sanitization
  validate.all_non_empty
    @ (items: list[string]) -> bool
    + returns true when every item has length > 0
    - returns false when the list is empty
    - returns false when any item is ""
    # validation
