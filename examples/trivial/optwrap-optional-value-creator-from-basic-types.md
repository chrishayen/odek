# Requirement: "functions for simplified creation of optional values from basic-type constants"

Single helper that lifts a constant into an optional wrapper.

std: (all units exist)

optwrap
  optwrap.of
    @ (value: string) -> optional[string]
    + wraps the given value in a present optional
    ? generic over basic types; string shown as the canonical case
    # optional_construction
