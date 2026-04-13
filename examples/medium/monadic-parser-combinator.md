# Requirement: "a monadic-style parser combinator"

A parser is a function from input to a result plus remaining input. Combinators build larger parsers from smaller ones.

std: (all units exist)

combinator
  combinator.run
    @ (p: parser, input: string) -> result[tuple[parser_value, string], string]
    + runs p against input and returns (value, rest) on success
    - returns error with position and message on failure
    # entry_point
  combinator.pure
    @ (value: parser_value) -> parser
    + returns a parser that yields value without consuming input
    # primitive
  combinator.char
    @ (c: i32) -> parser
    + matches exactly one codepoint equal to c
    - fails when the next codepoint differs or input is empty
    # primitive
  combinator.satisfy
    @ (predicate: char_predicate) -> parser
    + matches a single codepoint for which predicate returns true
    # primitive
  combinator.map
    @ (p: parser, f: value_map) -> parser
    + applies f to the result of p
    # combinator
  combinator.bind
    @ (p: parser, f: value_to_parser) -> parser
    + sequences p then f(value_of_p), threading input
    # combinator
  combinator.alt
    @ (left: parser, right: parser) -> parser
    + tries left first; on failure without consumption tries right
    # combinator
  combinator.many
    @ (p: parser) -> parser
    + repeatedly runs p, returning a list of values, possibly empty
    # combinator
  combinator.seq
    @ (ps: list[parser]) -> parser
    + runs each parser in order, collecting the values
    - fails when any child parser fails
    # combinator
