# Requirement: "an Excel-compatible spreadsheet reader and writer"

A workbook model with decoding and encoding for the zipped-XML spreadsheet format, plus helpers for cells, styles, and formulas.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the file contents
      - returns error when the path cannot be opened
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes the bytes to the file
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
      + returns the root node
      - returns error on malformed xml
      # serialization
    std.xml.serialize
      @ (root: xml_node) -> string
      + returns the xml text for the tree
      # serialization

xlsx
  xlsx.new_book
    @ () -> book_state
    + returns an empty workbook with no sheets
    # construction
  xlsx.add_sheet
    @ (book: book_state, name: string) -> result[book_state, string]
    + returns the book with the new sheet appended
    - returns error when the sheet name is duplicated
    # mutation
  xlsx.set_cell_number
    @ (book: book_state, sheet: string, row: i32, col: i32, value: f64) -> result[book_state, string]
    + returns the book with the numeric cell set
    - returns error when the sheet does not exist
    # mutation
  xlsx.set_cell_string
    @ (book: book_state, sheet: string, row: i32, col: i32, value: string) -> result[book_state, string]
    + returns the book with the string cell set
    - returns error when the sheet does not exist
    # mutation
  xlsx.set_cell_formula
    @ (book: book_state, sheet: string, row: i32, col: i32, formula: string) -> result[book_state, string]
    + returns the book with the formula stored for the cell
    - returns error when the sheet does not exist
    # mutation
  xlsx.get_cell
    @ (book: book_state, sheet: string, row: i32, col: i32) -> result[cell_value, string]
    + returns the value at the given cell
    - returns error when the sheet or cell does not exist
    # access
  xlsx.col_letter
    @ (index: i32) -> string
    + returns the one-based column letter (A, B, ..., Z, AA, AB, ...)
    # addressing
  xlsx.parse_cell_ref
    @ (ref: string) -> result[tuple[i32, i32], string]
    + returns the (row, col) pair for a cell reference such as "B12"
    - returns error when the reference is malformed
    # addressing
  xlsx.decode
    @ (data: bytes) -> result[book_state, string]
    + returns a workbook decoded from the container
    - returns error when required parts are missing
    - returns error when any sheet xml is malformed
    # codec
    -> std.compression.zip_extract
    -> std.xml.parse
  xlsx.encode
    @ (book: book_state) -> bytes
    + returns the workbook serialized as a container
    # codec
    -> std.compression.zip_pack
    -> std.xml.serialize
  xlsx.load
    @ (path: string) -> result[book_state, string]
    + returns a workbook read from the file
    - returns error when the file is unreadable or malformed
    # io
    -> std.fs.read_all
  xlsx.save
    @ (book: book_state, path: string) -> result[void, string]
    + writes the workbook to the file
    # io
    -> std.fs.write_all
  xlsx.evaluate_formulas
    @ (book: book_state) -> result[book_state, string]
    + returns the book with every formula cell replaced by its computed value
    - returns error when a formula references a missing cell
    - returns error on unknown functions
    # evaluation
