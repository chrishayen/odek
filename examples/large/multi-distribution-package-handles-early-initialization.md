# Requirement: "a cloud instance early-initialization library"

Reads instance metadata and a user config, then produces a sequence of configuration actions (set hostname, write files, run commands, add users) that a host-side runner applies.

std
  std.http
    std.http.get
      @ (url: string) -> result[string, string]
      + returns the response body
      - returns error on non-2xx status
      # http
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[yaml_node, string]
      + parses YAML text into a node tree
      - returns error on malformed input
      # serialization
  std.fs
    std.fs.write_file
      @ (path: string, content: bytes, mode: i32) -> result[void, string]
      + writes a file with the given permissions
      # filesystem
  std.process
    std.process.run
      @ (command: string, args: list[string]) -> result[i32, string]
      + runs a command and returns its exit code
      # process

cloudinit
  cloudinit.fetch_metadata
    @ (metadata_url: string) -> result[instance_metadata, string]
    + retrieves instance id, region, and public keys from the metadata service
    - returns error when the service is unreachable
    # metadata
    -> std.http.get
  cloudinit.parse_user_config
    @ (raw: string) -> result[user_config, string]
    + parses a user-data document into a structured config
    - returns error on invalid YAML
    # configuration
    -> std.yaml.parse
  cloudinit.detect_distribution
    @ (os_release: string) -> string
    + returns a normalized distribution identifier
    - returns "unknown" when no match is found
    # detection
  cloudinit.build_plan
    @ (meta: instance_metadata, config: user_config, distribution: string) -> list[init_action]
    + produces an ordered list of idempotent actions
    ? actions include set_hostname, write_file, add_user, run_command
    # planning
  cloudinit.apply_set_hostname
    @ (hostname: string) -> result[void, string]
    + writes and activates the hostname
    # action_runner
    -> std.fs.write_file
    -> std.process.run
  cloudinit.apply_write_file
    @ (path: string, content: bytes, mode: i32) -> result[void, string]
    + creates parent directories and writes the file
    # action_runner
    -> std.fs.write_file
  cloudinit.apply_add_user
    @ (username: string, ssh_keys: list[string]) -> result[void, string]
    + creates the user and installs authorized keys
    - returns error when the user already exists with a different uid
    # action_runner
    -> std.process.run
    -> std.fs.write_file
  cloudinit.apply_run_command
    @ (command: string, args: list[string]) -> result[i32, string]
    + runs a one-shot command and returns exit code
    # action_runner
    -> std.process.run
  cloudinit.run_plan
    @ (plan: list[init_action]) -> list[result[void, string]]
    + applies actions in order, returning one result per action
    ? continues past individual failures so operators see the full report
    # orchestration
