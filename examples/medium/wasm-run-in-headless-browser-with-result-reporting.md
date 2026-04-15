# Requirement: "a library that runs compiled wasm test binaries in a headless browser and reports results"

Launches a headless browser, serves the wasm artifact, injects a small harness, and scrapes structured results back out. The browser driver is pluggable.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file fully into memory
      - returns error when the file does not exist or is unreadable
      # filesystem
  std.http
    std.http.serve_static
      fn (root: string, port: i32) -> result[server_handle, string]
      + starts a static file server rooted at the given directory
      # http
    std.http.stop_server
      fn (handle: server_handle) -> void
      + stops a running server
      # http

wasmrun
  wasmrun.new_runner
    fn (driver: browser_driver, port: i32) -> runner_state
    + creates a runner bound to a browser driver and a port for the static server
    # construction
  wasmrun.prepare_bundle
    fn (wasm_path: string, out_dir: string) -> result[void, string]
    + copies the wasm artifact and a minimal harness HTML/glue into the output directory
    - returns error when the wasm file cannot be read
    # bundling
    -> std.fs.read_all
  wasmrun.run
    fn (state: runner_state, bundle_dir: string, timeout_ms: i64) -> result[test_report, string]
    + serves the bundle, drives the browser to the harness page, waits for completion, and returns structured results
    - returns error on driver failure or timeout
    # execution
    -> std.http.serve_static
    -> std.http.stop_server
  wasmrun.parse_report
    fn (raw: string) -> result[test_report, string]
    + parses the harness-emitted JSON into a structured report of passes, failures, and log lines
    - returns error on malformed payload
    # reporting
