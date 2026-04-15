# Requirement: "convert an arbitrary string into a valid filename"

Replaces reserved characters, trims length, and avoids reserved names.

std: (all units exist)

filenamify
  filenamify.sanitize
    fn (input: string, replacement: string) -> string
    + replaces /, \, ?, *, :, |, ", <, > with the replacement
    + collapses runs of replacements into a single replacement
    + strips control characters
    + trims trailing dots and spaces
    + truncates the stem to keep the total length at or below 255
    + appends "_" to reserved device names like "CON", "PRN", "AUX", "NUL"
    ? preserves the final extension when truncating
    # sanitization
