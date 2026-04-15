# Requirement: "a simple load balancer"

Round-robin selection over a pool of backends with health tracking. Networking is the caller's responsibility; this library only decides where to send next.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

load_balancer
  load_balancer.new
    fn (backends: list[string]) -> lb_state
    + creates a balancer with the given backend addresses, all marked healthy
    - returns an empty-pool state when the list is empty
    # construction
  load_balancer.next
    fn (state: lb_state) -> tuple[optional[string], lb_state]
    + returns the next healthy backend in round-robin order and the advanced state
    - returns none when no backends are healthy
    # selection
  load_balancer.mark_down
    fn (state: lb_state, backend: string, cooldown_ms: i64) -> lb_state
    + marks a backend unhealthy until cooldown expires
    # health
    -> std.time.now_millis
  load_balancer.mark_up
    fn (state: lb_state, backend: string) -> lb_state
    + marks a backend healthy and eligible for selection
    # health
  load_balancer.healthy_count
    fn (state: lb_state) -> i32
    + returns the number of backends currently eligible
    # introspection
    -> std.time.now_millis
