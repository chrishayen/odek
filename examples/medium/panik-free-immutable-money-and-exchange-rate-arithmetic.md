# Requirement: "a library for immutable monetary amounts and exchange rates with panic-free arithmetic"

Amounts are tagged with a currency code; exchange uses a rate table. All ops return result so bad inputs never panic.

std: (all units exist)

money
  money.new
    @ (units: i64, nanos: i32, currency: string) -> result[money_amount, string]
    + returns an amount where fractional value = nanos / 1e9
    - returns error when |nanos| >= 1e9
    - returns error when currency is not exactly 3 uppercase letters
    # construction
  money.from_string
    @ (s: string, currency: string) -> result[money_amount, string]
    + parses "12.34" into a money_amount with the given currency
    - returns error on non-numeric input
    # parsing
  money.to_string
    @ (m: money_amount) -> string
    + returns "123.45 USD" with trailing zeros trimmed but scale preserved
    # formatting
  money.add
    @ (a: money_amount, b: money_amount) -> result[money_amount, string]
    + returns a + b in the shared currency
    - returns error when currencies differ
    - returns error on overflow
    # arithmetic
  money.sub
    @ (a: money_amount, b: money_amount) -> result[money_amount, string]
    + returns a - b in the shared currency
    - returns error when currencies differ
    # arithmetic
  money.mul_scalar
    @ (m: money_amount, factor_num: i64, factor_den: i64) -> result[money_amount, string]
    + returns m * (num/den) with half-even rounding
    - returns error when den is zero
    # arithmetic
  money.compare
    @ (a: money_amount, b: money_amount) -> result[i32, string]
    + returns -1, 0, or 1
    - returns error when currencies differ
    # comparison
  money.rate_table_new
    @ () -> rate_table
    + returns an empty rate table
    # construction
  money.rate_table_set
    @ (t: rate_table, from: string, to: string, num: i64, den: i64) -> result[rate_table, string]
    + returns a new table with the rate from->to set to num/den
    - returns error when den is zero
    # rate_table
  money.exchange
    @ (m: money_amount, target: string, t: rate_table) -> result[money_amount, string]
    + returns the amount converted to target currency using the table
    - returns error when the rate is missing
    # exchange
