# Requirement: "a millisecond conversion utility"

Parses human-readable duration strings (e.g. "2 days", "1h", "500ms") into milliseconds and back.

std: (all units exist)

ms
  ms.parse
    @ (text: string) -> result[i64, string]
    + returns 2000 for "2s" and 3600000 for "1h"
    + accepts units: ms, s, m, h, d, w, y
    - returns error for unrecognized units or malformed numbers
    # duration_parsing
  ms.format
    @ (millis: i64) -> string
    + returns "2s" for 2000 and "1h" for 3600000
    + picks the largest unit that produces a whole number
    # duration_formatting
