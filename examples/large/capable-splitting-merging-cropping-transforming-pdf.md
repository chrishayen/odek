# Requirement: "a PDF manipulation library for splitting, merging, cropping, and transforming pages"

A PDF reader that parses the object tree and a writer that emits a new PDF. Page-level operations are implemented on top.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes (and overwrites) the file
      # filesystem
  std.compression
    std.compression.zlib_inflate
      @ (data: bytes) -> result[bytes, string]
      + decompresses a zlib stream
      - returns error on corrupt input
      # compression
    std.compression.zlib_deflate
      @ (data: bytes) -> bytes
      + compresses with zlib
      # compression

pdf
  pdf.parse
    @ (raw: bytes) -> result[pdf_doc, string]
    + parses the xref table and object tree from a PDF byte stream
    - returns error when the header is not "%PDF-"
    - returns error when the xref table is corrupt
    # parsing
    -> std.compression.zlib_inflate
  pdf.load
    @ (path: string) -> result[pdf_doc, string]
    + reads and parses a PDF from disk
    # io
    -> std.fs.read_all
  pdf.page_count
    @ (doc: pdf_doc) -> i32
    + returns the number of pages
    # introspection
  pdf.get_page
    @ (doc: pdf_doc, index: i32) -> result[pdf_page, string]
    + returns the page at the given zero-based index
    - returns error when index is out of range
    # introspection
  pdf.set_page_box
    @ (page: pdf_page, llx: f64, lly: f64, urx: f64, ury: f64) -> pdf_page
    + sets the crop box of the page
    # transformation
  pdf.apply_matrix
    @ (page: pdf_page, a: f64, b: f64, c: f64, d: f64, e: f64, f: f64) -> pdf_page
    + prepends a transformation matrix to the page content stream
    ? used for rotation, scaling, and translation
    # transformation
  pdf.split
    @ (doc: pdf_doc, ranges: list[tuple[i32,i32]]) -> result[list[pdf_doc], string]
    + produces one output document per (start, end) range
    - returns error when any range is out of bounds
    # page_ops
  pdf.merge
    @ (docs: list[pdf_doc]) -> pdf_doc
    + concatenates documents into a single output, rewriting object ids
    # page_ops
  pdf.select_pages
    @ (doc: pdf_doc, indices: list[i32]) -> result[pdf_doc, string]
    + returns a new document containing only the selected pages in order
    - returns error when any index is out of bounds
    # page_ops
  pdf.serialize
    @ (doc: pdf_doc) -> bytes
    + writes the document to PDF byte format with a fresh xref table
    # writing
    -> std.compression.zlib_deflate
  pdf.save
    @ (doc: pdf_doc, path: string) -> result[void, string]
    + writes the PDF to disk
    # io
    -> std.fs.write_all
