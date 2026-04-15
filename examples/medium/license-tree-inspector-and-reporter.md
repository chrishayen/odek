# Requirement: "a library that inspects a dependency tree and reports the licenses of every dependency"

Parses a manifest, walks the declared dependencies, reads each one's metadata, and produces a report. Allow/deny list filtering is included.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error when the path does not exist
      # filesystem
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns the names of entries within the directory
      - returns error when the path is not a directory
      # filesystem
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a flat JSON object
      - returns error on malformed input
      # serialization

license_report
  license_report.load_manifest
    fn (manifest_path: string) -> result[list[string], string]
    + returns the names of declared dependencies
    - returns error on malformed manifest
    # parsing
    -> std.fs.read_all
    -> std.json.parse_object
  license_report.read_dependency_license
    fn (module_root: string) -> result[license_record, string]
    + returns the license name and source file path for a single module
    - returns error when no license metadata is discoverable under the root
    # discovery
    -> std.fs.read_all
    -> std.fs.list_dir
  license_report.collect
    fn (manifest_path: string, modules_root: string) -> result[list[license_record], string]
    + returns one license record per declared dependency
    - returns error when the manifest cannot be parsed
    ? dependencies that cannot be resolved are reported with license "unknown"
    # orchestration
  license_report.filter_disallowed
    fn (records: list[license_record], allowed: list[string]) -> list[license_record]
    + returns only records whose license is not in the allowed list
    + returns empty list when every license is allowed
    # filtering
  license_report.format_table
    fn (records: list[license_record]) -> string
    + renders a plain-text table with dependency, license, and source columns
    # rendering
  license_report.format_csv
    fn (records: list[license_record]) -> string
    + renders records as CSV with a header row
    # rendering
