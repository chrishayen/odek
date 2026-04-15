# Requirement: "a balanced routing table for fast IP to CIDR longest-prefix lookups"

A prefix trie optimized for longest-prefix match over IPv4 and IPv6 address ranges.

std
  std.net
    std.net.parse_cidr
      fn (text: string) -> result[cidr_range, string]
      + parses a CIDR range such as "192.168.0.0/16" or "2001:db8::/32"
      - returns error on malformed text
      # networking
    std.net.parse_address
      fn (text: string) -> result[bytes, string]
      + parses an IPv4 or IPv6 address into its packed byte form
      - returns error on malformed text
      # networking

routing_table
  routing_table.new
    fn () -> route_state
    + returns an empty routing table
    # construction
  routing_table.insert
    fn (rt: route_state, cidr: string, value: string) -> result[route_state, string]
    + inserts a (prefix, value) mapping
    - returns error on invalid CIDR
    # updates
    -> std.net.parse_cidr
  routing_table.remove
    fn (rt: route_state, cidr: string) -> result[route_state, string]
    + removes the mapping for a prefix
    - returns error when the prefix is not present
    # updates
    -> std.net.parse_cidr
  routing_table.lookup
    fn (rt: route_state, address: string) -> result[optional[string], string]
    + returns the value for the longest prefix covering address
    - returns error on invalid address
    # query
    -> std.net.parse_address
  routing_table.contains
    fn (rt: route_state, cidr: string) -> result[bool, string]
    + returns whether an exact-match prefix exists in the table
    - returns error on invalid CIDR
    # query
    -> std.net.parse_cidr
  routing_table.size
    fn (rt: route_state) -> i32
    + returns the number of prefixes stored
    # inspection
