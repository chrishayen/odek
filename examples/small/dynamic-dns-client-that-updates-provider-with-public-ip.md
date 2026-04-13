# Requirement: "a dynamic DNS client that updates a DNS provider with the current public IP"

Detects the host's current public IP address and, when it has changed since the last update, rewrites a DNS record through a provider adapter.

std
  std.http
    std.http.get
      @ (url: string, headers: map[string, string]) -> result[http_response, string]
      + performs a GET and returns the status, headers, and body
      - returns error on network failure or malformed response
      # network
    std.http.put_json
      @ (url: string, headers: map[string, string], body: string) -> result[http_response, string]
      + performs a PUT with a JSON body
      - returns error on network failure
      # network

ddns_client
  ddns_client.detect_public_ip
    @ (resolver_url: string) -> result[string, string]
    + returns the current public IPv4 or IPv6 address as reported by the resolver endpoint
    - returns error when the response cannot be parsed as an IP
    # detection
    -> std.http.get
  ddns_client.has_changed
    @ (current: string, last_known: optional[string]) -> bool
    + returns true when last_known is none or differs from current
    # diff
  ddns_client.update_record
    @ (provider: provider_config, hostname: string, ip: string) -> result[void, string]
    + rewrites the A or AAAA record for hostname to ip using the provider's API
    - returns error when the provider rejects the update
    ? the record type is chosen from the address family of ip
    # update
    -> std.http.put_json
  ddns_client.tick
    @ (state: ddns_state) -> result[ddns_state, string]
    + detects the current IP and, if changed, updates the configured record
    + returns the new state with the IP cached
    # loop_step
