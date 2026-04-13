# Requirement: "a logfmt formatter"

Formats and parses log records as key=value pairs with the quoting rules of the logfmt convention.

std: (all units exist)

logfmt
  logfmt.format
    @ (fields: list[tuple[string, string]]) -> string
    + returns a space-separated sequence of key=value pairs in the given order
    + quotes values that contain spaces, equals signs or double quotes
    + escapes embedded double quotes and backslashes inside quoted values
    # encoding
  logfmt.parse
    @ (line: string) -> result[list[tuple[string, string]], string]
    + returns the pairs parsed from a logfmt line in order
    + accepts bare and quoted values, and values with escaped characters
    - returns error on an unterminated quoted value
    # decoding
