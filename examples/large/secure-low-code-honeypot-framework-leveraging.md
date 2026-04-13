# Requirement: "a honeypot framework that simulates services from declarative configuration"

Loads service definitions that describe listeners and canned responses, runs those listeners, and records attacker interactions to a pluggable sink.

std
  std.net
    std.net.listen_tcp
      @ (addr: string) -> result[listener, string]
      + binds to addr and returns a listener handle
      - returns error when the port is in use
      # network
    std.net.accept
      @ (l: listener) -> result[connection, string]
      + blocks until a new connection arrives
      # network
    std.net.read_bytes
      @ (c: connection, max: i32) -> result[bytes, string]
      + reads up to max bytes from the connection
      # network
    std.net.write_bytes
      @ (c: connection, data: bytes) -> result[i32, string]
      + writes bytes and returns the number written
      # network
    std.net.close
      @ (c: connection) -> void
      + releases the connection
      # network
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[yaml_value, string]
      + parses a YAML document into a tree
      - returns error on invalid YAML
      # serialization
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

honeypot
  honeypot.new
    @ (sink: fn(event) -> void) -> honeypot_state
    + creates an empty framework with the given event sink
    ? the event sink is pluggable so callers can route events anywhere
    # construction
  honeypot.load_config
    @ (state: honeypot_state, yaml_source: string) -> result[honeypot_state, string]
    + parses a service configuration document and registers the services
    - returns error on malformed YAML
    - returns error when a service entry is missing a required field
    # configuration
    -> std.yaml.parse
  honeypot.register_service
    @ (state: honeypot_state, name: string, proto: string, port: i32, script: response_script) -> result[honeypot_state, string]
    + adds a single service with its listen port and response script
    - returns error when the port is already claimed
    # configuration
  honeypot.start
    @ (state: honeypot_state) -> result[honeypot_state, string]
    + opens listeners for every registered service
    - returns error when any bind fails
    # lifecycle
    -> std.net.listen_tcp
  honeypot.stop
    @ (state: honeypot_state) -> honeypot_state
    + closes all listeners and drops active connections
    # lifecycle
    -> std.net.close
  honeypot.handle_connection
    @ (state: honeypot_state, service_name: string, c: connection) -> honeypot_state
    + runs the service script against the connection and emits events
    # session
    -> std.net.read_bytes
    -> std.net.write_bytes
    -> std.time.now_millis
  honeypot.script_step
    @ (script: response_script, input: bytes, cursor: i32) -> tuple[bytes, i32, bool]
    + returns (reply_bytes, next_cursor, session_done) for the given input
    + returns empty reply when the script has nothing queued for this step
    # scripting
  honeypot.emit_event
    @ (state: honeypot_state, ev: event) -> honeypot_state
    + passes the event to the configured sink
    # observability
  honeypot.record_session
    @ (state: honeypot_state, service_name: string, remote: string, transcript: list[bytes]) -> honeypot_state
    + appends a completed session to the in-memory audit log
    # audit
  honeypot.sessions_for
    @ (state: honeypot_state, service_name: string) -> list[session_record]
    + returns all recorded sessions for a service
    - returns empty list when the service is unknown
    # audit
