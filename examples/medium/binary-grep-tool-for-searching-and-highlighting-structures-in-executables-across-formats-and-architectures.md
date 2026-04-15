# Requirement: "a tool library that searches and highlights structures in executable binaries across formats and architectures"

Detect the binary format, parse headers and symbol tables, and return a structured listing with highlighted regions matching a user query.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file as bytes
      - returns error when the path does not exist
      # io

binary_grep
  binary_grep.detect_format
    fn (data: bytes) -> result[string, string]
    + returns "elf", "mach-o", "pe", or "coff" based on magic bytes
    - returns error when no known magic is recognized
    # detection
  binary_grep.parse_header
    fn (data: bytes, format: string) -> result[binary_header, string]
    + parses the format-specific header into a uniform structure
    - returns error on truncated or invalid header
    # parsing
  binary_grep.list_sections
    fn (header: binary_header) -> list[section_info]
    + returns one entry per section with name, offset, size, and flags
    # inspection
  binary_grep.list_symbols
    fn (data: bytes, header: binary_header) -> result[list[symbol_info], string]
    + returns every exported and imported symbol with its address
    - returns error when the symbol table is malformed
    # inspection
  binary_grep.find_matches
    fn (symbols: list[symbol_info], query: string) -> list[symbol_info]
    + returns symbols whose name contains query (case-insensitive)
    # query
  binary_grep.highlight
    fn (sections: list[section_info], matches: list[symbol_info]) -> list[highlight_span]
    + returns byte-offset spans suitable for colorized rendering
    # presentation
  binary_grep.open_path
    fn (path: string) -> result[binary_header, string]
    + reads the file and parses its header
    - returns error when the path is unreadable or format is unknown
    # convenience
    -> std.fs.read_all
