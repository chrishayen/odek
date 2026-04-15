# Requirement: "a library for extracting readable text and metadata from web pages"

Fetches a URL, parses the HTML into a tree, removes boilerplate (navigation, ads, scripts), identifies the main content region, and returns its plain text plus structured metadata.

std
  std.http
    std.http.get
      fn (url: string) -> result[http_response, string]
      + performs an HTTP GET and returns status, headers, and body
      - returns error on transport failure
      # http_client
  std.html
    std.html.parse
      fn (source: string) -> result[html_node, string]
      + returns the document root of a parsed HTML tree
      - returns error on a fundamentally broken document
      # parsing

content_extractor
  content_extractor.clean_tree
    fn (root: html_node) -> html_node
    + returns a copy with script, style, nav, header, footer, and aside elements removed
    # cleaning
  content_extractor.score_blocks
    fn (root: html_node) -> map[string, f64]
    + assigns a content score to each candidate block based on text density and link ratio
    ? nodes with many links relative to text are penalized
    # scoring
  content_extractor.pick_main_block
    fn (root: html_node, scores: map[string, f64]) -> optional[html_node]
    + returns the highest-scoring block
    - returns none when no block scores above the threshold
    # selection
    -> content_extractor.score_blocks
  content_extractor.extract_text
    fn (node: html_node) -> string
    + returns the concatenated text of the node with paragraphs separated by blank lines
    # extraction
  content_extractor.extract_title
    fn (root: html_node) -> optional[string]
    + returns the document title from <title>, falling back to the first h1
    - returns none when neither exists
    # metadata
  content_extractor.extract_meta
    fn (root: html_node) -> map[string, string]
    + returns name/content pairs from <meta> tags (description, author, og:*, article:*)
    # metadata
  content_extractor.extract_language
    fn (root: html_node) -> optional[string]
    + returns the value of the html lang attribute
    - returns none when absent
    # metadata
  content_extractor.extract
    fn (html_source: string) -> result[extracted, string]
    + parses, cleans, picks the main block, and returns text together with title, meta, and language
    - returns error when parsing fails
    # orchestration
    -> std.html.parse
    -> content_extractor.clean_tree
    -> content_extractor.pick_main_block
    -> content_extractor.extract_text
    -> content_extractor.extract_title
    -> content_extractor.extract_meta
    -> content_extractor.extract_language
  content_extractor.extract_from_url
    fn (url: string) -> result[extracted, string]
    + fetches the URL then runs extraction on the response body
    - returns error on transport failure or parse failure
    # orchestration
    -> std.http.get
    -> content_extractor.extract
