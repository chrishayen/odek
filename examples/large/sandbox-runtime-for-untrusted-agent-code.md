# Requirement: "a sandbox runtime library for running untrusted agent code in isolated environments"

Manages short-lived isolated workspaces where an agent can run commands, read and write files, and fetch a snapshot when done. The actual isolation backend is injected by the caller.

std
  std.fs
    std.fs.mkdir_all
      @ (path: string) -> result[void, string]
      + creates the directory and all parents
      # filesystem
    std.fs.remove_all
      @ (path: string) -> result[void, string]
      + recursively removes the path
      # filesystem
    std.fs.write_file
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating parent directories
      # filesystem
    std.fs.read_file
      @ (path: string) -> result[bytes, string]
      + reads the entire file
      - returns error when the file does not exist
      # filesystem
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the hex-encoded SHA-256 digest
      # hashing

sandbox_runtime
  sandbox_runtime.new
    @ (root: string, backend: isolation_backend) -> runtime_state
    + creates a runtime that stores sandboxes under root using the given backend
    ? the isolation backend is an injected interface so the library is backend-agnostic
    # construction
    -> std.fs.mkdir_all
  sandbox_runtime.create_sandbox
    @ (state: runtime_state, image: string, cpu_quota: i32, memory_mb: i32, network: bool) -> result[tuple[runtime_state, sandbox_id], string]
    + allocates a fresh sandbox with the given resource limits
    - returns error when cpu_quota or memory_mb is zero or negative
    # lifecycle
    -> std.hash.sha256_hex
    -> std.time.now_millis
  sandbox_runtime.destroy_sandbox
    @ (state: runtime_state, id: sandbox_id) -> result[runtime_state, string]
    + terminates the sandbox and removes its working directory
    - returns error when the sandbox is unknown
    # lifecycle
    -> std.fs.remove_all
  sandbox_runtime.start
    @ (state: runtime_state, id: sandbox_id) -> result[runtime_state, string]
    + brings the sandbox into the running state via the backend
    - returns error when already running
    # lifecycle
  sandbox_runtime.stop
    @ (state: runtime_state, id: sandbox_id) -> result[runtime_state, string]
    + stops the sandbox but preserves its filesystem
    # lifecycle
  sandbox_runtime.exec
    @ (state: runtime_state, id: sandbox_id, argv: list[string], stdin: bytes, timeout_ms: i32) -> result[exec_result, string]
    + runs a command inside the sandbox and returns stdout, stderr, exit code
    - returns error when the sandbox is not running
    - returns error with a timeout marker when runtime exceeds timeout_ms
    # exec
    -> std.time.now_millis
  sandbox_runtime.write_file
    @ (state: runtime_state, id: sandbox_id, path: string, data: bytes) -> result[runtime_state, string]
    + writes a file inside the sandbox workspace
    - returns error when path escapes the workspace
    # files
    -> std.fs.write_file
  sandbox_runtime.read_file
    @ (state: runtime_state, id: sandbox_id, path: string) -> result[bytes, string]
    + reads a file from the sandbox workspace
    - returns error when the file does not exist
    # files
    -> std.fs.read_file
  sandbox_runtime.list_files
    @ (state: runtime_state, id: sandbox_id, subpath: string) -> result[list[string], string]
    + returns file names relative to subpath inside the sandbox
    - returns error when subpath is outside the workspace
    # files
  sandbox_runtime.snapshot
    @ (state: runtime_state, id: sandbox_id) -> result[bytes, string]
    + returns a tar archive of the current sandbox workspace
    - returns error when the sandbox is unknown
    # snapshot
  sandbox_runtime.restore
    @ (state: runtime_state, image: string, archive: bytes) -> result[tuple[runtime_state, sandbox_id], string]
    + creates a new sandbox and populates its workspace from a snapshot archive
    - returns error on corrupt archive
    # snapshot
  sandbox_runtime.usage
    @ (state: runtime_state, id: sandbox_id) -> optional[sandbox_usage]
    + returns current cpu and memory usage for the sandbox
    - returns none when id is unknown
    # observability
