# Requirement: "generic option and result types that optionally contain a value"

Option and result helpers as simple constructors and inspectors over already-supported generic result/optional types.

std: (all units exist)

valor
  valor.some
    @ (value: string) -> optional[string]
    + wraps a value as a present optional
    # option_construction
  valor.none
    @ () -> optional[string]
    + returns an empty optional
    # option_construction
  valor.unwrap_or
    @ (opt: optional[string], default_value: string) -> string
    + returns the inner value when present, otherwise the default
    # option_inspection
  valor.ok
    @ (value: string) -> result[string, string]
    + wraps a value as a success result
    # result_construction
  valor.err
    @ (message: string) -> result[string, string]
    + wraps a message as an error result
    # result_construction
