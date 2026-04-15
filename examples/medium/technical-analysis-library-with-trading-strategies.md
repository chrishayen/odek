# Requirement: "a technical analysis library with trading strategies"

Builds a few classic indicators and composes them into entry/exit signals.

std: (all units exist)

technical_analysis
  technical_analysis.macd
    fn (prices: list[f64], fast: i32, slow: i32, signal: i32) -> macd_result
    + returns the macd line, signal line, and histogram
    - leading entries before slow are NaN
    # indicator
  technical_analysis.bollinger_bands
    fn (prices: list[f64], period: i32, std_multiplier: f64) -> bands_result
    + returns upper, middle, and lower bands
    ? middle band is the simple moving average
    # indicator
  technical_analysis.atr
    fn (highs: list[f64], lows: list[f64], closes: list[f64], period: i32) -> list[f64]
    + returns the average true range over period
    - returns empty list when input lengths differ
    # volatility
  technical_analysis.signal_from_macd
    fn (macd: macd_result) -> list[i8]
    + emits +1 when the macd line crosses above signal, -1 below, 0 otherwise
    # strategy
  technical_analysis.evaluate_strategy
    fn (prices: list[f64], signals: list[i8]) -> strategy_report
    + returns total return, win rate, and trade count
    - returns zero stats when signals never fire
    # evaluation
