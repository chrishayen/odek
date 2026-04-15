# Requirement: "a utility that sorts import statements"

Parses the import block at the top of a source file, groups imports, and re-emits them sorted.

std: (all units exist)

isort
  isort.classify
    fn (import_line: string) -> import_group
    + returns "standard", "third_party", or "local" for the import
    ? classification uses a simple prefix list
    # classification
  isort.split_file
    fn (source: string) -> tuple[list[string], string]
    + returns (import lines, remainder) splitting at the first non-import statement
    ? blank lines and comments before the split are kept with the imports
    # parsing
  isort.sort_block
    fn (import_lines: list[string]) -> list[string]
    + returns the imports sorted lexicographically within each group
    + groups are emitted in the order standard, third_party, local
    + a blank line separates non-empty groups
    # sorting
  isort.apply
    fn (source: string) -> string
    + returns source with its import block replaced by a sorted version
    + input with no imports is returned unchanged
    # rewriting
