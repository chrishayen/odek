# Requirement: "a platform-as-a-service core library for deploying and managing applications"

Tracks applications, their deployed versions, routes traffic, and streams lifecycle events. The std layer carries general-purpose primitives the project composes.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.ids
    std.ids.new_id
      @ () -> string
      + returns a fresh opaque identifier unique within the process
      # identifiers
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex SHA-256 digest of the input
      # hashing
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

paas
  paas.new_platform
    @ () -> platform_state
    + creates an empty platform with no applications and no events
    # construction
  paas.create_app
    @ (state: platform_state, name: string) -> result[tuple[string, platform_state], string]
    + returns (app_id, new_state) when the name is free
    - returns error when an application with the same name exists
    # applications
    -> std.ids.new_id
    -> std.time.now_seconds
  paas.push_version
    @ (state: platform_state, app_id: string, source: bytes) -> result[tuple[string, platform_state], string]
    + stores a new version keyed by content digest and returns the version id
    - returns error when the app_id is unknown
    # versioning
    -> std.hash.sha256_hex
    -> std.ids.new_id
    -> std.time.now_seconds
  paas.deploy
    @ (state: platform_state, app_id: string, version_id: string) -> result[platform_state, string]
    + marks the version as the active deployment for the app and emits a deploy event
    - returns error when the app or version is unknown
    # deployment
    -> std.time.now_seconds
  paas.rollback
    @ (state: platform_state, app_id: string) -> result[platform_state, string]
    + restores the previously active version and emits a rollback event
    - returns error when no previous version exists
    # deployment
    -> std.time.now_seconds
  paas.scale
    @ (state: platform_state, app_id: string, instances: i32) -> result[platform_state, string]
    + sets the desired instance count and emits a scale event
    - returns error when instances is negative
    # scaling
  paas.route
    @ (state: platform_state, host: string) -> optional[string]
    + returns the app_id bound to the given host, if any
    # routing
  paas.bind_route
    @ (state: platform_state, app_id: string, host: string) -> result[platform_state, string]
    + associates the host with the app
    - returns error when the host is already bound to a different app
    # routing
  paas.app_status
    @ (state: platform_state, app_id: string) -> result[string, string]
    + returns a JSON status document describing the app's active version and instance count
    - returns error when the app_id is unknown
    # status
    -> std.json.encode_object
  paas.recent_events
    @ (state: platform_state, app_id: string, limit: i32) -> list[string]
    + returns the most recent lifecycle events for the app as JSON strings, newest first
    # events
    -> std.json.encode_object
