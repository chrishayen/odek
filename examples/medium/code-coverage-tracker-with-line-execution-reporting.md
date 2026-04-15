# Requirement: "a code coverage tracking system that records which source lines were executed and produces a report"

Keeps a per-file bitmap of hit/miss counts and renders it as a summary.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads an entire file as text
      - returns error when the file is missing
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes contents to a file, creating or truncating it
      # filesystem

coverage
  coverage.new_profile
    fn () -> coverage_profile
    + creates an empty profile with no registered files
    # construction
  coverage.register_file
    fn (profile: coverage_profile, path: string, line_count: i32) -> coverage_profile
    + adds a file to the profile with an all-zero hit count vector
    # registration
  coverage.record_hit
    fn (profile: coverage_profile, path: string, line: i32) -> coverage_profile
    + increments the hit counter for a line
    - is a no-op when the file is not registered
    # recording
  coverage.file_percent
    fn (profile: coverage_profile, path: string) -> f64
    + returns the fraction of lines with at least one hit, in [0.0, 1.0]
    - returns 0.0 when the file is not registered
    # reporting
  coverage.total_percent
    fn (profile: coverage_profile) -> f64
    + returns the aggregate coverage across every registered file
    # reporting
  coverage.format_report
    fn (profile: coverage_profile) -> string
    + returns a human-readable summary with one line per file and a total
    # reporting
  coverage.save
    fn (profile: coverage_profile, path: string) -> result[void, string]
    + serializes the profile to disk
    # persistence
    -> std.fs.write_all
  coverage.load
    fn (path: string) -> result[coverage_profile, string]
    + reconstructs a profile from disk
    - returns error when the file is corrupt
    # persistence
    -> std.fs.read_all
