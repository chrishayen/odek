# Requirement: "an infrastructure automation library that runs declarative operations over remote hosts"

Hosts are inventory entries; an operation is a declarative step (package install, file write, service restart); a plan executes operations against hosts via a pluggable connector that shells out commands.

std
  std.text
    std.text.split_lines
      fn (raw: string) -> list[string]
      + splits on newlines and drops a trailing empty line
      # text
    std.text.join
      fn (parts: list[string], separator: string) -> string
      + joins parts with the separator
      # text
  std.io
    std.io.read_all
      fn (path: string) -> result[string, string]
      + reads an entire file as UTF-8 text
      - returns error when the file does not exist
      # io
  std.hash
    std.hash.sha256_hex
      fn (data: bytes) -> string
      + returns the hex-encoded SHA-256 digest of data
      # hashing

automate
  automate.host
    fn (name: string, address: string, user: string) -> host
    + returns a host record for the given connection details
    # inventory
  automate.load_inventory
    fn (path: string) -> result[list[host], string]
    + parses a newline-delimited "name address user" file into host records
    - returns error when the file cannot be read
    - returns error when a line has fewer than three fields
    # inventory
    -> std.io.read_all
    -> std.text.split_lines
  automate.op_install_package
    fn (name: string) -> operation
    + returns an operation that installs the named package when absent
    # operations
  automate.op_write_file
    fn (path: string, content: string, mode: i32) -> operation
    + returns an operation that writes content to path with the given permission bits
    + is idempotent: skipped when target already contains the content
    # operations
    -> std.hash.sha256_hex
  automate.op_restart_service
    fn (name: string) -> operation
    + returns an operation that restarts the named service
    # operations
  automate.op_run_shell
    fn (command: string) -> operation
    + returns an operation that runs a raw shell command
    # operations
  automate.plan
    fn (hosts: list[host], operations: list[operation]) -> plan
    + returns a plan that will apply every operation to every host in order
    # planning
  automate.execute
    fn (plan: plan, connector: host_connector) -> result[execution_report, string]
    + runs the plan via the connector and returns a per-host, per-operation report
    + reports changed vs unchanged operations
    - returns error when any operation on any host fails and no continue-on-error flag is set
    # execution
  automate.new_ssh_connector
    fn (ssh_key_path: string) -> host_connector
    + returns a connector that executes shell commands over SSH using the given key
    # connectors
  automate.new_local_connector
    fn () -> host_connector
    + returns a connector that executes shell commands on the local host
    # connectors
  automate.format_report
    fn (report: execution_report) -> string
    + returns a human-readable summary of an execution report
    # reporting
    -> std.text.join
