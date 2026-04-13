# Requirement: "a documentation comment extractor for source definitions"

Given source code and a symbol name, extracts the doc comment block immediately preceding the symbol's definition.

std
  std.strings
    std.strings.split_lines
      @ (s: string) -> list[string]
      + splits s on newline boundaries without trailing empties
      # text
    std.strings.starts_with
      @ (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # text

doc_extract
  doc_extract.find_symbol_line
    @ (source: string, symbol: string) -> optional[i32]
    + returns the zero-based line number where symbol is defined
    - returns none when symbol is not defined
    # locating
    -> std.strings.split_lines
  doc_extract.collect_doc_comment
    @ (source: string, line_index: i32) -> string
    + walks backwards from line_index collecting contiguous comment lines
    + returns "" when no comment precedes the line
    # extraction
    -> std.strings.split_lines
    -> std.strings.starts_with
  doc_extract.doc_for
    @ (source: string, symbol: string) -> result[string, string]
    + returns the doc comment associated with symbol
    - returns error when symbol is not found
    # api
