# Requirement: "a rich PDF document generation library"

A document builder that accumulates pages, text runs, and images, then serializes to the PDF file format. The real structure lives in std primitives for the wire format.

std
  std.encoding
    std.encoding.zlib_compress
      fn (data: bytes) -> bytes
      + returns a zlib-compressed stream of data
      # compression
    std.encoding.ascii85_encode
      fn (data: bytes) -> string
      + encodes bytes as ASCII85 text
      # encoding
  std.fs
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the named file
      - returns error when the parent directory does not exist
      # filesystem
  std.image
    std.image.decode_png
      fn (data: bytes) -> result[image_rgba, string]
      + decodes a PNG image into RGBA pixel data with width and height
      - returns error on invalid PNG
      # image

pdf
  pdf.new_document
    fn () -> pdf_doc
    + creates an empty PDF document with no pages
    # construction
  pdf.add_page
    fn (doc: pdf_doc, width_pt: f32, height_pt: f32) -> pdf_doc
    + appends a page with the given dimensions in points
    ? origin is bottom-left with y increasing upward
    # page_layout
  pdf.draw_text
    fn (doc: pdf_doc, page_index: i32, x: f32, y: f32, font: string, size: f32, text: string) -> result[pdf_doc, string]
    + places a text run on the named page at the given position
    - returns error when page_index is out of range
    # text_rendering
  pdf.draw_image
    fn (doc: pdf_doc, page_index: i32, x: f32, y: f32, w: f32, h: f32, data: bytes) -> result[pdf_doc, string]
    + places a decoded PNG image on the named page inside the given box
    - returns error when the image data cannot be decoded
    # image_embed
    -> std.image.decode_png
  pdf.draw_rect
    fn (doc: pdf_doc, page_index: i32, x: f32, y: f32, w: f32, h: f32, fill_gray: f32) -> result[pdf_doc, string]
    + fills a rectangle on the named page
    - returns error when page_index is out of range
    # shape_rendering
  pdf.render
    fn (doc: pdf_doc) -> bytes
    + serializes the document to a complete PDF byte stream with xref table
    # serialization
    -> std.encoding.zlib_compress
    -> std.encoding.ascii85_encode
  pdf.save
    fn (doc: pdf_doc, path: string) -> result[void, string]
    + renders and writes the document to the named file
    - returns error when the file cannot be written
    # persistence
    -> std.fs.write_all
