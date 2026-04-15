# Requirement: "an accurate currency computation library"

Represents money as an integer minor-unit amount plus a currency code. Avoids floating point entirely for arithmetic.

std: (all units exist)

money
  money.from_minor
    fn (amount: i64, currency: string) -> money_value
    + constructs a money value from an integer minor-unit amount
    # construction
  money.parse
    fn (text: string, currency: string) -> result[money_value, string]
    + parses "123.45" into minor units using the currency's fraction digits
    - returns error when more fraction digits are given than the currency supports
    - returns error on non-numeric input
    # parsing
  money.format
    fn (m: money_value) -> string
    + renders amount with the currency's fraction digits and an ISO code suffix
    # formatting
  money.add
    fn (a: money_value, b: money_value) -> result[money_value, string]
    + adds two amounts in the same currency
    - returns error when currencies differ
    # arithmetic
  money.allocate
    fn (m: money_value, ratios: list[i32]) -> list[money_value]
    + splits m in proportion to ratios, distributing remainder minor units so the total is exactly m
    - returns empty list when ratios is empty
    # allocation
  money.convert
    fn (m: money_value, to_currency: string, rate_num: i64, rate_den: i64) -> money_value
    + converts using the rational rate num/den, rounding half-to-even
    # conversion
