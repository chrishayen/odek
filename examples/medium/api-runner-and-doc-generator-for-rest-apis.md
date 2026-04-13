# Requirement: "an automated test runner and doc generator for REST APIs"

Loads a declarative spec of endpoints, runs each as a test case, and emits both a test report and reference docs.

std
  std.http
    std.http.request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an HTTP request and returns status, headers, and body
      - returns error on connection failure
      # http_client
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[yaml_value, string]
      + parses text into a generic YAML value
      - returns error on malformed input
      # serialization
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads a file fully as a string
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes contents to path
      # filesystem

api_runner
  api_runner.load_spec
    @ (path: string) -> result[api_spec, string]
    + loads an api spec file and returns the parsed structure
    - returns error when the file is missing or invalid
    # loading
    -> std.fs.read_all
    -> std.yaml.parse
  api_runner.run_case
    @ (base_url: string, case: endpoint_case) -> case_result
    + returns a passing result when status and body match the expectations
    - returns a failing result when status differs
    - returns a failing result when an expected body field is missing
    # execution
    -> std.http.request
  api_runner.run_all
    @ (base_url: string, spec: api_spec) -> run_report
    + runs every case in the spec and aggregates results
    # execution
  api_runner.render_test_report
    @ (report: run_report) -> string
    + returns a human-readable summary with per-case status
    # reporting
  api_runner.render_docs
    @ (spec: api_spec) -> string
    + returns markdown documentation for each endpoint with its request and response shapes
    # documentation
  api_runner.write_outputs
    @ (report_path: string, report: string, docs_path: string, docs: string) -> result[void, string]
    + writes the report and docs to disk
    # output
    -> std.fs.write_all
