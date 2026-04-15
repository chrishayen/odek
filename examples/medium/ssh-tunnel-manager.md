# Requirement: "a library for managing SSH tunnels"

Declarative tunnel configuration plus a small supervisor that opens, tracks, and closes tunnels. The ssh transport is an opaque dependency.

std
  std.net
    std.net.listen_tcp
      fn (host: string, port: i32) -> result[listener, string]
      + binds a TCP listener on host:port
      - returns error when the port is already bound
      # networking

tunnel
  tunnel.config_local
    fn (ssh_host: string, local_port: i32, remote_host: string, remote_port: i32) -> tunnel_config
    + describes a local port forward: listener on local_port targeting remote_host:remote_port through ssh_host
    # configuration
  tunnel.config_remote
    fn (ssh_host: string, remote_port: i32, local_host: string, local_port: i32) -> tunnel_config
    + describes a remote port forward: listener on ssh_host:remote_port targeting local_host:local_port
    # configuration
  tunnel.open
    fn (cfg: tunnel_config, ssh: ssh_session) -> result[tunnel_handle, string]
    + establishes the forward and begins accepting connections
    - returns error when the ssh session is closed
    - returns error when the local listener cannot bind
    # lifecycle
    -> std.net.listen_tcp
  tunnel.close
    fn (handle: tunnel_handle) -> result[void, string]
    + shuts down the listener and cancels in-flight forwards
    # lifecycle
  tunnel.supervisor_new
    fn () -> supervisor_state
    + constructs an empty supervisor
    # supervision
  tunnel.supervisor_add
    fn (state: supervisor_state, handle: tunnel_handle, name: string) -> supervisor_state
    + registers a tunnel with the supervisor under a human name
    # supervision
  tunnel.supervisor_close_all
    fn (state: supervisor_state) -> result[i32, string]
    + closes every tracked tunnel and returns how many were closed
    # supervision
    -> tunnel.close
