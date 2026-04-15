# Requirement: "an IT automation library with inventory, playbooks, and task execution over SSH"

Loads an inventory, parses declarative playbooks of tasks, and executes them against remote hosts. Modules are pluggable.

std
  std.yaml
    std.yaml.parse
      fn (raw: string) -> result[yaml_value, string]
      + parses YAML text into a generic yaml_value tree
      - returns error on invalid document
      # parsing
  std.ssh
    std.ssh.connect
      fn (host: string, user: string, key_path: string) -> result[ssh_session, string]
      + establishes an SSH connection using a private key
      - returns error on authentication failure
      # transport
    std.ssh.run
      fn (session: ssh_session, command: string) -> result[command_output, string]
      + runs command and returns exit code, stdout, and stderr
      - returns error on channel failure
      # execution
    std.ssh.put_file
      fn (session: ssh_session, local_bytes: bytes, remote_path: string, mode: i32) -> result[void, string]
      + uploads bytes to remote_path with the given permission bits
      # transfer
  std.template
    std.template.render
      fn (template: string, variables: map[string, string]) -> result[string, string]
      + substitutes {{name}} placeholders from variables
      - returns error on unknown placeholder
      # templating

automation
  automation.load_inventory
    fn (raw_yaml: string) -> result[inventory, string]
    + parses an inventory document into groups and hosts with variables
    - returns error when a host lacks a required address
    # inventory
    -> std.yaml.parse
  automation.parse_playbook
    fn (raw_yaml: string) -> result[playbook, string]
    + parses a playbook into plays, each with a host pattern and an ordered task list
    - returns error when a task lacks a known module name
    # playbook
    -> std.yaml.parse
  automation.register_module
    fn (name: string, handler: fn(map[string, string], ssh_session) -> result[module_result, string]) -> void
    + registers a task module handler by name
    # extensibility
  automation.select_hosts
    fn (inv: inventory, pattern: string) -> list[host]
    + returns hosts matching a group name or glob pattern
    # targeting
  automation.run_task
    fn (task: task, host: host, vars: map[string, string]) -> result[module_result, string]
    + renders task arguments, connects, and dispatches to the registered module
    - returns error when the module name is unknown
    # execution
    -> std.template.render
    -> std.ssh.connect
  automation.run_play
    fn (play: play, inv: inventory) -> result[list[module_result], string]
    + runs every task against every matching host in order, stopping a host on failure
    # orchestration
  automation.run_playbook
    fn (pb: playbook, inv: inventory) -> result[list[module_result], string]
    + runs every play in order and returns the aggregated results
    # orchestration
  automation.module_shell
    fn (args: map[string, string], session: ssh_session) -> result[module_result, string]
    + built-in module that runs a shell command and captures output
    - returns error when the command exits non-zero and ignore_errors is not set
    # builtin_module
    -> std.ssh.run
  automation.module_copy
    fn (args: map[string, string], session: ssh_session) -> result[module_result, string]
    + built-in module that uploads a file with a given mode
    - returns error when the source is missing
    # builtin_module
    -> std.ssh.put_file
