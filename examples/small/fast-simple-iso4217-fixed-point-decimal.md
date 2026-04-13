# Requirement: "an ISO 4217 fixed-point decimal money library"

Money values carry an amount (in minor units) and a currency code; arithmetic fails when currencies disagree.

std: (all units exist)

money
  money.from_minor
    @ (currency: string, minor_units: i64) -> result[money_value, string]
    + creates a money value from minor units (e.g. cents)
    - returns error when currency is not a known ISO 4217 code
    # construction
  money.from_decimal_string
    @ (currency: string, text: string) -> result[money_value, string]
    + parses a decimal string at the currency's fractional precision
    - returns error when the fractional digit count exceeds the currency's exponent
    - returns error on non-numeric input
    # parsing
  money.to_string
    @ (value: money_value) -> string
    + renders the amount with the currency's fractional digits and code (e.g. "12.50 USD")
    # formatting
  money.add
    @ (a: money_value, b: money_value) -> result[money_value, string]
    + sums two money values
    - returns error when currencies differ
    # arithmetic
  money.sub
    @ (a: money_value, b: money_value) -> result[money_value, string]
    + subtracts b from a
    - returns error when currencies differ
    # arithmetic
