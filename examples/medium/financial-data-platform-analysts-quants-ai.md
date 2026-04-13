# Requirement: "a financial data platform library with pluggable data sources"

Aggregates quotes, fundamentals, and news from one or more registered providers behind a uniform API.

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
      + parses a JSON array as raw object strings
      # serialization
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

openbb
  openbb.new
    @ () -> platform_state
    + creates a platform with no providers registered
    # construction
  openbb.register_provider
    @ (state: platform_state, name: string, base_url: string) -> platform_state
    + registers a named HTTP-backed data provider
    # configuration
  openbb.quote
    @ (state: platform_state, provider: string, symbol: string) -> result[quote, string]
    + returns the latest quote from the named provider
    - returns error when the provider is not registered
    # quotes
    -> std.http.get
    -> std.json.parse_object
  openbb.history
    @ (state: platform_state, provider: string, symbol: string, start: i64, end: i64) -> result[list[ohlcv], string]
    + returns historical bars for the given epoch-second range
    - returns error when the provider is not registered
    # history
    -> std.http.get
    -> std.json.parse_array
  openbb.fundamentals
    @ (state: platform_state, provider: string, symbol: string) -> result[map[string, string], string]
    + returns fundamental metrics as a map
    # fundamentals
    -> std.http.get
    -> std.json.parse_object
  openbb.news
    @ (state: platform_state, provider: string, symbol: string, since: i64) -> result[list[news_item], string]
    + returns news items published since the given epoch second
    # news
    -> std.http.get
    -> std.json.parse_array
    -> std.time.now_seconds
  openbb.cache_get
    @ (state: platform_state, key: string) -> optional[bytes]
    + returns a cached response body or none when absent or expired
    # caching
    -> std.time.now_seconds
  openbb.cache_put
    @ (state: platform_state, key: string, value: bytes, ttl_seconds: i64) -> platform_state
    + stores a response body with a TTL
    # caching
    -> std.time.now_seconds
