# Requirement: "a service that generates pdf documents from a json description"

Parses a structured document spec and emits a pdf byte stream.

std
  std.json
    std.json.parse
      @ (raw: bytes) -> result[json_value, string]
      + parses raw bytes as a json document
      - returns error on invalid json
      # serialization
  std.io
    std.io.read_file
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when the path does not exist
      # filesystem

pdf_service
  pdf_service.parse_spec
    @ (raw: bytes) -> result[document_spec, string]
    + parses a json object with title, pages, and page elements into a typed spec
    - returns error when a required field is missing
    - returns error on an unknown element kind
    # parsing
    -> std.json.parse
  pdf_service.layout
    @ (spec: document_spec, page_width: f64, page_height: f64) -> laid_out_document
    + assigns coordinates to every text and image element within page bounds
    + wraps text elements whose content overflows the page width
    # layout
  pdf_service.render
    @ (laid_out: laid_out_document) -> bytes
    + emits a pdf byte stream with the specified pages and elements
    + returns a valid empty-page pdf when the document has no elements
    # rendering
  pdf_service.load_image
    @ (path: string) -> result[image_resource, string]
    + loads an image from disk into a pdf-ready resource
    - returns error when the path does not exist
    # images
    -> std.io.read_file
