# Requirement: "a resizable semaphore"

A pure-value semaphore: acquire and release return a new state; the capacity can change at runtime.

std: (all units exist)

semaphore
  semaphore.new
    fn (capacity: i32) -> semaphore_state
    + creates a semaphore with the given initial permit count
    ? available == capacity on construction
    # construction
  semaphore.try_acquire
    fn (state: semaphore_state) -> tuple[bool, semaphore_state]
    + returns (true, new_state) when a permit was available
    - returns (false, unchanged_state) when no permits remain
    # acquire
  semaphore.release
    fn (state: semaphore_state) -> result[semaphore_state, string]
    + returns a state with one more available permit
    - returns error when releasing would exceed capacity
    # release
  semaphore.resize
    fn (state: semaphore_state, new_capacity: i32) -> result[semaphore_state, string]
    + changes the capacity; shrinking only affects future acquires
    - returns error when new_capacity is negative
    # resize
  semaphore.available
    fn (state: semaphore_state) -> i32
    + returns the number of permits currently free
    # inspection
