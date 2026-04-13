# Requirement: "a library for streaming realtime price quotes and rendering them as a terminal chart"

Fetches quotes from a pluggable data source and renders a scrolling ASCII chart. Project layer owns the polling loop and chart model; std provides http, json, and time.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      @ (ms: i64) -> void
      + blocks the calling task for ms milliseconds
      # time
  std.http
    std.http.get
      @ (url: string) -> result[string, string]
      + performs an http GET and returns the response body as a string
      - returns error on non-2xx status or network failure
      # networking
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

ticker
  ticker.new_series
    @ (symbol: string, capacity: i32) -> series_state
    + creates a bounded ring-buffer series for a ticker symbol
    # construction
  ticker.push_sample
    @ (series: series_state, price: f64, ts_millis: i64) -> series_state
    + appends a (ts, price) sample; oldest is dropped at capacity
    # data_ingestion
  ticker.fetch_latest
    @ (source_url: string, symbol: string) -> result[tuple[f64, i64], string]
    + fetches and parses the latest price and timestamp for the symbol
    - returns error when the response is missing expected fields
    # data_ingestion
    -> std.http.get
    -> std.json.parse_object
    -> std.time.now_millis
  ticker.poll_once
    @ (series: series_state, source_url: string, symbol: string) -> result[series_state, string]
    + performs one fetch and appends the sample to the series
    - returns error when the fetch fails
    # data_ingestion
  ticker.render_chart
    @ (series: series_state, width: i32, height: i32) -> string
    + renders the series as an ASCII chart with axis labels
    + uses min/max of the visible window for vertical scaling
    # rendering
  ticker.run_polling
    @ (series: series_state, source_url: string, symbol: string, interval_ms: i64) -> void
    + polls at the given interval indefinitely, updating the series each tick
    ? intended to run in a dedicated task; terminates on source error
    # control
    -> std.time.sleep_millis
