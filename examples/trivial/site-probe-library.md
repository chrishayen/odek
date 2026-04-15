# Requirement: "a library that checks whether a website is up"

A single probe that fetches a URL and reports reachability based on HTTP status.

std
  std.http
    std.http.get_status
      fn (url: string, timeout_millis: i32) -> result[i32, string]
      + returns the HTTP status code of a GET request
      - returns error on network failure or timeout
      # http_client

site_probe
  site_probe.is_up
    fn (url: string) -> bool
    + returns true when the URL responds with a status below 400 within a default timeout
    - returns false on connection error, timeout, or a 4xx/5xx status
    ? the default timeout is hardcoded
    # reachability
    -> std.http.get_status
