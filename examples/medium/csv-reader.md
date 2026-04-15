# Requirement: "a CSV reader"

CSV parsing is requirement-specific, not a reusable subsystem — it stays in the project package. Three runes split the pipeline: document → row → field.

std: (all units exist)

csv
  csv.parse
    fn (input: string) -> result[list[list[string]], string]
    + parses a CSV document into rows of string fields
    + handles quoted fields containing commas and newlines
    + handles doubled double-quotes inside quoted fields
    + trailing blank lines are ignored
    - returns error when a quoted field is unterminated at EOF
    ? consecutive separators produce empty-string fields
    # parsing
  csv.parse_line
    fn (line: string) -> result[list[string], string]
    + splits one unquoted row on commas
    - returns error on malformed quoting within the line
    # parsing
  csv.unescape_field
    fn (field: string) -> string
    + strips surrounding double quotes and collapses doubled inner quotes to single ones
    + returns unquoted fields unchanged
    # parsing
