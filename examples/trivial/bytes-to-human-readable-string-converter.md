# Requirement: "convert bytes to a human readable string"

A single function that formats a byte count into a compact SI-style string like "1.34 kB".

std: (all units exist)

bytes_format
  bytes_format.humanize
    fn (n: i64) -> string
    + formats 0 as "0 B"
    + formats 1337 as "1.34 kB"
    + formats values at or above 1000 bytes using kB, MB, GB, TB, PB
    - formats negative inputs with a leading minus sign
    ? uses SI (1000-based) units, two decimal places for non-byte magnitudes
    # formatting
