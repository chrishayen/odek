# Requirement: "a static analyzer that flags accidental logging of sensitive struct fields"

Walks a parsed source tree, tracks which struct fields are marked sensitive, and reports call sites where they flow into logging calls.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the path does not exist
      # filesystem
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns all file paths under root in depth-first order
      - returns error when root is not a directory
      # filesystem

leak_scanner
  leak_scanner.load_sources
    fn (root: string, extension: string) -> result[list[source_file], string]
    + returns each file under root with the given extension with its contents
    - returns error on any unreadable file
    # loading
    -> std.fs.walk
    -> std.fs.read_all
  leak_scanner.parse_file
    fn (src: source_file) -> result[ast_module, string]
    + returns a module AST
    - returns error on a parse failure with the offending position
    # parsing
  leak_scanner.collect_sensitive_fields
    fn (mods: list[ast_module], marker: string) -> map[string, list[string]]
    + returns struct name to list of field names that are tagged with marker
    + considers only fields whose declaration includes the marker annotation
    # indexing
  leak_scanner.find_logging_calls
    fn (mods: list[ast_module], logger_names: list[string]) -> list[call_site]
    + returns every call whose receiver or function name appears in logger_names
    # discovery
  leak_scanner.analyze_call
    fn (call: call_site, sensitive: map[string, list[string]]) -> list[finding]
    + returns a finding for each argument expression that reads a sensitive field
    + walks struct field accesses and member assignments through local aliases in the same function
    # analysis
  leak_scanner.scan
    fn (root: string, marker: string, logger_names: list[string]) -> result[list[finding], string]
    + returns every finding across the source tree
    - returns error when loading or parsing fails
    # orchestration
    -> leak_scanner.load_sources
    -> leak_scanner.parse_file
    -> leak_scanner.collect_sensitive_fields
    -> leak_scanner.find_logging_calls
    -> leak_scanner.analyze_call
  leak_scanner.format_finding
    fn (f: finding) -> string
    + returns a human-readable line with file, position, field path and logger call
    # reporting
