# Requirement: "a library returning colored symbols for various log levels"

Returns pre-styled symbol strings indexed by level name.

std: (all units exist)

logsymbols
  logsymbols.symbol
    fn (level: string) -> string
    + returns an ANSI-colored glyph for "info", "warn", "error", "success"
    - returns an uncolored fallback glyph for unknown levels
    ? colors use standard foreground escape sequences
    # presentation
