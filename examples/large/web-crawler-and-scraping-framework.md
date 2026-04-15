# Requirement: "a web crawling and scraping framework"

A crawler loop that feeds seed urls through a fetcher, runs user spider callbacks to extract items and next-page links, obeys robots rules, and schedules politely with per-host rate limits.

std
  std.http
    std.http.get
      fn (url: string, headers: map[string, string]) -> result[http_response, string]
      + issues a GET and returns status, headers, and body bytes
      - returns error on connection or tls failure
      # network
  std.html
    std.html.parse
      fn (source: string) -> html_document
      + returns a tree suitable for css-selector queries
      # parsing
    std.html.select
      fn (doc: html_document, selector: string) -> list[html_node]
      + returns nodes matching a css selector
      # traversal
  std.url
    std.url.resolve
      fn (base: string, ref: string) -> string
      + resolves a relative reference against a base url
      # url
    std.url.parse
      fn (url: string) -> result[url_parts, string]
      + splits a url into scheme, host, port, path, and query
      - returns error on malformed input
      # url
  std.time
    std.time.now_millis
      fn () -> i64
      + returns unix time in milliseconds
      # time

crawler
  crawler.new
    fn (start_urls: list[string], concurrency: i32, delay_ms: i32) -> crawler_state
    + creates a crawler with a frontier seeded from start_urls
    # construction
  crawler.register_spider
    fn (state: crawler_state, host: string, parser: spider_fn) -> crawler_state
    + associates a parser callback with a host
    # spider_registration
  crawler.fetch_robots
    fn (state: crawler_state, host: string) -> crawler_state
    + caches the robots.txt rules for a host
    # compliance
    -> std.http.get
  crawler.is_allowed
    fn (state: crawler_state, url: string) -> bool
    + returns true when the url passes robots rules and has not been visited
    # filtering
    -> std.url.parse
  crawler.step
    fn (state: crawler_state) -> tuple[crawler_state, list[scraped_item]]
    + dequeues one url, fetches it respecting per-host delay, runs the spider, and enqueues discovered links
    + deduplicates urls via a visited set
    - skips urls disallowed by robots
    # crawl_loop
    -> std.http.get
    -> std.html.parse
    -> std.url.resolve
    -> std.time.now_millis
  crawler.run
    fn (state: crawler_state, max_pages: i32) -> list[scraped_item]
    + repeatedly steps until the frontier is empty or max_pages is reached
    # driver
  crawler.extract
    fn (doc: html_document, rules: map[string, string]) -> map[string, string]
    + maps selector rules to text content for each matched node
    # extraction
    -> std.html.select
