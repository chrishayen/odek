# Requirement: "a web scraping framework that turns websites into structured data"

Fetches pages, extracts fields via CSS-like selectors, follows pagination, and emits structured records.

std
  std.http
    std.http.get
      @ (url: string, headers: map[string, string]) -> result[bytes, string]
      + returns body bytes on HTTP 2xx
      - returns error on network failure or non-2xx status
      # http
  std.html
    std.html.parse
      @ (source: bytes) -> result[html_doc, string]
      + parses HTML5 into a navigable document tree
      - returns error on malformed input that cannot be recovered
      # html_parsing
    std.html.select
      @ (doc: html_doc, selector: string) -> list[html_node]
      + returns nodes matching a CSS-like selector
      + returns empty list when nothing matches
      # html_query
    std.html.text
      @ (node: html_node) -> string
      + returns concatenated text content of a node
      # html_query
    std.html.attr
      @ (node: html_node, name: string) -> optional[string]
      + returns the attribute value when present
      # html_query
  std.url
    std.url.resolve
      @ (base: string, reference: string) -> result[string, string]
      + resolves a relative URL against a base URL
      - returns error when base is not absolute
      # url

scraper
  scraper.new_field
    @ (name: string, selector: string, attr: optional[string]) -> field_def
    + creates a field extractor bound to a selector and optional attribute
    ? when attr is absent the field captures text content
    # configuration
  scraper.new_schema
    @ (fields: list[field_def], follow: optional[string]) -> schema_def
    + creates a schema with field extractors and an optional follow selector
    # configuration
  scraper.extract_record
    @ (doc: html_doc, schema: schema_def) -> map[string, string]
    + returns a record mapping field names to extracted values
    + missing fields map to empty strings
    # extraction
    -> std.html.select
    -> std.html.text
    -> std.html.attr
  scraper.scrape_page
    @ (url: string, schema: schema_def) -> result[list[map[string, string]], string]
    + fetches a URL and returns all records matching the schema
    - returns error on fetch or parse failure
    # scraping
    -> std.http.get
    -> std.html.parse
  scraper.follow_links
    @ (doc: html_doc, base_url: string, selector: string) -> list[string]
    + returns absolute URLs from anchor hrefs matching the selector
    # link_following
    -> std.html.select
    -> std.html.attr
    -> std.url.resolve
  scraper.crawl
    @ (start_url: string, schema: schema_def, max_pages: i32) -> result[list[map[string, string]], string]
    + fetches pages following the schema's follow selector up to max_pages
    + deduplicates visited URLs
    - returns error when start_url cannot be fetched
    # crawling
    -> std.http.get
    -> std.html.parse
