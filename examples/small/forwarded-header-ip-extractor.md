# Requirement: "a library for parsing and resolving forwarded-address headers to extract the real client ip"

std: (all units exist)

forwarded
  forwarded.parse_for
    fn (header: string) -> list[string]
    + returns the comma-separated addresses from an x-forwarded-for header, trimmed
    + returns an empty list for an empty header
    # parsing
  forwarded.client_ip
    fn (header: string, trusted_proxies: list[string]) -> optional[string]
    + walks the parsed list from right to left, skipping trusted proxies, and returns the first remaining address
    - returns none when every address is trusted
    # resolution
  forwarded.parse_forwarded
    fn (header: string) -> list[map[string, string]]
    + parses an rfc 7239 forwarded header into a list of key-value maps
    + returns an empty list for an empty header
    # parsing
