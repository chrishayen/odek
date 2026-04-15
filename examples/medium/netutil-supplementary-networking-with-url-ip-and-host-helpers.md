# Requirement: "a supplementary networking library with helpers for URLs, IP addresses, and host lookups"

The input was a metapackage description. The library exposes the handful of utilities such a package actually contains.

std
  std.net
    std.net.resolve_host
      fn (host: string) -> result[list[string], string]
      + resolves a hostname to a list of IP address strings
      - returns error when no addresses are found
      # network
  std.text
    std.text.split_on
      fn (raw: string, sep: string) -> list[string]
      + splits a string on an exact separator
      # text

netutil
  netutil.parse_url
    fn (raw: string) -> result[parsed_url, string]
    + returns scheme, host, port, path, and query string
    - returns error when the scheme is missing
    # url
  netutil.join_url
    fn (base: string, ref: string) -> result[string, string]
    + resolves a relative reference against a base URL
    - returns error when the base is not absolute
    # url
  netutil.parse_ip
    fn (raw: string) -> result[ip_addr, string]
    + parses an IPv4 or IPv6 literal
    - returns error on malformed input
    # ip
  netutil.ip_in_cidr
    fn (addr: ip_addr, cidr: string) -> result[bool, string]
    + returns whether the address belongs to the CIDR block
    - returns error on malformed CIDR
    # ip
    -> std.text.split_on
  netutil.lookup_host
    fn (host: string) -> result[list[ip_addr], string]
    + resolves a hostname to parsed IP addresses
    - returns error when resolution fails
    # dns
    -> std.net.resolve_host
