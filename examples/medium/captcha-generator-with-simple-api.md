# Requirement: "a captcha generation library with a simple unopinionated api"

Generates a short random challenge, renders it as an image, and validates submitted answers against a store with expirations.

std
  std.random
    std.random.int_range
      @ (min: i32, max: i32) -> i32
      + returns a uniform random integer in [min, max]
      # randomness
    std.random.bytes
      @ (length: i32) -> bytes
      + returns cryptographically random bytes
      # randomness
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

captcha
  captcha.generate_challenge
    @ (length: i32, alphabet: string) -> tuple[string, string]
    + returns (challenge_id, solution) with a solution of the given length drawn from the alphabet
    # generation
    -> std.random.int_range
    -> std.random.bytes
  captcha.render_image
    @ (solution: string, width: i32, height: i32) -> bytes
    + returns a distorted image in PNG form displaying the solution
    # rendering
    -> std.random.int_range
  captcha.new_store
    @ (ttl_seconds: i64) -> captcha_store_state
    + returns an empty captcha store with the given entry lifetime
    # storage
  captcha.remember
    @ (store: captcha_store_state, id: string, solution: string) -> captcha_store_state
    + stores a challenge with the current timestamp
    # storage
    -> std.time.now_seconds
  captcha.verify
    @ (store: captcha_store_state, id: string, answer: string) -> tuple[bool, captcha_store_state]
    + returns (true, new_state) when the answer matches an unexpired entry and consumes it
    - returns (false, unchanged) when no entry matches
    - returns (false, new_state_with_entry_removed) when the entry was expired
    # verification
    -> std.time.now_seconds
