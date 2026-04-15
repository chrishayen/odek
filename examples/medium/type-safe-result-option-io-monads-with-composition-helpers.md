# Requirement: "a library of type-safe result, option, and IO monads with composition helpers"

Constructors, map/bind, and small composition helpers for three monadic containers.

std: (all units exist)

monads
  monads.some
    fn (value: string) -> option_state
    + wraps a value in a present option
    # option
  monads.none
    fn () -> option_state
    + returns an absent option
    # option
  monads.option_map
    fn (opt: option_state, fn: fn_string_to_string) -> option_state
    + applies fn when the option is present, returns none otherwise
    # option
  monads.option_bind
    fn (opt: option_state, fn: fn_string_to_option) -> option_state
    + chains option-returning operations
    # option
  monads.ok
    fn (value: string) -> result_state
    + wraps a value in a success result
    # result
  monads.err
    fn (error: string) -> result_state
    + wraps an error in a failure result
    # result
  monads.result_map
    fn (res: result_state, fn: fn_string_to_string) -> result_state
    + applies fn on success, propagates failure unchanged
    # result
  monads.result_bind
    fn (res: result_state, fn: fn_string_to_result) -> result_state
    + chains result-returning operations, short-circuiting on failure
    # result
  monads.io_pure
    fn (value: string) -> io_state
    + lifts a pure value into a deferred IO computation
    # io
  monads.io_run
    fn (action: io_state) -> string
    + executes a deferred IO computation and returns its result
    # io
  monads.pipe
    fn (fns: list[fn_string_to_string]) -> fn_string_to_string
    + returns a function that applies each fn in order left to right
    # composition
  monads.compose
    fn (fns: list[fn_string_to_string]) -> fn_string_to_string
    + returns a function that applies each fn in order right to left
    # composition
