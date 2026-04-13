# Requirement: "a library of named unicode symbols with ASCII fallbacks for terminals that cannot render them"

Callers ask for a symbol by name; the library returns the unicode glyph or an ASCII fallback depending on whether unicode is supported.

std: (all units exist)

symbols
  symbols.get
    @ (name: string, unicode_supported: bool) -> string
    + returns the unicode glyph when supported (e.g. "check" -> "\u2714")
    + returns the ASCII fallback when not supported (e.g. "check" -> "v")
    ? the name-to-glyph table is hardcoded
    # lookup
  symbols.names
    @ () -> list[string]
    + returns every known symbol name
    # introspection
  symbols.detect_unicode_supported
    @ (term_env: string, is_windows: bool) -> bool
    + returns true when the terminal environment indicates unicode support
    ? inputs are environment facts so the caller controls detection
    # detection
