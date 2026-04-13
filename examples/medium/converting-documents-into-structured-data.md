# Requirement: "convert documents in various formats into a uniform structured representation"

A small conversion pipeline: detect input format, dispatch to a format-specific extractor, then normalize the result into a shared document model.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire contents of a file
      - returns error when the file cannot be read
      # filesystem
  std.html
    std.html.tokenize
      @ (source: string) -> result[list[html_token], string]
      + tokenizes text, start tags, end tags
      - returns error on unclosed tags
      # parsing

docling
  docling.detect_format
    @ (filename: string, header: bytes) -> document_format
    + identifies PDF from the "%PDF" magic
    + identifies HTML from a doctype or leading "<html" tag
    + identifies plain text, markdown, and docx from extension and magic
    + returns "unknown" when no rule matches
    # detection
  docling.extract_text_blocks
    @ (raw: bytes, fmt: document_format) -> result[list[text_block], string]
    + extracts heading, paragraph, and list blocks with reading order preserved
    + preserves table rows as structured cells
    - returns error on formats it cannot handle
    # extraction
    -> std.html.tokenize
  docling.build_document
    @ (blocks: list[text_block]) -> structured_document
    + groups consecutive list items into a single list block
    + groups consecutive table rows into a single table block
    + assigns a nesting level to headings based on their tag
    # normalization
  docling.convert_file
    @ (path: string) -> result[structured_document, string]
    + reads, detects, extracts, and normalizes in one call
    - returns error when any stage fails
    # pipeline
    -> std.fs.read_all
    -> docling.detect_format
    -> docling.extract_text_blocks
    -> docling.build_document
