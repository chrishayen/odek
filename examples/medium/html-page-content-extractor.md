# Requirement: "a library for extracting the readable main content from an HTML page"

Given raw HTML, strip boilerplate and return the article's title and body text.

std
  std.html
    std.html.parse
      fn (source: string) -> result[html_node, string]
      + parses HTML into a node tree, tolerating minor malformation
      - returns error when the input is not recognizable HTML
      # html
    std.html.text_content
      fn (node: html_node) -> string
      + returns the concatenated text of node and its descendants
      # html
    std.html.find_all
      fn (node: html_node, tag: string) -> list[html_node]
      + returns every descendant (and self) with the given tag name
      # html
    std.html.get_attribute
      fn (node: html_node, name: string) -> optional[string]
      + returns the named attribute's value, if any
      # html

readable
  readable.extract
    fn (html: string) -> result[article, string]
    + returns an article with title and body text when a readable block is found
    - returns error when the HTML contains no candidate content
    # extraction
    -> std.html.parse
  readable.title
    fn (root: html_node) -> string
    + returns the best title using h1 when present, falling back to the title tag, then an empty string
    # extraction
    -> std.html.find_all
    -> std.html.text_content
  readable.score_candidates
    fn (root: html_node) -> list[scored_node]
    + returns candidate container nodes scored by text density and link ratio
    # scoring
    -> std.html.find_all
    -> std.html.text_content
  readable.pick_best
    fn (candidates: list[scored_node]) -> optional[html_node]
    + returns the highest-scoring candidate, or None when none exceed the minimum score
    # scoring
  readable.clean
    fn (node: html_node) -> string
    + returns the text of node with scripts, styles, navs, and aside elements removed
    # cleanup
    -> std.html.find_all
    -> std.html.text_content
