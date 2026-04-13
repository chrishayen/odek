# Requirement: "a collection of convenience string and array helpers familiar to web developers"

A small utility bundle: string padding, an implode helper, and an array slice with defaults.

std: (all units exist)

web_helpers
  web_helpers.str_pad_left
    @ (value: string, total_width: i32, pad: string) -> string
    + left-pads value with pad until its length reaches total_width
    + returns value unchanged when it is already at least total_width
    # strings
  web_helpers.join_with
    @ (parts: list[string], separator: string) -> string
    + concatenates parts with separator between each
    + returns "" for an empty list
    # strings
  web_helpers.array_slice
    @ (values: list[string], offset: i32, length: i32) -> list[string]
    + returns the subrange starting at offset with up to length elements
    + clamps length to the tail of the list
    - returns an empty list when offset is beyond the end
    # collections
  web_helpers.array_get
    @ (values: map[string, string], key: string, default_value: string) -> string
    + returns the value at key or default_value when missing
    # collections
