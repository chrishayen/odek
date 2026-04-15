# Requirement: "a utility for generating random data based on an input template"

Scans a template for placeholder tokens like {int}, {string}, {uuid} and substitutes random values.

std
  std.rand
    std.rand.next_int
      fn (min_val: i64, max_val: i64) -> i64
      + returns a uniformly random integer in [min_val, max_val)
      # randomness
    std.rand.next_string
      fn (length: i32) -> string
      + returns a random alphanumeric string of the requested length
      # randomness
    std.rand.uuid_v4
      fn () -> string
      + returns a random UUID v4
      # randomness

template_rand
  template_rand.render
    fn (template: string) -> result[string, string]
    + replaces each {int}, {string}, {uuid} token in the template with a fresh random value
    + leaves unknown tokens untouched
    - returns error when a token is unterminated (missing closing brace)
    # rendering
    -> std.rand.next_int
    -> std.rand.next_string
    -> std.rand.uuid_v4
