# Requirement: "a library of idioms for writing simple, maintainable, testable code"

A tiny library that returns a curated list of design-principle strings. The caller decides how to display them.

std: (all units exist)

idioms
  idioms.list_principles
    @ () -> list[string]
    + returns a non-empty list of short design-principle strings
    ? contents are hardcoded; no I/O or configuration
    # principles
