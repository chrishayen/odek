# Requirement: "a library for fixed-width text formatting with UTF-8 support"

Pads and truncates strings to a target display width, counting grapheme widths rather than bytes so UTF-8 content aligns correctly.

std
  std.unicode
    std.unicode.grapheme_width
      fn (grapheme: string) -> i32
      + returns display columns occupied by a grapheme (0, 1, or 2)
      + returns 2 for CJK and emoji
      # unicode
    std.unicode.grapheme_iter
      fn (s: string) -> list[string]
      + splits a string into grapheme clusters
      # unicode

fixed_width
  fixed_width.display_width
    fn (s: string) -> i32
    + returns the total display width of s in columns
    # measurement
    -> std.unicode.grapheme_iter
    -> std.unicode.grapheme_width
  fixed_width.pad_right
    fn (s: string, width: i32, fill: string) -> string
    + returns s padded on the right with fill until display width equals width
    - returns s unchanged when its display width already meets or exceeds width
    # formatting
  fixed_width.pad_left
    fn (s: string, width: i32, fill: string) -> string
    + returns s padded on the left with fill until display width equals width
    # formatting
  fixed_width.truncate
    fn (s: string, width: i32) -> string
    + returns a prefix of s whose display width does not exceed width
    ? never splits a grapheme cluster
    # formatting
  fixed_width.fit
    fn (s: string, width: i32, fill: string) -> string
    + pads or truncates s to exactly width columns
    # formatting
