# Requirement: "a library for creating HTML snippets and embeds from URLs using oEmbed, Open Graph, and meta tags"

Fetches a URL, extracts embed info from oEmbed endpoints, Open Graph tags, and generic meta tags, then renders an HTML snippet.

std
  std.http
    std.http.get
      @ (url: string) -> result[bytes, string]
      + returns the response body on 2xx
      - returns error on network failure or non-2xx status
      # http
    std.http.get_json
      @ (url: string) -> result[map[string, string], string]
      + fetches a URL and parses the body as a JSON object
      - returns error on invalid JSON
      # http
  std.html
    std.html.parse
      @ (markup: string) -> result[html_doc, string]
      + parses HTML into a traversable document
      - returns error on input that cannot be recovered
      # parsing
    std.html.find_meta_tags
      @ (doc: html_doc) -> map[string, string]
      + returns a map of name/property to content values
      # parsing
    std.html.find_link_rel
      @ (doc: html_doc, rel: string) -> optional[string]
      + returns the href of the first link element with the given rel
      # parsing
    std.html.escape
      @ (raw: string) -> string
      + escapes &, <, >, ", ' for safe HTML output
      # sanitization
  std.url
    std.url.parse
      @ (raw: string) -> result[url_parts, string]
      + splits a URL into scheme, host, path, query
      - returns error on malformed input
      # parsing
    std.url.is_absolute
      @ (raw: string) -> bool
      + true when the URL has a scheme and host
      # parsing

embed
  embed.fetch_info
    @ (url: string) -> result[embed_info, string]
    + returns structured info with title, description, thumbnail, provider
    + prefers oEmbed data when an oEmbed endpoint is discoverable
    + falls back to Open Graph tags when oEmbed is unavailable
    + falls back to generic meta tags when Open Graph is absent
    - returns error when the URL is not absolute
    - returns error when the page cannot be fetched
    # discovery
    -> std.url.is_absolute
    -> std.http.get
    -> std.html.parse
  embed.discover_oembed_endpoint
    @ (doc: html_doc) -> optional[string]
    + returns the href of a link[rel=alternate][type=application/json+oembed]
    # discovery
    -> std.html.find_link_rel
  embed.extract_open_graph
    @ (doc: html_doc) -> map[string, string]
    + collects meta[property^=og:] entries into a flat map
    # extraction
    -> std.html.find_meta_tags
  embed.extract_generic_meta
    @ (doc: html_doc) -> map[string, string]
    + collects title, description, and canonical image from meta tags
    # extraction
    -> std.html.find_meta_tags
  embed.fetch_oembed
    @ (endpoint: string) -> result[map[string, string], string]
    + fetches an oEmbed endpoint and returns the parsed JSON object
    - returns error when the endpoint returns malformed JSON
    # http
    -> std.http.get_json
  embed.render_snippet
    @ (info: embed_info) -> string
    + renders a self-contained HTML snippet from embed info
    + escapes all user-supplied strings
    ? output is a figure with title, thumbnail, and description
    # rendering
    -> std.html.escape
