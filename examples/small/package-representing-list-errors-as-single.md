# Requirement: "represent a list of errors as a single error"

An aggregator that collects individual error strings and exposes them as one flat error.

std: (all units exist)

multierr
  multierr.empty
    @ () -> multi_error
    + returns an empty multi-error
    # construction
  multierr.append
    @ (acc: multi_error, err: optional[string]) -> multi_error
    + returns acc with err appended when err is present
    + returns acc unchanged when err is absent
    ? absent errors are dropped so callers can append unconditionally
    # accumulation
  multierr.combine
    @ (a: multi_error, b: multi_error) -> multi_error
    + returns a multi-error whose entries are a's followed by b's
    # accumulation
  multierr.to_error
    @ (acc: multi_error) -> optional[string]
    + returns a joined error string separated by "; " when entries are present
    - returns absent when the multi-error is empty
    # presentation
  multierr.entries
    @ (acc: multi_error) -> list[string]
    + returns all error strings in insertion order
    # introspection
