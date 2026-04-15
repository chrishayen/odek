# Requirement: "a stock quote watcher and position tracker"

Holdings are a list of lots with cost basis. The library fetches quotes through a pluggable source and reports gain/loss per position.

std
  std.http
    std.http.get
      fn (url: string) -> result[bytes, string]
      + performs an HTTP GET and returns the body
      - returns error on non-2xx status
      # http
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses any json value
      - returns error on malformed input
      # serialization

portfolio
  portfolio.new
    fn () -> portfolio_state
    + creates an empty portfolio
    # construction
  portfolio.add_lot
    fn (state: portfolio_state, symbol: string, shares: f64, cost_per_share: f64) -> portfolio_state
    + records a purchase lot
    # mutation
  portfolio.remove_lot
    fn (state: portfolio_state, lot_id: string) -> result[portfolio_state, string]
    + removes a lot by id
    - returns error when lot_id is unknown
    # mutation
  portfolio.fetch_quote
    fn (source_url: string, symbol: string) -> result[f64, string]
    + fetches the latest price for symbol from a pluggable source endpoint
    - returns error on network or parse failure
    # data_fetch
    -> std.http.get
    -> std.json.parse
  portfolio.compute_positions
    fn (state: portfolio_state, quotes: map[string, f64]) -> list[position]
    + aggregates lots per symbol with shares, cost basis, market value, and unrealized gain
    - returns empty list when portfolio has no lots
    # aggregation
  portfolio.total_value
    fn (positions: list[position]) -> f64
    + sums market value across positions
    # aggregation
