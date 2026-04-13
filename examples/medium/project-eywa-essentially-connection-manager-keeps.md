# Requirement: "a connection manager that keeps track of connected devices"

Registers devices, tracks liveness via heartbeats, and emits events when devices connect or disconnect.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.uuid
    std.uuid.v4
      @ () -> string
      + returns a random UUIDv4 string
      # uuid

device_manager
  device_manager.new
    @ (heartbeat_timeout_ms: i64) -> manager_state
    + creates a manager with a liveness timeout
    # construction
  device_manager.connect
    @ (state: manager_state, device_id: string, metadata: map[string, string]) -> tuple[string, manager_state]
    + records a device as connected and assigns a session id
    + emits a "connected" event
    # lifecycle
    -> std.uuid.v4
    -> std.time.now_millis
  device_manager.heartbeat
    @ (state: manager_state, session_id: string) -> result[manager_state, string]
    + refreshes the last-seen timestamp for the session
    - returns error when session_id is unknown
    # lifecycle
    -> std.time.now_millis
  device_manager.disconnect
    @ (state: manager_state, session_id: string) -> result[manager_state, string]
    + removes the session and emits a "disconnected" event
    - returns error when session_id is unknown
    # lifecycle
  device_manager.sweep
    @ (state: manager_state) -> tuple[list[string], manager_state]
    + removes sessions whose last-seen is older than the timeout
    + returns the list of expired session ids
    # lifecycle
    -> std.time.now_millis
  device_manager.list_connected
    @ (state: manager_state) -> list[device]
    + returns every currently connected device
    # introspection
  device_manager.drain_events
    @ (state: manager_state) -> tuple[list[event], manager_state]
    + returns the buffered lifecycle events and clears the buffer
    # events
