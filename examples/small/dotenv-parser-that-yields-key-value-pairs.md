# Requirement: "a .env file parser that yields key/value pairs"

Parses the conventional KEY=value line format with comments and quoted values. Filesystem I/O is left to the caller.

std: (all units exist)

dotenv
  dotenv.parse
    fn (source: string) -> result[map[string, string], string]
    + parses lines of the form KEY=value into a map
    + strips surrounding single or double quotes from values
    + ignores blank lines and lines beginning with '#'
    + trims surrounding whitespace from keys and unquoted values
    - returns error on a line that has no '=' separator
    - returns error on a key containing whitespace
    # parsing
  dotenv.serialize
    fn (vars: map[string, string]) -> string
    + emits KEY=value lines, quoting values that contain whitespace or '#'
    # serialization
