# Requirement: "a server administration panel library exposing system metrics, service control, and user sessions"

A backend that powers a server admin UI: authenticated sessions, service lifecycle, metric snapshots. No rendering, just the core library.

std
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns 32 bytes of SHA-256 digest
      # cryptography
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      - returns empty bytes when n is zero or negative
      # cryptography
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + returns lowercase hex representation
      # encoding
  std.os
    std.os.read_cpu_percent
      @ () -> f64
      + returns system-wide CPU utilization 0.0-100.0
      # system
    std.os.read_memory_bytes
      @ () -> tuple[i64, i64]
      + returns (used_bytes, total_bytes)
      # system
    std.os.list_processes
      @ () -> list[process_info]
      + returns all running processes with pid and name
      # system
    std.os.signal_process
      @ (pid: i32, signal: string) -> result[void, string]
      + sends the named signal to the given pid
      - returns error when the signal name is unknown
      - returns error when the process does not exist
      # system

admin_panel
  admin_panel.new
    @ () -> panel_state
    + returns a panel with no users and no active sessions
    # construction
  admin_panel.add_user
    @ (state: panel_state, username: string, password: string) -> result[panel_state, string]
    + stores the user with a salted password hash
    - returns error when username already exists
    # users
    -> std.crypto.random_bytes
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  admin_panel.authenticate
    @ (state: panel_state, username: string, password: string) -> result[string, string]
    + returns a session token on success
    - returns error when credentials do not match
    # authentication
    -> std.crypto.sha256
    -> std.crypto.random_bytes
    -> std.encoding.hex_encode
    -> std.time.now_seconds
  admin_panel.check_session
    @ (state: panel_state, token: string) -> result[string, string]
    + returns the username owning a valid non-expired session
    - returns error when the token is unknown or expired
    # authentication
    -> std.time.now_seconds
  admin_panel.revoke_session
    @ (state: panel_state, token: string) -> panel_state
    + removes the session from the state
    # authentication
  admin_panel.snapshot_metrics
    @ () -> panel_metrics
    + returns current CPU, memory, and process count
    # metrics
    -> std.os.read_cpu_percent
    -> std.os.read_memory_bytes
    -> std.os.list_processes
  admin_panel.list_services
    @ (state: panel_state) -> list[service_info]
    + returns registered services with their last known status
    # services
  admin_panel.register_service
    @ (state: panel_state, name: string, pid: i32) -> panel_state
    + adds a service entry with running status
    # services
  admin_panel.stop_service
    @ (state: panel_state, name: string) -> result[panel_state, string]
    + sends a termination signal to the service's pid
    - returns error when the service is not registered
    # services
    -> std.os.signal_process
  admin_panel.audit_log
    @ (state: panel_state) -> list[audit_entry]
    + returns the chronological list of privileged actions
    # audit
