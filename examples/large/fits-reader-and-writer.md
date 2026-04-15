# Requirement: "a library for reading and writing FITS astronomy files"

Parses the FITS file format: a sequence of header-data units (HDUs), each with an ASCII header of 80-character cards followed by binary payload. The project layer wraps open/read/write operations for images and tables.

std
  std.fs
    std.fs.open_read
      fn (path: string) -> result[file_handle, string]
      + opens a file for reading
      - returns error when path does not exist
      # filesystem
    std.fs.open_write
      fn (path: string) -> result[file_handle, string]
      + opens a file for writing, truncating if present
      # filesystem
    std.fs.read_exact
      fn (handle: file_handle, n: i64) -> result[bytes, string]
      + reads exactly n bytes
      - returns error on short read
      # filesystem
    std.fs.write_all
      fn (handle: file_handle, data: bytes) -> result[void, string]
      + writes all bytes
      # filesystem
    std.fs.close
      fn (handle: file_handle) -> result[void, string]
      + closes the file
      # filesystem
  std.bytes
    std.bytes.read_big_endian_f64
      fn (data: bytes, offset: i64) -> f64
      + reads a big-endian IEEE 754 double
      # decoding
    std.bytes.read_big_endian_i32
      fn (data: bytes, offset: i64) -> i32
      + reads a big-endian signed 32-bit int
      # decoding

fits
  fits.parse_card
    fn (card: bytes) -> result[header_card, string]
    + parses an 80-byte card into key, value, and comment
    - returns error when card is not exactly 80 bytes
    - returns error when key area is malformed
    # header
  fits.parse_header
    fn (handle: file_handle) -> result[fits_header, string]
    + reads cards until an END card and returns the header
    + pads the read to the next 2880-byte block
    - returns error when END is missing
    # header
    -> std.fs.read_exact
  fits.data_unit_size
    fn (header: fits_header) -> result[i64, string]
    + returns the payload size in bytes from BITPIX and NAXIS keys
    - returns error when required keys are missing
    # header
  fits.read_data_unit
    fn (handle: file_handle, header: fits_header) -> result[bytes, string]
    + reads the data unit, padded to a 2880-byte boundary
    # io
    -> std.fs.read_exact
  fits.open
    fn (path: string) -> result[fits_file, string]
    + opens a FITS file and returns a fits_file with all HDUs parsed lazily
    - returns error on invalid primary header
    # api
    -> std.fs.open_read
  fits.next_hdu
    fn (file: fits_file) -> result[optional[hdu], string]
    + advances to the next HDU and returns it
    - returns none when EOF is reached
    # api
  fits.read_image_f64
    fn (file: fits_file, hdu_index: i32) -> result[list[f64], string]
    + returns the image pixels as f64 values
    - returns error when the HDU is not an image
    # image
    -> std.bytes.read_big_endian_f64
  fits.read_table_column
    fn (file: fits_file, hdu_index: i32, column: string) -> result[list[f64], string]
    + returns a numeric table column
    - returns error when the column does not exist
    # tables
    -> std.bytes.read_big_endian_i32
  fits.write_primary
    fn (path: string, header: fits_header, data: bytes) -> result[void, string]
    + writes a new FITS file with a single primary HDU
    - returns error on I/O failure
    # writing
    -> std.fs.open_write
    -> std.fs.write_all
    -> std.fs.close
  fits.format_header
    fn (header: fits_header) -> bytes
    + serializes a header to padded 2880-byte blocks of 80-byte cards
    # writing
  fits.close
    fn (file: fits_file) -> result[void, string]
    + closes the underlying handle
    # api
    -> std.fs.close
