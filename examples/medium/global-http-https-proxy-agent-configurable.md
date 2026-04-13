# Requirement: "a library that resolves HTTP or HTTPS proxy settings from environment variables and chooses the right proxy for a given request URL"

Reads http_proxy, https_proxy, and no_proxy style variables, normalizes them, and answers proxy-for-url queries honoring the no-proxy list.

std
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits a string by separator
      # strings
    std.strings.trim
      @ (s: string) -> string
      + removes leading and trailing whitespace
      # strings
    std.strings.to_lower
      @ (s: string) -> string
      + returns the lowercased form
      # strings
    std.strings.has_suffix
      @ (s: string, suffix: string) -> bool
      + returns true when s ends with suffix
      # strings

proxy_resolver
  proxy_resolver.parse_url
    @ (raw: string) -> result[parsed_url, string]
    + extracts scheme, host, and port from a URL
    - returns error when the scheme is missing or unsupported
    # parsing
    -> std.strings.split
  proxy_resolver.load_from_env
    @ (env: map[string,string]) -> proxy_config
    + reads http_proxy, https_proxy, all_proxy, and no_proxy using case-insensitive names
    + returns an empty config when no variables are set
    # config
    -> std.strings.to_lower
    -> std.strings.split
    -> std.strings.trim
  proxy_resolver.is_bypassed
    @ (cfg: proxy_config, host: string) -> bool
    + returns true when host matches any entry in the no-proxy list
    + treats a leading dot as a domain suffix match
    - returns false otherwise
    # matching
    -> std.strings.to_lower
    -> std.strings.has_suffix
  proxy_resolver.proxy_for_url
    @ (cfg: proxy_config, raw: string) -> result[optional[string], string]
    + returns the proxy URL to use for the request, or none when direct is appropriate
    + prefers https_proxy for https URLs, http_proxy for http URLs, falls back to all_proxy
    - returns error when the request URL cannot be parsed
    # resolution
