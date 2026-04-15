# Requirement: "a microservices framework with service discovery and load balancing"

Services register themselves under a name with one or more endpoints. Clients resolve a name to an endpoint via a selectable load-balancing strategy, and the registry prunes endpoints that fail health checks.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.random
    std.random.int_range
      fn (lo: i32, hi_exclusive: i32) -> i32
      + returns a uniformly random integer in [lo, hi_exclusive)
      # randomness
  std.hash
    std.hash.fnv1a
      fn (s: string) -> u64
      + returns the 64-bit FNV-1a hash of s
      # hashing

microsvc
  microsvc.new_registry
    fn () -> registry_state
    + creates an empty registry with no services
    # construction
  microsvc.register_instance
    fn (state: registry_state, service_name: string, instance_id: string, address: string) -> registry_state
    + adds an instance under service_name
    - returns unchanged state when instance_id is already registered for service_name
    # registration
  microsvc.deregister_instance
    fn (state: registry_state, service_name: string, instance_id: string) -> registry_state
    + removes an instance
    # registration
  microsvc.heartbeat
    fn (state: registry_state, service_name: string, instance_id: string) -> registry_state
    + records the current time as the last heartbeat for the instance
    # health
    -> std.time.now_millis
  microsvc.prune_unhealthy
    fn (state: registry_state, max_age_millis: i64) -> registry_state
    + removes instances whose last heartbeat is older than max_age_millis
    # health
    -> std.time.now_millis
  microsvc.list_instances
    fn (state: registry_state, service_name: string) -> list[instance]
    + returns all currently registered healthy instances for service_name
    - returns an empty list when the service has no instances
    # discovery
  microsvc.resolve_random
    fn (state: registry_state, service_name: string) -> result[instance, string]
    + returns a uniformly random instance
    - returns error "no instances" when none are registered
    # load_balancing
    -> std.random.int_range
  microsvc.resolve_round_robin
    fn (state: registry_state, service_name: string) -> tuple[registry_state, result[instance, string]]
    + returns the next instance in round-robin order and advances the cursor
    - returns error "no instances" when none are registered
    # load_balancing
  microsvc.resolve_consistent_hash
    fn (state: registry_state, service_name: string, routing_key: string) -> result[instance, string]
    + returns the instance whose hashed id is closest to the hashed routing key
    - returns error "no instances" when none are registered
    ? the same routing key resolves to the same instance across calls unless membership changes
    # load_balancing
    -> std.hash.fnv1a
  microsvc.watch
    fn (state: registry_state, service_name: string, watcher_id: string) -> registry_state
    + subscribes watcher_id to changes in the service's instance set
    # watches
  microsvc.drain_change_events
    fn (state: registry_state, watcher_id: string) -> tuple[registry_state, list[change_event]]
    + returns and clears the queued events for a watcher
    # watches
