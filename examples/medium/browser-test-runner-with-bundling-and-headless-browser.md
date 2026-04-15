# Requirement: "a browser test runner that bundles test sources and drives a headless browser"

Bundle test sources, serve them to a headless browser, and collect structured test results.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error when path does not exist
      # filesystem
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns filenames under a directory
      - returns error when path is not a directory
      # filesystem
  std.net
    std.net.serve_static
      fn (root: string, port: u16) -> result[server_handle, string]
      + starts an HTTP server that maps URLs to files under root
      - returns error when port is already bound
      # http
    std.net.stop_server
      fn (handle: server_handle) -> void
      + stops a running HTTP server
      # http

browser_test_runner
  browser_test_runner.bundle_tests
    fn (entry_files: list[string]) -> result[string, string]
    + concatenates test sources into a single script addressable by the browser
    - returns error when an entry file is missing
    # bundling
    -> std.fs.read_all
  browser_test_runner.launch_browser
    fn (url: string) -> result[browser_session, string]
    + spawns a headless browser pointing at url
    - returns error when no browser binary is found
    # browser_control
  browser_test_runner.collect_results
    fn (session: browser_session, timeout_ms: i64) -> result[test_report, string]
    + reads structured results the page posted back via a results endpoint
    - returns error when the session exits before posting results
    - returns error on timeout
    # result_collection
  browser_test_runner.run
    fn (entry_files: list[string], port: u16) -> result[test_report, string]
    + bundles, serves, launches, and returns a structured pass/fail report
    - returns error when any stage fails
    # orchestration
    -> std.net.serve_static
    -> std.net.stop_server
