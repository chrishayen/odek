# Requirement: "a parser building toolkit"

A parser combinator toolkit for building grammars: primitive matchers, combinators, and a runner that produces typed results with error locations.

std: (all units exist)

parser_kit
  parser_kit.literal
    fn (expected: string) -> parser
    + returns a parser that matches the exact literal at the input position
    - fails when the input does not start with expected at the current position
    # primitives
  parser_kit.char_class
    fn (pred: fn(string) -> bool) -> parser
    + returns a parser that matches a single character satisfying pred
    - fails on end of input
    # primitives
  parser_kit.sequence
    fn (parsers: list[parser]) -> parser
    + runs parsers in order, producing a list of their results
    - fails at the first child failure without consuming the next ones
    # combinators
  parser_kit.choice
    fn (parsers: list[parser]) -> parser
    + tries each parser in order and returns the first that succeeds
    - fails with the longest-reaching child error when all fail
    # combinators
  parser_kit.many
    fn (p: parser) -> parser
    + applies p zero or more times, returning a list
    + succeeds with [] when p never matches
    # combinators
  parser_kit.many1
    fn (p: parser) -> parser
    + applies p one or more times
    - fails when p does not match at least once
    # combinators
  parser_kit.optional_match
    fn (p: parser) -> parser
    + returns a parser that succeeds with some(result) or none without consuming
    # combinators
  parser_kit.map_result
    fn (p: parser, f: fn(parse_result) -> parse_result) -> parser
    + transforms the value produced by p without altering matching behavior
    # combinators
  parser_kit.run
    fn (p: parser, input: string) -> result[parse_result, parse_error]
    + runs p on input and returns the produced value
    - returns error with line and column of the farthest failure
    - returns error when input is not fully consumed
    # runner
