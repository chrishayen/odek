# Requirement: "a library for detecting outdated or insecure dependencies"

Parses a lockfile, queries a pluggable advisory and version source, and reports findings.

std
  std.http
    std.http.get_json
      fn (url: string) -> result[json_value, string]
      + performs an HTTP GET and parses the response body as JSON
      - returns error on non-2xx status or invalid JSON
      # http
  std.semver
    std.semver.parse
      fn (raw: string) -> result[semver, string]
      + parses a semantic version string
      - returns error on malformed input
      # versioning
    std.semver.compare
      fn (a: semver, b: semver) -> i32
      + returns -1, 0, or 1 comparing two versions
      # versioning

depcheck
  depcheck.parse_lockfile
    fn (source: string) -> result[list[dependency], string]
    + extracts name, version, and source for each locked dependency
    - returns error on malformed lockfile
    # parsing
  depcheck.fetch_latest
    fn (registry_url: string, name: string) -> result[semver, string]
    + queries a generic registry endpoint for a package's latest version
    - returns error when the package is not found
    # registry
    -> std.http.get_json
    -> std.semver.parse
  depcheck.fetch_advisories
    fn (advisory_url: string, name: string) -> result[list[advisory], string]
    + queries a pluggable advisory source for known vulnerabilities affecting a package
    - returns error on network failure
    # advisories
    -> std.http.get_json
  depcheck.classify
    fn (dep: dependency, latest: semver, advisories: list[advisory]) -> dependency_status
    + reports whether the dependency is up-to-date, outdated, or insecure
    + marks insecure when any advisory's affected range includes the dep's version
    # analysis
    -> std.semver.compare
  depcheck.scan
    fn (lockfile: string, registry_url: string, advisory_url: string) -> result[list[dependency_status], string]
    + orchestrates parse, fetch, and classify for every dependency in the lockfile
    - returns error when parsing fails
    # orchestration
