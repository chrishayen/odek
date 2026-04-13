# Requirement: "a library for managing packet-forwarding rules on a host firewall"

Add, list, and remove forwarding rules. The underlying firewall is a thin std seam.

std
  std.firewall
    std.firewall.add_forward
      @ (proto: string, src_port: i32, dst_host: string, dst_port: i32) -> result[string, string]
      + installs a forward and returns its rule id
      - returns error when the src_port is already forwarded
      # firewall
    std.firewall.remove_forward
      @ (rule_id: string) -> result[void, string]
      + removes a previously installed rule
      - returns error when the rule id is unknown
      # firewall
    std.firewall.list_forwards
      @ () -> result[list[firewall_rule], string]
      + returns all currently installed forwarding rules
      # firewall
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the file does not exist
      # io
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path
      # io

forwards
  forwards.new_manager
    @ (state_path: string) -> forwards_manager
    + creates a manager persisting rule metadata at state_path
    # construction
  forwards.add
    @ (mgr: forwards_manager, proto: string, src_port: i32, dst_host: string, dst_port: i32) -> result[forward_record, string]
    + installs a forward and records it in state
    - returns error when proto is not "tcp" or "udp"
    - returns error when src_port is out of [1, 65535]
    # management
    -> std.firewall.add_forward
    -> std.fs.read_all
    -> std.fs.write_all
  forwards.remove
    @ (mgr: forwards_manager, src_port: i32) -> result[void, string]
    + removes the forward for src_port
    - returns error when no forward exists for src_port
    # management
    -> std.firewall.remove_forward
    -> std.fs.read_all
    -> std.fs.write_all
  forwards.list
    @ (mgr: forwards_manager) -> result[list[forward_record], string]
    + returns all forwards known to the manager
    # management
    -> std.firewall.list_forwards
  forwards.sync
    @ (mgr: forwards_manager) -> result[void, string]
    + reinstalls any rules in state that are missing from the firewall
    # reconciliation
    -> std.firewall.list_forwards
    -> std.firewall.add_forward
