# Requirement: "a library to work with IP networks"

Parses CIDR strings, tests membership, and computes network/broadcast addresses for IPv4.

std: (all units exist)

ipnet
  ipnet.parse_ipv4
    @ (s: string) -> result[u32, string]
    + parses a dotted-quad string into a 32-bit integer
    - returns error when any octet is out of range or format is wrong
    # parsing
  ipnet.format_ipv4
    @ (addr: u32) -> string
    + returns the dotted-quad representation
    # formatting
  ipnet.parse_cidr
    @ (s: string) -> result[ipv4_network, string]
    + parses "a.b.c.d/prefix" into a network
    - returns error when prefix is not in 0..32
    - returns error when host bits are set outside strict mode
    # parsing
  ipnet.network_address
    @ (net: ipv4_network) -> u32
    + returns the network address (host bits cleared)
    # addressing
  ipnet.broadcast_address
    @ (net: ipv4_network) -> u32
    + returns the broadcast address (host bits set)
    # addressing
  ipnet.contains
    @ (net: ipv4_network, addr: u32) -> bool
    + returns true when addr falls within the network
    - returns false for addresses outside the prefix
    # membership
  ipnet.host_count
    @ (net: ipv4_network) -> u64
    + returns the number of addresses in the network
    + returns 1 for /32 and 2^32 for /0
    # addressing
  ipnet.overlaps
    @ (a: ipv4_network, b: ipv4_network) -> bool
    + returns true when two networks share any address
    # comparison
