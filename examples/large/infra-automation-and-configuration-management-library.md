# Requirement: "an infrastructure automation and configuration management library"

Describes target state for nodes, evaluates drift, and applies idempotent changes through pluggable executors.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the file contents
      - returns error when the file is missing
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the path
      # filesystem
  std.hash
    std.hash.sha256_hex
      fn (data: bytes) -> string
      + returns the hex-encoded sha256 digest
      # hashing
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

infra
  infra.new_inventory
    fn () -> inventory
    + creates an empty inventory of nodes
    # construction
  infra.add_node
    fn (inv: inventory, id: string, labels: map[string,string]) -> inventory
    + appends a node with the given labels
    # inventory
  infra.select_nodes
    fn (inv: inventory, label_key: string, label_value: string) -> list[string]
    + returns node ids whose label matches
    # inventory
  infra.declare_file_resource
    fn (path: string, content: bytes, mode: i32) -> resource
    + builds a file resource with desired content and permissions
    # resources
  infra.declare_service_resource
    fn (name: string, enabled: bool, running: bool) -> resource
    + builds a service resource with desired enabled and running flags
    # resources
  infra.declare_package_resource
    fn (name: string, version: string, present: bool) -> resource
    + builds a package resource at a desired version
    # resources
  infra.compose_playbook
    fn (resources: list[resource]) -> playbook
    + orders resources preserving declared sequence
    # composition
  infra.diff_resource
    fn (r: resource, observed: observation) -> diff
    + returns the set of changes needed to reach desired state
    + returns an empty diff when already converged
    # planning
    -> std.hash.sha256_hex
  infra.plan
    fn (pb: playbook, observations: map[string, observation]) -> list[diff]
    + returns diffs for all resources in the playbook
    # planning
  infra.apply
    fn (pb: playbook, exec: executor) -> result[apply_report, string]
    + applies each diff via the executor and records outcomes
    - stops and returns error on the first failing resource
    # execution
    -> std.time.now_seconds
  infra.register_executor
    fn (name: string, run_fn: string, read_fn: string) -> executor
    + builds a pluggable executor identified by name
    # executors
  infra.render_template
    fn (template: string, vars: map[string,string]) -> result[string, string]
    + substitutes {{name}} placeholders with values from vars
    - returns error on unclosed placeholders
    # templating
  infra.load_resource_file
    fn (path: string) -> result[list[resource], string]
    + reads a resource file and parses it into a list of resources
    # loading
    -> std.fs.read_all
  infra.write_report
    fn (path: string, report: apply_report) -> result[void, string]
    + writes the apply report to disk
    # reporting
    -> std.fs.write_all
