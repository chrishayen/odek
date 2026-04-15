# Requirement: "a behavior-driven test framework with a browser dashboard and live reload"

Test registry, runner, a minimal HTTP server, and a file-watcher loop. Networking and watching live in std.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns full file contents
      - returns error when path cannot be read
      # filesystem
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns absolute paths of all files under root
      - returns error when root cannot be opened
      # filesystem
    std.fs.mtime
      fn (path: string) -> result[i64, string]
      + returns the modification time as unix seconds
      - returns error when path cannot be stat'd
      # filesystem
  std.http
    std.http.serve
      fn (addr: string, handler: fn(http_request) -> http_response) -> result[void, string]
      + binds addr and dispatches each request to handler until shutdown
      - returns error when bind fails
      # networking
    std.http.response_ok
      fn (body: string, content_type: string) -> http_response
      + builds a 200 response with the given body and content type
      # networking
  std.json
    std.json.encode
      fn (value: json_value) -> string
      + encodes any json value as a compact string
      # serialization

bdd_web
  bdd_web.new_registry
    fn () -> bdd_registry
    + returns an empty test registry
    # construction
  bdd_web.register
    fn (reg: bdd_registry, name: string, body: fn() -> test_result) -> void
    + adds a named test case to the registry
    # registration
  bdd_web.run_all
    fn (reg: bdd_registry) -> run_report
    + runs every registered test and returns pass/fail counts and per-test results
    + captures panics as failed results with the panic message
    # execution
  bdd_web.report_to_json
    fn (report: run_report) -> string
    + encodes a run_report as JSON for the dashboard
    # serialization
    -> std.json.encode
  bdd_web.watch_sources
    fn (root: string, on_change: fn() -> void) -> result[void, string]
    + polls files under root and invokes on_change when any mtime advances
    - returns error when root cannot be walked
    # watching
    -> std.fs.walk
    -> std.fs.mtime
  bdd_web.serve_dashboard
    fn (addr: string, reg: bdd_registry) -> result[void, string]
    + serves an HTML page and a /results endpoint that returns the latest report
    - returns error when the address cannot be bound
    # dashboard
    -> std.http.serve
    -> std.http.response_ok
    -> std.fs.read_all
  bdd_web.run_with_live_reload
    fn (addr: string, root: string, reg: bdd_registry) -> result[void, string]
    + starts the dashboard and re-runs all tests whenever any source under root changes
    - returns error when the dashboard cannot start
    # orchestration
