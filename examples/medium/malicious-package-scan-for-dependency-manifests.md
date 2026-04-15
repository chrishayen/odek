# Requirement: "a library that scans dependency manifests for known-malicious packages"

Reads a manifest, resolves each declared dependency, looks each up against a vulnerability/malice database, and emits a report.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the full file contents as a string
      - returns error when the file does not exist or is unreadable
      # filesystem
  std.io
    std.io.http_get_json
      fn (url: string) -> result[json_value, string]
      + performs an HTTP GET and parses the body as JSON
      - returns error on non-2xx status or invalid JSON
      # http
  std.encoding
    std.encoding.json_decode
      fn (text: string) -> result[json_value, string]
      + parses JSON text into a JSON value
      - returns error on malformed JSON
      # serialization

malice_scan
  malice_scan.parse_manifest
    fn (text: string) -> result[list[package_ref], string]
    + returns the list of (name, version) pairs declared in the manifest
    - returns error when the manifest text is malformed
    # manifest_parsing
    -> std.encoding.json_decode
  malice_scan.lookup_package
    fn (pkg: package_ref, db_url: string) -> result[list[advisory], string]
    + queries the advisory database and returns matching advisories
    + returns an empty list when the package is clean
    - returns error when the lookup fails
    # advisory_lookup
    -> std.io.http_get_json
  malice_scan.classify
    fn (advisory: advisory) -> severity
    + maps a raw advisory record to a severity level (none/low/high/critical)
    # risk_classification
  malice_scan.scan_manifest
    fn (path: string, db_url: string) -> result[scan_report, string]
    + reads, parses, looks up each package, and returns a per-package report
    - returns error when reading or parsing the manifest fails
    # orchestration
    -> std.fs.read_all
  malice_scan.has_critical
    fn (report: scan_report) -> bool
    + returns true when any dependency has a critical-severity advisory
    # report_inspection
