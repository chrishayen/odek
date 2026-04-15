# Requirement: "a mock data generator driven by schemas"

Registers field schemas, then produces synthetic records for tests.

std
  std.random
    std.random.new_rng
      fn (seed: u64) -> rng_state
      + creates a deterministic pseudo-random generator seeded with seed
      # random
    std.random.next_u64
      fn (rng: rng_state) -> tuple[u64, rng_state]
      + returns the next 64-bit value and the advanced generator state
      # random

mock
  mock.new_registry
    fn () -> registry
    + creates an empty schema registry
    # construction
  mock.register
    fn (reg: registry, name: string, fields: list[field_spec]) -> registry
    + registers a named record schema with ordered field specs
    ? each field_spec carries a field name and a generator kind (int range, float range, string pattern, pick-from-list, nested-schema)
    # registration
  mock.generate
    fn (reg: registry, name: string, rng: rng_state) -> result[tuple[record, rng_state], string]
    + produces one record for the named schema using rng, returning the advanced state
    - returns error when the schema name is unknown
    # generation
    -> std.random.next_u64
  mock.generate_many
    fn (reg: registry, name: string, count: i32, rng: rng_state) -> result[tuple[list[record], rng_state], string]
    + produces count records from the named schema
    - returns error when the schema name is unknown
    # batch_generation
  mock.to_json
    fn (rec: record) -> string
    + encodes a record as JSON
    # serialization
