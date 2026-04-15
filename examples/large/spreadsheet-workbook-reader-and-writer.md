# Requirement: "a library for reading and writing spreadsheet workbook files"

A spreadsheet workbook is a zip archive containing XML parts for a shared string table, per-sheet cell data, and a workbook index. The library parses these parts into an in-memory workbook and serializes them back.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the file contents as bytes
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file
      - returns error when the path cannot be written
      # filesystem
  std.zip
    std.zip.read_archive
      fn (data: bytes) -> result[map[string, bytes], string]
      + returns a map from member path to decompressed bytes
      - returns error when the archive is malformed
      # archive
    std.zip.write_archive
      fn (members: map[string, bytes]) -> result[bytes, string]
      + returns a zip archive containing the given members
      # archive
  std.xml
    std.xml.parse
      fn (raw: string) -> result[xml_node, string]
      + returns a parsed XML tree
      - returns error on malformed XML
      # xml
    std.xml.serialize
      fn (node: xml_node) -> string
      + serializes an XML tree back to a string
      # xml
  std.text
    std.text.column_label_to_index
      fn (label: string) -> result[i32, string]
      + converts a letter-based column label (e.g. "AA") to a zero-based index
      - returns error when the label contains non-letters
      # text
    std.text.index_to_column_label
      fn (index: i32) -> string
      + converts a zero-based column index to its letter label
      # text

spreadsheet
  spreadsheet.parse_shared_strings
    fn (part: string) -> result[list[string], string]
    + extracts the shared string table entries in order
    - returns error when the XML does not contain an sst root
    # parsing
    -> std.xml.parse
  spreadsheet.parse_sheet
    fn (part: string, shared_strings: list[string]) -> result[sheet_data, string]
    + parses sheet XML into a map from (row, col) to cell value
    - returns error when a cell references a shared-string index that is out of range
    # parsing
    -> std.xml.parse
    -> std.text.column_label_to_index
  spreadsheet.parse_workbook_index
    fn (part: string) -> result[list[sheet_entry], string]
    + returns the list of sheets (name and rel id) in the workbook
    - returns error when the XML root is not a workbook element
    # parsing
    -> std.xml.parse
  spreadsheet.read
    fn (path: string) -> result[workbook, string]
    + loads a workbook from disk and assembles all sheets
    - returns error when the archive cannot be opened
    - returns error when any required part is missing
    # reading
    -> std.fs.read_all
    -> std.zip.read_archive
  spreadsheet.get_cell
    fn (book: workbook, sheet_name: string, row: i32, col: i32) -> optional[cell_value]
    + returns the cell value at the given row and column
    - returns none when the sheet does not exist
    - returns none when the cell is empty
    # access
  spreadsheet.set_cell
    fn (book: workbook, sheet_name: string, row: i32, col: i32, value: cell_value) -> result[workbook, string]
    + returns a workbook with the cell updated
    - returns error when the sheet does not exist
    # mutation
  spreadsheet.build_shared_strings
    fn (book: workbook) -> tuple[list[string], workbook]
    + returns a deduplicated shared string table and a workbook whose string cells index into it
    # building
  spreadsheet.serialize_sheet
    fn (sheet: sheet_data) -> string
    + returns XML for a sheet, using column letter labels
    # serialization
    -> std.xml.serialize
    -> std.text.index_to_column_label
  spreadsheet.write
    fn (path: string, book: workbook) -> result[void, string]
    + serializes the workbook and writes it to disk as an archive
    - returns error when the path cannot be written
    # writing
    -> std.zip.write_archive
    -> std.fs.write_all
