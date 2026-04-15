# Requirement: "an in-memory offender jailer with configurable rules"

Tracks infractions per actor, jails them when a rule threshold is crossed, and releases them after a sentence elapses.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

jailer
  jailer.new
    fn () -> jailer_state
    + creates an empty jailer with no rules and no offenders
    # construction
  jailer.add_rule
    fn (j: jailer_state, name: string, max_infractions: i32, window_seconds: i64, sentence_seconds: i64) -> jailer_state
    + registers a named rule with a sliding window and jail duration
    - max_infractions must be positive
    # rules
  jailer.infract
    fn (j: jailer_state, rule: string, actor: string) -> result[jailer_state, string]
    + records one infraction against the actor under the given rule
    + jails the actor when infractions within the window exceed the threshold
    - returns error when rule is unknown
    # infractions
    -> std.time.now_seconds
  jailer.is_jailed
    fn (j: jailer_state, actor: string) -> bool
    + returns true when the actor has an unexpired sentence
    + returns false when the sentence has elapsed
    # query
    -> std.time.now_seconds
  jailer.release
    fn (j: jailer_state, actor: string) -> jailer_state
    + clears any active sentence for the actor
    # release
  jailer.sweep_expired
    fn (j: jailer_state) -> jailer_state
    + removes expired sentences and infractions older than any rule window
    # maintenance
    -> std.time.now_seconds
