# Requirement: "a financial data interface library"

Fetches quotes and historical series from a generic HTTP endpoint and exposes them as typed records.

std
  std.http
    std.http.get
      @ (url: string) -> result[bytes, string]
      + issues an HTTP GET and returns the response body
      - returns error on non-2xx status
      # networking
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object
      - returns error on invalid JSON
      # serialization
    std.json.parse_array
      @ (raw: string) -> result[list[string], string]
      + parses a JSON array of objects as raw object strings
      - returns error on invalid JSON
      # serialization
  std.time
    std.time.parse_iso_date
      @ (text: string) -> result[i64, string]
      + parses a YYYY-MM-DD date into a unix epoch second
      - returns error on malformed input
      # time

akshare
  akshare.fetch_quote
    @ (symbol: string, base_url: string) -> result[quote, string]
    + fetches the latest quote for a symbol from the configured endpoint
    - returns error when the symbol is unknown
    # quotes
    -> std.http.get
    -> std.json.parse_object
  akshare.fetch_history
    @ (symbol: string, start: string, end: string, base_url: string) -> result[list[ohlcv], string]
    + fetches daily OHLCV bars for the inclusive date range
    - returns error when start is after end
    # history
    -> std.http.get
    -> std.json.parse_array
    -> std.time.parse_iso_date
  akshare.search_symbols
    @ (query: string, base_url: string) -> result[list[symbol_info], string]
    + searches for symbols matching the query string
    # discovery
    -> std.http.get
    -> std.json.parse_array
  akshare.to_returns
    @ (bars: list[ohlcv]) -> list[f64]
    + computes simple period-over-period returns from close prices
    + returns an empty list when the input has fewer than two bars
    # analytics
