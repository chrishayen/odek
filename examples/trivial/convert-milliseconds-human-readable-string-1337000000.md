# Requirement: "convert milliseconds to a human readable duration string"

One function that formats a millisecond count as "15d 11h 23m 20s".

std: (all units exist)

pretty_ms
  pretty_ms.humanize
    @ (ms: i64) -> string
    + formats 1337000000 as "15d 11h 23m 20s"
    + formats 0 as "0ms"
    + omits zero-valued units except when the entire duration is zero
    + uses the units d, h, m, s, ms in decreasing order
    - formats negative inputs with a leading "-"
    # formatting
