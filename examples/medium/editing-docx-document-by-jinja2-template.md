# Requirement: "a templated word document editor"

Loads a word-processing document, substitutes placeholders with values, and writes the result.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads file bytes
      - returns error when missing
      # fs
    std.fs.write_all
      @ (path: string, contents: bytes) -> result[void, string]
      + writes bytes to path
      - returns error on I/O failure
      # fs
  std.compression
    std.compression.zip_read
      @ (archive: bytes, entry: string) -> result[bytes, string]
      + returns the uncompressed bytes of the named entry
      - returns error when entry is absent
      # archive
    std.compression.zip_write
      @ (archive: bytes, entry: string, contents: bytes) -> result[bytes, string]
      + replaces the named entry inside the archive
      - returns error on malformed archive
      # archive
  std.xml
    std.xml.parse
      @ (raw: string) -> result[xml_node, string]
      + parses XML into a tree
      - returns error on malformed input
      # xml
    std.xml.serialize
      @ (node: xml_node) -> string
      + serializes a tree back to XML text
      # xml

docx
  docx.load
    @ (path: string) -> result[document, string]
    + reads a word document and extracts the main body XML
    - returns error when the archive lacks a document body
    # ingest
    -> std.fs.read_all
    -> std.compression.zip_read
    -> std.xml.parse
  docx.apply_template
    @ (doc: document, values: map[string, string]) -> result[document, string]
    + substitutes placeholders like {{name}} with provided values
    - returns error when a placeholder has no matching value
    # templating
  docx.render_text
    @ (doc: document) -> string
    + returns the plain text of the document for inspection
    # inspection
  docx.save
    @ (doc: document, path: string) -> result[void, string]
    + writes the edited document back to disk
    - returns error on I/O failure
    # persistence
    -> std.xml.serialize
    -> std.compression.zip_write
    -> std.fs.write_all
