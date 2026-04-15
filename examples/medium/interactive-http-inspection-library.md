# Requirement: "an interactive http inspection library"

The library builds, sends, and renders HTTP requests and responses. Interactivity (keystrokes, screen) is the caller's concern; this exposes the model and the transport.

std
  std.http
    std.http.send
      fn (method: string, url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
      + performs the request and returns status, headers, and body
      - returns error on network failure or malformed url
      # http
  std.url
    std.url.parse
      fn (raw: string) -> result[url_parts, string]
      + splits into scheme, host, path, query
      - returns error on missing scheme or host
      # url
  std.encoding
    std.encoding.form_encode
      fn (params: map[string,string]) -> string
      + encodes a flat map as application/x-www-form-urlencoded
      # encoding

http_inspector
  http_inspector.new_request
    fn (method: string, url: string) -> http_request
    + creates a blank request with the given method and url
    # construction
  http_inspector.set_header
    fn (req: http_request, name: string, value: string) -> http_request
    + returns a new request with the header set, replacing any prior value
    # request_building
  http_inspector.set_body
    fn (req: http_request, body: bytes, content_type: string) -> http_request
    + sets the raw body and content-type header together
    # request_building
  http_inspector.execute
    fn (req: http_request) -> result[http_response, string]
    + validates the request and sends it
    - returns error when method or url is empty
    # execution
    -> std.url.parse
    -> std.http.send
  http_inspector.format_response
    fn (resp: http_response) -> string
    + returns a human-readable status line, headers block, and body preview
    ? body is truncated at 4KB for display
    # rendering
  http_inspector.format_request
    fn (req: http_request) -> string
    + returns a human-readable method, url, headers, and body preview
    # rendering
