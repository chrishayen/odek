# Requirement: "a PDF reading and editing library"

The library parses PDF documents, exposes structural access to pages and metadata, and supports simple edits like page deletion and metadata updates.

std
  std.io
    std.io.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file into memory
      - returns error on missing or unreadable files
      # io
    std.io.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file
      # io

pdf
  pdf.parse
    @ (data: bytes) -> result[pdf_doc, string]
    + returns a parsed document with cross-reference table, pages, and metadata
    - returns error when the header is not "%PDF-"
    - returns error when the xref table is unreadable
    # parsing
  pdf.load
    @ (path: string) -> result[pdf_doc, string]
    + reads and parses a PDF file
    # parsing
    -> std.io.read_all
  pdf.page_count
    @ (doc: pdf_doc) -> i32
    + returns the number of pages
    # inspection
  pdf.get_metadata
    @ (doc: pdf_doc) -> map[string, string]
    + returns document-level metadata entries like Title, Author
    # inspection
  pdf.set_metadata
    @ (doc: pdf_doc, key: string, value: string) -> pdf_doc
    + updates a metadata entry and returns the updated document
    # editing
  pdf.delete_page
    @ (doc: pdf_doc, index: i32) -> result[pdf_doc, string]
    + removes the page at the given zero-based index
    - returns error when the index is out of range
    # editing
  pdf.serialize
    @ (doc: pdf_doc) -> bytes
    + produces a valid PDF byte stream from the document
    # serialization
  pdf.save
    @ (doc: pdf_doc, path: string) -> result[void, string]
    + writes the document to disk
    # serialization
    -> std.io.write_all
