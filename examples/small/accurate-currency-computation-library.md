# Requirement: "an accurate currency computation library"

Represents money as an integer minor-unit amount plus a currency code. Avoids floating point entirely for arithmetic.

std: (all units exist)

money
  money.from_minor
    @ (amount: i64, currency: string) -> money_value
    + constructs a money value from an integer minor-unit amount
    # construction
  money.parse
    @ (text: string, currency: string) -> result[money_value, string]
    + parses "123.45" into minor units using the currency's fraction digits
    - returns error when more fraction digits are given than the currency supports
    - returns error on non-numeric input
    # parsing
  money.format
    @ (m: money_value) -> string
    + renders amount with the currency's fraction digits and an ISO code suffix
    # formatting
  money.add
    @ (a: money_value, b: money_value) -> result[money_value, string]
    + adds two amounts in the same currency
    - returns error when currencies differ
    # arithmetic
  money.allocate
    @ (m: money_value, ratios: list[i32]) -> list[money_value]
    + splits m in proportion to ratios, distributing remainder minor units so the total is exactly m
    - returns empty list when ratios is empty
    # allocation
  money.convert
    @ (m: money_value, to_currency: string, rate_num: i64, rate_den: i64) -> money_value
    + converts using the rational rate num/den, rounding half-to-even
    # conversion
