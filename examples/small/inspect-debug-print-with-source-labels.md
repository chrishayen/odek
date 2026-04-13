# Requirement: "a debug print library that labels values with their source expression"

Given an expression text and its runtime value, format a human-readable inspection line.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

inspect
  inspect.format_value
    @ (expr: string, value: string) -> string
    + returns "expr = value" for single-value inspection
    # formatting
  inspect.format_many
    @ (pairs: list[tuple[string,string]]) -> string
    + returns a comma-separated inspection of many expression-value pairs
    # formatting
  inspect.format_with_location
    @ (file: string, line: i32, expr: string, value: string) -> string
    + prefixes the inspection with "file:line"
    # formatting
  inspect.format_timed
    @ (expr: string, value: string) -> string
    + prefixes the inspection with the current timestamp
    # formatting
    -> std.time.now_millis
