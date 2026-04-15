# Requirement: "a library for representing currency amounts with metadata and locale-aware formatting"

Amounts are stored as integer minor units (e.g. cents) to avoid floating-point rounding.

std: (all units exist)

money
  money.new
    fn (minor_units: i64, currency_code: string) -> result[money_state, string]
    + creates a money value in the given currency
    - returns error when the currency code is not recognized
    # construction
  money.add
    fn (a: money_state, b: money_state) -> result[money_state, string]
    + returns the sum when both values share the same currency
    - returns error when currencies differ
    # arithmetic
  money.subtract
    fn (a: money_state, b: money_state) -> result[money_state, string]
    + returns the difference when both values share the same currency
    - returns error when currencies differ
    # arithmetic
  money.currency_info
    fn (currency_code: string) -> result[currency_info, string]
    + returns metadata for a currency: symbol, name, and minor unit digits
    - returns error when the currency code is not recognized
    # metadata
  money.format
    fn (amount: money_state, locale: string) -> string
    + renders the amount using the locale's grouping, decimal, and symbol position conventions
    # formatting
