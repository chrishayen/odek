# Requirement: "a cancellation token for making execution flow interruptible"

A cooperative cancellation primitive: create a token, trigger it, observe it. Callers check the flag at safe points.

std: (all units exist)

breaker
  breaker.new
    @ () -> breaker_state
    + creates a fresh, unbroken breaker
    # construction
  breaker.trigger
    @ (state: breaker_state) -> breaker_state
    + marks the breaker as broken; subsequent checks return true
    ? triggering an already-broken breaker is a no-op
    # cancellation
  breaker.is_broken
    @ (state: breaker_state) -> bool
    + returns true after trigger has been called
    - returns false for a freshly created breaker
    # observation
  breaker.check
    @ (state: breaker_state) -> result[void, string]
    + returns ok when not broken, error "cancelled" when broken
    - returns error after trigger was called
    # cancellation
