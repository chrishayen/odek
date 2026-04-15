# Requirement: "convert a dash/dot/underscore/space separated string to camelCase"

A single function converting delimited input to camelCase.

std: (all units exist)

camelcase
  camelcase.to_camel
    fn (input: string) -> string
    + converts "foo-bar" to "fooBar"
    + converts "foo_bar_baz" to "fooBarBaz"
    + converts "foo.bar baz" to "fooBarBaz"
    + returns "" for empty input
    - collapses runs of delimiters rather than emitting empty segments
    ? first segment is lowercased, subsequent segments have their first letter uppercased
    # string_transform
