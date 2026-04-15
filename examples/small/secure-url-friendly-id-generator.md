# Requirement: "a secure url-friendly unique string id generator"

Generates compact URL-safe random IDs. Randomness comes from a thin std primitive.

std
  std.random
    std.random.fill_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness

secure_id
  secure_id.generate
    fn () -> string
    + returns a 21-character URL-safe id using the default alphabet "A-Za-z0-9_-"
    + successive calls return distinct values with overwhelming probability
    # id_generation
    -> std.random.fill_bytes
  secure_id.generate_sized
    fn (size: i32) -> string
    + returns a URL-safe id of the requested length
    - returns empty string when size <= 0
    ? uses rejection sampling so character frequency is uniform over the alphabet
    # id_generation
    -> std.random.fill_bytes
