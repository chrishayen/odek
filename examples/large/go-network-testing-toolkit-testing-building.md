# Requirement: "a toolkit for building and testing virtual network topologies"

Declarative node and link definitions get realized into isolated network namespaces, with helpers for running commands, capturing traffic, and tearing down cleanly.

std
  std.os
    std.os.run
      @ (argv: list[string]) -> result[process_result, string]
      + runs a command and returns exit code, stdout, and stderr
      - returns error when the executable cannot be launched
      # process
  std.netns
    std.netns.create
      @ (name: string) -> result[netns_handle, string]
      + creates a new network namespace with the given name
      - returns error when a namespace with that name already exists
      # namespace
    std.netns.destroy
      @ (handle: netns_handle) -> result[bool, string]
      + removes the namespace and releases its resources
      # namespace
    std.netns.exec
      @ (handle: netns_handle, argv: list[string]) -> result[process_result, string]
      + runs a command inside the namespace
      # namespace
  std.net
    std.net.create_veth
      @ (a_ns: netns_handle, b_ns: netns_handle) -> result[veth_pair, string]
      + creates a virtual ethernet pair placing one end in each namespace
      # networking
    std.net.assign_address
      @ (handle: netns_handle, iface: string, cidr: string) -> result[bool, string]
      + assigns an IP address in CIDR notation to an interface in the namespace
      - returns error when the CIDR is malformed
      # networking

nettest
  nettest.new_topology
    @ () -> topology
    + creates an empty topology with no nodes or links
    # construction
  nettest.add_node
    @ (topo: topology, name: string) -> topology
    + registers a node by name
    # topology
  nettest.add_link
    @ (topo: topology, a: string, b: string, a_cidr: string, b_cidr: string) -> topology
    + registers a point-to-point link with addresses at each end
    # topology
  nettest.realize
    @ (topo: topology) -> result[running_topology, string]
    + creates namespaces, veth pairs, and addresses for every declared node and link
    - returns error and rolls back on partial failure
    # lifecycle
    -> std.netns.create
    -> std.net.create_veth
    -> std.net.assign_address
  nettest.run_in_node
    @ (running: running_topology, node: string, argv: list[string]) -> result[process_result, string]
    + executes a command inside the named node
    - returns error when the node is unknown
    # execution
    -> std.netns.exec
  nettest.ping
    @ (running: running_topology, from_node: string, to_addr: string) -> result[bool, string]
    + returns true when an ICMP echo from from_node reaches to_addr
    # assertion
    -> std.netns.exec
  nettest.capture
    @ (running: running_topology, node: string, iface: string, duration_secs: i32) -> result[bytes, string]
    + runs a packet capture on the interface for duration_secs and returns pcap bytes
    # capture
    -> std.netns.exec
  nettest.teardown
    @ (running: running_topology) -> result[bool, string]
    + destroys every namespace and link created by realize
    # lifecycle
    -> std.netns.destroy
