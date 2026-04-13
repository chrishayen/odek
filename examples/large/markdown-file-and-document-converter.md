# Requirement: "a library that converts files and office documents to Markdown"

Dispatch by file type; each converter produces a canonical Markdown document. std hides format-specific extraction.

std
  std.fs
    std.fs.read_all_bytes
      @ (path: string) -> result[bytes, string]
      + reads a file into bytes
      - returns error when the file cannot be read
      # io
  std.docx
    std.docx.extract_text
      @ (data: bytes) -> result[list[docx_block], string]
      + returns an ordered list of paragraph and heading blocks with level and text
      - returns error when the document is not a valid office bundle
      # extraction
  std.pdf
    std.pdf.extract_text
      @ (data: bytes) -> result[list[pdf_page], string]
      + returns one text page per PDF page in document order
      - returns error when the PDF is encrypted
      # extraction
  std.xlsx
    std.xlsx.extract_sheets
      @ (data: bytes) -> result[list[sheet], string]
      + returns every worksheet as a name plus row-major cell grid
      # extraction
  std.html
    std.html.parse
      @ (source: string) -> result[html_node, string]
      + parses HTML into a DOM-like tree
      - returns error on unrecoverable markup
      # parsing

md
  md.detect_format
    @ (path: string, data: bytes) -> file_format
    + classifies input by extension and magic bytes
    ? returns "unknown" when no rule matches
    # detection
  md.convert
    @ (path: string) -> result[string, string]
    + dispatches to the matching converter by detected format
    - returns error when the format is unsupported
    # conversion
    -> std.fs.read_all_bytes
  md.convert_bytes
    @ (data: bytes, format: file_format) -> result[string, string]
    + converts an in-memory document to Markdown
    - returns error on extractor failure
    # conversion
  md.docx_to_markdown
    @ (data: bytes) -> result[string, string]
    + renders paragraphs as blank-line-separated text and headings with leading hashes
    # conversion_docx
    -> std.docx.extract_text
  md.pdf_to_markdown
    @ (data: bytes) -> result[string, string]
    + joins pages with a horizontal rule between them
    # conversion_pdf
    -> std.pdf.extract_text
  md.xlsx_to_markdown
    @ (data: bytes) -> result[string, string]
    + renders each sheet as a Markdown table with the sheet name as an H2
    # conversion_xlsx
    -> std.xlsx.extract_sheets
  md.html_to_markdown
    @ (source: string) -> result[string, string]
    + converts headings, paragraphs, lists, and links to Markdown equivalents
    - strips unsupported tags rather than erroring
    # conversion_html
    -> std.html.parse
  md.escape
    @ (text: string) -> string
    + escapes characters that would otherwise be interpreted by Markdown
    # escaping
  md.normalize
    @ (markdown: string) -> string
    + collapses runs of blank lines and trims trailing whitespace
    # postprocessing
