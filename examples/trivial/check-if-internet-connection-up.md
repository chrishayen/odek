# Requirement: "a library to check whether the internet connection is available"

The check delegates to a thin std primitive that opens a TCP connection to a well-known host.

std
  std.net
    std.net.tcp_reachable
      @ (host: string, port: i32, timeout_millis: i32) -> bool
      + returns true when a TCP connection to host:port completes within the timeout
      - returns false on timeout, refused, or DNS failure
      # networking

connectivity
  connectivity.is_online
    @ () -> bool
    + returns true when any of a small set of reliable public hosts is reachable
    - returns false when none of the probe targets can be reached
    ? the probe list is hardcoded and uses standard TCP ports
    # connectivity
    -> std.net.tcp_reachable
