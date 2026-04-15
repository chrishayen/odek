# Requirement: "an OSINT data-collection framework with a module package manager"

Modules are downloadable scripts that read a seed database, emit findings, and write them back. The framework exposes module management and a run loop.

std
  std.http
    std.http.get
      fn (url: string) -> result[bytes, string]
      + fetches a URL and returns the response body
      - returns error on non-2xx or network failure
      # http
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the path is missing
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating parents as needed
      # filesystem
  std.crypto
    std.crypto.sha256_hex
      fn (data: bytes) -> string
      + returns the lowercase hex SHA-256 digest
      # cryptography
  std.script
    std.script.eval
      fn (source: string, inputs: map[string,string]) -> result[map[string,string], string]
      + executes a scripted module with named inputs and returns its outputs
      - returns error on script failure
      # scripting

osint
  osint.new_workspace
    fn (root: string) -> workspace_state
    + creates a workspace rooted at a directory with empty module registry
    # construction
  osint.registry_sync
    fn (state: workspace_state, index_url: string) -> result[workspace_state, string]
    + downloads the remote module index and caches it locally
    - returns error when the index is unreachable or malformed
    # package_index
    -> std.http.get
  osint.install_module
    fn (state: workspace_state, name: string) -> result[workspace_state, string]
    + downloads a module by name, verifies its digest, and records it as installed
    - returns error when the digest does not match the index entry
    # package_install
    -> std.http.get
    -> std.crypto.sha256_hex
    -> std.fs.write_all
  osint.remove_module
    fn (state: workspace_state, name: string) -> workspace_state
    + removes a module from the workspace; no-op if not installed
    # package_remove
  osint.list_modules
    fn (state: workspace_state) -> list[string]
    + returns the names of installed modules
    # introspection
  osint.add_seed
    fn (state: workspace_state, kind: string, value: string) -> workspace_state
    + inserts a seed entity of kind (domain, ip, email, etc) into the workspace
    # seeding
  osint.run_module
    fn (state: workspace_state, name: string, filter: map[string,string]) -> result[list[finding], string]
    + runs an installed module against seeds matching filter and returns new findings
    - returns error when the module is not installed
    # execution
    -> std.script.eval
  osint.findings
    fn (state: workspace_state, kind: string) -> list[finding]
    + returns findings of the given entity kind
    # query
  osint.export_json
    fn (state: workspace_state) -> string
    + serializes all findings and seeds as a JSON document
    # export
