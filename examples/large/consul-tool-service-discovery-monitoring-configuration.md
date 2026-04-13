# Requirement: "a service discovery, health monitoring, and configuration library"

A registry for services with health checks and a key-value store for shared configuration. Persistence and networking are out of scope; this is the in-memory core.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.uuid
    std.uuid.new_v4
      @ () -> string
      + returns a random UUID string
      # identifiers

service_mesh
  service_mesh.new_registry
    @ () -> registry_state
    + creates an empty service registry with an embedded key-value store
    # construction
  service_mesh.register_service
    @ (state: registry_state, name: string, address: string, port: i32) -> tuple[string, registry_state]
    + assigns a fresh instance id, records the instance as healthy, and returns the id
    # registration
    -> std.uuid.new_v4
    -> std.time.now_millis
  service_mesh.deregister_service
    @ (state: registry_state, instance_id: string) -> result[registry_state, string]
    + removes the instance and any attached checks
    - returns error when no instance matches the id
    # registration
  service_mesh.list_healthy
    @ (state: registry_state, name: string) -> list[service_instance]
    + returns all currently healthy instances for a service name
    # discovery
  service_mesh.record_heartbeat
    @ (state: registry_state, instance_id: string, ttl_ms: i64) -> result[registry_state, string]
    + stamps the instance as healthy with an expiration ttl_ms in the future
    - returns error when the instance is unknown
    # health
    -> std.time.now_millis
  service_mesh.expire_stale
    @ (state: registry_state) -> registry_state
    + marks instances whose ttl has passed as unhealthy
    # health
    -> std.time.now_millis
  service_mesh.kv_put
    @ (state: registry_state, key: string, value: bytes) -> registry_state
    + stores a value under a hierarchical key like "app/db/password"
    # configuration
  service_mesh.kv_get
    @ (state: registry_state, key: string) -> optional[bytes]
    + returns the value for a key if present
    # configuration
  service_mesh.kv_list_prefix
    @ (state: registry_state, prefix: string) -> list[string]
    + returns all keys under a prefix in sorted order
    # configuration
  service_mesh.watch_key
    @ (state: registry_state, key: string, last_revision: i64) -> tuple[optional[bytes], i64]
    + returns the current value and revision if the key has changed since last_revision, else (none, last_revision)
    # configuration
