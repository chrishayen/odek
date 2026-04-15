# Requirement: "monads and functional programming abstractions (Option, Result, Either)"

Constructors and bind/map combinators for the core sum types. Values are carried opaquely so runtime wrapping and discriminators are hidden.

std: (all units exist)

fp
  fp.option_some
    fn (value: bytes) -> option_val
    + returns a Some wrapping the value
    # option
  fp.option_none
    fn () -> option_val
    + returns the None variant
    # option
  fp.option_map
    fn (opt: option_val, f: fn[bytes, bytes]) -> option_val
    + applies f to the inner value when present
    - returns None unchanged when the input is None
    # option
  fp.option_unwrap_or
    fn (opt: option_val, default: bytes) -> bytes
    + returns the inner value or the default
    # option
  fp.result_ok
    fn (value: bytes) -> result_val
    + returns an Ok variant
    # result
  fp.result_err
    fn (err: string) -> result_val
    + returns an Err variant
    # result
  fp.result_map
    fn (r: result_val, f: fn[bytes, bytes]) -> result_val
    + maps the Ok branch
    - returns the Err unchanged
    # result
  fp.result_and_then
    fn (r: result_val, f: fn[bytes, result_val]) -> result_val
    + chains a fallible step on the Ok branch
    # result
  fp.either_left
    fn (value: bytes) -> either_val
    + returns a Left variant
    # either
  fp.either_right
    fn (value: bytes) -> either_val
    + returns a Right variant
    # either
  fp.either_fold
    fn (e: either_val, on_left: fn[bytes, bytes], on_right: fn[bytes, bytes]) -> bytes
    + applies on_left or on_right depending on the variant
    # either
