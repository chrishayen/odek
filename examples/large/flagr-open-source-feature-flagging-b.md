# Requirement: "a feature flagging and A/B testing library"

Evaluates boolean and multi-variant flags based on rule predicates and deterministic user bucketing. Flags, segments, and rollouts live in a registry; evaluation consults the registry and a stable hash of the user key.

std
  std.hash
    std.hash.fnv1a_32
      @ (data: bytes) -> u32
      + returns 32-bit FNV-1a hash
      # hashing
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a json_value
      - returns error on malformed JSON
      # serialization

flags
  flags.new_registry
    @ () -> flag_registry
    + returns an empty flag registry
    # construction
  flags.upsert_flag
    @ (registry: flag_registry, flag: flag_definition) -> flag_registry
    + adds or replaces a flag by key
    # configuration
  flags.delete_flag
    @ (registry: flag_registry, key: string) -> flag_registry
    + removes a flag by key
    # configuration
  flags.bucket_user
    @ (flag_key: string, user_key: string) -> i32
    + returns a deterministic bucket 0..9999 for (flag, user)
    ? independent flags produce independent buckets for the same user
    # bucketing
    -> std.hash.fnv1a_32
  flags.match_segment
    @ (segment: segment_rule, attributes: map[string, string]) -> bool
    + returns true when all predicates in segment match the attributes
    # targeting
  flags.pick_variant
    @ (rollout: list[variant_weight], bucket: i32) -> string
    + returns the variant whose cumulative weight range contains bucket
    ? the sum of weights must be 10000
    # rollout
  flags.evaluate
    @ (registry: flag_registry, key: string, user_key: string, attributes: map[string, string]) -> flag_result
    + returns the chosen variant and whether an overriding rule matched
    - returns the default variant when the flag is disabled or absent
    # evaluation
    -> std.time.now_seconds
  flags.evaluate_all
    @ (registry: flag_registry, user_key: string, attributes: map[string, string]) -> map[string, flag_result]
    + evaluates every flag in the registry for the user
    # evaluation
  flags.record_exposure
    @ (stats: exposure_stats, key: string, variant: string) -> exposure_stats
    + increments the exposure counter for (flag, variant)
    # telemetry
  flags.variant_rates
    @ (stats: exposure_stats, key: string) -> map[string, f64]
    + returns the fraction of exposures observed for each variant
    # telemetry
  flags.load_definitions
    @ (raw: string) -> result[list[flag_definition], string]
    + parses a JSON document into flag definitions
    - returns error on malformed JSON or missing fields
    # serialization
    -> std.json.parse
