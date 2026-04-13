# Requirement: "a higher-level http client with retries, timeouts, and query building"

A small facade around a primitive transport: builds requests, retries on transient failures, decodes responses.

std
  std.http
    std.http.send
      @ (req: http_request) -> result[http_response, string]
      + sends a request over the network and returns the response
      - returns error on transport failure
      # transport
  std.time
    std.time.sleep_millis
      @ (ms: i64) -> void
      + blocks the current flow for the given milliseconds
      # time
  std.encoding
    std.encoding.url_encode
      @ (s: string) -> string
      + percent-encodes a string for safe use in URLs
      # encoding

httpx
  httpx.new_client
    @ () -> client_state
    + creates a client with default retry and timeout settings
    # construction
  httpx.with_timeout
    @ (c: client_state, ms: i64) -> client_state
    + sets the per-request timeout
    # configuration
  httpx.with_retries
    @ (c: client_state, max: i32, base_delay_ms: i64) -> client_state
    + sets the retry count and exponential backoff base delay
    # configuration
  httpx.build_url
    @ (base: string, query: map[string, string]) -> string
    + appends encoded query parameters to a base URL
    + returns base unchanged when query is empty
    # url_building
    -> std.encoding.url_encode
  httpx.get
    @ (c: client_state, url: string) -> result[http_response, string]
    + sends a GET request and retries on transient failures with backoff
    - returns error after max retries on persistent failure
    # requests
    -> std.http.send
    -> std.time.sleep_millis
  httpx.post_json
    @ (c: client_state, url: string, body: string) -> result[http_response, string]
    + sends a POST request with JSON content type and retries
    - returns error after max retries on persistent failure
    # requests
    -> std.http.send
    -> std.time.sleep_millis
  httpx.is_retryable_status
    @ (status: i32) -> bool
    + returns true for 5xx and 429 responses
    # retry_policy
