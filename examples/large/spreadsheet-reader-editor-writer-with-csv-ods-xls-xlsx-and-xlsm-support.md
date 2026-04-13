# Requirement: "one API for reading, editing, and writing spreadsheet-family formats (csv, ods, xls, xlsx, xlsm)"

A uniform sheet model with per-format codecs.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the file contents
      - returns error when the path cannot be opened
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes the bytes to the file, replacing existing content
      # filesystem
  std.compression
    std.compression.zip_extract
      @ (archive: bytes) -> result[map[string, bytes], string]
      + returns a map from entry name to entry contents
      - returns error when the archive is malformed
      # compression
    std.compression.zip_pack
      @ (entries: map[string, bytes]) -> bytes
      + returns a zip archive containing the given entries
      # compression
  std.xml
    std.xml.parse
      @ (raw: string) -> result[xml_node, string]
      + returns the root node of the document
      - returns error on malformed xml
      # serialization
    std.xml.serialize
      @ (root: xml_node) -> string
      + returns the xml text for the tree
      # serialization

spreadsheet
  spreadsheet.new_book
    @ () -> book_state
    + returns an empty workbook containing no sheets
    # construction
  spreadsheet.add_sheet
    @ (book: book_state, name: string) -> result[book_state, string]
    + returns the book with a new empty sheet appended
    - returns error when a sheet with that name already exists
    # mutation
  spreadsheet.set_cell
    @ (book: book_state, sheet: string, row: i32, col: i32, value: cell_value) -> result[book_state, string]
    + stores the value at the given cell
    - returns error when the sheet does not exist
    # mutation
  spreadsheet.get_cell
    @ (book: book_state, sheet: string, row: i32, col: i32) -> result[cell_value, string]
    + returns the value at the given cell
    - returns error when the sheet or cell does not exist
    # access
  spreadsheet.load
    @ (path: string) -> result[book_state, string]
    + returns a book decoded from the file
    + dispatches to the codec matching the file extension
    - returns error when the extension is not supported
    - returns error when the file contents are malformed for that format
    # io
    -> std.fs.read_all
    -> std.compression.zip_extract
    -> std.xml.parse
  spreadsheet.save
    @ (book: book_state, path: string) -> result[void, string]
    + writes the book using the codec for the file extension
    - returns error when the extension is unsupported
    # io
    -> std.fs.write_all
    -> std.compression.zip_pack
    -> std.xml.serialize
  spreadsheet.decode_csv
    @ (raw: string) -> result[book_state, string]
    + returns a single-sheet book with rows and cells from the csv text
    - returns error on unterminated quoted fields
    # codec
  spreadsheet.encode_csv
    @ (book: book_state, sheet: string) -> result[string, string]
    + returns the sheet serialized as csv text
    - returns error when the sheet does not exist
    # codec
  spreadsheet.decode_xlsx
    @ (raw: bytes) -> result[book_state, string]
    + returns a book decoded from the xlsx container
    - returns error when required parts are missing
    # codec
    -> std.compression.zip_extract
    -> std.xml.parse
  spreadsheet.encode_xlsx
    @ (book: book_state) -> bytes
    + returns the workbook serialized as an xlsx container
    # codec
    -> std.compression.zip_pack
    -> std.xml.serialize
  spreadsheet.decode_ods
    @ (raw: bytes) -> result[book_state, string]
    + returns a book decoded from the ods container
    - returns error when the archive lacks content.xml
    # codec
    -> std.compression.zip_extract
    -> std.xml.parse
  spreadsheet.encode_ods
    @ (book: book_state) -> bytes
    + returns the workbook serialized as an ods container
    # codec
    -> std.compression.zip_pack
    -> std.xml.serialize
  spreadsheet.decode_xls
    @ (raw: bytes) -> result[book_state, string]
    + returns a book decoded from the legacy binary format
    - returns error when the record stream is truncated
    # codec
