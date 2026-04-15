# Requirement: "a privacy-preserving front-end library for a social link aggregator"

Fetches posts and comments from an upstream social-link-aggregator API, strips tracking, rewrites media URLs to a configurable proxy, and renders minimal HTML. No user accounts.

std
  std.net
    std.net.http_get
      fn (url: string) -> result[http_response, string]
      + performs an HTTP GET and returns status, headers, and body
      - returns error on connection failure
      # http
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document into a tagged value tree
      - returns error on malformed JSON
      # serialization
  std.html
    std.html.escape
      fn (text: string) -> string
      + returns text with HTML-significant characters escaped
      # html
  std.url
    std.url.parse
      fn (raw: string) -> result[url_parts, string]
      + parses a URL into scheme, host, path, and query
      - returns error when raw is not a valid URL
      # url

private_frontend
  private_frontend.configure
    fn (upstream_base: string, media_proxy_base: string) -> frontend_config
    + builds a configuration with upstream API base and media proxy base
    # configuration
  private_frontend.fetch_listing
    fn (cfg: frontend_config, community: string, sort: string) -> result[list[post], string]
    + retrieves a community listing and returns a normalized list of posts
    - returns error when the community does not exist
    # listing
    -> std.net.http_get
    -> std.json.parse
  private_frontend.fetch_thread
    fn (cfg: frontend_config, thread_id: string) -> result[thread, string]
    + retrieves a thread with its comments as a tree
    - returns error when the thread does not exist
    # thread_fetch
    -> std.net.http_get
    -> std.json.parse
  private_frontend.rewrite_media_url
    fn (cfg: frontend_config, url: string) -> result[string, string]
    + rewrites an upstream media URL to one served by the configured media proxy
    - returns error when the input is not a valid URL
    # privacy
    -> std.url.parse
  private_frontend.render_listing
    fn (cfg: frontend_config, posts: list[post]) -> string
    + renders a listing as minimal HTML with escaped content and proxied media
    # rendering
    -> std.html.escape
  private_frontend.render_thread
    fn (cfg: frontend_config, t: thread) -> string
    + renders a thread with its comment tree as minimal HTML
    # rendering
    -> std.html.escape
  private_frontend.strip_tracking
    fn (url: string) -> string
    + returns the URL with tracking query parameters removed
    # privacy
    -> std.url.parse
