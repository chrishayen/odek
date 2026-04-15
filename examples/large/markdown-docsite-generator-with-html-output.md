# Requirement: "a documentation site generator that serves markdown without pre-built HTML"

Reads markdown files from a directory, serves them dynamically, and renders them to HTML on demand. Navigation is derived from the directory tree.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file as UTF-8
      - returns error when the file is missing
      # filesystem
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns every file path under root in deterministic order
      - returns error when root does not exist
      # filesystem
  std.markdown
    std.markdown.render_html
      fn (source: string) -> string
      + renders CommonMark source to HTML
      # rendering
  std.http
    std.http.serve
      fn (addr: string, handler: fn(http_request) -> http_response) -> result[void, string]
      + starts an HTTP server that dispatches every request to handler
      - returns error when the address cannot be bound
      # networking

docsite
  docsite.load
    fn (root: string) -> result[docsite_state, string]
    + scans root, collects every .md file, and builds a route-to-path index
    - returns error when root cannot be walked
    # loading
    -> std.fs.walk
  docsite.route_to_path
    fn (state: docsite_state, route: string) -> optional[string]
    + maps a URL route to an on-disk markdown file
    + the empty route resolves to index.md when present
    # routing
  docsite.render_page
    fn (state: docsite_state, route: string) -> result[string, string]
    + returns the HTML for the requested route, including sidebar navigation
    - returns error when the route is unknown
    # rendering
    -> std.fs.read_all
    -> std.markdown.render_html
  docsite.build_sidebar
    fn (state: docsite_state) -> string
    + returns the HTML sidebar listing every page grouped by directory
    # navigation
  docsite.handle
    fn (state: docsite_state, req: http_request) -> http_response
    + returns a 200 rendered page for known routes and a 404 otherwise
    - returns a 404 response for unknown routes
    # request_handling
  docsite.serve
    fn (state: docsite_state, addr: string) -> result[void, string]
    + starts an HTTP server using handle as the request handler
    - returns error when the address cannot be bound
    # serving
    -> std.http.serve
