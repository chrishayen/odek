# Requirement: "a technical analysis library with indicators, strategies, and a backtest runner"

Indicators compute rolling metrics over price series. Strategies turn indicator output into trade signals. The backtest runner walks a price series and reports realized P&L.

std: (all units exist)

technical_analysis
  technical_analysis.sma
    fn (prices: list[f64], period: i32) -> list[f64]
    + returns the simple moving average; first (period-1) entries are NaN
    - returns an empty list when prices is empty
    # indicator
  technical_analysis.ema
    fn (prices: list[f64], period: i32) -> list[f64]
    + returns the exponential moving average with smoothing factor 2/(period+1)
    # indicator
  technical_analysis.rsi
    fn (prices: list[f64], period: i32) -> list[f64]
    + returns the relative strength index in [0, 100]
    - leading entries before period are NaN
    # indicator
  technical_analysis.crossover_strategy
    fn (fast: list[f64], slow: list[f64]) -> list[i8]
    + returns +1 on bullish cross, -1 on bearish cross, 0 otherwise
    ? output length equals the shorter of the two inputs
    # strategy
  technical_analysis.backtest
    fn (prices: list[f64], signals: list[i8], starting_cash: f64) -> backtest_result
    + walks the series: +1 buys with all cash, -1 sells all holdings
    + returns final equity, number of trades, and max drawdown
    - returns the starting balance unchanged when signals are all zero
    # backtesting
