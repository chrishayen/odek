# Requirement: "a vpn management library for multiple tunnel backends with telemetry, leak detection, and a kill switch"

Supports multiple tunnel backends via a pluggable adapter. Collects telemetry, checks for DNS/IP leaks, and enforces a kill switch that blocks traffic when the tunnel drops.

std
  std.net
    std.net.list_interfaces
      fn () -> result[list[interface_info], string]
      + returns interface metadata
      - returns error when the platform has no interface API
      # networking
    std.net.resolve_hostname
      fn (hostname: string) -> result[list[string], string]
      + returns ip addresses for a hostname
      - returns error on resolution failure
      # networking
    std.net.http_get
      fn (url: string) -> result[bytes, string]
      + performs an HTTP GET
      - returns error on non-2xx
      # http
  std.firewall
    std.firewall.add_block_rule
      fn (interface: string) -> result[rule_handle, string]
      + blocks all outbound traffic on interface
      - returns error when insufficient privileges
      # firewall
    std.firewall.remove_rule
      fn (handle: rule_handle) -> result[void, string]
      + removes a previously added rule
      # firewall

vpn_manager
  vpn_manager.register_backend
    fn (state: manager_state, name: string, adapter: tunnel_adapter) -> manager_state
    + registers a tunnel adapter under a name
    # registration
  vpn_manager.new
    fn () -> manager_state
    + creates an empty manager with no registered backends
    # construction
  vpn_manager.connect
    fn (state: manager_state, backend_name: string, config: map[string,string]) -> result[tunnel_session, string]
    + dispatches to the named backend to bring up a tunnel
    - returns error when backend_name is not registered
    - returns error when the adapter fails to connect
    # connect
  vpn_manager.disconnect
    fn (session: tunnel_session) -> result[void, string]
    + tears down the tunnel
    # disconnect
  vpn_manager.tunnel_status
    fn (session: tunnel_session) -> tunnel_status_info
    + returns connection state, assigned ip, peer, and uptime seconds
    # status
  vpn_manager.read_telemetry
    fn (session: tunnel_session) -> telemetry_sample
    + returns current rx/tx bytes and packet counts
    # telemetry
  vpn_manager.detect_dns_leak
    fn (expected_dns: list[string]) -> result[bool, string]
    + true if the system would resolve through a non-tunnel DNS server
    - returns error when resolution fails entirely
    # leak_detection
    -> std.net.resolve_hostname
  vpn_manager.detect_ip_leak
    fn (session: tunnel_session, check_url: string) -> result[bool, string]
    + true if an external check returns an IP outside the tunnel's assigned address
    - returns error on network failure
    # leak_detection
    -> std.net.http_get
  vpn_manager.enable_kill_switch
    fn (state: manager_state, tunnel_interface: string) -> result[kill_switch_handle, string]
    + installs firewall rules blocking non-tunnel traffic
    - returns error when privileges are insufficient
    # kill_switch
    -> std.firewall.add_block_rule
  vpn_manager.disable_kill_switch
    fn (handle: kill_switch_handle) -> result[void, string]
    + removes kill-switch firewall rules
    # kill_switch
    -> std.firewall.remove_rule
  vpn_manager.list_tunnel_interfaces
    fn () -> result[list[string], string]
    + returns names of interfaces created by registered tunnel backends
    # inspection
    -> std.net.list_interfaces
