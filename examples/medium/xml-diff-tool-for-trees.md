# Requirement: "an XML diff tool that reports differences between two XML trees"

Parses both inputs, walks them structurally, and emits a list of typed differences.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file into memory
      - returns error when the file does not exist
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns entries in the directory
      # filesystem
  std.xml
    std.xml.parse
      @ (raw: bytes) -> result[xml_node, string]
      + parses an XML document into a tree
      - returns error on malformed XML
      # parsing

xml_diff
  xml_diff.compare_nodes
    @ (left: xml_node, right: xml_node) -> list[xml_difference]
    + returns tag, attribute, and text differences between two nodes recursively
    + reports insertions and deletions in element order
    # diffing
  xml_diff.compare_files
    @ (left_path: string, right_path: string) -> result[list[xml_difference], string]
    + parses both files and returns their structural differences
    - returns error when either file fails to parse
    # diffing
    -> std.fs.read_all
    -> std.xml.parse
  xml_diff.compare_folders
    @ (left_dir: string, right_dir: string) -> result[folder_diff, string]
    + walks both directories, matches files by relative path, and diffs XML content
    - returns error when either directory cannot be listed
    # diffing
    -> std.fs.list_dir
  xml_diff.format_difference
    @ (diff: xml_difference) -> string
    + returns a one-line human-readable description of a single difference
    # reporting
  xml_diff.format_report
    @ (diffs: list[xml_difference]) -> string
    + returns a multi-line report grouped by element path
    # reporting
