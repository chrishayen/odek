# Requirement: "a reference library for HTTP methods, status codes, headers, and TCP/UDP ports"

Lookup-only data library. Each reference category is a simple query against embedded data.

std: (all units exist)

httpref
  httpref.lookup_method
    @ (name: string) -> optional[string]
    + returns the description for a known HTTP method like "GET"
    - returns none for an unknown method name
    # method_reference
  httpref.lookup_status
    @ (code: i32) -> optional[string]
    + returns the name and description for a known status code
    - returns none for codes outside the standard ranges
    # status_reference
  httpref.lookup_header
    @ (name: string) -> optional[string]
    + returns the description for a known header; case-insensitive
    - returns none for unknown headers
    # header_reference
  httpref.lookup_port
    @ (port: i32, protocol: string) -> optional[string]
    + returns the service name for a well-known port and protocol ("tcp" or "udp")
    - returns none for unassigned ports
    # port_reference
