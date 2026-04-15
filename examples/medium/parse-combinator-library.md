# Requirement: "a parser combinator library"

Combinators produce parsers that consume a string position and either succeed with a value and a new position or fail with a message.

std: (all units exist)

parse
  parse.literal
    fn (text: string) -> parser
    + returns a parser that matches text exactly at the current position
    - fails when the input at the current position does not start with text
    # primitives
  parse.one_of
    fn (chars: string) -> parser
    + returns a parser that matches any single character in chars
    - fails when the next character is not in chars or the input is exhausted
    # primitives
  parse.sequence
    fn (parsers: list[parser]) -> parser
    + runs parsers in order and returns their results as a list
    - fails when any parser in the sequence fails
    # combinators
  parse.choice
    fn (parsers: list[parser]) -> parser
    + tries each parser in order and returns the first success
    - fails with the last error when every alternative fails
    # combinators
  parse.many
    fn (p: parser) -> parser
    + runs p repeatedly until it fails and returns the list of results
    ? always succeeds; the result list may be empty
    # combinators
  parse.map
    fn (p: parser, f: any) -> parser
    + transforms the result of p through a user-supplied function
    # combinators
  parse.run
    fn (p: parser, input: string) -> result[any, string]
    + runs p against input and returns the final value
    - returns error when p fails or when input has trailing unconsumed characters
    # execution
