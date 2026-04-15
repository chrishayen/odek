# Requirement: "a service registry with load balancing and failover for mid-tier services"

Clients register instances under service names with health and weight; consumers pick an instance using weighted selection and mark failing instances so they rotate out.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.random
    std.random.next_u64
      fn () -> u64
      + returns a uniformly distributed 64-bit value
      # randomness

registry
  registry.new
    fn () -> registry_state
    + creates an empty registry
    # construction
  registry.register
    fn (state: registry_state, service: string, instance_id: string, address: string, weight: i32) -> registry_state
    + adds or replaces an instance under the given service name
    # registration
    -> std.time.now_millis
  registry.heartbeat
    fn (state: registry_state, service: string, instance_id: string) -> result[registry_state, string]
    + refreshes the instance's last-seen timestamp
    - returns error when the instance is not registered
    # health
    -> std.time.now_millis
  registry.expire
    fn (state: registry_state, ttl_millis: i64) -> registry_state
    + removes instances whose last heartbeat is older than ttl_millis
    # health
    -> std.time.now_millis
  registry.pick
    fn (state: registry_state, service: string) -> result[service_instance, string]
    + returns an instance chosen by weighted random among healthy instances
    - returns error when no healthy instance exists for the service
    # load_balancing
    -> std.random.next_u64
  registry.mark_failed
    fn (state: registry_state, service: string, instance_id: string) -> registry_state
    + marks the instance as failed so it is skipped until its cooldown elapses
    # failover
    -> std.time.now_millis
