# Requirement: "a library that generates fake data"

Deterministic, seedable generators for common synthetic fields: names, emails, addresses, numbers, dates. Randomness is a std primitive.

std
  std.random
    std.random.new
      fn (seed: i64) -> rng_state
      + creates a deterministic rng from a seed
      # randomness
    std.random.next_int
      fn (state: rng_state, min: i64, max: i64) -> tuple[i64, rng_state]
      + returns an integer in [min, max] inclusive and the advanced state
      # randomness

fake
  fake.new
    fn (seed: i64) -> fake_state
    + creates a fake generator with a seeded rng
    # construction
    -> std.random.new
  fake.first_name
    fn (state: fake_state) -> tuple[string, fake_state]
    + returns a random first name and the advanced state
    # people
    -> std.random.next_int
  fake.last_name
    fn (state: fake_state) -> tuple[string, fake_state]
    + returns a random last name and the advanced state
    # people
    -> std.random.next_int
  fake.email
    fn (state: fake_state) -> tuple[string, fake_state]
    + returns "first.last@domain.tld" with randomized parts
    # internet
  fake.phone
    fn (state: fake_state) -> tuple[string, fake_state]
    + returns a formatted phone number
    # contact
    -> std.random.next_int
  fake.address
    fn (state: fake_state) -> tuple[string, fake_state]
    + returns a single-line street address
    # geography
  fake.int_between
    fn (state: fake_state, low: i64, high: i64) -> tuple[i64, fake_state]
    + returns an integer in [low, high] inclusive
    - returns an error-valued tuple when low > high
    # numeric
    -> std.random.next_int
  fake.date_between
    fn (state: fake_state, start_unix: i64, end_unix: i64) -> tuple[i64, fake_state]
    + returns a unix timestamp in the given range
    # temporal
    -> std.random.next_int
  fake.uuid
    fn (state: fake_state) -> tuple[string, fake_state]
    + returns a random uuid string derived from the rng
    # identifiers
    -> std.random.next_int
