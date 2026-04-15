# Requirement: "a client for a microblogging service API"

Typed wrappers over status, timeline, and search endpoints of a microblogging service.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + issues an HTTP request and returns status and body
      - returns error on network failure
      # http
  std.crypto
    std.crypto.hmac_sha1
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA1 of data under key
      # cryptography
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.parse_array_of_objects
      fn (raw: string) -> result[list[map[string, string]], string]
      + parses a JSON array of objects into a list of maps
      - returns error when the root is not an array
      # serialization

microblog_client
  microblog_client.new
    fn (consumer_key: string, consumer_secret: string, access_token: string, access_secret: string) -> microblog_client_state
    + stores credentials used to sign every request
    # construction
  microblog_client.sign_request
    fn (state: microblog_client_state, method: string, url: string, params: map[string, string]) -> string
    + returns an OAuth1 authorization header value for the request
    # auth
    -> std.crypto.hmac_sha1
  microblog_client.post_status
    fn (state: microblog_client_state, text: string) -> result[map[string, string], string]
    + returns the created status object
    - returns error when text exceeds the length limit
    # posting
    -> std.http.request
    -> std.json.parse_object
  microblog_client.home_timeline
    fn (state: microblog_client_state, count: i32) -> result[list[map[string, string]], string]
    + returns the most recent count items from the authenticated user's home timeline
    # reading
    -> std.http.request
    -> std.json.parse_array_of_objects
  microblog_client.search
    fn (state: microblog_client_state, query: string) -> result[list[map[string, string]], string]
    + returns statuses matching the query string
    - returns error when the query is empty
    # search
    -> std.http.request
    -> std.json.parse_array_of_objects
