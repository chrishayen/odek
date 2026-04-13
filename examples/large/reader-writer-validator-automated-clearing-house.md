# Requirement: "a reader, writer, and validator for Automated Clearing House (ACH) files"

ACH files are fixed-width records grouped into batches inside a file header/control envelope. The library parses records, emits them back, and enforces the standard's structural and checksum rules.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes contents to a file, creating or overwriting it
      - returns error when the path cannot be written
      # filesystem
  std.text
    std.text.pad_left
      @ (s: string, width: i32, pad: string) -> string
      + returns s padded on the left to width with pad
      # text
    std.text.pad_right
      @ (s: string, width: i32, pad: string) -> string
      + returns s padded on the right to width with pad
      # text

ach
  ach.parse_file_header
    @ (line: string) -> result[file_header, string]
    + parses a 94-character file header record
    - returns error when the record type is not "1"
    - returns error when the line length is wrong
    # parsing
  ach.parse_batch_header
    @ (line: string) -> result[batch_header, string]
    + parses a 94-character batch header record
    - returns error when the record type is not "5"
    # parsing
  ach.parse_entry_detail
    @ (line: string) -> result[entry_detail, string]
    + parses a 94-character entry detail record
    - returns error when the record type is not "6"
    - returns error when the amount field is not numeric
    # parsing
  ach.parse_batch_control
    @ (line: string) -> result[batch_control, string]
    + parses a 94-character batch control record
    - returns error when the record type is not "8"
    # parsing
  ach.parse_file_control
    @ (line: string) -> result[file_control, string]
    + parses a 94-character file control record
    - returns error when the record type is not "9"
    # parsing
  ach.read_file
    @ (path: string) -> result[ach_file, string]
    + reads and assembles a full ACH file from disk
    - returns error when any record fails to parse
    - returns error when records are not in the expected order
    # reading
    -> std.fs.read_all
  ach.compute_entry_hash
    @ (entries: list[entry_detail]) -> i64
    + returns the 10-digit entry hash (sum of routing numbers, mod 10^10)
    # checksums
  ach.validate
    @ (file: ach_file) -> result[void, list[string]]
    + returns ok when counts, totals, and entry hashes match the control records
    - returns error list when the batch entry count does not match
    - returns error list when debit or credit totals do not match
    - returns error list when the entry hash does not match
    # validation
  ach.format_file_header
    @ (header: file_header) -> string
    + returns the 94-character header line
    # formatting
    -> std.text.pad_left
    -> std.text.pad_right
  ach.format_entry_detail
    @ (entry: entry_detail) -> string
    + returns the 94-character entry detail line
    # formatting
    -> std.text.pad_left
    -> std.text.pad_right
  ach.write_file
    @ (path: string, file: ach_file) -> result[void, string]
    + writes a validated ACH file, padding to a multiple of 10 lines with filler
    - returns error when validation fails
    - returns error when the path cannot be written
    # writing
    -> std.fs.write_all
