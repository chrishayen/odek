# Requirement: "a library for downloading structured regulatory filings and parsing xbrl financial statements"

Fetching is done over a pluggable http primitive; parsing turns a filing document into typed facts.

std
  std.net
    std.net.http_get
      @ (url: string) -> result[bytes, string]
      + fetches a URL and returns the response body
      - returns error on non-2xx status
      # http
  std.xml
    std.xml.parse
      @ (raw: bytes) -> result[xml_node, string]
      + parses XML into a tree node
      - returns error on malformed XML
      # parsing
    std.xml.find_all
      @ (root: xml_node, tag: string) -> list[xml_node]
      + returns every descendant element matching the tag name
      # parsing
    std.xml.attr
      @ (node: xml_node, name: string) -> optional[string]
      + returns the named attribute value when present
      # parsing
    std.xml.text
      @ (node: xml_node) -> string
      + returns the concatenated text content of the node
      # parsing

edgar
  edgar.fetch_filing
    @ (accession: string) -> result[bytes, string]
    + downloads the primary document for an accession number
    - returns error when the transport fails
    # retrieval
    -> std.net.http_get
  edgar.fetch_filing_index
    @ (cik: string) -> result[list[string], string]
    + returns the list of accession numbers for the given filer id
    # retrieval
    -> std.net.http_get
  edgar.parse_xbrl_facts
    @ (document: bytes) -> result[list[xbrl_fact], string]
    + returns facts with (concept, unit, period, value) extracted from the document
    - returns error on malformed input
    # xbrl_parsing
    -> std.xml.parse
    -> std.xml.find_all
    -> std.xml.attr
    -> std.xml.text
  edgar.select_fact
    @ (facts: list[xbrl_fact], concept: string) -> optional[xbrl_fact]
    + returns the first fact matching the concept tag
    # xbrl_query
