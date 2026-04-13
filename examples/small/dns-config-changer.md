# Requirement: "a library for changing the system DNS resolver configuration"

Exposes read/write of the resolver list. File IO is delegated to thin std primitives.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file contents
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes atomically to path
      - returns error on permission failure
      # filesystem

dns_config
  dns_config.parse
    @ (raw: string) -> list[string]
    + returns the nameserver entries in declaration order
    + ignores comment and non-nameserver lines
    # parsing
  dns_config.render
    @ (servers: list[string]) -> string
    + returns a resolver file body with one nameserver line per server
    # rendering
  dns_config.load
    @ (path: string) -> result[list[string], string]
    + loads the current resolver list from disk
    # io
    -> std.fs.read_all
  dns_config.save
    @ (path: string, servers: list[string]) -> result[void, string]
    + writes the resolver list to disk
    - returns error when any server is not a valid IP address
    # io
    -> std.fs.write_all
  dns_config.is_valid_address
    @ (value: string) -> bool
    + returns true for well-formed IPv4 or IPv6 addresses
    - returns false for hostnames or malformed strings
    # validation
