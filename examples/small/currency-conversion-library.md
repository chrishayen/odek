# Requirement: "a currency conversion library"

A rate table keyed by ISO currency code, with a conversion function that triangulates through a base currency.

std: (all units exist)

currency_convert
  currency_convert.new_rate_table
    fn (base: string) -> rate_table
    + creates an empty table whose rates are expressed relative to the base currency
    # construction
  currency_convert.set_rate
    fn (table: rate_table, code: string, rate_per_base: f64) -> result[rate_table, string]
    + records that one unit of base equals rate_per_base units of code
    - returns error when rate_per_base is not positive
    # rates
  currency_convert.convert
    fn (table: rate_table, amount: f64, from_code: string, to_code: string) -> result[f64, string]
    + returns the amount converted by triangulating from_code to base to to_code
    - returns error when either code is missing from the table
    - returns error when amount is negative
    # conversion
  currency_convert.list_codes
    fn (table: rate_table) -> list[string]
    + returns all known currency codes in sorted order
    # introspection
