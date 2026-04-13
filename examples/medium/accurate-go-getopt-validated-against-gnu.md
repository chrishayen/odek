# Requirement: "a getopt-style command-line argument parser matching GNU semantics"

The library parses an argv-like list against a flag specification. It does not touch the process environment or print usage on its own — those are the caller's choice.

std: (all units exist)

getopt
  getopt.spec_new
    @ () -> spec
    + returns an empty flag specification
    # construction
  getopt.spec_add
    @ (s: spec, short: string, long: string, has_arg: i32) -> result[spec, string]
    + adds a flag definition carrying a short name, long name, and argument mode
    ? has_arg is 0 (no arg), 1 (required arg), 2 (optional arg)
    - returns error when short is not a single character
    - returns error when short or long collides with an existing entry
    # specification
  getopt.parse
    @ (s: spec, args: list[string]) -> result[parse_result, string]
    + returns a parse result with flag values and leftover positional arguments
    + supports "--long", "--long=value", "-s", "-svalue", combined short flags "-abc"
    + treats "--" as the end-of-flags marker
    - returns error on an unknown flag
    - returns error when a required argument is missing
    # parsing
    -> getopt.spec_add
  getopt.value
    @ (r: parse_result, name: string) -> optional[string]
    + returns the most recent value supplied for a flag, identified by short or long name
    - returns none when the flag was not present on the command line
    # access
  getopt.count
    @ (r: parse_result, name: string) -> i32
    + returns how many times the flag appeared
    # access
  getopt.positionals
    @ (r: parse_result) -> list[string]
    + returns the arguments that were not consumed as flags
    # access
