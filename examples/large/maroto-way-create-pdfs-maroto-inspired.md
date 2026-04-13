# Requirement: "a grid-based PDF document builder"

Rows and columns like a CSS grid layout, with cells that hold text, images, barcodes, or nested components. Emits a byte-buffer PDF.

std
  std.pdf
    std.pdf.new_document
      @ (page_size: string) -> pdf_doc
      + creates an empty PDF document with the given page size
      # pdf
    std.pdf.draw_text
      @ (doc: pdf_doc, x: f64, y: f64, text: string, font: string, size: f64) -> pdf_doc
      + draws text at the given point
      # pdf
    std.pdf.draw_image
      @ (doc: pdf_doc, x: f64, y: f64, w: f64, h: f64, image: bytes) -> result[pdf_doc, string]
      + draws an embedded image at the given rectangle
      - returns error on unsupported image formats
      # pdf
    std.pdf.draw_line
      @ (doc: pdf_doc, x1: f64, y1: f64, x2: f64, y2: f64) -> pdf_doc
      + draws a stroked line
      # pdf
    std.pdf.new_page
      @ (doc: pdf_doc) -> pdf_doc
      + starts a new page
      # pdf
    std.pdf.to_bytes
      @ (doc: pdf_doc) -> bytes
      + serializes the document to a PDF byte buffer
      # pdf
  std.barcode
    std.barcode.encode_code128
      @ (value: string) -> bytes
      + encodes a Code 128 barcode as a monochrome image
      # barcode

pdfgrid
  pdfgrid.new
    @ (page_size: string) -> pdfgrid_state
    + creates a builder with a blank page and default 12-column grid
    # construction
    -> std.pdf.new_document
  pdfgrid.row
    @ (state: pdfgrid_state, height: f64, cells: list[cell_spec]) -> pdfgrid_state
    + adds a row whose cells fill 12 columns left to right
    + wraps to a new page when the row would overflow
    # layout
    -> std.pdf.new_page
  pdfgrid.text_cell
    @ (columns: i32, text: string, font: string, size: f64) -> cell_spec
    + constructs a text cell spanning the given columns
    # cells
  pdfgrid.image_cell
    @ (columns: i32, image: bytes) -> cell_spec
    + constructs an image cell spanning the given columns
    # cells
  pdfgrid.barcode_cell
    @ (columns: i32, value: string) -> cell_spec
    + constructs a Code 128 barcode cell spanning the given columns
    # cells
    -> std.barcode.encode_code128
  pdfgrid.line
    @ (state: pdfgrid_state, thickness: f64) -> pdfgrid_state
    + draws a horizontal separator across the current page width
    # layout
    -> std.pdf.draw_line
  pdfgrid.render_cell
    @ (state: pdfgrid_state, cell: cell_spec, x: f64, y: f64, w: f64, h: f64) -> pdfgrid_state
    + renders a cell at the given rectangle
    # rendering
    -> std.pdf.draw_text
    -> std.pdf.draw_image
  pdfgrid.page_break
    @ (state: pdfgrid_state) -> pdfgrid_state
    + starts a fresh page and resets the vertical cursor
    # layout
    -> std.pdf.new_page
  pdfgrid.build
    @ (state: pdfgrid_state) -> bytes
    + finalizes the document and returns PDF bytes
    # build
    -> std.pdf.to_bytes
