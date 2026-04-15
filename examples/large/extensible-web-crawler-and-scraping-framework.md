# Requirement: "an extensible web crawling and scraping framework"

A crawler holds a frontier queue, a per-host politeness clock, and a pipeline of user-supplied extractors.

std
  std.http
    std.http.get
      fn (url: string, headers: map[string,string]) -> result[http_response, string]
      + fetches a URL and returns status, headers, and body
      - returns error on connection failure or invalid URL
      # network
  std.url
    std.url.parse
      fn (raw: string) -> result[parsed_url, string]
      + parses a URL into scheme, host, path, and query components
      - returns error on malformed URL
      # url
    std.url.resolve
      fn (base: parsed_url, ref: string) -> result[parsed_url, string]
      + resolves a relative reference against a base URL
      - returns error when the reference cannot be resolved
      # url
  std.html
    std.html.parse
      fn (body: bytes) -> result[html_doc, string]
      + parses HTML bytes into a DOM tree
      - returns error on severely malformed input
      # parsing
    std.html.select
      fn (doc: html_doc, css_selector: string) -> list[html_node]
      + returns nodes matching the CSS selector
      # dom
  std.time
    std.time.now_millis
      fn () -> i64
      + returns unix time in milliseconds
      # time
    std.time.sleep_millis
      fn (ms: i64) -> void
      + blocks the caller for the given duration
      # time

crawler
  crawler.new
    fn (user_agent: string, politeness_ms: i64) -> crawler_state
    + creates a crawler with an empty frontier and the given per-host delay
    # construction
  crawler.enqueue
    fn (state: crawler_state, url: string) -> result[crawler_state, string]
    + adds a URL to the frontier if not already seen
    - returns error when the URL is malformed
    # frontier
    -> std.url.parse
  crawler.fetch_next
    fn (state: crawler_state) -> result[tuple[crawler_state, fetched_page], string]
    + dequeues a URL, waits out the host's politeness window, and fetches it
    - returns error when the frontier is empty
    - returns error when fetch fails
    # fetching
    -> std.time.now_millis
    -> std.time.sleep_millis
    -> std.http.get
  crawler.extract_links
    fn (page: fetched_page) -> list[string]
    + returns absolute URLs for every anchor in the page
    + resolves relative hrefs against the page URL
    # link_extraction
    -> std.html.parse
    -> std.html.select
    -> std.url.resolve
  crawler.register_extractor
    fn (state: crawler_state, name: string, selector: string) -> crawler_state
    + registers a named CSS selector whose matches will be captured
    # pipeline
  crawler.run_extractors
    fn (state: crawler_state, page: fetched_page) -> map[string, list[string]]
    + runs every registered extractor over the page's DOM and returns their text results
    # extraction
    -> std.html.parse
    -> std.html.select
  crawler.mark_visited
    fn (state: crawler_state, url: string) -> crawler_state
    + marks a URL as visited so enqueue will skip it
    # deduplication
  crawler.set_robots
    fn (state: crawler_state, host: string, disallowed: list[string]) -> crawler_state
    + records disallowed path prefixes for a host
    # robots
  crawler.is_allowed
    fn (state: crawler_state, url: string) -> bool
    + returns true when the URL's path is not blocked by the host's disallow list
    - returns false when the URL matches a disallowed prefix
    # robots
    -> std.url.parse
