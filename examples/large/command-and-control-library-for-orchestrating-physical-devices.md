# Requirement: "a command and control library for orchestrating physical devices"

Devices register with a central hub, receive dispatched commands, and report telemetry. Transport and auth are pluggable; the hub tracks device state and routes commands.

std
  std.net
    std.net.tcp_listen
      fn (host: string, port: i32) -> result[listener_state, string]
      + binds a TCP listener
      # networking
    std.net.tcp_accept
      fn (listener: listener_state) -> result[conn_state, string]
      + blocks for a client
      # networking
    std.net.conn_read
      fn (conn: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes
      # networking
    std.net.conn_write
      fn (conn: conn_state, data: bytes) -> result[void, string]
      + writes the full buffer
      # networking
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.crypto
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + cryptographically random bytes
      # cryptography
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + returns 32-byte HMAC
      # cryptography
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (m: map[string,string]) -> string
      + encodes a string map as JSON
      # serialization

control
  control.new_hub
    fn () -> hub_state
    + creates a hub with no registered devices
    # construction
  control.register_device
    fn (hub: hub_state, device_id: string, shared_secret: bytes) -> hub_state
    + adds the device and stores its HMAC secret for auth
    # registration
  control.authenticate
    fn (hub: hub_state, device_id: string, nonce: bytes, mac: bytes) -> result[void, string]
    + verifies mac equals HMAC(secret, nonce) for the device
    - returns error when device is unknown
    - returns error when the MAC does not match
    # auth
    -> std.crypto.hmac_sha256
  control.device_seen
    fn (hub: hub_state, device_id: string) -> hub_state
    + updates the device's last_seen timestamp to now
    # presence
    -> std.time.now_millis
  control.is_online
    fn (hub: hub_state, device_id: string, stale_after_ms: i64) -> bool
    + returns true when last_seen is newer than now - stale_after_ms
    # presence
    -> std.time.now_millis
  control.issue_command
    fn (hub: hub_state, device_id: string, command: string, payload: map[string,string]) -> result[command_envelope, string]
    + queues a command for the device with a fresh random id
    - returns error when device is not registered
    # command_dispatch
    -> std.crypto.random_bytes
    -> std.json.encode_object
  control.pop_commands
    fn (hub: hub_state, device_id: string) -> tuple[list[command_envelope], hub_state]
    + drains and returns queued commands for the device
    # command_dispatch
  control.report_telemetry
    fn (hub: hub_state, device_id: string, raw: string) -> result[hub_state, string]
    + parses raw as a JSON telemetry object and records it
    - returns error when JSON is invalid
    # telemetry
    -> std.json.parse_object
  control.query_telemetry
    fn (hub: hub_state, device_id: string) -> optional[map[string,string]]
    + returns the latest telemetry for the device
    # telemetry
  control.serve
    fn (hub: hub_state, host: string, port: i32) -> result[void, string]
    + accepts device connections and runs the auth+command+telemetry loop
    # server_loop
    -> std.net.tcp_listen
    -> std.net.tcp_accept
    -> std.net.conn_read
    -> std.net.conn_write
