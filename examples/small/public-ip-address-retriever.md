# Requirement: "returns the host's public-facing IPv4 address"

Queries an external echo endpoint and parses the returned address. The HTTP fetch is a thin std primitive so tests can substitute a fake transport.

std
  std.http
    std.http.get
      fn (url: string) -> result[string, string]
      + returns the response body as a string on HTTP 200
      - returns error on non-200 or transport failure
      # http_client

publicip
  publicip.lookup
    fn (endpoint: string) -> result[string, string]
    + returns the trimmed IPv4 string returned by the endpoint
    - returns error when the endpoint response is not a valid dotted-quad
    # ip_lookup
    -> std.http.get
  publicip.parse_ipv4
    fn (raw: string) -> result[string, string]
    + returns the canonical "a.b.c.d" form with leading/trailing whitespace removed
    - returns error when any octet is missing, non-numeric, or out of 0..255
    # parsing
