# Requirement: "an http proxy for a columnar analytics database"

Authenticates clients, routes queries to configured upstreams, and enforces simple rate limits.

std
  std.http
    std.http.forward
      fn (url: string, body: bytes, headers: map[string,string]) -> result[http_response, string]
      + forwards a request to an upstream and returns the response
      - returns error on connection failure
      # http
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

db_proxy
  db_proxy.new_config
    fn () -> proxy_config
    + returns an empty configuration with no users or upstreams
    # construction
  db_proxy.add_user
    fn (cfg: proxy_config, username: string, password_hash: string, max_qps: i32) -> proxy_config
    + registers a user with credentials and a per-second quota
    # configuration
  db_proxy.add_upstream
    fn (cfg: proxy_config, name: string, url: string) -> proxy_config
    + registers an upstream by logical name
    # configuration
  db_proxy.authenticate
    fn (cfg: proxy_config, username: string, password: string) -> result[user_id, string]
    + returns the user id when credentials match
    - returns error when the user does not exist or the password does not match
    # authentication
  db_proxy.check_quota
    fn (state: quota_state, uid: user_id, limit: i32) -> tuple[bool, quota_state]
    + returns (true, updated_state) when the user is under their per-second limit
    - returns (false, unchanged_state) when the limit would be exceeded
    # rate_limiting
    -> std.time.now_seconds
  db_proxy.handle_query
    fn (cfg: proxy_config, state: quota_state, uid: user_id, upstream: string, query: string) -> result[bytes, string]
    + forwards the query to the named upstream and returns its response body
    - returns error when the upstream name is not configured
    - returns error when the quota check fails
    # dispatch
    -> std.http.forward
