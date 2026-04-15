# Requirement: "a csv record encoder and decoder"

Encodes and decodes CSV records as string rows with proper quoting and escaping per RFC 4180.

std
  std.text
    std.text.split
      fn (s: string, sep: string) -> list[string]
      + splits a string on a separator
      # text
    std.text.contains
      fn (s: string, sub: string) -> bool
      + returns true when sub appears anywhere in s
      # text

csv
  csv.encode_row
    fn (fields: list[string]) -> string
    + joins fields with commas, quoting fields containing commas, quotes, or newlines
    + doubles embedded quotes inside quoted fields
    # encoding
    -> std.text.contains
  csv.encode_all
    fn (rows: list[list[string]]) -> string
    + emits multiple rows separated by CRLF
    # encoding
  csv.decode_row
    fn (line: string) -> result[list[string], string]
    + parses a single CSV line honoring quoted fields and escaped quotes
    - returns error on an unterminated quoted field
    # decoding
  csv.decode_all
    fn (input: string) -> result[list[list[string]], string]
    + parses multi-row input, handling quoted fields that span lines
    - returns error on an unterminated quoted field
    # decoding
    -> std.text.split
