# Requirement: "common structured-log attribute formatters and a helper to compose new ones"

Formatters transform a (key, value) attribute into a (key, value) pair for a structured logger. The library provides a few standard ones plus a combinator.

std: (all units exist)

log_format
  log_format.redact
    @ (keys: list[string]) -> attr_formatter
    + returns a formatter that replaces the value of any attribute whose key is in keys with "***"
    + passes through attributes whose key is not in keys
    # redaction
  log_format.rename
    @ (mapping: map[string, string]) -> attr_formatter
    + returns a formatter that rewrites attribute keys according to mapping
    + leaves unknown keys unchanged
    # renaming
  log_format.error_to_string
    @ () -> attr_formatter
    + returns a formatter that stringifies error-typed values to their message
    # error_shaping
  log_format.chain
    @ (formatters: list[attr_formatter]) -> attr_formatter
    + returns a formatter that applies each formatter in order, feeding the output of one into the next
    - when formatters is empty, returns an identity formatter
    # composition
  log_format.apply
    @ (f: attr_formatter, key: string, value: string) -> tuple[string, string]
    + runs formatter f on (key, value) and returns the resulting pair
    # invocation
