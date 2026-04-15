# Requirement: "a skiplist data structure"

An ordered key-value skiplist supporting insert, lookup, and delete. A thin std random source enables deterministic tests.

std
  std.random
    std.random.next_u32
      fn (state: random_state) -> tuple[u32, random_state]
      + returns a uniformly distributed u32 and the updated state
      # randomness

skiplist
  skiplist.new
    fn (seed: u64) -> skiplist_state
    + creates an empty skiplist with a seeded rng for level generation
    # construction
  skiplist.insert
    fn (state: skiplist_state, key: string, value: i64) -> skiplist_state
    + inserts a new key-value pair in sorted order
    + replaces the value when key already exists
    # writes
    -> std.random.next_u32
  skiplist.get
    fn (state: skiplist_state, key: string) -> optional[i64]
    + returns the value associated with key
    - returns none when key is absent
    # reads
  skiplist.delete
    fn (state: skiplist_state, key: string) -> skiplist_state
    + removes key from the skiplist
    - is a no-op when key is absent
    # writes
  skiplist.range
    fn (state: skiplist_state, from_key: string, to_key: string) -> list[tuple[string, i64]]
    + returns all entries with from_key <= key < to_key in sorted order
    # reads
