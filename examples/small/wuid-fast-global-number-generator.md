# Requirement: "a fast globally unique number generator"

Each generator is initialized with a shard-reserved high section; the low section counts up locally. When the low section is exhausted the caller must reserve a new high section.

std: (all units exist)

wuid
  wuid.new
    fn (reserve_high: fn() -> result[u64, string]) -> result[wuid_state, string]
    + reserves an initial high section and returns a ready generator
    - returns error when reserve_high fails
    # construction
  wuid.next
    fn (state: wuid_state) -> result[tuple[u64, wuid_state], string]
    + returns the next identifier combining the current high section with an incremented low counter
    - returns error when the low counter has reached its maximum; caller should call renew
    # generation
  wuid.renew
    fn (state: wuid_state) -> result[wuid_state, string]
    + reserves a new high section and resets the low counter
    - returns error when the reserve function fails
    # renewal
  wuid.should_renew
    fn (state: wuid_state, watermark: u32) -> bool
    + returns true when the low counter has passed the watermark, signaling the caller to renew soon
    # renewal
