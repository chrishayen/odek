# Requirement: "a client library for an uptime monitoring HTTP API"

A thin wrapper around a monitors-and-alerts REST API.

std
  std.http
    std.http.post_form
      fn (url: string, form: map[string, string]) -> result[string, string]
      + POSTs form-encoded fields and returns the response body
      - returns error on non-2xx status
      # http_client
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

uptime
  uptime.new_client
    fn (api_key: string) -> uptime_client
    + holds the api key and base url for subsequent calls
    # construction
  uptime.list_monitors
    fn (client: uptime_client) -> result[list[map[string, string]], string]
    + returns every monitor configured on the account
    - returns error when the api key is rejected
    # query
    -> std.http.post_form
    -> std.json.parse_object
  uptime.create_monitor
    fn (client: uptime_client, name: string, url: string, interval_sec: i32) -> result[string, string]
    + returns the new monitor id
    - returns error when the url is malformed
    # mutation
    -> std.http.post_form
    -> std.json.parse_object
  uptime.delete_monitor
    fn (client: uptime_client, monitor_id: string) -> result[void, string]
    + removes the monitor from the account
    # mutation
    -> std.http.post_form
