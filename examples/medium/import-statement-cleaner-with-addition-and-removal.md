# Requirement: "a library that adds missing and removes unused import statements in source files"

Parses a source file, figures out which imported names are actually used in the body, drops the unused ones, and looks up missing names in an import index.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the entire file contents as a string
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to the file, creating or overwriting
      # filesystem

import_fixer
  import_fixer.parse_imports
    fn (source: string) -> result[parsed_source, string]
    + returns the import block span, the imported names, and the rest of the file
    - returns error when the import block is malformed
    # parsing
  import_fixer.collect_used_names
    fn (body: string) -> list[string]
    + returns every identifier prefix that appears before a '.' in the body
    ? this is the set of module aliases actually referenced
    # analysis
  import_fixer.drop_unused
    fn (parsed: parsed_source, used: list[string]) -> parsed_source
    + removes imports whose alias does not appear in used
    # rewriting
  import_fixer.add_missing
    fn (parsed: parsed_source, used: list[string], index: map[string, string]) -> parsed_source
    + adds an import for every used name that is not yet imported and that the index resolves
    ? index maps unqualified names to canonical import paths
    # rewriting
  import_fixer.render
    fn (parsed: parsed_source) -> string
    + returns the source with the import block rewritten and sorted alphabetically
    # rewriting
  import_fixer.fix_source
    fn (source: string, index: map[string, string]) -> result[string, string]
    + parses, drops unused imports, adds missing imports, and renders the result
    - returns error when parsing fails
    # orchestration
    -> import_fixer.parse_imports
    -> import_fixer.collect_used_names
    -> import_fixer.drop_unused
    -> import_fixer.add_missing
    -> import_fixer.render
  import_fixer.fix_file
    fn (path: string, index: map[string, string]) -> result[void, string]
    + reads, fixes, and writes the file in place
    - returns error when reading, fixing, or writing fails
    # orchestration
    -> std.fs.read_all
    -> import_fixer.fix_source
    -> std.fs.write_all
